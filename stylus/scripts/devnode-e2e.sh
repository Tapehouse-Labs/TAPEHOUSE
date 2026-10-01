#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-e2e.sh RPC_URL PRIVATE_KEY BAND_ADDRESS DEPLOYMENTS_JSON
set -euo pipefail

rpc=$1 key=$2 band=$3 registry=$4
root=$(cd "$(dirname "$0")/../.." && pwd)
payload="$(dirname "$0")/redstone-payload.py"
symbol_nvda=$(cast format-bytes32-string NVDA)
symbol_tsla=$(cast format-bytes32-string TSLA)
symbol_aapl=$(cast format-bytes32-string AAPL)
symbol_spy=$(cast format-bytes32-string SPY)
node_interface=0x00000000000000000000000000000000000000C8
nvda=$(cast format-bytes32-string NVDA---24_7)
usa500=$(cast format-bytes32-string USA500.Y---24_7)
zero_id=0x0000000000000000000000000000000000000000000000000000000000000000
zero_address=0x0000000000000000000000000000000000000000
token=$(jq -r .tokens.NVDA "$registry")
token_px() { python3 -c "import sys; print(int(sys.argv[1]) * int(sys.argv[2]) // 10**18)" "$1" "$(cast call --rpc-url "$rpc" "$token" "uiMultiplier()(uint256)" | cut -d' ' -f1)"; }
status_feeds="[$(cast format-bytes32-string NY_MARKET_CURRENT_STATUS),$(cast format-bytes32-string NY_MARKET_NEXT_STATUS),$(cast format-bytes32-string NY_MARKET_NEXT_CHANGE_TIME)]"
price_written=$(cast keccak "PriceWritten(bytes32,uint256,uint64)")

fail() { echo "FAIL: $*"; exit 1; }

write() {
  local hash
  hash=$(cast send --rpc-url "$rpc" --private-key "$key" "$band" "writePrices(bytes32[],bytes)" "$1" "$2" --json |
    jq -r .transactionHash) || return 1
  cast rpc --rpc-url "$rpc" eth_getTransactionReceipt "$hash" |
    jq -r '"\(.status) \(.gasUsed) \(.gasUsedForL1) \(.blockNumber) \(.logs[0].topics[0]) \(.logs[0].topics[1])"'
}

l2_gas() {
  cast call --rpc-url "$rpc" $node_interface "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" \
    "${2:-$band}" false "$1" | sed -n 1,2p | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}

gas_log=$root/stylus/target/devnode-gas.txt
: > "$gas_log"
record() {
  echo "$1 $2" >> "$gas_log"
  [[ $2 =~ ^[0-9]+$ ]] || { echo "FAIL: no L2 gas for $1" >&2; exit 1; }
  echo "$2"
}
measure() { record "$1" "$(l2_gas "$(cast calldata "$2" "${@:3}")")"; }

legs() {
  cast call --rpc-url "$rpc" "$band" "legs(bytes32)(uint256,uint64,uint256,uint64)" "$1" | cut -d' ' -f1 | tr '\n' ' '
}

expect_revert() {
  local out
  if out=$(cast call --rpc-url "$rpc" "$band" "writePrices(bytes32[],bytes)" "$1" "$2" 2>&1); then
    fail "expected revert $3, got $out"
  fi
  grep -q "$3" <<<"$out" || fail "expected revert $3, got $out"
}

price_slot() {
  local base
  base=$(cast index bytes32 "$1" 0)
  cast storage --rpc-url "$rpc" "$band" "$(python3 -c "print(hex(int('$base', 16) + 1))")"
}
sample_ms() {
  local word
  word=$(price_slot "$1")
  echo $((16#${word:18:16}))
}
tracked() {
  local word
  word=$(price_slot "$1")
  echo $((16#${word:16:2}))
}

[ "$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$nvda" | tr '\n' ' ')" = "0 0 0 " ] ||
  fail "unwritten feed does not read as zero"

nvda_payload=$(python3 "$payload" NVDA---24_7)
expected=$(cast call --rpc-url "$rpc" "$band" "writePrices(bytes32[],bytes)(uint256[])" "[$nvda]" "$nvda_payload" | tr -d '[]' | cut -d' ' -f1)
read -r status gas l1 block topic feed <<<"$(write "[$nvda]" "$nvda_payload")"
[ "$status" = 0x1 ] || fail "writePrices reverted for NVDA---24_7"
[ "$topic" = "$price_written" ] && [ "$feed" = "$nvda" ] || fail "PriceWritten not emitted for NVDA---24_7"
read -r value package_ms written <<<"$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$nvda" | cut -d' ' -f1 | tr '\n' ' ')"
[ "$value" = "$expected" ] || fail "stored $value, verified $expected"
[ "$written" = "$(cast block --rpc-url "$rpc" "$((block))" -f timestamp)" ] || fail "write timestamp is not the block timestamp"
echo "NVDA---24_7 = $value at $package_ms ms, L2 gas $((gas - l1))"

expect_revert "[$nvda]" "$nvda_payload" 0x5d1d5d09

BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-band.sh" "$rpc" "$band" "$registry"
jq '.chainlink |= {NVDA_USD: .TSLA_USD, TSLA_USD: .NVDA_USD}' "$registry" > "$registry.swapped"
if out=$(BAND_ASSETS="NVDA TSLA" "$(dirname "$0")/check-band.sh" "$rpc" "$band" "$registry.swapped"); then
  fail "check-band.sh accepted a band that differs from the registry"
fi
grep -q 'FAIL: band holds' <<<"$out" || fail "check-band.sh failed for another reason: $out"
jq 'del(.chainlink.TSLA_USD)' "$registry" > "$registry.missing"
if out=$(BAND_ASSETS="NVDA TSLA" "$(dirname "$0")/check-band.sh" "$rpc" "$band" "$registry.missing" 2>&1); then
  fail "check-band.sh accepted a registry that band-args.sh rejects: $out"
fi
jq '.tapehouse.HaltSigner = "0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"' "$registry" > "$registry.signer"
if out=$(BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-band.sh" "$rpc" "$band" "$registry.signer"); then
  fail "check-band.sh accepted another halt signer"
fi
grep -q "halt signer" <<<"$out" || fail "check-band.sh failed for another reason: $out"
words=$(BAND_ASSETS="NVDA TSLA" "$(dirname "$0")/band-args.sh" "$registry.swapped")
read -r -a args <<<"$words"
swapped=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "${args[@]}" 2>&1 | grep 'deployed code at address' | grep -o '0x[0-9a-f]\{40\}')
if out=$(BAND_ASSETS="NVDA TSLA" "$(dirname "$0")/check-band.sh" "$rpc" "$swapped" "$registry.swapped"); then
  fail "check-band.sh accepted feeds that describe other assets"
fi
grep -q 'describes itself as' <<<"$out" || fail "check-band.sh failed for another reason: $out"
read -r cl_price cl_updated_at price_247 price_247_ms <<<"$(legs "$symbol_nvda")"
[ "$cl_price" = 22900000000 ] && [ "$cl_updated_at" -gt 0 ] && [ "$price_247" = "$(token_px "$value")" ] &&
  [ "$price_247_ms" = "$package_ms" ] || fail "legs(NVDA) = $cl_price $cl_updated_at $price_247 $price_247_ms, not the token's price"
read -r cl_price cl_updated_at price_247 price_247_ms <<<"$(legs "$symbol_tsla")"
[ "$cl_price" = 37800000000 ] && [ "$cl_updated_at" -gt 0 ] && [ "$price_247" = 0 ] && [ "$price_247_ms" = 0 ] ||
  fail "legs(TSLA) = $cl_price $cl_updated_at $price_247 $price_247_ms before any TSLA write"
[ "$(legs "$symbol_aapl")" = "0 0 0 0 " ] || fail "legs(AAPL) is not empty for an unknown symbol"
read -r cl_price cl_updated_at price_247 price_247_ms <<<"$(legs "$symbol_spy")"
[ "$cl_price" = 77232802713 ] && [ "$price_247" = 0 ] || fail "legs(SPY) = $cl_price $cl_updated_at $price_247 $price_247_ms before its first anchor"
echo "legs L2 gas: NVDA $(measure "legs(NVDA)" "legs(bytes32)" "$symbol_nvda"), unknown symbol $(measure "legs(unknown)" "legs(bytes32)" "$symbol_aapl")"

again=0x5585258d$(cast abi-encode "f(bytes32[],address[],bytes32[],bytes32[],address[])" "[$symbol_aapl]" \
  "[$zero_address]" "[$(cast format-bytes32-string AAPL---24_7)]" "[$zero_id]" "[$zero_address]" | cut -c3-)
if out=$(cast call --rpc-url "$rpc" "$band" "$again" 2>&1); then
  fail "the constructor ran a second time: $out"
fi
grep -q 'data: "0x"$' <<<"$out" || fail "the second constructor call reverted for another reason: $out"
extra=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "[$symbol_aapl,$symbol_spy]" "[$zero_address,$zero_address]" \
  "[$(cast format-bytes32-string AAPL---24_7),$usa500]" "[$zero_id,$zero_id]" "[$zero_address,$zero_address]" "${args[5]}" "${args[6]}" \
  "$zero_address" false 2>&1 |
  grep 'deployed code at address' | grep -o '0x[0-9a-f]\{40\}')
jq -n --arg signer "${args[5]}" --arg owner "${args[6]}" '{chainId: 412346, tapehouse: {HaltSigner: $signer, Owner: $owner}}' \
  > "$registry.redstone"
if out=$(BAND_ASSETS="AAPL SPY" "$(dirname "$0")/check-band.sh" "$rpc" "$extra" "$registry.redstone"); then
  fail "check-band.sh accepted a band that configures an asset band-args.sh leaves out"
fi
grep -q 'FAIL: band configures SPY' <<<"$out" || fail "check-band.sh failed for another reason: $out"

stub_18=$(forge create --root "$root/contracts" test/devnode/StubAggregator.sol:StubAggregator \
  --rpc-url "$rpc" --private-key "$key" --broadcast --json --constructor-args 18 1 1 "NVDA / USD" | jq -r .deployedTo)
if out=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "[$symbol_nvda]" "[$stub_18]" "[$nvda]" "[$zero_id]" "[$zero_address]" "${args[5]}" "${args[6]}" \
  "$zero_address" false 2>&1); then
  fail "a feed with 18 decimals was accepted"
fi
if ! grep -q 0x88d8f57d <<<"$out" || ! grep -q 6787b555 <<<"$out"; then
  fail "expected InvalidFeed inside ContractInitializationError, got: $(tail -n 1 <<<"$out")"
fi
if out=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "[$symbol_nvda]" "[$zero_address]" "[$nvda]" "[$zero_id]" "[$zero_address]" "$zero_address" "${args[6]}" \
  "$zero_address" false 2>&1); then
  fail "a zero halt signer was accepted"
fi
if ! grep -q 0x88d8f57d <<<"$out" || ! grep -q 126917ee <<<"$out"; then
  fail "expected InvalidHaltSigner inside ContractInitializationError, got: $(tail -n 1 <<<"$out")"
fi

deploy_tx=$(grep 'deployment tx hash' "$root/stylus/target/devnode-band.log" | grep -o '0x[0-9a-f]\{64\}')
cast rpc --rpc-url "$rpc" eth_getTransactionReceipt "$deploy_tx" | jq -r '"\(.gasUsed) \(.gasUsedForL1)"' |
  { read -r gas l1; echo "deploy, activation and constructor through StylusDeployer: L2 gas $((gas - l1))"; }

flip=$((${#nvda_payload} - 186))
tampered=${nvda_payload:0:flip}$(printf '%02x' $((0x${nvda_payload:flip:2} ^ 1)))${nvda_payload:flip+2}
expect_revert "[$nvda]" "$tampered" 0xec459bc0
value_size_at=$((${#nvda_payload} - 172))
oversized=${nvda_payload:0:value_size_at}ffffffe0${nvda_payload:value_size_at+8}
expect_revert "[$nvda]" "$oversized" 0x5796f78a

[ "$(cast call --rpc-url "$rpc" "$band" "variance(bytes32)(uint128)" "$nvda")" = 0 ] || fail "variance before a second sample"
[ "$(tracked "$nvda")" = 1 ] && [ "$(sample_ms "$nvda")" = "$package_ms" ] ||
  fail "the first write of a tracked feed is not its first sample"

first_ms=$package_ms first_px=$value attempts=0
while :; do
  attempts=$((attempts + 1))
  [ "$attempts" -le 20 ] || fail "no NVDA---24_7 package 50 s newer than the first within 200 s"
  sleep 10
  next_payload=$(python3 "$payload" NVDA---24_7)
  cast call --rpc-url "$rpc" "$band" "writePrices(bytes32[],bytes)" "[$nvda]" "$next_payload" > /dev/null 2>&1 || continue
  read -r status gas l1 _ <<<"$(write "[$nvda]" "$next_payload")"
  [ "$status" = 0x1 ] || fail "writePrices reverted for NVDA---24_7"
  read -r px ms _ <<<"$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$nvda" | cut -d' ' -f1 | tr '\n' ' ')"
  if [ $((ms - first_ms)) -lt 50000 ]; then
    [ "$(sample_ms "$nvda")" = "$first_ms" ] || fail "a write within 50 s was sampled"
    continue
  fi
  if [ "$px" -ge "$first_px" ]; then r=$(((px - first_px) * 1000000 / first_px)); else r=$((((first_px - px) * 1000000 + first_px - 1) / first_px)); fi
  expected_var=$((6 * (r * r * 60 / ((ms - first_ms) / 1000)) / 100))
  [ "$(cast call --rpc-url "$rpc" "$band" "variance(bytes32)(uint128)" "$nvda" | cut -d' ' -f1)" = "$expected_var" ] ||
    fail "variance is not the EWMA of $first_px and $px over $((ms - first_ms)) ms"
  [ "$(sample_ms "$nvda")" = "$ms" ] || fail "the sample did not move to the package at $ms ms"
  echo "NVDA---24_7 sample after $((ms - first_ms)) ms: variance $expected_var, L2 gas $((gas - l1))"
  break
done

read -r status gas l1 _ <<<"$(write "$status_feeds" "$(python3 "$payload" NY_MARKET_STATUS)")"
[ "$status" = 0x1 ] || fail "writePrices reverted for NY_MARKET_STATUS"
change_time=$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$(cast format-bytes32-string NY_MARKET_NEXT_CHANGE_TIME)" | sed -n 1p | cut -d' ' -f1)
[ "$(cast to-base "$change_time" 16 | wc -c)" -gt 19 ] || fail "NY_MARKET_NEXT_CHANGE_TIME $change_time fits in 64 bits"
status_current=$(cast format-bytes32-string NY_MARKET_CURRENT_STATUS)
[ "$(tracked "$status_current")" = 0 ] && [ "$(sample_ms "$status_current")" = 0 ] || fail "an untracked feed was sampled"
echo "NY_MARKET_STATUS: 3 feeds, NEXT_CHANGE_TIME = $change_time, L2 gas $((gas - l1))"
expect_revert "[$(cast format-bytes32-string NY_MARKET_CURRENT_STATUS)]" "$(python3 "$payload" NY_MARKET_STATUS)" 0x140a6d4d

spy_feed=$(jq -r .chainlink.SPY_USD "$registry")
keeper_feeds="[$nvda,$usa500,${status_feeds#[}"
attempts=0
until keeper_payload=$(python3 "$payload" NVDA---24_7 USA500.Y---24_7 NY_MARKET_STATUS) &&
  cast call --rpc-url "$rpc" "$band" "writePrices(bytes32[],bytes)" "$keeper_feeds" "$keeper_payload" > /dev/null 2>&1; do
  attempts=$((attempts + 1))
  [ "$attempts" -le 20 ] || fail "no newer NVDA---24_7, USA500.Y---24_7 and NY_MARKET_STATUS packages within 200 s"
  sleep 10
done
print_at=$(date +%s)
cast send --rpc-url "$rpc" --private-key "$key" "$spy_feed" "setRound(int256,uint256)" 77232802713 "$print_at" > /dev/null
read -r status gas l1 _ <<<"$(write "$keeper_feeds" "$keeper_payload")"
[ "$status" = 0x1 ] || fail "writePrices reverted for NVDA---24_7, USA500.Y---24_7 and NY_MARKET_STATUS"
read -r index_px _ <<<"$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$usa500" | cut -d' ' -f1 | tr '\n' ' ')"
[ "$(cast call --rpc-url "$rpc" "$band" "anchor(bytes32)(uint64,uint64,uint64)" "$symbol_spy" | cut -d' ' -f1 | tr '\n' ' ')" = "77232802713 $index_px $print_at " ] ||
  fail "the index write did not anchor SPY's print at $print_at"
[ "$(cast call --rpc-url "$rpc" "$band" "variance(bytes32)(uint128)" "$usa500")" = 0 ] || fail "variance of the index after one sample"
echo "keeper write of NVDA---24_7, USA500.Y---24_7 and NY_MARKET_STATUS, anchoring SPY: L2 gas $((gas - l1))"

read -r session nyse nyse_next change_ms boundary_ms <<<"$(cast call --rpc-url "$rpc" "$band" "session()(uint8,uint8,uint8,uint64,uint64)" | cut -d' ' -f1 | tr '\n' ' ')"
now_ms=$(($(cast block --rpc-url "$rpc" latest -f timestamp) * 1000))
case "$session-$nyse" in
  2-1 | 2-2 | 2-3) expected="3 2" spy_half=55 ;;
  1-3) [ "$boundary_ms" -gt "$now_ms" ] || fail "a closed session reopens at $boundary_ms, before now"
    expected="2 1" spy_half=80 ;;
  1-2) expected="2 1" spy_half=80 ;;
  *) fail "session() = $session $nyse $nyse_next $change_ms $boundary_ms" ;;
esac
quote() {
  cast call --rpc-url "$rpc" "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$1" | cut -d' ' -f1 | tr '\n' ' '
}
read -r nvda_share _ <<<"$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$nvda" | cut -d' ' -f1 | tr '\n' ' ')"
nvda_px=$(token_px "$nvda_share")
read -r state live mid _ <<<"$(quote "$symbol_nvda")"
[ "$state $live" = "$expected" ] && [ "$mid" = "$nvda_px" ] || fail "quote(NVDA) = $(quote "$symbol_nvda") with session $session"
if [ "$session" = 2 ]; then
  read -r state live mid _ <<<"$(quote "$symbol_tsla")"
  [ "$state $live $mid" = "1 1 37800000000" ] || fail "quote(TSLA) = $(quote "$symbol_tsla") with Chainlink alone in session"
else
  [ "$(quote "$symbol_tsla")" = "0 0 0 0 0 0 " ] || fail "quote(TSLA) = $(quote "$symbol_tsla") with Chainlink asleep"
fi
read -r state live mid half _ <<<"$(quote "$symbol_spy")"
[ "$state $live" = "$expected" ] && [ "$mid" = 77232802713 ] && [ "$half" -ge "$spy_half" ] ||
  fail "quote(SPY) = $(quote "$symbol_spy") with session $session"
echo "session $session, NYSE $nyse; quote(NVDA) $(quote "$symbol_nvda")"
echo "L2 gas: quote(NVDA) $(measure "quote(NVDA)" "quote(bytes32)" "$symbol_nvda"), quote(SPY) $(measure "quote(SPY)" "quote(bytes32)" "$symbol_spy"), session() $(measure "session()" "session()")"
echo "L2 gas: price(NVDA---24_7) $(measure "price(NVDA---24_7)" "price(bytes32)" "$nvda"), anchor(SPY) $(measure "anchor(SPY)" "anchor(bytes32)" "$symbol_spy"), halt(NVDA) $(measure "halt(NVDA)" "halt(bytes32)" "$symbol_nvda")"

attempts=0
until keeper_payload=$(python3 "$payload" NVDA---24_7 USA500.Y---24_7 NY_MARKET_STATUS) &&
  cast call --rpc-url "$rpc" "$band" "writePrices(bytes32[],bytes)" "$keeper_feeds" "$keeper_payload" > /dev/null 2>&1; do
  attempts=$((attempts + 1))
  [ "$attempts" -le 20 ] || fail "no newer NVDA---24_7, USA500.Y---24_7 and NY_MARKET_STATUS packages within 200 s"
  sleep 10
done
read -r status gas l1 _ <<<"$(write "$keeper_feeds" "$keeper_payload")"
[ "$status" = 0x1 ] || fail "the second keeper write reverted"
echo "keeper update of the same five feeds: L2 gas $((gas - l1))"

feeds_of() {
  jq --arg band "$1" '.tapehouse.Band = $band' "$registry" > "$registry.feeds"
  BAND_ASSETS=$2 "$root/contracts/script/deploy-band-feeds.sh" "$rpc" "$registry.feeds" --private-key "$key"
}
description() { cast call --rpc-url "$rpc" "$1" "description()(string)"; }
round_of() { cast call --rpc-url "$rpc" "$1" "latestRoundData()(uint80,int256,uint256,uint256,uint80)" | cut -d' ' -f1 | tr '\n' ' '; }
expect_feed_revert() {
  local out
  if out=$(cast call --rpc-url "$rpc" "$1" "latestRoundData()(uint80,int256,uint256,uint256,uint80)" 2>&1); then
    fail "expected revert $2, got $out"
  fi
  grep -q "$2" <<<"$out" || fail "expected revert $2, got $out"
}
band_view="latestBand()((uint8,uint8,uint64,uint64,uint64,uint128,uint128,uint8,uint8,uint8,uint64,uint64,bool,uint64,uint64,bool,bool,bool,uint256,int256))"
feeds=$(feeds_of "$band" "NVDA TSLA SPY")
feed_address() { awk -v asset="$1" '$1 == asset { print $2 }' <<<"$feeds"; }
feed=$(feed_address NVDA)
[ "$(description "$feed")" = '"NVDA / USD Tapehouse band, low side (Robinhood NVDA Stock Token)"' ] &&
  [ "$(description "$(feed_address TSLA)")" = '"TSLA / USD Tapehouse band, low side (TSLA share)"' ] &&
  [ "$(description "$(feed_address SPY)")" = '"SPY / USD Tapehouse band, low side (Robinhood SPY Stock Token)"' ] ||
  fail "deploy-band-feeds.sh did not describe what each feed prices: $feeds"
read -r state _ mid _ low _ <<<"$(quote "$symbol_nvda")"
read -r round answer started updated _ <<<"$(round_of "$feed")"
now=$(cast block --rpc-url "$rpc" latest -f timestamp)
[ "$answer" = "$low" ] && [ "$round $started $updated" = "$now $now $now" ] || fail "the feed answered $answer at $round, quote's low is $low"
[ "$(cast call --rpc-url "$rpc" "$feed" "decimals()(uint8)")" = 8 ] || fail "the feed does not have 8 decimals"
[ "$(cast call --rpc-url "$rpc" "$feed" "$band_view" | sed 's/ \[[^]]*\]//g' | tr -d '() ' | cut -d, -f1,3)" = "$state,$mid" ] ||
  fail "latestBand() is not the band: $(cast call --rpc-url "$rpc" "$feed" "$band_view")"
if out=$(cast call --rpc-url "$rpc" "$feed" "seal()" 2>&1); then
  echo "the band is within ten minutes of its reopen: seal() would store it"
else
  grep -q 0x7c193929 <<<"$out" || fail "expected NotSealWindow, got $out"
fi
echo "BandFeed: answers quote's low; L2 gas latestRoundData $(record "BandFeed.latestRoundData(NVDA)" "$(l2_gas "$(cast calldata "latestRoundData()")" "$feed")"), latestBand $(record "BandFeed.latestBand(NVDA)" "$(l2_gas "$(cast calldata "latestBand()")" "$feed")")"

nvda_feed=$(jq -r .chainlink.NVDA_USD "$registry")
math() { python3 -c "print($1)"; }
mine() { cast send --rpc-url "$rpc" --private-key "$key" --value 0 "$(cast wallet address "$key")" > /dev/null; }
action() {
  cast call --rpc-url "$rpc" "${2:-$band}" "corporateAction(bytes32)(uint8,uint64,uint128,uint128)" "${1:-$symbol_nvda}" |
    cut -d' ' -f1 | tr '\n' ' '
}
sync_multiplier() {
  local hash
  hash=$(cast send --rpc-url "$rpc" --private-key "$key" "${2:-$band}" "syncMultiplier(bytes32)" "${1:-$symbol_nvda}" --json |
    jq -r .transactionHash)
  cast rpc --rpc-url "$rpc" eth_getTransactionReceipt "$hash" | jq -r '"\(.status) \(.gasUsed) \(.gasUsedForL1) \(.logs | length)"' |
    { read -r status gas l1 logs; echo "$status $((gas - l1)) $logs"; }
}
schedule() {
  local at=$(($(date +%s) + 20))
  cast send --rpc-url "$rpc" --private-key "$key" "${2:-$token}" "updateMultiplier(uint256,uint256)" "$1" "$at" > /dev/null
  echo "$at"
}
after() {
  local s=$(($1 + 1 - $(date +%s)))
  [ "$s" -le 0 ] || sleep "$s"
  mine
}
fresh_nvda() {
  local attempts=0 nvda_payload
  until nvda_payload=$(python3 "$payload" NVDA---24_7) &&
    cast call --rpc-url "$rpc" "$1" "writePrices(bytes32[],bytes)" "[$nvda]" "$nvda_payload" > /dev/null 2>&1; do
    attempts=$((attempts + 1))
    [ "$attempts" -le 20 ] || fail "no newer NVDA---24_7 package within 200 s"
    sleep 10
  done
  cast send --rpc-url "$rpc" --private-key "$key" "$1" "writePrices(bytes32[],bytes)" "[$nvda]" "$nvda_payload" > /dev/null
  cast call --rpc-url "$rpc" "$1" "price(bytes32)(uint256,uint64,uint64)" "$nvda" | sed -n 1p | cut -d' ' -f1
}
halted="0 0 0 0 0 0 "

halt_key=0x59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d
chain_id=$(cast chain-id --rpc-url "$rpc")
block_time() { cast block --rpc-url "$rpc" latest -f timestamp; }
next_second() {
  local t
  until t=$(block_time) && [ "$t" -gt "$1" ]; do
    sleep 0.2
    mine
  done
  echo "$t"
}
halt_signature() {
  jq -n --arg band "${5:-$band}" --argjson chain "${6:-$chain_id}" --arg symbol "$symbol_nvda" --argjson halted "$1" \
    --argjson issued "$2" --argjson expires "$3" '{
    types: {
      EIP712Domain: [{name: "name", type: "string"}, {name: "version", type: "string"},
        {name: "chainId", type: "uint256"}, {name: "verifyingContract", type: "address"}],
      HaltState: [{name: "symbol", type: "bytes32"}, {name: "halted", type: "bool"},
        {name: "issuedAt", type: "uint64"}, {name: "expiresAt", type: "uint64"}]
    },
    primaryType: "HaltState",
    domain: {name: "Tapehouse Band", version: "1", chainId: $chain, verifyingContract: $band},
    message: {symbol: $symbol, halted: $halted, issuedAt: $issued, expiresAt: $expires}
  }' > "$registry.halt.json"
  cast wallet sign --private-key "${4:-$halt_key}" --data --from-file "$registry.halt.json"
}
write_halt() {
  local signature hash
  signature=$(halt_signature "$@")
  hash=$(cast send --rpc-url "$rpc" --private-key "$key" "$band" "writeHalt(bytes32,bool,uint64,uint64,bytes)" \
    "$symbol_nvda" "$1" "$2" "$3" "$signature" --json | jq -r .transactionHash) || fail "writeHalt($1, $2, $3) reverted"
  cast rpc --rpc-url "$rpc" eth_getTransactionReceipt "$hash" | jq -r '"\(.status) \(.gasUsed) \(.gasUsedForL1)"' |
    { read -r status gas l1; echo "$status $((gas - l1))"; }
}
expect_halt_revert() {
  local out signature
  signature=$(halt_signature "$1" "$2" "$3" "${5:-}" "${6:-}" "${7:-}")
  if out=$(cast call --rpc-url "$rpc" "$band" "writeHalt(bytes32,bool,uint64,uint64,bytes)" \
    "$symbol_nvda" "$1" "$2" "$3" "$signature" 2>&1); then
    fail "expected revert $4, got $out"
  fi
  grep -q "$4" <<<"$out" || fail "expected revert $4, got $out"
}
halt_view() {
  cast call --rpc-url "$rpc" "$band" "halt(bytes32)(bool,uint64,uint64,bool)" "$symbol_nvda" | cut -d' ' -f1 | tr '\n' ' '
}

[ "$(cast call --rpc-url "$rpc" "$band" "haltSigner()(address)")" = "$(jq -r .tapehouse.HaltSigner "$registry")" ] ||
  fail "haltSigner() is not the registry's halt signer"
live_state=$(quote "$symbol_nvda" | cut -d' ' -f1)
[ "$live_state" != 0 ] || fail "NVDA is halted before any halt"
issued=$(block_time)
expect_halt_revert true "$issued" $((issued + 30)) 0xec459bc0 0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
expect_halt_revert true "$issued" $((issued + 3601)) 0x1d538dfe
expect_halt_revert true "$issued" $((issued + 30)) 0xec459bc0 "" "$swapped"
expect_halt_revert true "$issued" $((issued + 30)) 0xec459bc0 "" "" 46630
read -r status halt_gas <<<"$(write_halt true "$issued" $((issued + 30)))"
[ "$status" = 0x1 ] && [ "$(halt_view)" = "true $((issued + 30)) $issued false " ] &&
  [ "$(quote "$symbol_nvda")" = "0 0 0 0 0 0 " ] || fail "a signed halt did not halt NVDA: $(halt_view)"
expect_halt_revert true "$issued" $((issued + 30)) 0xedd6aa50
expect_feed_revert "$feed" 0x39c5542d
signed_quote_gas=$(record "quote(NVDA,signedHalt)" "$(l2_gas "$(cast calldata "quote(bytes32)" "$symbol_nvda")")")
lifted=$(next_second "$issued")
write_halt false "$lifted" $((lifted + 30)) > /dev/null
[ "$(halt_view)" = "false 0 $lifted false " ] && [ "$(quote "$symbol_nvda" | cut -d' ' -f1)" != 0 ] ||
  fail "a newer lift did not end the halt: $(halt_view)"
issued=$(next_second "$lifted")
write_halt true "$issued" $((issued + 20)) > /dev/null
[ "$(quote "$symbol_nvda" | cut -d' ' -f1)" = 0 ] || fail "the second halt did not hold"
after $((issued + 20))
[ "$(halt_view)" = "false $((issued + 20)) $issued false " ] && [ "$(quote "$symbol_nvda" | cut -d' ' -f1)" = 1 ] ||
  fail "a lapsed halt did not degrade NVDA: $(halt_view) $(quote "$symbol_nvda")"
lifted=$(block_time)
write_halt false "$lifted" $((lifted + 30)) > /dev/null
[ "$(quote "$symbol_nvda" | cut -d' ' -f1)" = "$live_state" ] || fail "a lift after a lapse did not restore NVDA's band"
cast send --rpc-url "$rpc" --private-key "$key" "$token" "pauseOracle()" > /dev/null
[ "$(halt_view | cut -d' ' -f4)" = true ] && [ "$(quote "$symbol_nvda")" = "0 0 0 0 0 0 " ] ||
  fail "the issuer's oracle pause did not halt NVDA: $(halt_view)"
halted_quote_gas=$(record "quote(NVDA,oraclePaused)" "$(l2_gas "$(cast calldata "quote(bytes32)" "$symbol_nvda")")")
cast send --rpc-url "$rpc" --private-key "$key" "$token" "unpauseOracle()" > /dev/null
[ "$(quote "$symbol_nvda" | cut -d' ' -f1)" != 0 ] || fail "NVDA stayed halted after the oracle pause ended"
echo "signed halt: written with L2 gas $(record "writeHalt(NVDA)" "$halt_gas"); replay, other key, band and chain rejected; lifted; lapsed to degraded"
echo "L2 gas: quote(NVDA) under a signed halt $signed_quote_gas, while its oracle is paused $halted_quote_gas"

m0=$(cast call --rpc-url "$rpc" "$token" "uiMultiplier()(uint256)" | cut -d' ' -f1)
m1=$(math "$m0 * 10017 // 10000")
at=$(schedule "$m1")
read -r status sync_gas logs <<<"$(sync_multiplier)"
[ "$status $logs" = "0x1 1" ] && [ "$(action)" = "0 $at $m0 $m1 " ] || fail "a 17 bps dividend was not recorded, or needs a window: $(action)"
after "$at"
read -r nvda_share _ <<<"$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$nvda" | cut -d' ' -f1 | tr '\n' ' ')"
read -r cl_price _ price_247 _ <<<"$(legs "$symbol_nvda")"
[ "$cl_price" = "$(math "22900000000 * $m1 // $m0")" ] && [ "$price_247" = "$(token_px "$nvda_share")" ] ||
  fail "after a dividend, legs(NVDA) = $cl_price $price_247"
[ "$(quote "$symbol_nvda" | cut -d' ' -f1)" != 0 ] || fail "a dividend halted the band"
echo "dividend of 17 bps: recorded with L2 gas $sync_gas, Chainlink scaled to $cl_price"

m2=$(math "$m1 * 4")
at=$(schedule "$m2")
[ "$(action)" = "1 $at $m1 $m2 " ] && [ "$(quote "$symbol_nvda")" != "$halted" ] || fail "a split scheduled and not synced is not status 1: $(action)"
after "$at"
[ "$(action)" = "2 $at 0 $m2 " ] && [ "$(quote "$symbol_nvda")" = "$halted" ] || fail "a split past its step and not synced did not halt NVDA: $(action)"
read -r status _ logs <<<"$(sync_multiplier)"
[ "$logs $(action | cut -d' ' -f1)" = "1 2" ] || fail "a split was confirmed without a Chainlink round after it"
split_at=$at
m3=$(math "$m2 * 10010 // 10000")
at=$(schedule "$m3")
read -r status _ logs <<<"$(sync_multiplier)"
[ "$logs $(action)" = "0 2 $split_at 0 $m2 " ] && [ "$(quote "$symbol_nvda")" = "$halted" ] || fail "a dividend scheduled over an unconfirmed split ended its halt: $(action)"
after "$at"
[ "$(action)" = "2 $split_at 0 $m2 " ] && [ "$(quote "$symbol_nvda")" = "$halted" ] || fail "a dividend's step ended an unconfirmed split's halt: $(action)"
nvda_share=$(fresh_nvda "$band")
cast send --rpc-url "$rpc" --private-key "$key" "$nvda_feed" "setRound(int256,uint256)" "$(math "$(token_px "$nvda_share") * 2")" "$(date +%s)" > /dev/null
read -r status _ logs <<<"$(sync_multiplier)"
[ "$logs" = 0 ] && [ "$(quote "$symbol_nvda")" = "$halted" ] || fail "a Chainlink round outside the 24/7 band confirmed the split: $(action)"
cast send --rpc-url "$rpc" --private-key "$key" "$nvda_feed" "setRound(int256,uint256)" "$(token_px "$nvda_share")" "$(date +%s)" > /dev/null
cast send --rpc-url "$rpc" --private-key "$key" "$token" "pauseOracle()" > /dev/null
read -r status _ logs <<<"$(sync_multiplier)"
[ "$logs $(action | cut -d' ' -f1)" = "0 2" ] || fail "a split was confirmed while the issuer paused the oracle: $(action)"
cast send --rpc-url "$rpc" --private-key "$key" "$token" "unpauseOracle()" > /dev/null
read -r status confirm_gas logs <<<"$(sync_multiplier)"
[ "$logs $(action)" = "3 0 $at 0 $m3 " ] && [ "$(quote "$symbol_nvda")" != "$halted" ] ||
  fail "a Chainlink round inside the 24/7 band did not confirm the split and the dividend after it: $(action)"
echo "4:1 split, then a dividend: halted until confirmed; confirming both, L2 gas $confirm_gas"

m4=$(math "$m3 * 2")
at=$(schedule "$m4")
after "$at"
m5=$(math "$m4 * 10010 // 10000")
m5_at=$(schedule "$m5")
now=$(cast block --rpc-url "$rpc" latest -f timestamp)
[ "$(action)" = "2 $now 0 $m4 " ] && [ "$(quote "$symbol_nvda")" = "$halted" ] || fail "a step nobody recorded did not halt NVDA: $(action)"
read -r status _ logs <<<"$(sync_multiplier)"
[ "$status $logs $(action | cut -d' ' -f1)" = "0x1 1 2" ] || fail "a step nobody recorded was not recorded as unconfirmed: $(action)"
echo "a 2:1 split nobody recorded: halted from the first read"

spy_token=$(jq -r .tokens.SPY "$registry")
s0=$(cast call --rpc-url "$rpc" "$spy_token" "uiMultiplier()(uint256)" | cut -d' ' -f1)
s1=$(math "$s0 * 10017 // 10000")
read -r spy_cl _ <<<"$(legs "$symbol_spy")"
at=$(schedule "$s1" "$spy_token")
read -r status _ logs <<<"$(sync_multiplier "$symbol_spy")"
[ "$status $logs" = "0x1 1" ] || fail "SPY's dividend was not recorded: $(action "$symbol_spy")"
after "$at"
read -r anchor_cl anchor_index _ <<<"$(cast call --rpc-url "$rpc" "$band" "anchor(bytes32)(uint64,uint64,uint64)" "$symbol_spy" | cut -d' ' -f1 | tr '\n' ' ')"
read -r index_px _ <<<"$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$usa500" | sed -n 1p)"
read -r cl_price _ price_247 _ <<<"$(legs "$symbol_spy")"
[ "$cl_price" = "$(math "$spy_cl * $s1 // $s0")" ] &&
  [ "$price_247" = "$(math "($anchor_cl * $s1 // $s0) * $index_px // $anchor_index")" ] ||
  fail "after SPY's dividend, legs(SPY) = $cl_price $price_247 with anchor $anchor_cl $anchor_index"
echo "SPY's dividend: Chainlink and the index leg's anchor scaled to $cl_price and $price_247"

after "$m5_at"
words=$(BAND_ASSETS="NVDA" "$(dirname "$0")/band-args.sh" "$registry")
read -r -a args <<<"$words"
fresh=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "${args[@]}" 2>&1 | grep 'deployed code at address' | grep -o '0x[0-9a-f]\{40\}')
[ "$(action "$symbol_nvda" "$fresh")" = "2 $m5_at 0 $m5 " ] || fail "a band deployed after a step does not start halted: $(action "$symbol_nvda" "$fresh")"
nvda_share=$(fresh_nvda "$fresh")
cast send --rpc-url "$rpc" --private-key "$key" "$nvda_feed" "setRound(int256,uint256)" "$(token_px "$nvda_share")" "$(date +%s)" > /dev/null
read -r status _ logs <<<"$(sync_multiplier "$symbol_nvda" "$fresh")"
[ "$logs $(action "$symbol_nvda" "$fresh")" = "1 0 $m5_at 0 $m5 " ] || fail "a fresh band's first sync did not confirm: $(action "$symbol_nvda" "$fresh")"
echo "a band deployed after a step: halted until its first sync confirms"
echo "L2 gas: quote(NVDA) with its Stock Token $(l2_gas "$(cast calldata "quote(bytes32)" "$symbol_nvda")")"

owner=$(jq -r .tapehouse.Owner "$registry")
next_key=0x5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a
next=$(cast wallet address --private-key "$next_key")
[ "$(cast call --rpc-url "$rpc" "$band" "owner()(address)")" = "$owner" ] &&
  [ "$(cast call --rpc-url "$rpc" "$band" "chainConfig()(address,bool)" | tr '\n' ' ')" = "$zero_address false " ] ||
  fail "the band's owner or chain configuration is not the registry's"
cast send --rpc-url "$rpc" --private-key "$key" --value 1ether "$next" > /dev/null
if out=$(cast call --rpc-url "$rpc" --from "$next" "$band" "transferOwnership(address)" "$next" 2>&1); then
  fail "a stranger started an ownership transfer"
fi
grep -q 0x118cdaa7 <<<"$out" || fail "expected OwnableUnauthorizedAccount, got $out"
cast send --rpc-url "$rpc" --private-key "$key" "$band" "transferOwnership(address)" "$next" > /dev/null
[ "$(cast call --rpc-url "$rpc" "$band" "owner()(address)") $(cast call --rpc-url "$rpc" "$band" "pendingOwner()(address)")" = "$owner $next" ] ||
  fail "a started transfer moved ownership before it was accepted"
cast send --rpc-url "$rpc" --private-key "$next_key" "$band" "acceptOwnership()" > /dev/null
[ "$(cast call --rpc-url "$rpc" "$band" "owner()(address)")" = "$next" ] || fail "an accepted transfer did not move ownership"
if out=$(cast call --rpc-url "$rpc" --from "$owner" "$band" "setHaltSigner(address)" "$owner" 2>&1); then
  fail "the previous owner rotated the halt signer"
fi
grep -q 0x118cdaa7 <<<"$out" || fail "expected OwnableUnauthorizedAccount, got $out"
if out=$(BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-band.sh" "$rpc" "$band" "$registry"); then
  fail "check-band.sh accepted a band whose owner is not the registry's"
fi
grep -q "FAIL: band's owner" <<<"$out" || fail "check-band.sh failed for another reason: $out"
hash=$(cast send --rpc-url "$rpc" --private-key "$next_key" "$band" "setHaltSigner(address)" "$next" --json | jq -r .transactionHash)
rotate_gas=$(cast rpc --rpc-url "$rpc" eth_getTransactionReceipt "$hash" | jq -r '"\(.gasUsed) \(.gasUsedForL1)"' | { read -r gas l1; echo $((gas - l1)); })
[ "$(cast call --rpc-url "$rpc" "$band" "haltSigner()(address)")" = "$next" ] || fail "the owner did not rotate the halt signer"
issued=$(next_second "$(block_time)")
halt_key=$next_key
write_halt true "$issued" $((issued + 30)) > /dev/null
[ "$(halt_view | cut -d' ' -f1)" = true ] || fail "the rotated signer's halt was not accepted"
echo "ownership: transferred in two steps; halt signer rotated with L2 gas $(record "setHaltSigner(address)" "$rotate_gas")"

sequencer=$(forge create --root "$root/contracts" test/devnode/StubAggregator.sol:StubAggregator --rpc-url "$rpc" \
  --private-key "$key" --broadcast --json --constructor-args 0 0 $(($(block_time) - 7200)) "L2 Sequencer Uptime Status Feed" |
  jq -r .deployedTo)
words=$(BAND_ASSETS="NVDA" "$(dirname "$0")/band-args.sh" "$registry")
read -r -a args <<<"$words"
args[4]="[$zero_address]" args[7]=$sequencer args[8]=true
arbitrum=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "${args[@]}" 2>&1 | grep 'deployed code at address' | grep -o '0x[0-9a-f]\{40\}')
[ "$(cast call --rpc-url "$rpc" "$arbitrum" "chainConfig()(address,bool)" | tr '\n' ' ')" = "$(cast to-check-sum-address "$sequencer") true " ] ||
  fail "a band deployed for Arbitrum One does not hold its chain configuration"
fresh_nvda "$arbitrum" > /dev/null
cast send --rpc-url "$rpc" --private-key "$key" "$arbitrum" "writePrices(bytes32[],bytes)" "$status_feeds" \
  "$(python3 "$payload" NY_MARKET_STATUS)" > /dev/null
arbitrum_state() { cast call --rpc-url "$rpc" "$arbitrum" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$symbol_nvda" | sed -n 1p; }
settled_view() { cast call --rpc-url "$rpc" "$arbitrum" "sequencerSettled()(bool)"; }
if out=$(BAND_ASSETS="NVDA" "$(dirname "$0")/check-band.sh" "$rpc" "$arbitrum" "$registry" 2>&1); then
  fail "check-band.sh accepted a band whose chain configuration is not the registry's"
fi
grep -q "FAIL: band's chain configuration" <<<"$out" || fail "check-band.sh failed for another reason: $out"
settled=$(arbitrum_state)
[ "$(settled_view)" = true ] && [ "$settled" != 0 ] || fail "NVDA's band on a settled sequencer is $settled"
record "quote(NVDA,sequencer)" "$(l2_gas "$(cast calldata "quote(bytes32)" "$symbol_nvda")" "$arbitrum")" > /dev/null
arbitrum_feed=$(feeds_of "$arbitrum" NVDA | awk '{ print $2 }')
[ "$(description "$arbitrum_feed")" = '"NVDA / USD Tapehouse band, low side (NVDA share)"' ] || fail "the Arbitrum One feed does not price the share"
round_of "$arbitrum_feed" > /dev/null
cast send --rpc-url "$rpc" --private-key "$key" "$sequencer" "setRound(int256,uint256)" 1 "$(block_time)" > /dev/null
[ "$(settled_view) $(arbitrum_state)" = "false 1" ] || fail "a down sequencer did not degrade the band"
expect_feed_revert "$arbitrum_feed" 0xc6b5066d
cast send --rpc-url "$rpc" --private-key "$key" "$sequencer" "setRound(int256,uint256)" 0 $(($(block_time) - 3500)) > /dev/null
[ "$(settled_view) $(arbitrum_state)" = "false 1" ] || fail "a sequencer back for less than an hour did not keep the band degraded"
expect_feed_revert "$arbitrum_feed" 0xc6b5066d
fresh_nvda "$arbitrum" > /dev/null
cast send --rpc-url "$rpc" --private-key "$key" "$sequencer" "setRound(int256,uint256)" 0 $(($(block_time) - 3700)) > /dev/null
[ "$(settled_view) $(arbitrum_state)" = "true $settled" ] || fail "a sequencer back for over an hour did not restore the band"
round=$(round_of "$arbitrum_feed") || fail "the feed did not answer once the sequencer settled"
read -r _ answer _ <<<"$round"
arbitrum_low=$(cast call --rpc-url "$rpc" "$arbitrum" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$symbol_nvda" | sed -n 5p | cut -d' ' -f1)
[ "$answer" = "$arbitrum_low" ] || fail "the feed answered $answer once the sequencer settled, quote's low is $arbitrum_low"
echo "Arbitrum One configuration: a down or recent sequencer degrades the band, and the feed refuses to answer"

tsla=$(cast format-bytes32-string TSLA---24_7)
read -r status _ <<<"$(write "[$tsla]" "$(python3 "$payload" TSLA---24_7)")"
[ "$status" = 0x1 ] || fail "writePrices reverted for TSLA---24_7"
tsla_first_ms=$(sample_ms "$tsla")
[ "$tsla_first_ms" -gt 0 ] || fail "the first TSLA---24_7 write is not a sample"
attempts=0
until tsla_payload=$(python3 "$payload" TSLA---24_7) &&
  cast call --rpc-url "$rpc" "$band" "writePrices(bytes32[],bytes)" "[$tsla]" "$tsla_payload" > /dev/null 2>&1; do
  attempts=$((attempts + 1))
  [ "$attempts" -le 8 ] || fail "no newer TSLA---24_7 package within 40 s"
  sleep 5
done
read -r status gas l1 _ <<<"$(write "[$tsla]" "$tsla_payload")"
[ "$status" = 0x1 ] || fail "writePrices reverted for TSLA---24_7"
read -r _ tsla_ms _ <<<"$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$tsla" | cut -d' ' -f1 | tr '\n' ' ')"
[ $((tsla_ms - tsla_first_ms)) -lt 50000 ] || fail "the gateway gave no TSLA---24_7 package within 50 s"
[ "$(sample_ms "$tsla")" = "$tsla_first_ms" ] || fail "a write within 50 s was sampled"
echo "TSLA---24_7 written again within 50 s of its first sample: not a sample, L2 gas $((gas - l1))"
unmeasured=$(grep -v ' [0-9][0-9]*$' "$gas_log" || true)
[ -z "$unmeasured" ] || fail "no L2 gas for: $unmeasured"
echo "PASS"
