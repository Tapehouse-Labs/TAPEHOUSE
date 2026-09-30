#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-margin-accounts.sh RPC_URL DEPLOYMENTS_JSON
# Verifies the registry's margin accounts on Sourcify, with the constructor arguments read back from them. The owner
# is read as the initial owner, and the premium rate as the initial one, so this runs before any ownership transfer
# or rate change.
set -euo pipefail

rpc=$1 registry=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
accounts=$(jq -er .tapehouse.MarginAccounts "$registry")
read_() { cast call --rpc-url "$rpc" "$accounts" "$1()(address)"; }
caps=
for symbol in $(cast call --rpc-url "$rpc" "$(read_ engine)" "assets()(bytes32[])" | tr -d '[],'); do
  caps=${caps:+$caps,}$(cast call --rpc-url "$rpc" "$accounts" "holding(bytes32)(uint256,uint256,uint256)" "$symbol" | sed -n 3p | cut -d' ' -f1)
done
uint_() { cast call --rpc-url "$rpc" "$accounts" "$1()($2)" | cut -d' ' -f1; }
arguments=$(cast abi-encode "constructor(address,address,address,address,address,uint256[],uint256,uint256,uint32,uint16)" \
  "$(read_ band)" "$(read_ engine)" "$(read_ vault)" "$(read_ weth)" "$(read_ owner)" "[$caps]" \
  "$(uint_ debtCap uint256)" "$(uint_ weekendDebtCap uint256)" "$(uint_ premiumRate uint32)" "$(uint_ reserveShare uint16)")
forge verify-contract --root "$root/contracts" --chain-id "$(jq -r .chainId "$registry")" --verifier sourcify \
  --constructor-args "$arguments" --watch "$accounts" src/MarginAccounts.sol:MarginAccounts
