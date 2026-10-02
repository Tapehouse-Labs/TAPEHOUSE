#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-short-positions-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks the short positions on the dev node: they are SPY's lending vault's borrower and short SPY through its USDG
# pool. A position lends 10 SPY; a short adds 1,000 USDG of margin, borrows 2 SPY and sells them through the stub
# router at the band's centre, and the margin program margins it as a negative quantity at the band's high edge; it buys
# them back and repays them, which closes it, and takes its USDG out. Prints the L2 gas of the sale and of the buy-back.
set -euo pipefail

rpc=$1 key=$2 registry=$3
accounts=$(jq -r .tapehouse.MarginAccounts "$registry")
shorts=$(jq -r .tapehouse.ShortPositions "$registry")
lending=$(jq -r .tapehouse.StockLending.SPY "$registry")
band=$(jq -r .tapehouse.Band "$registry")
spy=$(jq -r .tokens.SPY "$registry")
usdg=$(jq -r .tokens.USDG "$registry")
me=$(cast wallet address --private-key "$key")
symbol=$(cast format-bytes32-string SPY)
node_interface=0x00000000000000000000000000000000000000C8
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$key" "$@" > /dev/null; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }
calc() { python3 -c "print($1)"; }
l2_gas() {
  cast call --rpc-url "$rpc" $node_interface "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" \
    "$shorts" false "$1" --from "$me" | sed -n 1,2p | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}
position() {
  cast call --rpc-url "$rpc" "$shorts" "position(address,bytes32)(int256,uint256,uint256)" "$me" "$symbol" | cut -d' ' -f1 | tr '\n' ' '
}

[ "$(read_ "$shorts" "accounts()(address)")" = "$accounts" ] || fail "the shorts are not over the registry's accounts"
[ "$(read_ "$lending" "borrower()(address)")" = "$shorts" ] || fail "the shorts do not borrow SPY's lending vault"
[ "$(read_ "$shorts" "fee(bytes32)(uint24)" "$symbol")" = 500 ] || fail "the shorts do not short SPY through its 0.05% pool"
read -r state _ mid _ _ high <<<"$(cast call --rpc-url "$rpc" "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$symbol" | cut -d' ' -f1 | tr '\n' ' ')"
[ "$state" != 0 ] || fail "SPY's band is halted"
send "$(jq -r .uniswapV3.SwapRouter02 "$registry")" "setPrice(address,uint256)" "$spy" "$mid"
ten=10000000000000000000
two=2000000000000000000
margin=1000000000
send "$spy" "mint(address,uint256)" "$me" "$ten"
send "$spy" "approve(address,uint256)" "$accounts" "$ten"
send "$accounts" "deposit(bytes32,address,uint256,address)" "$symbol" "$spy" "$ten" "$me"
send "$accounts" "lend(bytes32,address,uint256,address)" "$symbol" "$spy" "$ten" "$me"
send "$usdg" "mint(address,uint256)" "$me" "$margin"
send "$usdg" "approve(address,uint256)" "$shorts" "$margin"
send "$shorts" "deposit(bytes32,uint256,address)" "$symbol" "$margin" "$me"
sell_gas=$(l2_gas "$(cast calldata "sell(bytes32,uint256,uint256,address)" "$symbol" "$two" 0 "$me")")
send "$shorts" "sell(bytes32,uint256,uint256,address)" "$symbol" "$two" 0 "$me"
[ "$(read_ "$lending" "debt()(uint256)")" = "$two" ] || fail "the vault did not lend the shorts 2 SPY"
read -r held debt _ <<<"$(position)"
[ "$debt" = "$two" ] || fail "the short owes $debt, not 2 SPY"
[ "$held" = "$(calc "$margin + $two * $mid // 10**20")" ] || fail "the short holds $held USDG, not its margin and 2 SPY at the band's centre"
read -r equity requirement missing regime <<<"$(cast call --rpc-url "$rpc" "$shorts" "health(address,bytes32)(int256,uint256,uint8,uint8)" "$me" "$symbol" | cut -d' ' -f1 | tr '\n' ' ')"
read -r quantities prices <<<"$(python3 -c "
import sys
assets = sys.argv[1].strip('[]').split(', ')
spy = [i for i, a in enumerate(assets) if bytes.fromhex(a[2:]).rstrip(b'\0') == b'SPY'][0]
print('[' + ','.join('-$two' if i == spy else '0' for i in range(len(assets))) + ']',
      '[' + ','.join('$high' if i == spy else '0' for i in range(len(assets))) + ']')
" "$(cast call --rpc-url "$rpc" "$accounts" "stocks()(bytes32[],address[])" | sed -n 1p)")"
expected=$(read_ "$(jq -r .tapehouse.Margin "$registry")" "currentRequirement(int256[],uint256[])(uint256,uint8,uint8)" "$quantities" "$prices" | sed -n 1p)
[ "$requirement" = "$expected" ] || fail "the short's requirement $requirement is not the margin program's $expected for -2 SPY at the high edge"
[ "$(calc "$requirement > 0")" = True ] || fail "the margin program required nothing of the short"
[ "$equity" = "$(calc "$held * 10**12 - $two * $high // 10**8")" ] || fail "the short's equity is not its USDG less 2 SPY at the high edge"
cover_gas=$(l2_gas "$(cast calldata "cover(bytes32,uint256,uint256,address)" "$symbol" "$(cast max-uint)" "$(cast max-uint)" "$me")")
send "$shorts" "cover(bytes32,uint256,uint256,address)" "$symbol" "$(cast max-uint)" "$(cast max-uint)" "$me"
[ "$(read_ "$lending" "debt()(uint256)")" = 0 ] || fail "the buy-back did not repay the vault"
read -r left _ shares <<<"$(position)"
[ "$shares" = 0 ] || fail "the short is still open"
send "$shorts" "withdraw(bytes32,uint256,address,address)" "$symbol" "$left" "$me" "$me"
[ "$(read_ "$usdg" "balanceOf(address)(uint256)" "$shorts")" = 0 ] || fail "the shorts kept USDG"
echo "short positions $shorts: sold 2 SPY at $mid, margined at $requirement for -2 SPY at $high (missing $missing, regime $regime) on $equity of equity, bought them back and took $left USDG out; sell $sell_gas L2 gas, cover $cover_gas"
echo "PASS"
