#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-erc4337.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Deploys ERC-4337's EntryPoint v0.7 and its SimpleAccountFactory, built from eth-infinitism's v0.7.0 sources, on the
# dev node, which has neither, and records them under .erc4337. Robinhood Chain has both at their canonical addresses.
set -euo pipefail

rpc=$1 key=$2 registry=$3
root=$(cd "$(dirname "$0")/../.." && pwd)
create() {
  forge create --root "$root/contracts" --rpc-url "$rpc" --private-key "$key" --broadcast --json "$@" | jq -er .deployedTo
}
entry_point=$(create lib/account-abstraction/contracts/core/EntryPoint.sol:EntryPoint)
factory=$(create lib/account-abstraction/contracts/samples/SimpleAccountFactory.sol:SimpleAccountFactory \
  --constructor-args "$entry_point")
jq --arg entryPoint "$entry_point" --arg factory "$factory" \
  '.erc4337 = {EntryPoint: $entryPoint, SimpleAccountFactory: $factory}' "$registry" > "$registry.tmp"
mv "$registry.tmp" "$registry"
echo "EntryPoint $entry_point, SimpleAccountFactory $factory"
