#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-reference-e2e.sh RPC_URL PRIVATE_KEY MARGIN_ADDRESS DEPLOYMENTS_JSON
# Deploys the Solidity reference of the margin engine, contracts/test/reference/MarginReference.sol, with the margin
# program's constructor arguments for NVDA, TSLA and SPY, checks that it answers every view the program answers to the
# bit, and measures both on the same calls. Deploys it again for the six launch assets and measures it there. Run it
# before devnode-margin-e2e.sh, which updates the program's parameters. Appends its L2 gas to
# stylus/target/devnode-gas.txt.
set -euo pipefail

rpc=$1 key=$2 margin=$3 registry=$4
root=$(cd "$(dirname "$0")/../.." && pwd)
node_interface=0x00000000000000000000000000000000000000C8
gas_log=$root/stylus/target/devnode-gas.txt
owner=$(cast wallet address --private-key "$key")
vectors=$root/stylus/contracts/margin/testdata/requirement-vectors.json
requirement="requirement(int256[],uint256[],uint64,bool)"
current="currentRequirement(int256[],uint256[])"

fail() { echo "FAIL: $*"; exit 1; }
record() {
  echo "$1 $2" >> "$gas_log"
  [[ $2 =~ ^[0-9]+$ ]] || { echo "FAIL: no L2 gas for $1" >&2; exit 1; }
  echo "$2"
}
l2_gas() {
  cast call --rpc-url "$rpc" --from "$owner" $node_interface \
    "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" "$2" false "$1" |
    head -n 2 | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}
measure() { record "$1" "$(l2_gas "$(cast calldata "$3" "${@:4}")" "$2")"; }
list() { jq -r "$1"' | map(tostring) | "[" + join(",") + "]"' "$vectors"; }
deploy_reference() {
  forge create --root "$root/contracts" test/reference/MarginReference.sol:MarginReference --rpc-url "$rpc" \
    --private-key "$key" --broadcast --json --constructor-args "$@" | jq -r .deployedTo
}
same() {
  local program solidity
  program=$(cast call --rpc-url "$rpc" "$margin" "$1" "${@:2}") || fail "$1 ${*:2}: the program reverted"
  solidity=$(cast call --rpc-url "$rpc" "$reference" "$1" "${@:2}") || fail "$1 ${*:2}: the reference reverted"
  [ "$program" = "$solidity" ] || fail "$1 ${*:2}: the program answers $program, the reference $solidity"
}

words=$(BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/margin-args.sh" "$registry")
read -r -a args <<<"$words"
reference=$(deploy_reference "${args[@]}")
nvda=$(cast format-bytes32-string NVDA)
spy=$(cast format-bytes32-string SPY)
quantities=$(list .devnode.quantities)
prices=$(list .devnode.prices)
same "assets()(bytes32[])"
same "market()(bytes32)"
same "band()(address)"
same "ethUsdFeed()(address)"
same "owner()(address)"
same "pendingOwner()(address)"
same "lastUpdate()(uint64)"
same "weekendLeverage()(uint32)"
same "volatility(bytes32)(uint32,uint32)" "$nvda"
same "correlation(bytes32,bytes32)(uint16,uint16)" "$nvda" "$spy"
same "weekendGap(bytes32)(uint32,uint32)" "$spy"
same "pool(bytes32)(address)" "$spy"
same "depth(bytes32)(uint32,uint32,uint32,uint32)" "$nvda"
for index in 0 255 256 767 768 771 772 777; do
  same "scenario(uint16,uint64)(int32[])" "$index" 172800
done
for size in 32 64 128 256; do
  same "scenarioDigest(uint16,uint64)(bytes32)" "$size" 172800
done
same "scenarioDigest(uint16,uint64)(bytes32)" 256 0
for horizon in 0 172800 302400; do
  for closure in false true; do
    same "$requirement(uint256,uint8)" "$quantities" "$prices" "$horizon" "$closure"
  done
done
same "$current(uint256,uint8,uint8)" "$quantities" "$prices"
same "$current(uint256,uint8,uint8)" "$(list '.devnode.quantities | [.[0], "0", "0"]')" "$prices"
echo "reference: the Solidity reference answers every view of the margin program to the bit"

for size in 32 64 128 256; do
  sizes+=" $size $(measure "reference.scenarioDigest($size)" "$reference" "scenarioDigest(uint16,uint64)" "$size" 172800),"
done
echo "reference scenarios: L2 gas by lattice size${sizes%,}; one scenario $(measure "reference.scenario(0)" "$reference" \
  "scenario(uint16,uint64)" 0 172800)"
echo "reference requirement: L2 gas $(measure "reference.requirement(3)" "$reference" "$requirement" "$quantities" "$prices" \
  172800 false) for three assets; current requirement $(measure "reference.currentRequirement(3)" "$reference" "$current" \
  "$quantities" "$prices")"
echo "reference views: L2 gas assets $(measure "reference.assets()" "$reference" "assets()"), volatility $(measure \
  "reference.volatility(NVDA)" "$reference" "volatility(bytes32)" "$nvda"), correlation $(measure \
  "reference.correlation(NVDA,SPY)" "$reference" "correlation(bytes32,bytes32)" "$nvda" "$spy")"

jq '.chainlink += {AAPL_USD: .chainlink.NVDA_USD, MSFT_USD: .chainlink.NVDA_USD, GOOGL_USD: .chainlink.NVDA_USD}' \
  "$registry" > "$registry.launch"
words=$("$(dirname "$0")/margin-args.sh" "$registry.launch")
read -r -a args <<<"$words"
reference=$(deploy_reference "${args[@]}")
launch=$(list '[.vectors[] | select(.label == "100k each of the six")][0].quantities')
isolated=$(list '[.vectors[] | select(.label == "100k each of the six")][0].quantities | [.[0]] + [range(5) | 0]')
echo "reference for the six launch assets: L2 gas scenarios $(measure "reference.scenarioDigest(256,launch)" "$reference" \
  "scenarioDigest(uint16,uint64)" 256 172800), requirement $(measure "reference.requirement(6)" "$reference" "$requirement" \
  "$launch" "$(list .prices)" 172800 false), current requirement of one position $(measure \
  "reference.currentRequirement(1 of 6)" "$reference" "$current" "$isolated" "$(list .prices)")"
unmeasured=$(grep -v ' [0-9][0-9]*$' "$gas_log" || true)
[ -z "$unmeasured" ] || fail "no L2 gas for: $unmeasured"
echo "PASS"
