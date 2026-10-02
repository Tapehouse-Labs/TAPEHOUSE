#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-basket-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks the registry's basket PAIR on the dev node: it holds NVDA and SPY at its first target, owned by the registry's
# owner and taken by the margin accounts. A fresh address, its key made at run time, mints two shares in kind and
# deposits them in its cross position, where the margin program margins them as the 2 NVDA and 1 SPY they redeem for,
# at their bands' low edges; while NVDA's band is halted they back no loan. Unwrapped, the position holds those tokens
# with the same equity and requirement; with NVDA taken out, it borrows against SPY and repays. A share minted and
# redeemed returns what it took. Prints the L2 gas of a mint, a deposit and an unwrap.
set -euo pipefail

rpc=$1 key=$2 registry=$3
basket=$(jq -er .tapehouse.Baskets.PAIR "$registry")
accounts=$(jq -r .tapehouse.MarginAccounts "$registry")
margin=$(jq -r .tapehouse.Margin "$registry")
band=$(jq -r .tapehouse.Band "$registry")
nvda=$(jq -r .tokens.NVDA "$registry")
spy=$(jq -r .tokens.SPY "$registry")
usdg=$(jq -r .tokens.USDG "$registry")
user_key=$(cast wallet new --json | jq -er '.data[0].private_key')
me=$(cast wallet address --private-key "$user_key")
cross=0x0000000000000000000000000000000000000000000000000000000000000000
node_interface=0x00000000000000000000000000000000000000C8
one=1000000000000000000
two=2000000000000000000
half=500000000000000000
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$user_key" "$@" > /dev/null; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }
list() { cast call --rpc-url "$rpc" "$@" | sed -E 's/ \[[^]]*\]//g'; }
health() { cast call --rpc-url "$rpc" "$accounts" "health(address,bytes32)(int256,uint256,uint8,uint8)" "$me" $cross | cut -d' ' -f1 | tr '\n' ' '; }
l2_gas() {
  cast call --rpc-url "$rpc" $node_interface "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" \
    "$1" false "$2" --from "$me" | sed -n 1,2p | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}

components=$(cast call --rpc-url "$rpc" "$basket" "components()(bytes32[],address[])" | sed -n 2p)
[ "$components" = "[$nvda, $spy]" ] || fail "the basket holds $components, not NVDA and SPY"
[ "$(list "$basket" "target()(uint256[])")" = "[$one, $half]" ] || fail "the basket's target is not 1 NVDA and 0.5 SPY a share"
[ "$(read_ "$basket" "owner()(address)")" = "$(jq -r .tapehouse.Owner "$registry")" ] || fail "the basket's owner is not the registry's"
grep -qi "$basket" <<<"$(cast call --rpc-url "$rpc" "$accounts" "baskets()(address[])")" || fail "the accounts do not take the basket"
cast send --rpc-url "$rpc" --private-key "$key" --value 1ether "$me" > /dev/null

send "$nvda" "mint(address,uint256)" "$me" 3000000000000000000
send "$spy" "mint(address,uint256)" "$me" 1500000000000000000
send "$nvda" "approve(address,uint256)" "$basket" 3000000000000000000
send "$spy" "approve(address,uint256)" "$basket" 1500000000000000000
[ "$(list "$basket" "previewMint(uint256)(uint256[])" $two)" = "[$two, $one]" ] || fail "two shares do not take 2 NVDA and 1 SPY"
mint_gas=$(l2_gas "$basket" "$(cast calldata "mint(uint256,address,uint256[])" $two "$me" "[$two,$one]")")
send "$basket" "mint(uint256,address,uint256[])" $two "$me" "[$two,$one]"
send "$basket" "approve(address,uint256)" "$accounts" $two
deposit_gas=$(l2_gas "$accounts" "$(cast calldata "deposit(bytes32,address,uint256,address)" $cross "$basket" $two "$me")")
send "$accounts" "deposit(bytes32,address,uint256,address)" $cross "$basket" $two "$me"

stocks=$(cast call --rpc-url "$rpc" "$accounts" "stocks()(bytes32[],address[])")
read -ra symbols <<<"$(sed -n 1p <<<"$stocks" | tr -d '[],')"
read -ra tokens <<<"$(sed -n 2p <<<"$stocks" | tr -d '[],')"
read -ra held <<<"$(list "$accounts" "inBaskets(address,bytes32)(uint256[])" "$me" $cross | tr -d '[],')"
quantities=() prices=()
for i in "${!symbols[@]}"; do
  case "${tokens[$i]}" in
    "$nvda") [ "${held[$i]}" = $two ] || fail "the basket holds ${held[$i]} NVDA for the position, not 2" ;;
    "$spy") [ "${held[$i]}" = $one ] || fail "the basket holds ${held[$i]} SPY for the position, not 1" ;;
    *) [ "${held[$i]}" = 0 ] || fail "the basket holds ${held[$i]} of ${symbols[$i]} for the position" ;;
  esac
  low=0
  if [ "${held[$i]}" != 0 ]; then
    low=$(cast call --rpc-url "$rpc" "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "${symbols[$i]}" | sed -n 5p | cut -d' ' -f1)
  fi
  quantities+=("${held[$i]}") prices+=("$low")
done
expected=$(read_ "$margin" "currentRequirement(int256[],uint256[])(uint256,uint8,uint8)" \
  "[$(IFS=,; echo "${quantities[*]}")]" "[$(IFS=,; echo "${prices[*]}")]" | sed -n 1p)
read -r equity requirement _ _ <<<"$(health)"
[ "$requirement" = "$expected" ] || fail "the basket's requirement $requirement is not the program's $expected for its tokens"
nvda_state=$(read_ "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$(cast format-bytes32-string NVDA)" | sed -n 1p)
if [ "$nvda_state" = 0 ]; then
  refused=$(cast call --rpc-url "$rpc" --from "$me" "$accounts" "borrow(bytes32,uint256,address,address)" $cross 1 "$me" "$me" 2>&1) &&
    fail "a loan went through against a basket with a halted token"
  grep -q "$(cast calldata "AssetHalted(bytes32)" "$(cast format-bytes32-string NVDA)")" <<<"$refused" ||
    fail "the loan reverted for another reason: $refused"
  loan="refused while NVDA's band is halted"
else
  loan="open while NVDA's band is in state $nvda_state"
fi

unwrap_gas=$(l2_gas "$accounts" "$(cast calldata "unwrap(address,address,uint256)" "$me" "$basket" $two)")
send "$accounts" "unwrap(address,address,uint256)" "$me" "$basket" $two
[ "$(read_ "$accounts" "collateral(address,bytes32,address)(uint256)" "$me" $cross "$basket")" = 0 ] || fail "shares are left"
[ "$(read_ "$accounts" "collateral(address,bytes32,address)(uint256)" "$me" $cross "$nvda")" = $two ] || fail "the position does not hold 2 NVDA"
[ "$(read_ "$accounts" "collateral(address,bytes32,address)(uint256)" "$me" $cross "$spy")" = $one ] || fail "the position does not hold 1 SPY"
read -r equity_after requirement_after _ _ <<<"$(health)"
[ "$equity_after $requirement_after" = "$equity $requirement" ] ||
  fail "unwrapping moved the position from $equity $requirement to $equity_after $requirement_after"
send "$accounts" "withdraw(bytes32,address,uint256,address,address)" $cross "$nvda" $two "$me" "$me"
send "$accounts" "borrow(bytes32,uint256,address,address)" $cross 100000000 "$me" "$me"
send "$usdg" "mint(address,uint256)" "$me" 1000000
send "$usdg" "approve(address,uint256)" "$accounts" 101000000
send "$accounts" "repay(bytes32,uint256,address)" $cross 101000000 "$me"
[ "$(read_ "$accounts" "debt(address,bytes32)(uint256)" "$me" $cross)" = 0 ] || fail "the loan against SPY was not repaid"
send "$accounts" "withdraw(bytes32,address,uint256,address,address)" $cross "$spy" $one "$me" "$me"

send "$basket" "mint(uint256,address,uint256[])" $one "$me" "[$one,$half]"
[ "$(list "$basket" "previewRedeem(uint256)(uint256[])" $one)" = "[$one, $half]" ] || fail "a share does not redeem for what it took"
send "$basket" "redeem(uint256,address,address)" $one "$me" "$me"
[ "$(read_ "$basket" "totalSupply()(uint256)")" = 0 ] || fail "shares are left in the basket"
echo "basket PAIR $basket: 2 shares margined as 2 NVDA and 1 SPY, requirement $requirement, the program's own; a loan $loan; unwrapped at the same equity and requirement; minted in $mint_gas, deposited in $deposit_gas and unwrapped in $unwrap_gas L2 gas"
echo "PASS"
