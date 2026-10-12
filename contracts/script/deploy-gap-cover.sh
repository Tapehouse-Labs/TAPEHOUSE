#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: deploy-gap-cover.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys the gap cover over the registry's margin program and USDG, and records it under .tapehouse.GapCover. An asset
# settles on the registry's RedStone round feed, .roundFeeds.<ASSET>, where it names one, which follow NYSE's regular
# hours, and otherwise on the Chainlink feed of the margin program's band. Checks first that the margin program's band
# is the registry's and that at least one asset has a feed of either kind. SIGNER is forge's wallet flags, such as
# `--account NAME --password-file FILE`.
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
coverable=0 feeds="" round_feeds=0
for symbol in $(cast call --rpc-url "$rpc" "$margin" "assets()(bytes32[])" | tr -d '[],'); do
  feed=$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "$symbol" | sed -n 1p)
  round=$(jq -r --arg asset "$(cast parse-bytes32-string "$symbol")" '.roundFeeds[$asset] // empty' "$registry")
  [ "$feed" = "$zero" ] && [ -z "$round" ] || coverable=$((coverable + 1))
  [ -z "$round" ] || round_feeds=$((round_feeds + 1))
  feeds=${feeds:+$feeds,}${round:-$zero}
done
[ "$coverable" -gt 0 ] || { echo "the band names no Chainlink feed, and the registry no round feed, for any of the margin program's assets" >&2; exit 1; }
[ "$round_feeds" -gt 0 ] || feeds=
created=$(forge create --root "$root/contracts" src/GapCover.sol:GapCover --rpc-url "$rpc" --broadcast --json "$@" \
  --constructor-args "$margin" "$usdg" "[$feeds]" "$([ "$round_feeds" -gt 0 ] && echo true || echo false)")
cover=$(jq -er .deployedTo <<<"$created")
echo "GapCover $cover, deployed in $(jq -r .transactionHash <<<"$created"), covering $coverable assets, $round_feeds of them by a RedStone round feed"
jq --arg cover "$cover" '.tapehouse.GapCover = $cover' "$registry" > "$registry.tmp"
mv "$registry.tmp" "$registry"
