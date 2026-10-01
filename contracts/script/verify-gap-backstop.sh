#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-gap-backstop.sh RPC_URL DEPLOYMENTS_JSON
# Verifies the registry's gap backstop on Sourcify, with the constructor arguments read back from it. The owner is read
# as the initial owner, and the exposure limits as the initial ones, so this runs before any ownership transfer or limit
# change.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
backstop=$(jq -er .tapehouse.GapBackstop "$registry")
accounts=$(cast call --rpc-url "$rpc" "$backstop" "accounts()(address)")
limit() { cast call --rpc-url "$rpc" "$backstop" "exposureLimit(bytes32)(uint256)" "$1" | cut -d' ' -f1; }
limits=
for symbol in $(cast call --rpc-url "$rpc" "$accounts" "stocks()(bytes32[],address[])" | head -n 1 | tr -d '[],'); do
  limits=${limits:+$limits,}$(limit "$symbol")
done
arguments=$(cast abi-encode "constructor(address,address,uint256,uint256[])" "$accounts" \
  "$(cast call --rpc-url "$rpc" "$backstop" "owner()(address)")" "$(limit 0x0000000000000000000000000000000000000000000000000000000000000000)" "[$limits]")
forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
  --constructor-args "$arguments" --watch "$backstop" src/GapBackstop.sol:GapBackstop
