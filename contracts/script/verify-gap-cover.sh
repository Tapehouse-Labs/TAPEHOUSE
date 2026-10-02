#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-gap-cover.sh RPC_URL DEPLOYMENTS_JSON
# Verifies the registry's gap cover on Sourcify, with the constructor arguments read back from it.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
cover=$(jq -er .tapehouse.GapCover "$registry")
arguments=$(cast abi-encode "constructor(address,address)" \
  "$(cast call --rpc-url "$rpc" "$cover" "engine()(address)")" "$(cast call --rpc-url "$rpc" "$cover" "asset()(address)")")
forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
  --constructor-args "$arguments" --watch "$cover" src/GapCover.sol:GapCover
