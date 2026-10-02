#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: deploy-morpho-oracles.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys one MorphoBandOracle over the registry's band for each launch asset the band names a Stock Token for,
# pricing it in the registry's USDG and owned by the registry's owner, and records each under .morphoOracles.<ASSET>.
# Skips the assets the band prices as a share or does not configure, and stops if the band names another token than
# the registry. SIGNER is forge's wallet flags, such as `--account NAME --password-file FILE`. MORPHO_ASSETS overrides
# the launch assets.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
zero=0x0000000000000000000000000000000000000000
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
band=$(field .tapehouse.Band)
owner=$(field .tapehouse.Owner)
usdg=$(field .tokens.USDG)

for asset in ${MORPHO_ASSETS:-NVDA TSLA AAPL MSFT GOOGL SPY}; do
  symbol=$(cast format-bytes32-string "$asset")
  token=$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "$symbol" | sed -n 4p)
  [ "$token" != "$zero" ] || { echo "$asset: the band names no Stock Token for it, skipped" >&2; continue; }
  [ "$(lower "$token")" = "$(lower "$(jq -r --arg asset "$asset" '.tokens[$asset] // ""' "$registry")")" ] ||
    { echo "$asset: the band names $token, not the registry's token" >&2; exit 1; }
  created=$(forge create --root "$root/contracts" src/MorphoBandOracle.sol:MorphoBandOracle --rpc-url "$rpc" --broadcast \
    --json "$@" --constructor-args "$band" "$symbol" "$usdg" "$owner")
  oracle=$(jq -er .deployedTo <<<"$created")
  echo "MorphoBandOracle $oracle for $asset, deployed in $(jq -r .transactionHash <<<"$created")"
  jq --arg asset "$asset" --arg oracle "$oracle" '.morphoOracles[$asset] = $oracle' "$registry" > "$registry.tmp"
  mv "$registry.tmp" "$registry"
done
