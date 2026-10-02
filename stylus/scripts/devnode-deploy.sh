#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-deploy.sh RPC_URL PRIVATE_KEY
# Deploys stub Chainlink aggregators for NVDA, TSLA, SPY and ETH, each with a round a day over the nine days before
# its current one, a few tenths of a percent apart, stub Stock Tokens for NVDA and SPY, a stub
# USDG and WETH, stub NVDA/USDG and SPY/WETH pools holding the real pools' mean ticks and liquidity on
# 28 September 2026, the SPY pool with the Stock Token as token0, a stub SPY/USDG pool at SPY's stub price, whose entry
# gives the shorts SPY's fee tier, a stub swap router that trades NVDA and SPY for USDG at their stub prices and a stub
# QuoterV2 that quotes its swaps, writes them to stylus/target/devnode-registry.json with Anvil's second test account
# as the halt signer and the deploying account as the owner, and deploys the band and margin programs configured from
# that file through StylusDeployer, writing each program's address into the file as it is deployed. Writes each
# program's address to stylus/target/devnode-<program>.
set -euo pipefail

rpc=$1 key=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
target=$root/stylus/target
mkdir -p "$target"

stub() {
  local feed day
  feed=$(forge create --root "$root/contracts" test/devnode/StubAggregator.sol:StubAggregator --rpc-url "$rpc" \
    --private-key "$key" --broadcast --json --constructor-args 8 "$1" $((now - 9 * 86400)) "$2" | jq -r .deployedTo)
  for day in 8 7 6 5 4 3 2 1; do
    cast send --rpc-url "$rpc" --private-key "$key" "$feed" "setRound(int256,uint256)" \
      $(($1 * (10000 + (day % 2 * 2 - 1) * (20 + day * 5)) / 10000)) $((now - day * 86400)) > /dev/null
  done
  cast send --rpc-url "$rpc" --private-key "$key" "$feed" "setRound(int256,uint256)" "$1" "$now" > /dev/null
  echo "$feed"
}

token() {
  forge create --root "$root/contracts" test/devnode/StubStockToken.sol:StubStockToken --rpc-url "$rpc" \
    --private-key "$key" --broadcast --json --constructor-args "$1" | jq -r .deployedTo
}

create() {
  forge create --root "$root/contracts" "test/devnode/$1.sol:$1" --rpc-url "$rpc" --private-key "$key" --broadcast \
    --json --constructor-args "${@:2}" | jq -r .deployedTo
}

now=$(cast block --rpc-url "$rpc" latest -f timestamp)
nvda_token=$(token 1000775159164630595)
spy_token=$(token 1001717991187472003)
usdg=$(forge create --root "$root/contracts" test/devnode/StubUsdg.sol:StubUsdg --rpc-url "$rpc" --private-key "$key" \
  --broadcast --json | jq -r .deployedTo)
weth=$(create StubToken 18)
router=$(create StubSwapRouter "$usdg")
cast send --rpc-url "$rpc" --private-key "$key" "$router" "setPrice(address,uint256)" "$nvda_token" 22900000000 > /dev/null
cast send --rpc-url "$rpc" --private-key "$key" "$router" "setPrice(address,uint256)" "$spy_token" 77232802713 > /dev/null
jq -n --arg nvda "$(stub 22900000000 "NVDA / USD")" --arg tsla "$(stub 37800000000 "TSLA / USD")" \
  --arg spy "$(stub 77232802713 "SPY / USD")" --arg eth "$(stub 268330550000 "ETH / USD")" \
  --arg nvda_token "$nvda_token" --arg spy_token "$spy_token" --arg usdg "$usdg" --arg weth "$weth" \
  --arg nvda_pool "$(create StubPool "$usdg" "$nvda_token" 500 221989 11245526858841909681)" \
  --arg spy_pool "$(create StubPool "$spy_token" "$weth" 500 -12513 16029297629534329325587)" \
  --arg spy_usdg_pool "$(create StubPool "$spy_token" "$usdg" 500 -209827 823471432023715773)" \
  --arg router "$router" --arg quoter "$(create StubQuoterV2 "$router")" \
  --arg owner "$(cast wallet address --private-key "$key")" \
  '{chainId: 412346, chainlink: {NVDA_USD: $nvda, TSLA_USD: $tsla, SPY_USD: $spy, ETH_USD: $eth},
    tokens: {NVDA: $nvda_token, SPY: $spy_token, USDG: $usdg, WETH: $weth},
    uniswapV3: {NVDA_USDG_500: $nvda_pool, SPY_WETH_500: $spy_pool, SPY_USDG_500: $spy_usdg_pool, SwapRouter02: $router,
      QuoterV2: $quoter},
    tapehouse: {HaltSigner: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", Owner: $owner}}' \
  > "$target/devnode-registry.json"

deploy() {
  local program=$1 words args logged tx
  words=$(BAND_ASSETS="NVDA TSLA SPY" "$root/stylus/scripts/$program-args.sh" "$target/devnode-registry.json")
  read -r -a args <<<"$words"
  (cd "$root/stylus/contracts/$program" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
    --constructor-args "${args[@]}") > "$target/devnode-$program.log"
  cast to-check-sum-address "$(grep 'deployed code at address' "$target/devnode-$program.log" | grep -o '0x[0-9a-f]\{40\}')" \
    > "$target/devnode-$program"
  tx=$(grep 'deployment tx hash' "$target/devnode-$program.log" | grep -o '0x[0-9a-f]\{64\}')
  logged=$(cast receipt --rpc-url "$rpc" "$tx" --json | jq -r --arg topic "$(cast keccak 'ContractDeployed(address)')" \
    '.logs[] | select(.address == "0xcecba2f1dc234f70dd89f2041029807f8d03a990" and .topics[0] == $topic) | .data')
  [ "$(cast to-check-sum-address "0x${logged: -40}")" = "$(cat "$target/devnode-$program")" ] ||
    { echo "StylusDeployer's ContractDeployed log names 0x${logged: -40}, not $(cat "$target/devnode-$program")" >&2; exit 1; }
  echo "$program deployed on the dev node at $(cat "$target/devnode-$program")."
}

record() {
  jq --arg address "$(cat "$target/devnode-$1")" ".tapehouse.$2 = \$address" "$target/devnode-registry.json" > "$target/registry.tmp"
  mv "$target/registry.tmp" "$target/devnode-registry.json"
}

deploy band
record band Band
deploy margin
record margin Margin
