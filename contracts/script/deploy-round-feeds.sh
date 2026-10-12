#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: ROUND_FEEDS="TSLA AMZN" deploy-round-feeds.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# For a chain without Chainlink feeds of its assets: deploys a RedStone round feed for each asset in ROUND_FEEDS, which
# follows RedStone's data feed of the asset's name, its price in regular hours, and records it under
# .roundFeeds.<ASSET>. An asset the registry names already is left as it is. The rounds are recorded by the keepers
# (see `rounds`), and the gap cover then settles on them. SIGNER is forge's wallet flags, such as
# `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
: "${ROUND_FEEDS:?ROUND_FEEDS, the assets to deploy a round feed for}"
for asset in $ROUND_FEEDS; do
  [ -z "$(jq -r --arg asset "$asset" '.roundFeeds[$asset] // empty' "$registry")" ] || { echo "$asset: already in the registry"; continue; }
  created=$(forge create --root "$root/contracts" src/RedStoneRoundFeed.sol:RedStoneRoundFeed --rpc-url "$rpc" --broadcast --json "$@" \
    --constructor-args "$(cast format-bytes32-string "$asset")" "$asset / USD, RedStone regular hours")
  feed=$(jq -er .deployedTo <<<"$created")
  echo "$asset $feed, deployed in $(jq -r .transactionHash <<<"$created")"
  jq --arg asset "$asset" --arg feed "$feed" '.roundFeeds[$asset] = $feed' "$registry" > "$registry.tmp"
  mv "$registry.tmp" "$registry"
done
