#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: deploy-gap-cover.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys the gap cover over the registry's margin program and USDG, and records it under .tapehouse.GapCover. Checks
# first that the margin program's band is the registry's and that the band names a Chainlink feed for at least one of
# its assets. SIGNER is forge's wallet flags, such as `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
zero=0x0000000000000000000000000000000000000000
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
margin=$(field .tapehouse.Margin)
band=$(field .tapehouse.Band)
usdg=$(field .tokens.USDG)
[ "$(lower "$(cast call --rpc-url "$rpc" "$margin" "band()(address)")")" = "$(lower "$band")" ] ||
  { echo "the margin program's band is not the registry's $band" >&2; exit 1; }
coverable=0
for symbol in $(cast call --rpc-url "$rpc" "$margin" "assets()(bytes32[])" | tr -d '[],'); do
  feed=$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "$symbol" | sed -n 1p)
  [ "$feed" = "$zero" ] || coverable=$((coverable + 1))
done
[ "$coverable" -gt 0 ] || { echo "the band names no Chainlink feed for any of the margin program's assets" >&2; exit 1; }
created=$(forge create --root "$root/contracts" src/GapCover.sol:GapCover --rpc-url "$rpc" --broadcast --json "$@" \
  --constructor-args "$margin" "$usdg")
cover=$(jq -er .deployedTo <<<"$created")
echo "GapCover $cover, deployed in $(jq -r .transactionHash <<<"$created"), covering $coverable assets"
jq --arg cover "$cover" '.tapehouse.GapCover = $cover' "$registry" > "$registry.tmp"
mv "$registry.tmp" "$registry"
