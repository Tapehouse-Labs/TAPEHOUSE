#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: BASKET=<KEY> NAME=<share name> SYMBOL=<share symbol> ASSETS=<ASSET>,... UNITS=<raw units>,... \
#   deploy-basket.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys a Basket of the registry's band's Stock Tokens, UNITS[i] raw units of ASSETS[i] per share as its first
# target, owned by the registry's owner, records it under .tapehouse.Baskets.<KEY>, and makes the margin accounts take
# its shares. The assets are listed in the engine's order whatever order ASSETS gives. Checks first that the signer is
# the registry's owner and owns the accounts, and that each asset is one of theirs with a Stock Token.
# SIGNER is forge's wallet flags for the accounts' owner, such as `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
zero=0x0000000000000000000000000000000000000000
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
: "${BASKET:?}" "${NAME:?}" "${SYMBOL:?}" "${ASSETS:?}" "${UNITS:?}"
accounts=$(field .tapehouse.MarginAccounts)
band=$(field .tapehouse.Band)
owner=$(field .tapehouse.Owner)
jq -e --arg key "$BASKET" '.tapehouse.Baskets[$key] // empty' "$registry" > /dev/null &&
  { echo "the registry has a basket $BASKET already" >&2; exit 1; }
signer=$(cast wallet address "$@")
[ "$(lower "$signer")" = "$(lower "$owner")" ] || { echo "the signer $signer is not the registry's owner $owner" >&2; exit 1; }
[ "$(lower "$(cast call --rpc-url "$rpc" "$accounts" "owner()(address)")")" = "$(lower "$owner")" ] ||
  { echo "the accounts' owner is not $owner" >&2; exit 1; }
IFS=, read -ra assets <<<"$ASSETS"
IFS=, read -ra units <<<"$UNITS"
[ "${#assets[@]}" = "${#units[@]}" ] || { echo "ASSETS and UNITS differ in length" >&2; exit 1; }
stocks=$(cast call --rpc-url "$rpc" "$accounts" "stocks()(bytes32[],address[])")
read -ra symbols <<<"$(sed -n 1p <<<"$stocks" | tr -d '[],')"
read -ra tokens <<<"$(sed -n 2p <<<"$stocks" | tr -d '[],')"
ordered=() amounts=()
for i in "${!symbols[@]}"; do
  name=$(cast parse-bytes32-string "${symbols[$i]}")
  for k in "${!assets[@]}"; do
    [ "${assets[$k]}" = "$name" ] || continue
    [ "${tokens[$i]}" != "$zero" ] || { echo "$name has no Stock Token in the accounts" >&2; exit 1; }
    ordered+=("${symbols[$i]}") amounts+=("${units[$k]}")
  done
done
[ "${#ordered[@]}" = "${#assets[@]}" ] || { echo "ASSETS names an asset the accounts do not take" >&2; exit 1; }
created=$(forge create --root "$root/contracts" src/Basket.sol:Basket --rpc-url "$rpc" --broadcast --json "$@" \
  --constructor-args "$NAME" "$SYMBOL" "$band" "[$(IFS=,; echo "${ordered[*]}")]" "[$(IFS=,; echo "${amounts[*]}")]" "$owner")
basket=$(jq -er .deployedTo <<<"$created")
echo "Basket $basket, $SYMBOL, deployed in $(jq -r .transactionHash <<<"$created")"
jq --arg key "$BASKET" --arg basket "$basket" '.tapehouse.Baskets[$key] = $basket' "$registry" > "$registry.tmp"
mv "$registry.tmp" "$registry"
cast send --rpc-url "$rpc" "$@" "$accounts" "addBasket(address)" "$basket" > /dev/null
echo "the margin accounts take $SYMBOL in their cross positions"
