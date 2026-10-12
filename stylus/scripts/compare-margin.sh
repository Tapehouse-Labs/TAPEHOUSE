#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: compare-margin.sh RPC_URL DEPLOYMENTS_JSON
# Calls the margin program, .tapehouse.Margin, and its Solidity reference, .references.Margin, which the
# registry names for a chain where both are deployed, with the same arguments; stops at the first answer that
# differs by a single bit; then prints the L2 gas of each call as a Markdown table. It reads the chain and
# sends nothing, so it needs no key.
set -euo pipefail

rpc=$1 registry=$2
node_interface=0x00000000000000000000000000000000000000C8
margin=$(jq -er '.tapehouse.Margin' "$registry")
reference=$(jq -er '.references.Margin' "$registry")
band=$(jq -er '.tapehouse.Band' "$registry")
requirement="requirement(int256[],uint256[],uint64,bool)"
current="currentRequirement(int256[],uint256[])"
horizon=172800
checked=0

fail() { echo "FAIL: $*" >&2; exit 1; }
same() {
  local program solidity
  program=$(cast call --rpc-url "$rpc" "$margin" "$1" "${@:2}") || fail "$1 ${*:2}: the program reverted"
  solidity=$(cast call --rpc-url "$rpc" "$reference" "$1" "${@:2}") || fail "$1 ${*:2}: the reference reverted"
  [ "$program" = "$solidity" ] || fail "$1 ${*:2}: the program answers $program, the reference $solidity"
  checked=$((checked + 1))
}
l2_gas() {
  cast call --rpc-url "$rpc" $node_interface \
    "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" "$1" false "$2" |
    sed -n 1,2p | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}
row() {
  local label=$1 program solidity
  program=$(l2_gas "$margin" "$(cast calldata "${@:2}")")
  solidity=$(l2_gas "$reference" "$(cast calldata "${@:2}")")
  awk -v label="$label" -v p="$program" -v s="$solidity" 'function commas(n,  t, out) {
    t = n ""; out = ""; while (length(t) > 3) { out = "," substr(t, length(t) - 2) out; t = substr(t, 1, length(t) - 3) } return t out }
    BEGIN { printf "| `%s` | %s | %s | %.1f× |\n", label, commas(p), commas(s), s / p }'
}

assets=$(cast call --rpc-url "$rpc" "$margin" "assets()(bytes32[])" | tr -d '[] ')
IFS=, read -r -a symbols <<<"$assets"
[ "${#symbols[@]}" -ge 1 ] || fail "the margin program holds no assets"
first=${symbols[0]}
second=${symbols[$((${#symbols[@]} > 1 ? 1 : 0))]}

for view in "assets()(bytes32[])" "market()(bytes32)" "band()(address)" "ethUsdFeed()(address)" "owner()(address)" \
  "pendingOwner()(address)" "lastUpdate()(uint64)" "weekendLeverage()(uint32)"; do same "$view"; done
same "volatility(bytes32)(uint32,uint32)" "$first"
same "correlation(bytes32,bytes32)(uint16,uint16)" "$first" "$second"
same "weekendGap(bytes32)(uint32,uint32)" "$first"
same "pool(bytes32)(address)" "$first"
same "depth(bytes32)(uint32,uint32,uint32,uint32)" "$first"
for index in 0 1 255 256; do same "scenario(uint16,uint64)(int32[])" "$index" "$horizon"; done
for size in 32 64 128 256; do
  same "scenarioDigest(uint16,uint64)(bytes32)" "$size" "$horizon"
  same "scenarioDigest(uint16,uint64)(bytes32)" "$size" 0
done

quantities="["
prices="["
for i in "${!symbols[@]}"; do
  mid=$(cast call --rpc-url "$rpc" "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint64)" "${symbols[$i]}" | sed -n 3p | cut -d' ' -f1)
  [ "$mid" != 0 ] || fail "the band has no price for ${symbols[$i]} now"
  quantities+="$([ "$i" = 0 ] || echo ,)$([ $((i % 2)) = 0 ] && echo 1000000000000000000 || echo -500000000000000000)"
  prices+="$([ "$i" = 0 ] || echo ,)$mid"
done
quantities+="]"
prices+="]"
for h in 0 "$horizon" 302400; do
  for closure in false true; do same "$requirement(uint256,uint8)" "$quantities" "$prices" "$h" "$closure"; done
done
same "$current(uint256,uint8,uint8)" "$quantities" "$prices"
echo "PASS: the Solidity reference answers $checked calls of the margin program to the bit, on a portfolio of ${#symbols[@]} assets"
echo
echo "L2 gas, from the node's own estimate (NodeInterface), the same call on each:"
echo
echo "| Call | Stylus | Solidity | Solidity / Stylus |"
echo "|---|---|---|---|"
row "scenario(0)" "scenario(uint16,uint64)" 0 "$horizon"
for size in 32 64 128 256; do row "scenarioDigest($size)" "scenarioDigest(uint16,uint64)" "$size" "$horizon"; done
row "requirement(${#symbols[@]})" "$requirement" "$quantities" "$prices" "$horizon" false
row "currentRequirement(${#symbols[@]})" "$current" "$quantities" "$prices"
row "assets()" "assets()"
