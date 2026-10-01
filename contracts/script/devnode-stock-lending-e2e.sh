#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-stock-lending-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks SPY's stock lending vault on the dev node: it is the accounts' lending vault for SPY, takes deposits from them
# alone, and has no borrower yet. A position lends 4 of its 10 SPY: the tokens leave its holding for the vault, and its
# equity falls by the recall haircut on what it lent, while the engine still sees all 10. It takes them back, and its
# equity is what it was. Prints the L2 gas of the loan and of its return.
set -euo pipefail

rpc=$1 key=$2 registry=$3
accounts=$(jq -r .tapehouse.MarginAccounts "$registry")
lending=$(jq -r .tapehouse.StockLending.SPY "$registry")
spy=$(jq -r .tokens.SPY "$registry")
me=$(cast wallet address --private-key "$key")
position=$(cast format-bytes32-string SPY)
node_interface=0x00000000000000000000000000000000000000C8
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$key" "$@" > /dev/null; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }
calc() { python3 -c "print($1)"; }
health() { cast call --rpc-url "$rpc" "$accounts" "health(address,bytes32)(int256,uint256,uint8,uint8)" "$me" "$position" | sed -n 1,2p | cut -d' ' -f1 | tr '\n' ' '; }
l2_gas() {
  cast call --rpc-url "$rpc" $node_interface "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" \
    "$accounts" false "$1" --from "$me" | sed -n 1,2p | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}

[ "$(read_ "$accounts" "lending(bytes32)(address)" "$position")" = "$lending" ] || fail "the accounts do not lend SPY through the registry's vault"
[ "$(read_ "$lending" "asset()(address)")" = "$spy" ] || fail "the vault does not lend SPY"
[ "$(read_ "$lending" "depositor()(address)")" = "$accounts" ] || fail "the vault's depositor is not the accounts"
[ "$(read_ "$lending" "borrower()(address)")" = 0x0000000000000000000000000000000000000000 ] || fail "the vault has a borrower already"
[ "$(read_ "$lending" "maxDeposit(address)(uint256)" "$me")" = 0 ] || fail "the vault takes deposits from $me"
ten=10000000000000000000
four=4000000000000000000
held=$(read_ "$accounts" "collateral(address,bytes32,address)(uint256)" "$me" "$position" "$spy")
if [ "$(calc "$held < $ten")" = True ]; then
  send "$spy" "mint(address,uint256)" "$me" "$ten"
  send "$spy" "approve(address,uint256)" "$accounts" "$ten"
  send "$accounts" "deposit(bytes32,address,uint256,address)" "$position" "$spy" "$ten" "$me"
  held=$(read_ "$accounts" "collateral(address,bytes32,address)(uint256)" "$me" "$position" "$spy")
fi
read -r before requirement <<<"$(health)"
lend_gas=$(l2_gas "$(cast calldata "lend(bytes32,address,uint256,address)" "$position" "$spy" "$four" "$me")")
send "$accounts" "lend(bytes32,address,uint256,address)" "$position" "$spy" "$four" "$me"
[ "$(read_ "$accounts" "collateral(address,bytes32,address)(uint256)" "$me" "$position" "$spy")" = "$(calc "$held - $four")" ] ||
  fail "the holding did not fall by the 4 SPY lent"
[ "$(read_ "$accounts" "lent(address,bytes32,address)(uint256)" "$me" "$position" "$spy")" = "$four" ] || fail "the position did not lend 4 SPY"
[ "$(read_ "$spy" "balanceOf(address)(uint256)" "$lending")" = "$four" ] || fail "the vault does not hold the 4 SPY"
read -r lent_equity lent_requirement <<<"$(health)"
haircut=$(calc "round(($before - $lent_equity) * 10000 * $held / ($before * $four))")
[ "$lent_requirement" = "$requirement" ] || fail "the engine no longer sees all the SPY: $requirement became $lent_requirement"
[ "$haircut" = 500 ] || fail "lending 4 SPY took $(calc "$before - $lent_equity") off the equity, not 5% of their value"
unlend_gas=$(l2_gas "$(cast calldata "unlend(bytes32,address,uint256,address)" "$position" "$spy" "$four" "$me")")
send "$accounts" "unlend(bytes32,address,uint256,address)" "$position" "$spy" "$four" "$me"
read -r after _ <<<"$(health)"
[ "$after" = "$before" ] || fail "taking the SPY back left the equity at $after, not $before"
echo "stock lending $lending: SPY's, lent 4 SPY at a 5% haircut on $before of equity, the engine unchanged, and took them back; lend $lend_gas L2 gas, unlend $unlend_gas"
echo "PASS"
