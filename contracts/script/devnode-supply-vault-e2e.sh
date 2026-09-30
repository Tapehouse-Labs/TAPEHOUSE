#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-supply-vault-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks the supply vault the registry names on the dev node, freshly deployed with the margin accounts as its
# borrower, against its stub USDG: a deposit, a loan through a WETH-backed account at the model's rate, its
# repayment three seconds later, a USDG pause that closes the vault, and the lender's redemption with interest.
set -euo pipefail

rpc=$1 key=$2 registry=$3
vault=$(jq -r .tapehouse.SupplyVault "$registry")
accounts=$(jq -r .tapehouse.MarginAccounts "$registry")
usdg=$(jq -r .tokens.USDG "$registry")
weth=$(jq -r .tokens.WETH "$registry")
me=$(cast wallet address --private-key "$key")
cross=0x0000000000000000000000000000000000000000000000000000000000000000
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$key" "$@" > /dev/null; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }

[ "$(read_ "$vault" "asset()(address)")" = "$usdg" ] || fail "the vault's asset is not the registry's USDG"
[ "$(read_ "$vault" "owner()(address)")" = "$(jq -r .tapehouse.Owner "$registry")" ] || fail "the vault's owner is not the registry's"
[ "$(read_ "$vault" "borrower()(address)")" = "$accounts" ] || fail "the vault's borrower is not the margin accounts"
send "$usdg" "mint(address,uint256)" "$me" 2000000000
send "$usdg" "approve(address,uint256)" "$vault" 1000000000
send "$vault" "deposit(uint256,address)" 1000000000 "$me"
[ "$(read_ "$vault" "balanceOf(address)(uint256)" "$me")" = 1000000000000000 ] || fail "1,000 USDG did not mint 1e15 shares"
send "$weth" "mint(address,uint256)" "$me" 1000000000000000000
send "$weth" "approve(address,uint256)" "$accounts" 1000000000000000000
send "$accounts" "deposit(bytes32,address,uint256,address)" $cross "$weth" 1000000000000000000 "$me"
send "$accounts" "borrow(bytes32,uint256,address,address)" $cross 900000000 "$me" "$me"
[ "$(read_ "$vault" "utilization()(uint256)")" = 900000000000000000 ] || fail "900 of 1,000 lent is not 90% utilization"
[ "$(read_ "$vault" "borrowRate()(uint256)")" = 60000000000000000 ] || fail "the borrow rate at 90% is not 6%"
sleep 3
send "$usdg" "approve(address,uint256)" "$accounts" 1000000000
send "$accounts" "repay(bytes32,uint256,address)" $cross 1000000000 "$me"
[ "$(read_ "$vault" "debt()(uint256)")" = 0 ] || fail "the debt was not repaid"
idle=$(read_ "$vault" "idle()(uint256)")
[ "$idle" -gt 1000000000 ] || fail "the vault holds $idle, no interest over 1,000 USDG"
send "$usdg" "pause()"
[ "$(read_ "$vault" "maxDeposit(address)(uint256)" "$me")" = 0 ] || fail "a paused USDG left deposits open"
send "$usdg" "unpause()"
shares=$(read_ "$vault" "balanceOf(address)(uint256)" "$me")
before=$(read_ "$usdg" "balanceOf(address)(uint256)" "$me")
send "$vault" "redeem(uint256,address,address)" "$shares" "$me" "$me"
after=$(read_ "$usdg" "balanceOf(address)(uint256)" "$me")
echo "supply vault $vault: lent 900 of 1,000 USDG at 6%, repaid, paused and redeemed $((after - before)) with $((idle - 1000000000)) of interest"
echo "PASS"
