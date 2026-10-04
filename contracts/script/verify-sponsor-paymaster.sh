#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-sponsor-paymaster.sh RPC_URL DEPLOYMENTS_JSON
# Verifies the registry's sponsor paymaster on Sourcify, with the constructor arguments read back from it. The owner,
# the signer and the cost limit are read as the initial ones, so this runs before any of them changes.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
paymaster=$(jq -er .tapehouse.SponsorPaymaster "$registry")
arguments=$(cast abi-encode "constructor(address,address,address,uint128)" \
  "$(cast call --rpc-url "$rpc" "$paymaster" "entryPoint()(address)")" \
  "$(cast call --rpc-url "$rpc" "$paymaster" "owner()(address)")" \
  "$(cast call --rpc-url "$rpc" "$paymaster" "signer()(address)")" \
  "$(cast call --rpc-url "$rpc" "$paymaster" "maxCostPerOperation()(uint128)" | cut -d' ' -f1)")
forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
  --constructor-args "$arguments" --watch "$paymaster" src/SponsorPaymaster.sol:SponsorPaymaster
