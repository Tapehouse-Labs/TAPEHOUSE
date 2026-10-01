#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-liquidator-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks the liquidator the registry names on the dev node against the band and margin programs: it is the accounts'
# liquidator; a position of WETH borrowed to the accounts' 80% is not short at the liquidator's 84%, and is after WETH
# falls 10% on the dev node's ETH/USD stub; its auction starts at WETH's price; a buyer takes a quarter of its WETH,
# within the half of the debt a purchase may repay, paying USDG that repays the position, less the fee the reserve
# keeps. SPY held alone is judged against the margin program's requirement as the session stands: the current one, or
# the buffered open one while the session is not known. Prints the L2 gas of the start and the purchase.
set -euo pipefail

rpc=$1 key=$2 registry=$3
accounts=$(jq -r .tapehouse.MarginAccounts "$registry")
liquidator=$(jq -r .tapehouse.Liquidator "$registry")
usdg=$(jq -r .tokens.USDG "$registry")
weth=$(jq -r .tokens.WETH "$registry")
feed=$(jq -r .chainlink.ETH_USD "$registry")
me=$(cast wallet address --private-key "$key")
position=$(cast format-bytes32-string NVDA)
node_interface=0x00000000000000000000000000000000000000C8
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$key" "$@" > /dev/null; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }
short() { cast call --rpc-url "$rpc" "$liquidator" "shortfall(address,bytes32)(int256,uint256,bool,bool)" "$me" "$position" | sed -n 3p; }
l2_gas() {
  cast call --rpc-url "$rpc" $node_interface "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" \
    "$liquidator" false "$1" --from "$me" | sed -n 1,2p | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}

[ "$(read_ "$accounts" "liquidator()(address)")" = "$liquidator" ] || fail "the accounts' liquidator is not the registry's"
[ "$(read_ "$liquidator" "accounts()(address)")" = "$accounts" ] || fail "the liquidator's accounts are not the registry's"
price=$(cast call --rpc-url "$rpc" "$feed" "latestRoundData()(uint80,int256,uint256,uint256,uint80)" | sed -n 2p | cut -d' ' -f1)
now=$(cast block --rpc-url "$rpc" latest --field timestamp)
send "$feed" "setRound(int256,uint256)" "$price" "$now"
send "$weth" "mint(address,uint256)" "$me" 1000000000000000000
send "$weth" "approve(address,uint256)" "$accounts" 1000000000000000000
send "$accounts" "deposit(bytes32,address,uint256,address)" "$position" "$weth" 1000000000000000000 "$me"
send "$accounts" "borrow(bytes32,uint256,address,address)" "$position" $((price * 8 / 1000 - 1)) "$me" "$me"
[ "$(short)" = false ] || fail "a WETH position borrowed to 80% is short at 84%"
refused=$(cast call --rpc-url "$rpc" --from "$me" "$liquidator" "start(address,bytes32)" "$me" "$position" 2>&1) &&
  fail "a healthy position's auction started"
grep -q "$(cast sig "NotLiquidatable(address,bytes32)")" <<<"$refused" || fail "a healthy position's start reverted otherwise: $refused"
fallen=$((price * 9 / 10))
send "$feed" "setRound(int256,uint256)" "$fallen" "$(cast block --rpc-url "$rpc" latest --field timestamp)"
[ "$(short)" = true ] || fail "a WETH position is not short after WETH fell 10%"
start_gas=$(l2_gas "$(cast calldata "start(address,bytes32)" "$me" "$position")")
send "$liquidator" "start(address,bytes32)" "$me" "$position"
[ "$(read_ "$liquidator" "price(address,bytes32,address)(uint256)" "$me" "$position" "$weth")" = "$fallen" ] ||
  fail "the auction does not start at WETH's price"
debt=$(read_ "$accounts" "debt(address,bytes32)(uint256)" "$me" "$position")
reserve=$(read_ "$accounts" "reserve()(uint128)")
send "$usdg" "mint(address,uint256)" "$me" 10000000000
send "$usdg" "approve(address,uint256)" "$liquidator" 10000000000
calldata=$(cast calldata "buy(address,bytes32,address,uint256,uint256,address)" "$me" "$position" "$weth" 250000000000000000 10000000000 "$me")
buy_gas=$(l2_gas "$calldata")
topic=$(cast keccak "Bought(address,bytes32,address,uint256,uint256,uint256,address,address)")
data=$(cast send --rpc-url "$rpc" --private-key "$key" --json "$liquidator" "$calldata" |
  jq -r --arg topic "$topic" '.logs[] | select(.topics[0] == $topic) | .data')
cost=$(cast to-dec "0x${data:66:64}")
fee=$(cast to-dec "0x${data:130:64}")
kept=$(( $(read_ "$accounts" "reserve()(uint128)") - reserve ))
repaid=$(( debt - $(read_ "$accounts" "debt(address,bytes32)(uint256)" "$me" "$position") ))
[ "$(read_ "$accounts" "collateral(address,bytes32,address)(uint256)" "$me" "$position" "$weth")" = 750000000000000000 ] ||
  fail "the buyer did not take a quarter of the WETH"
[ "$fee" = $((cost * 50 / 10000)) ] && [ "$kept" -ge "$fee" ] && [ "$repaid" -gt 0 ] ||
  fail "the purchase cost $cost with a fee of $fee; the reserve kept $kept and the position repaid $repaid"
[ "$(read_ "$usdg" "balanceOf(address)(uint256)" "$liquidator")" = 0 ] || fail "the liquidator kept USDG"
send "$feed" "setRound(int256,uint256)" "$price" "$(cast block --rpc-url "$rpc" latest --field timestamp)"
spy=$(jq -r .tokens.SPY "$registry")
margin=$(jq -r .tapehouse.Margin "$registry")
alone=$(cast format-bytes32-string SPY)
send "$spy" "mint(address,uint256)" "$me" 10000000000000000000
send "$spy" "approve(address,uint256)" "$accounts" 10000000000000000000
send "$accounts" "deposit(bytes32,address,uint256,address)" "$alone" "$spy" 10000000000000000000 "$me"
low=$(cast call --rpc-url "$rpc" "$(jq -r .tapehouse.Band "$registry")" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$alone" | sed -n 5p | cut -d' ' -f1)
read -r quantities prices <<<"$(python3 -c "
import sys
assets = sys.argv[1].strip('[]').split(', ')
spy = [i for i, a in enumerate(assets) if bytes.fromhex(a[2:]).rstrip(b'\0') == b'SPY'][0]
q = ['10000000000000000000' if i == spy else '0' for i in range(len(assets))]
p = [sys.argv[2] if i == spy else '0' for i in range(len(assets))]
print('[' + ','.join(q) + '] [' + ','.join(p) + ']')" "$(cast call --rpc-url "$rpc" "$margin" "assets()(bytes32[])")" "$low")"
read -r current _ regime <<<"$(cast call --rpc-url "$rpc" "$margin" "currentRequirement(int256[],uint256[])(uint256,uint8,uint8)" "$quantities" "$prices" | cut -d' ' -f1 | tr '\n' ' ')"
open=$(cast call --rpc-url "$rpc" "$margin" "requirement(int256[],uint256[],uint64,bool)(uint256,uint8)" "$quantities" "$prices" 172800 false | sed -n 1p | cut -d' ' -f1)
judged=$(cast call --rpc-url "$rpc" "$liquidator" "shortfall(address,bytes32)(int256,uint256,bool,bool)" "$me" "$alone" | sed -n 2p | cut -d' ' -f1)
python3 -c "import sys; sys.exit($judged != ((($open * 5 - 1) // 4 + 1) if $regime == 0 else $current))" ||
  fail "SPY alone is judged against $judged in regime $regime; the program's current requirement is $current, its open one $open"
send "$accounts" "withdraw(bytes32,address,uint256,address,address)" "$alone" "$spy" 10000000000000000000 "$me" "$me"
echo "liquidator $liquidator: WETH borrowed to 80% is not short at 84%, is after a 10% fall, auctioned from $fallen; a quarter of its WETH cost $cost, $fee of it to the reserve; SPY alone judged against $judged in regime $regime; start $start_gas L2 gas, buy $buy_gas"
echo "PASS"
