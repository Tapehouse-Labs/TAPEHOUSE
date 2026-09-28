#!/usr/bin/env bash
# Usage: verify-band-feeds.sh RPC_URL DEPLOYMENTS_JSON
# Verifies every BandFeed in the registry's .bandFeeds group on Sourcify, with the constructor arguments read
# back from each feed.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
chain_id=$(jq -r .chainId "$registry")
log=$(mktemp)
trap 'rm -f "$log"' EXIT
for asset in $(jq -r '.bandFeeds // {} | keys[]' "$registry"); do
  feed=$(jq -r --arg asset "$asset" '.bandFeeds[$asset]' "$registry")
  band=$(cast call --rpc-url "$rpc" "$feed" "band()(address)")
  symbol=$(cast call --rpc-url "$rpc" "$feed" "symbol()(bytes32)")
  side=$(cast call --rpc-url "$rpc" "$feed" "side()(uint8)")
  pool=$(cast call --rpc-url "$rpc" "$feed" "pool()(address)")
  description=$(cast call --rpc-url "$rpc" "$feed" "description()(string)" | jq -r .)
  arguments=$(cast abi-encode "constructor(address,bytes32,uint8,address,string)" "$band" "$symbol" "$side" "$pool" "$description")
  forge verify-contract --root "$root/contracts" --chain-id "$chain_id" --verifier sourcify \
    --constructor-args "$arguments" --watch "$feed" src/BandFeed.sol:BandFeed > "$log" 2>&1 || { tail -n 5 "$log"; exit 1; }
  echo "$asset $feed verified"
done
