#!/usr/bin/env bash
# Usage: band-args.sh DEPLOYMENTS_JSON
# Prints the band constructor arguments as nine shell words: symbols, Chainlink feeds, RedStone feed IDs,
# index feed IDs, Stock Tokens, the halt signer, the owner, the L2 sequencer-uptime feed and whether
# Chainlink follows NYSE regular hours. Assets are BAND_ASSETS (default: the launch set). Feeds come from
# the file's .chainlink group. A file without that group configures no Chainlink legs; a launch asset
# missing from the group is an error. SPY has no RedStone 24/7 feed: its 24/7 leg is the S&P 500 index,
# USA500.Y---24_7, anchored to its Chainlink feed. Where SPY has no Chainlink feed it has no leg. Stock
# Tokens come from the file's .tokens group; a file without it, or without the asset, gives no token. An
# asset left with no leg at all is omitted. The halt signer and the owner are the file's
# .tapehouse.HaltSigner and .tapehouse.Owner, which must be set. The sequencer-uptime feed is
# .chainlinkSequencer.Uptime, zero where absent and refused when malformed. Arbitrum One's Chainlink feeds
# follow NYSE regular hours, so there the ninth word is true.
set -euo pipefail

registry=$1
assets=${BAND_ASSETS:-NVDA TSLA AAPL MSFT GOOGL SPY}
zero_address=0x0000000000000000000000000000000000000000
zero_id=0x0000000000000000000000000000000000000000000000000000000000000000
has_chainlink=$(jq 'has("chainlink")' "$registry")
chain_id=$(jq -r .chainId "$registry")
symbols="" feeds="" feed_ids="" index_ids="" tokens=""

for asset in $assets; do
  feed=$zero_address
  if [ "$has_chainlink" = true ]; then
    feed=$(jq -er --arg key "${asset}_USD" '.chainlink[$key]' "$registry") ||
      { echo "$registry has no .chainlink.${asset}_USD" >&2; exit 1; }
  fi
  feed_id=$zero_id index_id=$zero_id
  if [ "$asset" != SPY ]; then
    feed_id=$(cast format-bytes32-string "${asset}---24_7")
  elif [ "$feed" != $zero_address ]; then
    index_id=$(cast format-bytes32-string USA500.Y---24_7)
  fi
  [ "$feed" = $zero_address ] && [ "$feed_id" = $zero_id ] && continue
  token=$(jq -r --arg asset "$asset" --arg zero $zero_address '.tokens[$asset] // $zero' "$registry")
  symbols+=",$(cast format-bytes32-string "$asset")" feeds+=",$feed" feed_ids+=",$feed_id" index_ids+=",$index_id" tokens+=",$token"
done

[ -n "$symbols" ] || { echo "$registry configures no asset with a leg" >&2; exit 1; }
address() {
  jq -er --arg key "$1" 'getpath($key | split(".")) // "" | select(test("^0x[0-9a-fA-F]{40}$") and (test("^0x0{40}$") | not))' "$registry" ||
    { echo "$registry has no .$1" >&2; exit 1; }
}
halt_signer=$(address tapehouse.HaltSigner)
owner=$(address tapehouse.Owner)
sequencer=$zero_address
[ "$(jq 'has("chainlinkSequencer")' "$registry")" = false ] || sequencer=$(address chainlinkSequencer.Uptime)
regular_hours=false
[ "$chain_id" != 42161 ] || regular_hours=true
echo "[${symbols#,}]" "[${feeds#,}]" "[${feed_ids#,}]" "[${index_ids#,}]" "[${tokens#,}]" "$halt_signer" "$owner" "$sequencer" "$regular_hours"
