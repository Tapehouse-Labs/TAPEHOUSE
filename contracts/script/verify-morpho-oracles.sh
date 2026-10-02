#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-morpho-oracles.sh RPC_URL DEPLOYMENTS_JSON
# Verifies every MorphoBandOracle in the registry's .morphoOracles group on Sourcify, with the constructor arguments
# read back from each oracle. Run it before an oracle is re-pointed or its ownership moves, while its band and owner
# are still the ones it was deployed with.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
chain_id=$(jq -r .chainId "$registry")
log=$(mktemp)
trap 'rm -f "$log"' EXIT
for asset in $(jq -r '.morphoOracles // {} | keys[]' "$registry"); do
  oracle=$(jq -r --arg asset "$asset" '.morphoOracles[$asset]' "$registry")
  arguments=$(cast abi-encode "constructor(address,bytes32,address,address)" \
    "$(cast call --rpc-url "$rpc" "$oracle" "band()(address)")" \
    "$(cast call --rpc-url "$rpc" "$oracle" "symbol()(bytes32)")" \
    "$(cast call --rpc-url "$rpc" "$oracle" "loanToken()(address)")" \
    "$(cast call --rpc-url "$rpc" "$oracle" "owner()(address)")")
  forge verify-contract --root "$root/contracts" --chain-id "$chain_id" --verifier sourcify \
    --constructor-args "$arguments" --watch "$oracle" src/MorphoBandOracle.sol:MorphoBandOracle > "$log" 2>&1 ||
    { tail -n 5 "$log"; exit 1; }
  echo "$asset $oracle verified"
done
