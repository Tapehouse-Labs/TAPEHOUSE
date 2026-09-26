#!/usr/bin/env bash
# Usage: devnode-e2e.sh RPC_URL PRIVATE_KEY BAND_ADDRESS
set -euo pipefail

rpc=$1 key=$2 band=$3
payload="$(dirname "$0")/redstone-payload.py"
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

flip=$((${#nvda_payload} - 186))
tampered=${nvda_payload:0:flip}$(printf '%02x' $((0x${nvda_payload:flip:2} ^ 1)))${nvda_payload:flip+2}
expect_revert "[$nvda]" "$tampered" 0xec459bc0

read -r status gas l1 _ <<<"$(write "$status_feeds" "$(python3 "$payload" NY_MARKET_STATUS)")"
[ "$status" = 0x1 ] || fail "writePrices reverted for NY_MARKET_STATUS"
change_time=$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$(cast format-bytes32-string NY_MARKET_NEXT_CHANGE_TIME)" | head -n 1 | cut -d' ' -f1)
[ "$(cast to-base "$change_time" 16 | wc -c)" -gt 19 ] || fail "NY_MARKET_NEXT_CHANGE_TIME $change_time fits in 64 bits"
echo "NY_MARKET_STATUS: 3 feeds, NEXT_CHANGE_TIME = $change_time, L2 gas $((gas - l1))"
echo "PASS"
