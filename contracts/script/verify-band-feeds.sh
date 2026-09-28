#!/usr/bin/env bash
# Usage: verify-band-feeds.sh RPC_URL DEPLOYMENTS_JSON
# Verifies every BandFeed in the registry's .bandFeeds group on Sourcify, with the constructor arguments read
# back from each feed.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
chain_id=$(jq -r .chainId "$registry")
for asset in $(jq -r '.bandFeeds // {} | keys[]' "$registry"); do
  feed=$(jq -r --arg asset "$asset" '.bandFeeds[$asset]' "$registry")
  log=$(mktemp)
  arguments=$(cast abi-encode "constructor(address,bytes32,uint8,address,string)" \
    "$(cast call --rpc-url "$rpc" "$feed" "band()(address)")" \
    "$(cast call --rpc-url "$rpc" "$feed" "symbol()(bytes32)")" \
    "$(cast call --rpc-url "$rpc" "$feed" "side()(uint8)")" \
    "$(cast call --rpc-url "$rpc" "$feed" "pool()(address)")" \
    "$(cast call --rpc-url "$rpc" "$feed" "description()(string)" | jq -r .)")
  forge verify-contract --root "$root/contracts" --chain-id "$chain_id" --verifier sourcify \
    --constructor-args "$arguments" --watch "$feed" src/BandFeed.sol:BandFeed > "$log" 2>&1 || { tail -n 5 "$log"; exit 1; }
  rm -f "$log"
  echo "$asset $feed verified"
done
