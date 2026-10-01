#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: deploy-liquidator.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys the liquidator over the registry's margin accounts, makes it their liquidator, then records it in the
# registry under .tapehouse.Liquidator. Checks first that the signer is the registry's owner and owns the accounts.
# The accounts take a liquidator once: if they have one over them already, it records that one and deploys nothing.
# Uses forge create, as the accounts' script does.
# SIGNER is forge's wallet flags for the accounts' owner, such as `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
accounts=$(field .tapehouse.MarginAccounts)
owner=$(field .tapehouse.Owner)
signer=$(cast wallet address "$@")
[ "$(lower "$signer")" = "$(lower "$owner")" ] || { echo "the signer $signer is not the registry's owner $owner" >&2; exit 1; }
[ "$(lower "$(cast call --rpc-url "$rpc" "$accounts" "owner()(address)")")" = "$(lower "$owner")" ] ||
  { echo "the accounts' owner is not $owner" >&2; exit 1; }
record() {
  jq --arg liquidator "$1" '.tapehouse.Liquidator = $liquidator' "$registry" > "$registry.tmp"
  mv "$registry.tmp" "$registry"
}
current=$(cast call --rpc-url "$rpc" "$accounts" "liquidator()(address)")
if [ "$current" != 0x0000000000000000000000000000000000000000 ]; then
  [ "$(lower "$(cast call --rpc-url "$rpc" "$current" "accounts()(address)" 2>/dev/null || true)")" = "$(lower "$accounts")" ] ||
    { echo "the accounts already have a liquidator, $current, that is not over them" >&2; exit 1; }
  record "$current"
  echo "the accounts' liquidator is $current, already set; recorded"
  exit 0
fi
created=$(forge create --root "$root/contracts" src/Liquidator.sol:Liquidator --rpc-url "$rpc" --broadcast --json "$@" \
  --constructor-args "$accounts")
liquidator=$(jq -er .deployedTo <<<"$created")
echo "Liquidator $liquidator, deployed in $(jq -r .transactionHash <<<"$created")"
cast send --rpc-url "$rpc" "$@" "$accounts" "setLiquidator(address)" "$liquidator" > /dev/null
record "$liquidator"
echo "the accounts' liquidator is $liquidator"
