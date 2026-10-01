#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-reopening-auction.sh RPC_URL DEPLOYMENTS_JSON
# Verifies the registry's reopening auction on Sourcify, with the constructor arguments read back from it: its
# liquidator, and each Stock Token in .bandFeeds with the feed it holds for it.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
auction=$(jq -er .tapehouse.ReopeningAuction "$registry")
symbols=
feeds=
while IFS=' ' read -r name; do
  [ -n "$name" ] || continue
  symbol=$(cast format-bytes32-string "$name")
  symbols=${symbols:+$symbols,}$symbol
  feeds=${feeds:+$feeds,}$(cast call --rpc-url "$rpc" "$auction" "feeds(bytes32)(address)" "$symbol")
done <<<"$(jq -r '(.bandFeeds // {}) | keys_unsorted[]' "$registry")"
arguments=$(cast abi-encode "constructor(address,bytes32[],address[])" \
  "$(cast call --rpc-url "$rpc" "$auction" "liquidator()(address)")" "[$symbols]" "[$feeds]")
forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
  --constructor-args "$arguments" --watch "$auction" src/ReopeningAuction.sol:ReopeningAuction
