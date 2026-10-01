#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-liquidator.sh RPC_URL DEPLOYMENTS_JSON
# Verifies the registry's liquidator on Sourcify, with its constructor argument, the accounts, read back from it.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
liquidator=$(jq -er .tapehouse.Liquidator "$registry")
arguments=$(cast abi-encode "constructor(address)" "$(cast call --rpc-url "$rpc" "$liquidator" "accounts()(address)")")
forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
  --constructor-args "$arguments" --watch "$liquidator" src/Liquidator.sol:Liquidator
