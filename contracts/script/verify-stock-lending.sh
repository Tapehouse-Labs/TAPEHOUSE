#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-stock-lending.sh RPC_URL DEPLOYMENTS_JSON
# Verifies each of the registry's stock lending vaults on Sourcify, with the constructor arguments read back from it.
# The owner and the rate model are read as the initial ones, so this runs before any ownership transfer or rate change.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
for lending in $(jq -r '.tapehouse.StockLending // {} | .[]' "$registry"); do
  model=$(cast call --rpc-url "$rpc" "$lending" "rateModel()((uint16,uint32,uint32,uint32))")
  arguments=$(cast abi-encode "constructor(address,address,(uint16,uint32,uint32,uint32),uint16)" \
    "$(cast call --rpc-url "$rpc" "$lending" "asset()(address)")" "$(cast call --rpc-url "$rpc" "$lending" "owner()(address)")" \
    "$model" "$(cast call --rpc-url "$rpc" "$lending" "feeShare()(uint16)")")
  forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
    --constructor-args "$arguments" --watch "$lending" src/StockLendingVault.sol:StockLendingVault
done
