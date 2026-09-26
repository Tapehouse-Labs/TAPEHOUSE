#!/usr/bin/env bash
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
args=$(BAND_ASSETS="NVDA TSLA" "$(dirname "$0")/band-args.sh" "$registry.swapped")
read -r symbols feeds feed_ids index_ids <<<"$args"
swapped=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "$symbols" "$feeds" "$feed_ids" "$index_ids" 2>&1 | grep 'deployed code at address' | grep -o '0x[0-9a-f]\{40\}')
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
read -r cl_price cl_updated_at price_247 price_247_ms <<<"$(legs "$symbol_spy")"
[ "$cl_price" = 77232802713 ] && [ "$price_247" = 0 ] || fail "legs(SPY) = $cl_price $cl_updated_at $price_247 $price_247_ms before its first anchor"
echo "legs L2 gas: NVDA $(l2_gas "$(cast calldata "legs(bytes32)" "$symbol_nvda")"), unknown symbol $(l2_gas "$(cast calldata "legs(bytes32)" "$symbol_aapl")")"

again=0x5585258d$(cast abi-encode "f(bytes32[],address[],bytes32[],bytes32[])" "[$symbol_aapl,$symbol_spy]" \
  "[0x0000000000000000000000000000000000000000,0x0000000000000000000000000000000000000000]" \
  "[$(cast format-bytes32-string AAPL---24_7),$usa500]" "[$zero_id,$zero_id]" | cut -c3-)
(cd "$root/stylus/contracts/band" && cargo stylus get-initcode --output "$root/stylus/target/initcode.hex" > /dev/null 2>&1)
plain=$(cast send --rpc-url "$rpc" --private-key "$key" --create "0x$(tr -d '\n' < "$root/stylus/target/initcode.hex")" --json |
  jq -r .contractAddress)
[ "$(cast send --rpc-url "$rpc" --private-key "$key" "$plain" "$again" --json | jq -r .status)" = 0x1 ] ||
  fail "the first constructor call on a plain deploy failed"
[ "$(cast call --rpc-url "$rpc" "$plain" "asset(bytes32)(address,bytes32,bytes32)" "$symbol_aapl" | sed -n 2p)" = \
  "$(cast format-bytes32-string AAPL---24_7)" ] || fail "the first constructor call did not configure AAPL"
for program in "$plain" "$band"; do
  if out=$(cast call --rpc-url "$rpc" "$program" "$again" 2>&1); then
    fail "the constructor of $program ran a second time: $out"
  fi
  grep -q 'data: "0x"$' <<<"$out" || fail "the second constructor call of $program reverted for another reason: $out"
done
jq -n '{chainId: 412346}' > "$registry.redstone"
if out=$(BAND_ASSETS="AAPL SPY" "$(dirname "$0")/check-band.sh" "$rpc" "$plain" "$registry.redstone"); then
  fail "check-band.sh accepted a band that configures an asset band-args.sh leaves out"
fi
grep -q 'FAIL: band configures SPY' <<<"$out" || fail "check-band.sh failed for another reason: $out"

stub_18=$(forge create --root "$root/contracts" test/devnode/StubAggregator.sol:StubAggregator \
  --rpc-url "$rpc" --private-key "$key" --broadcast --json --constructor-args 18 1 1 "NVDA / USD" | jq -r .deployedTo)
if out=$(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "[$symbol_nvda]" "[$stub_18]" "[$nvda]" "[$zero_id]" 2>&1); then
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
change_time=$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$(cast format-bytes32-string NY_MARKET_NEXT_CHANGE_TIME)" | head -n 1 | cut -d' ' -f1)
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
  *) fail "session() = $session $nyse $nyse_next $change_ms $boundary_ms" ;;
esac
quote() {
  cast call --rpc-url "$rpc" "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$1" | cut -d' ' -f1 | tr '\n' ' '
}
read -r nvda_px _ <<<"$(cast call --rpc-url "$rpc" "$band" "price(bytes32)(uint256,uint64,uint64)" "$nvda" | cut -d' ' -f1 | tr '\n' ' ')"
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
echo "L2 gas: quote(NVDA) $(l2_gas "$(cast calldata "quote(bytes32)" "$symbol_nvda")"), quote(SPY) $(l2_gas "$(cast calldata "quote(bytes32)" "$symbol_spy")"), session() $(l2_gas "$(cast calldata "session()")")"

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
echo "PASS"
