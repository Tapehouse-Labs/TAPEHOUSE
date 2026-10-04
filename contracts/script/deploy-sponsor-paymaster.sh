#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: SPONSOR_SIGNER=<address> MAX_COST=<wei> DEPOSIT=<wei> STAKE=<wei> [UNSTAKE_DELAY=<seconds>] deploy-sponsor-paymaster.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys the sponsor paymaster on the registry's ERC-4337 EntryPoint, owned by the registry's owner, paying only for
# operations SPONSOR_SIGNER, the sponsor service's key, signs, and letting no operation cost it more than MAX_COST wei.
# Deposits DEPOSIT wei in the EntryPoint for it to pay with, and stakes STAKE wei, which bundlers require of a paymaster
# that reads its own storage, withdrawable UNSTAKE_DELAY seconds, a day by default, after the owner unlocks it. Records
# it under .tapehouse.SponsorPaymaster. Checks first that the signer is the registry's owner and that the registry's
# SimpleAccountFactory makes accounts on that EntryPoint. SIGNER is forge's wallet flags for the owner, such as
# `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
: "${SPONSOR_SIGNER:?SPONSOR_SIGNER, the address the sponsor service signs with}"
: "${MAX_COST:?MAX_COST, the most an operation may cost the paymaster, in wei}"
: "${DEPOSIT:?DEPOSIT, the wei the paymaster pays operations with}"
: "${STAKE:?STAKE, the wei the paymaster stakes}"
delay=${UNSTAKE_DELAY:-86400}
entry_point=$(field .erc4337.EntryPoint)
factory=$(field .erc4337.SimpleAccountFactory)
owner=$(field .tapehouse.Owner)
signer=$(cast wallet address "$@")
[ "$(lower "$signer")" = "$(lower "$owner")" ] || { echo "the signer $signer is not the registry's owner $owner" >&2; exit 1; }
implementation=$(cast call --rpc-url "$rpc" "$factory" "accountImplementation()(address)")
[ "$(lower "$(cast call --rpc-url "$rpc" "$implementation" "entryPoint()(address)")")" = "$(lower "$entry_point")" ] ||
  { echo "the registry's SimpleAccountFactory makes accounts on another EntryPoint than $entry_point" >&2; exit 1; }
created=$(forge create --root "$root/contracts" src/SponsorPaymaster.sol:SponsorPaymaster --rpc-url "$rpc" --broadcast \
  --json "$@" --constructor-args "$entry_point" "$owner" "$SPONSOR_SIGNER" "$MAX_COST")
paymaster=$(jq -er .deployedTo <<<"$created")
echo "SponsorPaymaster $paymaster, deployed in $(jq -r .transactionHash <<<"$created")"
cast send --rpc-url "$rpc" "$@" --value "$DEPOSIT" "$paymaster" "deposit()" > /dev/null
cast send --rpc-url "$rpc" "$@" --value "$STAKE" "$paymaster" "addStake(uint32)" "$delay" > /dev/null
jq --arg paymaster "$paymaster" '.tapehouse.SponsorPaymaster = $paymaster' "$registry" > "$registry.tmp"
mv "$registry.tmp" "$registry"
echo "it pays for what $SPONSOR_SIGNER signs, at most $MAX_COST wei an operation, from a deposit of $DEPOSIT wei and a stake of $STAKE wei unlockable in $delay seconds"
