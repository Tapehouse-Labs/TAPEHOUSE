#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-margin-accounts-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks the margin accounts the registry names on the dev node against the band and margin programs: their caps
# (SPY's within the 1% that dividends paid since the deployment may move it) and guardian, SPY in its isolated
# position margined alone by the engine at its band's low edge, a loan refused under the guardian's pause, a loan
# against it, its liquidation price previewed before the loan and found after it, an untouched cross position, and
# the repayment of the debt and its premium, and the withdrawal; then the premium's rate, the closure an accrual
# records from the band's session as it stands, and that the accounts hold in USDG only their reserve and the premium
# set aside. Prints the L2 gas of the loan and of an accrual, and the gas of the liquidation-price search.
set -euo pipefail

rpc=$1 key=$2 registry=$3
accounts=$(jq -r .tapehouse.MarginAccounts "$registry")
vault=$(jq -r .tapehouse.SupplyVault "$registry")
band=$(jq -r .tapehouse.Band "$registry")
usdg=$(jq -r .tokens.USDG "$registry")
spy=$(jq -r .tokens.SPY "$registry")
me=$(cast wallet address --private-key "$key")
cross=0x0000000000000000000000000000000000000000000000000000000000000000
isolated=$(cast format-bytes32-string SPY)
node_interface=0x00000000000000000000000000000000000000C8
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$key" "$@" > /dev/null; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }
health() { cast call --rpc-url "$rpc" "$accounts" "health(address,bytes32)(int256,uint256,uint8,uint8)" "$me" "$1" | cut -d' ' -f1 | tr '\n' ' '; }
l2_gas() {
  cast call --rpc-url "$rpc" $node_interface "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" \
    "$accounts" false "$1" --from "$me" | head -n 2 | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}

[ "$(read_ "$accounts" "engine()(address)")" = "$(jq -r .tapehouse.Margin "$registry")" ] || fail "the accounts' engine is not the registry's"
[ "$(read_ "$accounts" "band()(address)")" = "$band" ] || fail "the accounts' band is not the registry's"
[ "$(read_ "$accounts" "vault()(address)")" = "$vault" ] || fail "the accounts' vault is not the registry's"
[ "$(read_ "$accounts" "debtCap()(uint256)")" = 1000000000000 ] && [ "$(read_ "$accounts" "weekendDebtCap()(uint256)")" = 500000000000 ] ||
  fail "the debt caps are not the deployment's"
[ "$(read_ "$accounts" "guardian()(address)")" = "$(jq -r .tapehouse.Owner "$registry")" ] || fail "the guardian is not the owner"
[ "$(read_ "$accounts" "premiumRate()(uint32)")" = 500 ] && [ "$(read_ "$accounts" "reserveShare()(uint16)")" = 1000 ] ||
  fail "the premium is not the deployment's"
depth=$(cast call --rpc-url "$rpc" "$(jq -r .tapehouse.Margin "$registry")" "depth(bytes32)(uint32,uint32,uint32,uint32)" "$isolated" | head -n 1 | cut -d' ' -f1)
feed_price=$(cast call --rpc-url "$rpc" "$(jq -r .chainlink.SPY_USD "$registry")" "latestRoundData()(uint80,int256,uint256,uint256,uint80)" | sed -n 2p | cut -d' ' -f1)
multiplier=$(cast call --rpc-url "$rpc" "$spy" "uiMultiplier()(uint256)" | cut -d' ' -f1)
cap=$(cast call --rpc-url "$rpc" "$accounts" "holding(bytes32)(uint256,uint256,uint256)" "$isolated" | sed -n 3p | cut -d' ' -f1)
python3 -c "import sys; sys.exit(abs($cap * $feed_price * $multiplier // 10**44 - $depth) * 100 > $depth)" ||
  fail "SPY's cap $cap is not its selling depth over its token's price, within the dividends paid since"
send "$usdg" "mint(address,uint256)" "$me" 10200000000
send "$usdg" "approve(address,uint256)" "$vault" 10000000000
send "$vault" "deposit(uint256,address)" 10000000000 "$me"
send "$spy" "mint(address,uint256)" "$me" 10000000000000000000
send "$spy" "approve(address,uint256)" "$accounts" 10000000000000000000
send "$accounts" "deposit(bytes32,address,uint256,address)" "$isolated" "$spy" 10000000000000000000 "$me"
low=$(cast call --rpc-url "$rpc" "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$isolated" | sed -n 5p | cut -d' ' -f1)
[ "$low" -gt 0 ] || fail "SPY's band is halted"
read -r equity requirement missing regime <<<"$(health "$isolated")"
[ "$equity" = "$((10 * low))0000000000" ] || fail "10 SPY at $low are not worth $equity"
[ "$requirement" != 0 ] && [ "$missing" = 0 ] && [ "$regime" != 0 ] || fail "health = $equity $requirement $missing $regime"
send "$accounts" "setBorrowingPaused(bool)" true
paused=$(cast call --rpc-url "$rpc" --from "$me" "$accounts" "borrow(bytes32,uint256,address,address)" "$isolated" 1 "$me" "$me" 2>&1) &&
  fail "a loan went through the guardian's pause"
grep -q "$(cast sig "BorrowingIsPaused()")" <<<"$paused" || fail "a paused loan reverted for another reason: $paused"
send "$accounts" "setBorrowingPaused(bool)" false
cross_health=$(health $cross)
preview=$(read_ "$accounts" "liquidationPrice(address,bytes32,bytes32,uint256)(uint256)" "$me" "$isolated" "$isolated" 100000000)
gas=$(l2_gas "$(cast calldata "borrow(bytes32,uint256,address,address)" "$isolated" 100000000 "$me" "$me")")
send "$accounts" "borrow(bytes32,uint256,address,address)" "$isolated" 100000000 "$me" "$me"
debt=$(read_ "$accounts" "debt(address,bytes32)(uint256)" "$me" "$isolated")
[ "$debt" -ge 100000000 ] && [ "$debt" -lt 100001000 ] || fail "the position owes $debt, not 100 USDG and a few seconds of interest"
[ "$(health $cross)" = "$cross_health" ] || fail "the cross position moved from $cross_health to $(health $cross)"
price=$(read_ "$accounts" "liquidationPrice(address,bytes32,bytes32,uint256)(uint256)" "$me" "$isolated" "$isolated" 0)
[ "$price" -gt 0 ] && [ "$price" -lt "$low" ] || fail "liquidation price $price against $low"
[ "$price" = "$preview" ] || fail "the liquidation price before the loan, $preview, is not $price after it"
search_gas=$(cast estimate --rpc-url "$rpc" --from "$me" "$accounts" "liquidationPrice(address,bytes32,bytes32,uint256)" "$me" "$isolated" "$isolated" 0)
send "$usdg" "approve(address,uint256)" "$accounts" 200000000
send "$accounts" "repay(bytes32,uint256,address)" "$isolated" 200000000 "$me"
[ "$(read_ "$accounts" "debt(address,bytes32)(uint256)" "$me" "$isolated")" = 0 ] || fail "the debt was not repaid"
[ "$(read_ "$accounts" "premium(address,bytes32)(uint256)" "$me" "$isolated")" = 0 ] || fail "the premium was not paid"
send "$accounts" "withdraw(bytes32,address,uint256,address,address)" "$isolated" "$spy" 10000000000000000000 "$me" "$me"
accrue_gas=$(l2_gas "$(cast calldata "accruePremium()")")
read -r state _ _ _ boundary <<<"$(cast call --rpc-url "$rpc" "$band" "session()(uint8,uint8,uint8,uint64,uint64)" | cut -d' ' -f1 | tr '\n' ' ')"
block=$(cast send --rpc-url "$rpc" --private-key "$key" --json "$accounts" "accruePremium()" | jq -r .blockNumber)
at=$(( $(cast block --rpc-url "$rpc" "$(cast to-dec "$block")" --field timestamp) * 1000 ))
read -r closes reopens accrued <<<"$(cast call --rpc-url "$rpc" "$accounts" "closure()(uint64,uint64,uint64)" | cut -d' ' -f1 | tr '\n' ' ')"
case $state in
  1) [ "$accrued" = "$at" ] && [ "$closes" -le "$at" ] && { [ "$boundary" = 0 ] || [ "$reopens" = "$boundary" ]; } ;;
  2) [ "$boundary" = 0 ] && { [ "$closes" = 0 ] || [ "$reopens" -le "$at" ]; } ||
    { [ "$closes" = "$boundary" ] && [ "$reopens" = 0 ]; } ;;
  *) [ "$closes" = 0 ] ;;
esac || fail "session $state $boundary left the closure at $closes $reopens $accrued"
held=$(read_ "$usdg" "balanceOf(address)(uint256)" "$accounts")
kept=$(( $(read_ "$accounts" "reserve()(uint128)") + $(read_ "$accounts" "backstopPremium()(uint128)") ))
[ "$held" = "$kept" ] || fail "the accounts hold $held USDG, not the $kept of their reserve and the premium set aside"
echo "margin accounts $accounts: SPY capped at $cap, 10 SPY at $low isolated with a requirement of $requirement in regime $regime, paused and resumed by the guardian, lent 100 USDG in $gas L2 gas, liquidation at $price, found in $search_gas gas, repaid and withdrawn; premium accrued in $accrue_gas L2 gas in session $state, $kept USDG kept"
echo "PASS"
