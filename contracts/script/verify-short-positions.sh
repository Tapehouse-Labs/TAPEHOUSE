#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-short-positions.sh RPC_URL DEPLOYMENTS_JSON
# Verifies the registry's short positions on Sourcify, with the constructor arguments read back from them.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
shorts=$(jq -er .tapehouse.ShortPositions "$registry")
accounts=$(cast call --rpc-url "$rpc" "$shorts" "accounts()(address)")
symbols=$(cast call --rpc-url "$rpc" "$accounts" "stocks()(bytes32[],address[])" | sed -n 1p | tr -d '[],')
fees=()
for symbol in $symbols; do
  fees+=("$(cast call --rpc-url "$rpc" "$shorts" "fee(bytes32)(uint24)" "$symbol" 2>/dev/null || echo 0)")
done
arguments=$(cast abi-encode "constructor(address,address,uint24[])" "$accounts" \
  "$(cast call --rpc-url "$rpc" "$shorts" "router()(address)")" "[$(IFS=,; echo "${fees[*]}")]")
forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
  --constructor-args "$arguments" --watch "$shorts" src/ShortPositions.sol:ShortPositions
