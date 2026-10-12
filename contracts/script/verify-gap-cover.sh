#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-gap-cover.sh RPC_URL DEPLOYMENTS_JSON
# Verifies the registry's gap cover on Sourcify, with the constructor arguments read back from it.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
cover=$(jq -er .tapehouse.GapCover "$registry")
zero=0x0000000000000000000000000000000000000000
engine=$(cast call --rpc-url "$rpc" "$cover" "engine()(address)")
band=$(cast call --rpc-url "$rpc" "$cover" "band()(address)")
feeds="" round_feeds=0
for symbol in $(cast call --rpc-url "$rpc" "$engine" "assets()(bytes32[])" | tr -d '[],'); do
  feed=$(cast call --rpc-url "$rpc" "$cover" "feed(bytes32)(address)" "$symbol")
  own=$feed
  [ "$feed" != "$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "$symbol" | sed -n 1p)" ] || own=$zero
  [ "$own" = "$zero" ] || round_feeds=$((round_feeds + 1))
  feeds=${feeds:+$feeds,}$own
done
[ "$round_feeds" -gt 0 ] || feeds=
arguments=$(cast abi-encode "constructor(address,address,address[],bool)" "$engine" \
  "$(cast call --rpc-url "$rpc" "$cover" "asset()(address)")" "[$feeds]" \
  "$([ "$round_feeds" -gt 0 ] && cast call --rpc-url "$rpc" "$cover" "regularHours()(bool)" || echo false)")
forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
  --constructor-args "$arguments" --watch "$cover" src/GapCover.sol:GapCover
