#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: EXPOSURE_LIMITS=<cross>,<asset>,... SEED=<USDG units> deploy-gap-backstop.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys the gap backstop over the registry's margin accounts, owned by the registry's owner, with its exposure limits
# in USDG units a closure: the cross positions' first, then each asset's isolated positions' in the accounts' order.
# Makes it the accounts' backstop, records it in the registry under .tapehouse.GapBackstop, and deposits SEED USDG units
# from the signer for the owner: the team's seed, public on chain. Checks first that the signer is the registry's owner,
# owns the accounts and holds the seed. The accounts take a backstop once: if they have one over them already, it
# records that one, seeds it if nothing was deposited yet, and deploys nothing. Uses forge create, as the accounts'
# script does. SIGNER is forge's wallet flags for the accounts' owner, such as `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
: "${EXPOSURE_LIMITS:?EXPOSURE_LIMITS, the limit of the cross positions and then of each asset, in USDG units a closure}"
: "${SEED:?SEED, the first deposit of the team, in USDG units}"
accounts=$(field .tapehouse.MarginAccounts)
owner=$(field .tapehouse.Owner)
signer=$(cast wallet address "$@")
[ "$(lower "$signer")" = "$(lower "$owner")" ] || { echo "the signer $signer is not the registry's owner $owner" >&2; exit 1; }
[ "$(lower "$(cast call --rpc-url "$rpc" "$accounts" "owner()(address)")")" = "$(lower "$owner")" ] ||
  { echo "the accounts' owner is not $owner" >&2; exit 1; }
usdg=$(cast call --rpc-url "$rpc" "$accounts" "usdg()(address)")
holds_seed() {
  [ "$(cast call --rpc-url "$rpc" "$usdg" "balanceOf(address)(uint256)" "$signer" | cut -d' ' -f1)" -ge "$SEED" ] ||
    { echo "the signer holds less than the seed of $SEED USDG units" >&2; exit 1; }
}
record() {
  jq --arg backstop "$1" '.tapehouse.GapBackstop = $backstop' "$registry" > "$registry.tmp"
  mv "$registry.tmp" "$registry"
}
seed() {
  local backstop=$1
  shift
  [ "$SEED" != 0 ] && [ "$(cast call --rpc-url "$rpc" "$backstop" "totalSupply()(uint256)" | cut -d' ' -f1)" = 0 ] || return 0
  holds_seed
  cast send --rpc-url "$rpc" "$@" "$usdg" "approve(address,uint256)" "$backstop" "$SEED" > /dev/null
  local tx
  tx=$(cast send --rpc-url "$rpc" "$@" --json "$backstop" "deposit(uint256,address)" "$SEED" "$owner" | jq -r .transactionHash)
  echo "the team seeded it with $SEED USDG units in $tx"
}
current=$(cast call --rpc-url "$rpc" "$accounts" "backstop()(address)")
if [ "$current" != 0x0000000000000000000000000000000000000000 ]; then
  [ "$(lower "$(cast call --rpc-url "$rpc" "$current" "accounts()(address)" 2>/dev/null || true)")" = "$(lower "$accounts")" ] ||
    { echo "the accounts already have a backstop, $current, that is not over them" >&2; exit 1; }
  record "$current"
  echo "the accounts' backstop is $current, already set; recorded"
  seed "$current" "$@"
  exit 0
fi
cross=${EXPOSURE_LIMITS%%,*}
assets=${EXPOSURE_LIMITS#*,}
[ "$assets" != "$EXPOSURE_LIMITS" ] || assets=
count=$(cast call --rpc-url "$rpc" "$accounts" "stocks()(bytes32[],address[])" | sed -n 1p | tr -d '[] ' | tr ',' '\n' | grep -c . || true)
[ "$(tr ',' '\n' <<<"$assets" | grep -c . || true)" = "$count" ] ||
  { echo "EXPOSURE_LIMITS needs the cross positions' limit and $count assets' limits" >&2; exit 1; }
holds_seed
created=$(forge create --root "$root/contracts" src/GapBackstop.sol:GapBackstop --rpc-url "$rpc" --broadcast --json "$@" \
  --constructor-args "$accounts" "$owner" "$cross" "[$assets]")
backstop=$(jq -er .deployedTo <<<"$created")
echo "GapBackstop $backstop, deployed in $(jq -r .transactionHash <<<"$created")"
cast send --rpc-url "$rpc" "$@" "$accounts" "setBackstop(address)" "$backstop" > /dev/null
record "$backstop"
echo "the accounts' backstop is $backstop, covering at most $cross USDG units a closure for the cross positions and [$assets] for each asset's isolated positions"
seed "$backstop" "$@"
