#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: deploy-short-positions.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys the short positions over the registry's margin accounts and its Uniswap v3 SwapRouter02, each Stock Token
# shortable through the registry's one <SYMBOL>_USDG_<fee> pool and not shortable without one, records them under
# .tapehouse.ShortPositions, and makes them the borrower of each of the accounts' lending vaults that has none. Checks
# first that the signer is the registry's owner and owns the accounts, and so every lending vault the accounts took.
# SIGNER is forge's wallet flags for the accounts' owner, such as `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
zero=0x0000000000000000000000000000000000000000
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
lower() { tr '[:upper:]' '[:lower:]' <<<"$1"; }
accounts=$(field .tapehouse.MarginAccounts)
owner=$(field .tapehouse.Owner)
router=$(field .uniswapV3.SwapRouter02)
signer=$(cast wallet address "$@")
[ "$(lower "$signer")" = "$(lower "$owner")" ] || { echo "the signer $signer is not the registry's owner $owner" >&2; exit 1; }
[ "$(lower "$(cast call --rpc-url "$rpc" "$accounts" "owner()(address)")")" = "$(lower "$owner")" ] ||
  { echo "the accounts' owner is not $owner" >&2; exit 1; }
stocks=$(cast call --rpc-url "$rpc" "$accounts" "stocks()(bytes32[],address[])")
read -ra symbols <<<"$(sed -n 1p <<<"$stocks" | tr -d '[],')"
read -ra tokens <<<"$(sed -n 2p <<<"$stocks" | tr -d '[],')"
fees=()
for i in "${!symbols[@]}"; do
  name=$(cast parse-bytes32-string "${symbols[$i]}")
  pools=$(jq -r --arg prefix "${name}_USDG_" '.uniswapV3 | keys[] | select(startswith($prefix))' "$registry")
  [ "$(grep -c . <<<"$pools")" -le 1 ] || { echo "the registry names more than one ${name}/USDG pool" >&2; exit 1; }
  if [ -n "$pools" ] && [ "${tokens[$i]}" != "$zero" ]; then fees+=("${pools##*_}"); else fees+=(0); fi
done
list=$(IFS=,; echo "[${fees[*]}]")
created=$(forge create --root "$root/contracts" src/ShortPositions.sol:ShortPositions --rpc-url "$rpc" --broadcast \
  --json "$@" --constructor-args "$accounts" "$router" "$list")
shorts=$(jq -er .deployedTo <<<"$created")
echo "ShortPositions $shorts, deployed in $(jq -r .transactionHash <<<"$created"), at fees $list"
jq --arg shorts "$shorts" '.tapehouse.ShortPositions = $shorts' "$registry" > "$registry.tmp"
mv "$registry.tmp" "$registry"
borrowing=0
for i in "${!symbols[@]}"; do
  name=$(cast parse-bytes32-string "${symbols[$i]}")
  [ "${tokens[$i]}" != "$zero" ] || continue
  lending=$(cast call --rpc-url "$rpc" "$accounts" "lending(bytes32)(address)" "${symbols[$i]}")
  [ "$lending" != "$zero" ] || continue
  borrower=$(cast call --rpc-url "$rpc" "$lending" "borrower()(address)")
  if [ "$borrower" != "$zero" ]; then
    echo "$name's lending vault $lending has the borrower $borrower already; left as it is"
    continue
  fi
  cast send --rpc-url "$rpc" "$@" "$lending" "setBorrower(address)" "$shorts" > /dev/null
  borrowing=$((borrowing + 1))
done
echo "the shorts borrow from $borrowing lending vaults"
