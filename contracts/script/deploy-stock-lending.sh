#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: RATE_MODEL=<optimal>,<base>,<slope1>,<slope2> FEE_SHARE=<bps> deploy-stock-lending.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys a stock lending vault for each Stock Token of the registry's margin accounts, owned by the registry's owner,
# with the fee curve in basis points a year and the owner's share of the fee in basis points. Makes the accounts each
# vault's depositor and the vault the accounts' lending vault for that asset, and records it in the registry under
# .tapehouse.StockLending.<SYMBOL>. The borrower is set later. Checks first that the signer is the registry's owner and
# owns the accounts. The accounts take an asset's lending vault once: an asset that has one already is recorded as it is.
# SIGNER is forge's wallet flags for the accounts' owner, such as `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
: "${RATE_MODEL:?RATE_MODEL, the optimal utilisation, base, slope1 and slope2 in basis points}"
: "${FEE_SHARE:?FEE_SHARE, the share of the fee the owner takes, in basis points}"
accounts=$(field .tapehouse.MarginAccounts)
owner=$(field .tapehouse.Owner)
signer=$(cast wallet address "$@")
[ "$(lower "$signer")" = "$(lower "$owner")" ] || { echo "the signer $signer is not the registry's owner $owner" >&2; exit 1; }
[ "$(lower "$(cast call --rpc-url "$rpc" "$accounts" "owner()(address)")")" = "$(lower "$owner")" ] ||
  { echo "the accounts' owner is not $owner" >&2; exit 1; }
record() {
  jq --arg symbol "$1" --arg lending "$2" '.tapehouse.StockLending[$symbol] = $lending' "$registry" > "$registry.tmp"
  mv "$registry.tmp" "$registry"
}
stocks=$(cast call --rpc-url "$rpc" "$accounts" "stocks()(bytes32[],address[])")
read -ra symbols <<<"$(sed -n 1p <<<"$stocks" | tr -d '[],')"
read -ra tokens <<<"$(sed -n 2p <<<"$stocks" | tr -d '[],')"
deployed=0
for i in "${!symbols[@]}"; do
  symbol=${symbols[$i]} token=${tokens[$i]}
  [ "$token" != 0x0000000000000000000000000000000000000000 ] || continue
  name=$(cast parse-bytes32-string "$symbol")
  current=$(cast call --rpc-url "$rpc" "$accounts" "lending(bytes32)(address)" "$symbol")
  if [ "$current" != 0x0000000000000000000000000000000000000000 ]; then
    record "$name" "$current"
    echo "$name lends through $current, already set; recorded"
    continue
  fi
  created=$(forge create --root "$root/contracts" src/StockLendingVault.sol:StockLendingVault --rpc-url "$rpc" --broadcast \
    --json "$@" --constructor-args "$token" "$owner" "($RATE_MODEL)" "$FEE_SHARE")
  lending=$(jq -er .deployedTo <<<"$created")
  echo "StockLendingVault $lending for $name, deployed in $(jq -r .transactionHash <<<"$created")"
  cast send --rpc-url "$rpc" "$@" "$lending" "setDepositor(address)" "$accounts" > /dev/null
  cast send --rpc-url "$rpc" "$@" "$accounts" "setLending(bytes32,address)" "$symbol" "$lending" > /dev/null
  record "$name" "$lending"
  deployed=$((deployed + 1))
done
echo "the accounts lend $deployed Stock Tokens through new vaults, at a fee of ($RATE_MODEL) bps a year and $FEE_SHARE bps of it to the owner"
