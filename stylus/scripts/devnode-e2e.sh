#!/usr/bin/env bash
# Usage: devnode-e2e.sh RPC_URL PRIVATE_KEY BAND_ADDRESS DEPLOYMENTS_JSON
set -euo pipefail

rpc=$1 key=$2 band=$3 registry=$4
root=$(cd "$(dirname "$0")/../.." && pwd)
payload="$(dirname "$0")/redstone-payload.py"
symbol_nvda=$(cast format-bytes32-string NVDA)
symbol_tsla=$(cast format-bytes32-string TSLA)
symbol_aapl=$(cast format-bytes32-string AAPL)
node_interface=0x00000000000000000000000000000000000000C8
nvda=$(cast format-bytes32-string NVDA---24_7)
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
    "$band" false "$1" | head -n 2 | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}

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

BAND_ASSETS="NVDA TSLA" "$(dirname "$0")/check-band.sh" "$rpc" "$band" "$registry"
jq '.chainlink |= {NVDA_USD: .TSLA_USD, TSLA_USD: .NVDA_USD}' "$registry" > "$registry.swapped"
if out=$(BAND_ASSETS="NVDA TSLA" "$(dirname "$0")/check-band.sh" "$rpc" "$band" "$registry.swapped"); then
  fail "check-band.sh accepted a band that differs from the registry"
fi
grep -q 'FAIL: band holds' <<<"$out" || fail "check-band.sh failed for another reason: $out"
read -r symbols feeds feed_ids <<<"$(BAND_ASSETS="NVDA TSLA" "$(dirname "$0")/band-args.sh" "$registry.swapped")"
swapped=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "$symbols" "$feeds" "$feed_ids" 2>&1 | grep 'deployed code at address' | grep -o '0x[0-9a-f]\{40\}')
if out=$(BAND_ASSETS="NVDA TSLA" "$(dirname "$0")/check-band.sh" "$rpc" "$swapped" "$registry.swapped"); then
  fail "check-band.sh accepted feeds that describe other assets"
fi
grep -q 'describes itself as' <<<"$out" || fail "check-band.sh failed for another reason: $out"
read -r cl_price cl_updated_at price_247 price_247_ms <<<"$(legs "$symbol_nvda")"
[ "$cl_price" = 22900000000 ] && [ "$cl_updated_at" -gt 0 ] && [ "$price_247" = "$value" ] && [ "$price_247_ms" = "$package_ms" ] ||
  fail "legs(NVDA) = $cl_price $cl_updated_at $price_247 $price_247_ms"
read -r cl_price cl_updated_at price_247 price_247_ms <<<"$(legs "$symbol_tsla")"
[ "$cl_price" = 37800000000 ] && [ "$cl_updated_at" -gt 0 ] && [ "$price_247" = 0 ] && [ "$price_247_ms" = 0 ] ||
  fail "legs(TSLA) = $cl_price $cl_updated_at $price_247 $price_247_ms before any TSLA write"
[ "$(legs "$symbol_aapl")" = "0 0 0 0 " ] || fail "legs(AAPL) is not empty for an unknown symbol"
echo "legs L2 gas: NVDA $(l2_gas "$(cast calldata "legs(bytes32)" "$symbol_nvda")"), unknown symbol $(l2_gas "$(cast calldata "legs(bytes32)" "$symbol_aapl")")"

again=0x5585258d$(cast abi-encode "f(bytes32[],address[],bytes32[])" "[$symbol_aapl]" \
  "[0x0000000000000000000000000000000000000000]" "[$(cast format-bytes32-string AAPL---24_7)]" | cut -c3-)
(cd "$root/stylus/contracts/band" && cargo stylus get-initcode --output "$root/stylus/target/initcode.hex" > /dev/null 2>&1)
plain=$(cast send --rpc-url "$rpc" --private-key "$key" --create "0x$(tr -d '\n' < "$root/stylus/target/initcode.hex")" --json |
  jq -r .contractAddress)
[ "$(cast send --rpc-url "$rpc" --private-key "$key" "$plain" "$again" --json | jq -r .status)" = 0x1 ] ||
  fail "the first constructor call on a plain deploy failed"
[ "$(cast call --rpc-url "$rpc" "$plain" "asset(bytes32)(address,bytes32)" "$symbol_aapl" | tail -n 1)" = \
  "$(cast format-bytes32-string AAPL---24_7)" ] || fail "the first constructor call did not configure AAPL"
for program in "$plain" "$band"; do
  if out=$(cast call --rpc-url "$rpc" "$program" "$again" 2>&1); then
    fail "the constructor of $program ran a second time: $out"
  fi
done

stub_18=$(forge create --root "$root/contracts" test/devnode/StubAggregator.sol:StubAggregator \
  --rpc-url "$rpc" --private-key "$key" --broadcast --json --constructor-args 18 1 1 "NVDA / USD" | jq -r .deployedTo)
if out=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "[$symbol_nvda]" "[$stub_18]" "[$nvda]" 2>&1); then
  fail "a feed with 18 decimals was accepted"
fi
if ! grep -q 0x88d8f57d <<<"$out" || ! grep -q 6787b555 <<<"$out"; then
  fail "expected InvalidFeed inside ContractInitializationError, got: $(tail -n 1 <<<"$out")"
fi

deploy_tx=$(grep 'deployment tx hash' "$root/stylus/target/devnode-band.log" | grep -o '0x[0-9a-f]\{64\}')
cast rpc --rpc-url "$rpc" eth_getTransactionReceipt "$deploy_tx" | jq -r '"\(.gasUsed) \(.gasUsedForL1)"' |
  { read -r gas l1; echo "deploy, activation and constructor through StylusDeployer: L2 gas $((gas - l1))"; }

flip=$((${#nvda_payload} - 186))
tampered=${nvda_payload:0:flip}$(printf '%02x' $((0x${nvda_payload:flip:2} ^ 1)))${nvda_payload:flip+2}
expect_revert "[$nvda]" "$tampered" 0xec459bc0

read -r status gas l1 _ <<<"$(write "$status_feeds" "$(python3 "$payload" NY_MARKET_STATUS)")"
[ "$status" = 0x1 ] || fail "writePrices reverted for NY_MARKET_STATUS"
change_time=$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$(cast format-bytes32-string NY_MARKET_NEXT_CHANGE_TIME)" | head -n 1 | cut -d' ' -f1)
[ "$(cast to-base "$change_time" 16 | wc -c)" -gt 19 ] || fail "NY_MARKET_NEXT_CHANGE_TIME $change_time fits in 64 bits"
echo "NY_MARKET_STATUS: 3 feeds, NEXT_CHANGE_TIME = $change_time, L2 gas $((gas - l1))"
echo "PASS"
