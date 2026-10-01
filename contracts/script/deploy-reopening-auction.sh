#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: deploy-reopening-auction.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys the reopening auction over the registry's liquidator, with each Stock Token's band feed from .bandFeeds, whose
# sealed band sets a round's floor; makes it the liquidator's auction; and records it in the registry under
# .tapehouse.ReopeningAuction. Checks first that the signer is the registry's owner and owns the accounts. The liquidator
# takes an auction once: if it has one over it already, it records that one and deploys nothing. Uses forge create, as
# the accounts' script does. SIGNER is forge's wallet flags for the accounts' owner, such as
# `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
liquidator=$(field .tapehouse.Liquidator)
owner=$(field .tapehouse.Owner)
signer=$(cast wallet address "$@")
[ "$(lower "$signer")" = "$(lower "$owner")" ] || { echo "the signer $signer is not the registry's owner $owner" >&2; exit 1; }
accounts=$(cast call --rpc-url "$rpc" "$liquidator" "accounts()(address)")
[ "$(lower "$(cast call --rpc-url "$rpc" "$accounts" "owner()(address)")")" = "$(lower "$owner")" ] ||
  { echo "the accounts' owner is not $owner" >&2; exit 1; }
record() {
  jq --arg auction "$1" '.tapehouse.ReopeningAuction = $auction' "$registry" > "$registry.tmp"
  mv "$registry.tmp" "$registry"
}
current=$(cast call --rpc-url "$rpc" "$liquidator" "auction()(address)")
if [ "$current" != 0x0000000000000000000000000000000000000000 ]; then
  [ "$(lower "$(cast call --rpc-url "$rpc" "$current" "liquidator()(address)" 2>/dev/null || true)")" = "$(lower "$liquidator")" ] ||
    { echo "the liquidator already has an auction, $current, that is not over it" >&2; exit 1; }
  record "$current"
  echo "the liquidator's auction is $current, already set; recorded"
  exit 0
fi
symbols=
feeds=
while IFS=' ' read -r name feed; do
  [ -n "$name" ] || continue
  symbols=${symbols:+$symbols,}$(cast format-bytes32-string "$name")
  feeds=${feeds:+$feeds,}$feed
done <<<"$(jq -r '(.bandFeeds // {}) | to_entries[] | "\(.key) \(.value)"' "$registry")"
created=$(forge create --root "$root/contracts" src/ReopeningAuction.sol:ReopeningAuction --rpc-url "$rpc" --broadcast --json \
  "$@" --constructor-args "$liquidator" "[$symbols]" "[$feeds]")
auction=$(jq -er .deployedTo <<<"$created")
echo "ReopeningAuction $auction, deployed in $(jq -r .transactionHash <<<"$created")"
cast send --rpc-url "$rpc" "$@" "$liquidator" "setAuction(address)" "$auction" > /dev/null
record "$auction"
echo "the liquidator's auction is $auction, with the sealed bands of [$(jq -r '(.bandFeeds // {}) | keys_unsorted | join(",")' "$registry")]"
