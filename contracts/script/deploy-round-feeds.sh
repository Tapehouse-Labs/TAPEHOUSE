#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: ROUND_FEEDS="TSLA AMZN" deploy-round-feeds.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# For a chain without Chainlink feeds of its assets: deploys, for each asset in ROUND_FEEDS, a RedStone round feed of
# the Stock Token's price, which follows RedStone's data feed of the asset's name, its price in regular hours, scaled by
# the registry's Stock Token's multiplier, and records it as the asset's Chainlink feed, under .chainlink.<ASSET>_USD,
# with .chainlinkRegularHours set, which the band and the gap cover read as they do Arbitrum One's. An asset the registry
# names already is left as it is. The rounds are recorded by the keepers (see `rounds`). SIGNER is forge's wallet flags,
# such as `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
: "${ROUND_FEEDS:?ROUND_FEEDS, the assets to deploy a round feed for}"
for asset in $ROUND_FEEDS; do
  [ -z "$(jq -r --arg key "${asset}_USD" '.chainlink[$key] // empty' "$registry")" ] || { echo "$asset: already in the registry"; continue; }
  token=$(jq -er --arg asset "$asset" '.tokens[$asset]' "$registry") || { echo "$registry has no .tokens.$asset" >&2; exit 1; }
  created=$(forge create --root "$root/contracts" src/RedStoneRoundFeed.sol:RedStoneRoundFeed --rpc-url "$rpc" --broadcast --json "$@" \
    --constructor-args "$(cast format-bytes32-string "$asset")" "$asset / USD, RedStone regular hours" "$token")
  feed=$(jq -er .deployedTo <<<"$created")
  echo "$asset $feed, deployed in $(jq -r .transactionHash <<<"$created")"
  jq --arg key "${asset}_USD" --arg feed "$feed" '.chainlink[$key] = $feed | .chainlinkRegularHours = true' "$registry" > "$registry.tmp"
  mv "$registry.tmp" "$registry"
done
