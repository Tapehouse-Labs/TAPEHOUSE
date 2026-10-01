#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-gap-backstop-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks the gap backstop the registry names on the dev node: it is the accounts' backstop, over them, holding the
# team's seed; it claims the premium set aside. A position of WETH borrowed to the accounts' 80% becomes worth less than
# it owes when WETH falls 40% on the dev node's ETH/USD stub; the liquidator sells all its WETH, and the backstop repays
# what the emptied position still owes, so the vault writes nothing off. Prints the L2 gas of the claim and the cover.
set -euo pipefail

rpc=$1 key=$2 registry=$3
accounts=$(jq -r .tapehouse.MarginAccounts "$registry")
liquidator=$(jq -r .tapehouse.Liquidator "$registry")
backstop=$(jq -r .tapehouse.GapBackstop "$registry")
usdg=$(jq -r .tokens.USDG "$registry")
weth=$(jq -r .tokens.WETH "$registry")
feed=$(jq -r .chainlink.ETH_USD "$registry")
owner=$(jq -r .tapehouse.Owner "$registry")
me=$(cast wallet address --private-key "$key")
position=$(cast format-bytes32-string SPY)
node_interface=0x00000000000000000000000000000000000000C8
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$key" "$@" > /dev/null; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }
event() { jq -r --arg topic "$(cast keccak "$1")" '.logs[] | select(.topics[0] == $topic) | .data'; }
word() { cast to-dec "0x${1:$((2 + 64 * $2)):64}"; }
l2_gas() {
  cast call --rpc-url "$rpc" $node_interface "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" \
    "$backstop" false "$1" --from "$me" | head -n 2 | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}
now() { cast block --rpc-url "$rpc" latest --field timestamp; }

[ "$(read_ "$accounts" "backstop()(address)")" = "$backstop" ] || fail "the accounts' backstop is not the registry's"
[ "$(read_ "$backstop" "accounts()(address)")" = "$accounts" ] || fail "the backstop's accounts are not the registry's"
seed=$(read_ "$backstop" "convertToAssets(uint256)(uint256)" "$(read_ "$backstop" "balanceOf(address)(uint256)" "$owner")")
[ "$seed" -gt 0 ] || fail "the owner holds no seed in the backstop"
claim_gas=$(l2_gas "$(cast calldata "claim()")")
held=$(read_ "$backstop" "held()(uint256)")
claimed=$(cast send --rpc-url "$rpc" --private-key "$key" --json "$backstop" "claim()" | event "PremiumAdded(uint256)")
claimed=$(word "$claimed" 0)
[ "$(read_ "$backstop" "held()(uint256)")" = $((held + claimed)) ] || fail "the claim of $claimed did not reach the backstop"
price=$(cast call --rpc-url "$rpc" "$feed" "latestRoundData()(uint80,int256,uint256,uint256,uint80)" | sed -n 2p | cut -d' ' -f1)
send "$feed" "setRound(int256,uint256)" "$price" "$(now)"
send "$weth" "mint(address,uint256)" "$me" 1000000000000000000
send "$weth" "approve(address,uint256)" "$accounts" 1000000000000000000
send "$accounts" "deposit(bytes32,address,uint256,address)" "$position" "$weth" 1000000000000000000 "$me"
send "$accounts" "borrow(bytes32,uint256,address,address)" "$position" $((price * 8 / 1000 - 1)) "$me" "$me"
send "$feed" "setRound(int256,uint256)" $((price * 6 / 10)) "$(now)"
send "$liquidator" "start(address,bytes32)" "$me" "$position"
send "$usdg" "mint(address,uint256)" "$me" 10000000000
send "$usdg" "approve(address,uint256)" "$liquidator" 10000000000
for _ in 1 2 3 4; do
  [ "$(read_ "$accounts" "collateral(address,bytes32,address)(uint256)" "$me" "$position" "$weth")" = 0 ] && break
  send "$liquidator" "buy(address,bytes32,address,uint256,uint256,address)" "$me" "$position" "$weth" 1000000000000000000 10000000000 "$me"
done
[ "$(read_ "$accounts" "collateral(address,bytes32,address)(uint256)" "$me" "$position" "$weth")" = 0 ] ||
  fail "the liquidator did not sell all the WETH"
owed=$(read_ "$accounts" "debt(address,bytes32)(uint256)" "$me" "$position")
[ "$owed" -gt 0 ] || fail "the emptied position owes nothing"
held=$(read_ "$backstop" "held()(uint256)")
cover_gas=$(l2_gas "$(cast calldata "cover(address,bytes32)" "$me" "$position")")
receipt=$(cast send --rpc-url "$rpc" --private-key "$key" --json "$backstop" "cover(address,bytes32)" "$me" "$position")
covered=$(event "Covered(address,bytes32,uint64,uint256,uint256)" <<<"$receipt")
paid=$(word "$covered" 0)
written=$(word "$covered" 1)
[ "$written" = 0 ] && [ "$paid" -ge "$owed" ] && [ "$paid" -le $((owed + 1000)) ] ||
  fail "the backstop paid $paid of $owed, and $written was written off"
[ "$(read_ "$accounts" "debt(address,bytes32)(uint256)" "$me" "$position")" = 0 ] || fail "the position still owes"
[ "$(read_ "$backstop" "held()(uint256)")" = $((held - paid)) ] || fail "the backstop's count did not fall by $paid"
[ "$(read_ "$usdg" "balanceOf(address)(uint256)" "$backstop")" = $((held - paid)) ] || fail "the backstop's USDG is not its count"
send "$feed" "setRound(int256,uint256)" "$price" "$(now)"
echo "gap backstop $backstop: seeded with $seed, claimed $claimed of premium in $claim_gas L2 gas; a WETH position emptied after a 40% fall owed $owed, and the backstop paid $paid, $written to lenders, in $cover_gas L2 gas"
echo "PASS"
