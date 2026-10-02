#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-baskets.sh RPC_URL DEPLOYMENTS_JSON
# Verifies each of the registry's baskets on Sourcify, with the constructor arguments read back from it. The target and
# the owner are read as the first ones, so this runs before any target change or ownership transfer.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
for basket in $(jq -r '.tapehouse.Baskets // {} | .[]' "$registry"); do
  symbols=$(cast call --rpc-url "$rpc" "$basket" "components()(bytes32[],address[])" | sed -n 1p)
  arguments=$(cast abi-encode "constructor(string,string,address,bytes32[],uint256[],address)" \
    "$(cast call --rpc-url "$rpc" "$basket" "name()(string)" | jq -r .)" \
    "$(cast call --rpc-url "$rpc" "$basket" "symbol()(string)" | jq -r .)" \
    "$(cast call --rpc-url "$rpc" "$basket" "band()(address)")" "$symbols" \
    "$(cast call --rpc-url "$rpc" "$basket" "target()(uint256[])" | sed -E 's/ \[[^]]*\]//g')" \
    "$(cast call --rpc-url "$rpc" "$basket" "owner()(address)")")
  forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
    --constructor-args "$arguments" --watch "$basket" src/Basket.sol:Basket
done
