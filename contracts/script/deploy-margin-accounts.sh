#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: DEBT_CAP=<USDG units> WEEKEND_DEBT_CAP=<USDG units> PREMIUM_RATE=<bps a year> RESERVE_SHARE=<bps> \
#   deploy-margin-accounts.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys the margin accounts over the registry's band, margin engine, supply vault and WETH, owned by the registry's
# owner, makes them the vault's borrower, then records them in the registry under .tapehouse.MarginAccounts, and makes
# the owner their guardian. Each asset's cap, in raw tokens, is the engine's selling depth in USD over the token's
# price, its Chainlink price times its ERC-8056 multiplier; zero for an asset without a Stock Token or a feed. The
# weekend premium starts at PREMIUM_RATE, and RESERVE_SHARE of each premium paid goes to the reserve.
# Checks first that the signer is the registry's owner and owns the vault, that the vault has no borrower, and that
# each feed answers a positive price at most a day and a minute old. Uses forge create: the
# constructor reads the Stylus programs, which forge script cannot run. SIGNER is forge's wallet flags for the
# vault's owner, such as `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
vault=$(field .tapehouse.SupplyVault)
owner=$(field .tapehouse.Owner)
band=$(field .tapehouse.Band)
margin=$(field .tapehouse.Margin)
: "${DEBT_CAP:?DEBT_CAP, the most USDG the accounts may owe, in USDG units}"
: "${WEEKEND_DEBT_CAP:?WEEKEND_DEBT_CAP, the most they may owe across a closure, in USDG units}"
: "${PREMIUM_RATE:?PREMIUM_RATE, the weekend premium, in basis points a year}"
: "${RESERVE_SHARE:?RESERVE_SHARE, the part of each premium paid kept as a reserve, in basis points}"
signer=$(cast wallet address "$@")
[ "$(tr '[:upper:]' '[:lower:]' <<<"$signer")" = "$(tr '[:upper:]' '[:lower:]' <<<"$owner")" ] ||
  { echo "the signer $signer is not the registry's owner $owner" >&2; exit 1; }
now=$(cast block --rpc-url "$rpc" latest --field timestamp)
caps=
for symbol in $(cast call --rpc-url "$rpc" "$margin" "assets()(bytes32[])" | tr -d '[],'); do
  read -r feed _ _ token <<<"$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "$symbol" | tr '\n' ' ')"
  cap=0
  if [ "$token" != 0x0000000000000000000000000000000000000000 ] && [ "$feed" != 0x0000000000000000000000000000000000000000 ]; then
    depth=$(cast call --rpc-url "$rpc" "$margin" "depth(bytes32)(uint32,uint32,uint32,uint32)" "$symbol" | head -n 1 | cut -d' ' -f1)
    read -r _ price _ updated _ <<<"$(cast call --rpc-url "$rpc" "$feed" "latestRoundData()(uint80,int256,uint256,uint256,uint80)" | cut -d' ' -f1 | tr '\n' ' ')"
    [ "$price" -gt 0 ] && [ "$updated" -ge $((now - 86460)) ] ||
      { echo "$(cast parse-bytes32-string "$symbol")'s feed answers $price, updated at $updated" >&2; exit 1; }
    multiplier=$(cast call --rpc-url "$rpc" "$token" "uiMultiplier()(uint256)" | cut -d' ' -f1)
    cap=$(python3 -c "import sys; print(int(sys.argv[1]) * 10**44 // (int(sys.argv[2]) * int(sys.argv[3])))" "$depth" "$price" "$multiplier")
  fi
  caps=${caps:+$caps,}$cap
done
[ "$(cast call --rpc-url "$rpc" "$vault" "owner()(address)")" = "$owner" ] || { echo "the vault's owner is not $owner" >&2; exit 1; }
[ "$(cast call --rpc-url "$rpc" "$vault" "borrower()(address)")" = 0x0000000000000000000000000000000000000000 ] ||
  { echo "the vault already has a borrower" >&2; exit 1; }
created=$(forge create --root "$root/contracts" src/MarginAccounts.sol:MarginAccounts --rpc-url "$rpc" --broadcast --json "$@" \
  --constructor-args "$band" "$margin" "$vault" "$(field .tokens.WETH)" "$owner" "[$caps]" "$DEBT_CAP" "$WEEKEND_DEBT_CAP" \
  "$PREMIUM_RATE" "$RESERVE_SHARE")
accounts=$(jq -er .deployedTo <<<"$created")
echo "MarginAccounts $accounts, deployed in $(jq -r .transactionHash <<<"$created")"
cast send --rpc-url "$rpc" "$@" "$vault" "setBorrower(address)" "$accounts" > /dev/null
echo "the vault's borrower is $accounts"
jq --arg accounts "$accounts" '.tapehouse.MarginAccounts = $accounts' "$registry" > "$registry.tmp"
mv "$registry.tmp" "$registry"
cast send --rpc-url "$rpc" "$@" "$accounts" "setGuardian(address)" "$owner" > /dev/null
echo "caps [$caps], debt $DEBT_CAP, weekend debt $WEEKEND_DEBT_CAP, premium $PREMIUM_RATE bps a year with $RESERVE_SHARE bps to the reserve; the guardian is $owner"
