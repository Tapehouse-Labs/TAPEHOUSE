#!/usr/bin/env bash
# Usage: band-args.sh DEPLOYMENTS_JSON
# Prints the band constructor arguments as three shell words: symbols, Chainlink feeds, RedStone feed IDs.
# Assets are BAND_ASSETS (default: the launch set). Feeds come from the file's .chainlink group. A file
# without that group configures no Chainlink legs; a launch asset missing from the group is an error.
# SPY has no RedStone 24/7 feed. An asset left with no leg at all is omitted.
set -euo pipefail

registry=$1
assets=${BAND_ASSETS:-NVDA TSLA AAPL MSFT GOOGL SPY}
zero_address=0x0000000000000000000000000000000000000000
zero_id=0x0000000000000000000000000000000000000000000000000000000000000000
has_chainlink=$(jq 'has("chainlink")' "$registry")
symbols="" feeds="" feed_ids=""

for asset in $assets; do
  feed=$zero_address
  if [ "$has_chainlink" = true ]; then
    feed=$(jq -er --arg key "${asset}_USD" '.chainlink[$key]' "$registry") ||
      { echo "$registry has no .chainlink.${asset}_USD" >&2; exit 1; }
  fi
  feed_id=$zero_id
  [ "$asset" = SPY ] || feed_id=$(cast format-bytes32-string "${asset}---24_7")
  [ "$feed" = $zero_address ] && [ "$feed_id" = $zero_id ] && continue
  symbols+=",$(cast format-bytes32-string "$asset")" feeds+=",$feed" feed_ids+=",$feed_id"
done

[ -n "$symbols" ] || { echo "$registry configures no asset with a leg" >&2; exit 1; }
echo "[${symbols#,}]" "[${feeds#,}]" "[${feed_ids#,}]"
