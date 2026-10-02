#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-morpho-oracle-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks the registry's Morpho oracles against the band program on the dev node: each prices the Stock Token the band
# names for its asset in USDG, at the band's low edge with Morpho's 36 + 6 - 18 decimals, or has no answer while the
# band is halted; it reports the band's halt for its asset, and only the registry's owner may re-point it. Then, for
# each asset nothing else halts, it puts Chainlink's round past the age the band accepts and lets the 24/7 leg go
# stale: the band has no live leg and the oracle no answer, until a fresh RedStone package is written. Prints each
# oracle's answer and the L2 gas of reading it.
set -euo pipefail

rpc=$1 key=$2 registry=$3
root=$(cd "$(dirname "$0")/../.." && pwd)
band=$(jq -r .tapehouse.Band "$registry")
owner=$(jq -r .tapehouse.Owner "$registry")
usdg=$(jq -r .tokens.USDG "$registry")
me=$(cast wallet address --private-key "$key")
node_interface=0x00000000000000000000000000000000000000C8
zero_id=0x0000000000000000000000000000000000000000000000000000000000000000
chainlink_max_age_s=86460
live_max_age_ms=120000
no_answer=$(cast sig "NoAnswer(bytes32)")
package_not_newer=$(cast sig "PackageNotNewer(bytes32,uint64,uint64)")
unauthorized=$(cast sig "OwnableUnauthorizedAccount(address)")
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$key" "$@" > /dev/null; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
same() { [ "$(lower "$1")" = "$(lower "$2")" ]; }
words() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1 | tr '\n' ' '; }
block_time() { cast block --rpc-url "$rpc" latest -f timestamp; }
quote() { words "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$1"; }
l2_gas() {
  cast call --rpc-url "$rpc" $node_interface "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" \
    "$1" false "$(cast calldata "price()")" --from "$me" | sed -n 1,2p | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}
expect_no_answer() {
  local out
  if out=$(cast call --rpc-url "$rpc" "$1" "price()(uint256)" 2>&1); then fail "$2, yet its oracle answered $out"; fi
  grep -q "${no_answer#0x}" <<<"$out" || fail "expected NoAnswer from $1, got $out"
}

oracles=$(jq -r '.morphoOracles // {} | keys[]' "$registry")
[ -n "$oracles" ] || fail "the registry has no Morpho oracles"
for asset in $oracles; do
  oracle=$(jq -r --arg asset "$asset" '.morphoOracles[$asset]' "$registry")
  symbol=$(cast format-bytes32-string "$asset")
  token=$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "$symbol" | sed -n 4p)
  same "$(read_ "$oracle" "band()(address)")" "$band" || fail "$asset's oracle does not read the registry's band"
  [ "$(read_ "$oracle" "symbol()(bytes32)")" = "$symbol" ] || fail "$asset's oracle prices another asset"
  same "$(read_ "$oracle" "collateralToken()(address)")" "$token" &&
    same "$token" "$(jq -r --arg asset "$asset" '.tokens[$asset]' "$registry")" ||
    fail "$asset's oracle does not price the Stock Token the band and the registry name"
  same "$(read_ "$oracle" "loanToken()(address)")" "$usdg" || fail "$asset's oracle does not price in the registry's USDG"
  [ "$(read_ "$oracle" "scaleFactor()(uint256)")" = 10000000000000000 ] || fail "$asset's oracle does not scale by 1e16"
  same "$(read_ "$oracle" "owner()(address)")" "$owner" || fail "$asset's oracle is not the registry owner's"
  [ "$(cast call --rpc-url "$rpc" "$oracle" "halt()(bool,uint64,uint64,bool)")" = \
    "$(cast call --rpc-url "$rpc" "$band" "halt(bytes32)(bool,uint64,uint64,bool)" "$symbol")" ] || fail "$asset's oracle reports another halt"
  read -r state _ _ _ low _ <<<"$(quote "$symbol")"
  if [ "$state" = 0 ]; then
    expect_no_answer "$oracle" "$asset's band is halted"
    echo "morpho oracle $oracle: $asset's band is halted, and the oracle has no answer"
  else
    price=$(read_ "$oracle" "price()(uint256)")
    [ "$price" = "$(python3 -c "print($low * 10**16)")" ] || fail "$asset's oracle answered $price, the band's low edge is $low"
    echo "morpho oracle $oracle: $asset at $price, the band's low edge $low in state $state; price $(l2_gas "$oracle") L2 gas"
  fi
  if out=$(cast call --rpc-url "$rpc" --from 0x000000000000000000000000000000000000dEaD "$oracle" "setBand(address)" "$band" 2>&1); then
    fail "a stranger may re-point $asset's oracle"
  fi
  grep -q "${unauthorized#0x}" <<<"$out" || fail "expected OwnableUnauthorizedAccount, got $out"
  cast call --rpc-url "$rpc" --from "$owner" "$oracle" "setBand(address)" "$band" > /dev/null || fail "the owner may not re-point $asset's oracle"
done

staled=0
for asset in $oracles; do
  oracle=$(jq -r --arg asset "$asset" '.morphoOracles[$asset]' "$registry")
  symbol=$(cast format-bytes32-string "$asset")
  read -r signed _ _ paused <<<"$(words "$band" "halt(bytes32)(bool,uint64,uint64,bool)" "$symbol")"
  read -r action _ <<<"$(words "$band" "corporateAction(bytes32)(uint8,uint64,uint128,uint128)" "$symbol")"
  [ "$signed $paused" = "false false" ] && [ "$action" != 2 ] || continue
  read -r feed live_id index_id _ <<<"$(words "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "$symbol")"
  feed_id=$live_id
  [ "$feed_id" != $zero_id ] || feed_id=$index_id
  package=$(cast parse-bytes32-string "$feed_id")
  read -r _ answer _ updated _ <<<"$(words "$feed" "latestRoundData()(uint80,int256,uint256,uint256,uint80)")"
  send "$feed" "setRound(int256,uint256)" "$answer" $(($(block_time) - chainlink_max_age_s - 1))
  read -r _ package_ms _ <<<"$(words "$band" "price(bytes32)(uint256,uint64,uint64)" "$feed_id")"
  until [ $(($(block_time) * 1000 - package_ms)) -gt $((live_max_age_ms + 1000)) ]; do
    sleep 5
    send --value 0 "$me"
  done
  [ "$(quote "$symbol")" = "0 0 0 0 0 0 " ] || fail "$asset's band has a price with no live leg: $(quote "$symbol")"
  expect_no_answer "$oracle" "$asset's band has no live leg"
  stale_s=$(($(block_time) - package_ms / 1000))
  attempts=0
  until payload=$(python3 "$root/stylus/scripts/redstone-payload.py" "$package") &&
    cast call --rpc-url "$rpc" "$band" "writePrices(bytes32[],bytes)" "[$feed_id]" "$payload" > /dev/null 2>&1; do
    attempts=$((attempts + 1))
    [ "$attempts" -le 20 ] || fail "no newer $package package within 200 s"
    sleep 10
  done
  send "$band" "writePrices(bytes32[],bytes)" "[$feed_id]" "$payload"
  read -r state _ _ _ low _ <<<"$(quote "$symbol")"
  [ "$state" != 0 ] || fail "$asset's band has no price after a fresh $package package"
  price=$(read_ "$oracle" "price()(uint256)")
  [ "$price" = "$(python3 -c "print($low * 10**16)")" ] || fail "$asset's oracle answered $price, the band's low edge is $low"
  if out=$(cast call --rpc-url "$rpc" "$band" "writePrices(bytes32[],bytes)" "[$feed_id]" "$payload" 2>&1); then
    fail "the same $package package was written twice"
  fi
  grep -q "${package_not_newer#0x}" <<<"$out" || fail "expected PackageNotNewer, got $out"
  send "$feed" "setRound(int256,uint256)" "$answer" "$updated"
  echo "morpho oracle $oracle: $asset with Chainlink's round older than the band takes and its 24/7 leg $stale_s s old" \
    "has no answer; after a fresh $package package, $price, the low edge in state $state; the same package again" \
    "reverts PackageNotNewer"
  staled=$((staled + 1))
done
[ "$staled" -gt 0 ] || fail "every oracle's band is halted by something other than its legs"
echo "PASS"
