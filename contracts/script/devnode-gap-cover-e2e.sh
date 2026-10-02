#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-gap-cover-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks the gap cover the registry names on the dev node: it is over the margin program and its band, and prices SPY
# cover from the margin program's weekend gap, moved during the sales to follow the week's realised move on the feed,
# and the band's centre, its smallest deductible and its premium matching the fitted tail to the unit. The dev node's
# stub feeds keep a round a day over the nine days before their deployment, so during the sales the week is read from
# them: the gap is twice its move, at least 55% of the margin program's and at most eight times it. A writer deposits
# 10,000 USDG. Then, as the cover's sales stand: while they are open, a keeper measures SPY's week, a buyer pays the
# premium of 10,000 USDG of SPY cover from its smallest deductible to 10% beyond, and the cover reserves the layer's
# whole payout; otherwise no cover is sold and the writer takes its USDG back. Prints the L2 gas of the measurement and
# the purchase, or of recording the session.
set -euo pipefail

rpc=$1 key=$2 registry=$3
cover=$(jq -r .tapehouse.GapCover "$registry")
margin=$(jq -r .tapehouse.Margin "$registry")
band=$(jq -r .tapehouse.Band "$registry")
usdg=$(jq -r .tokens.USDG "$registry")
me=$(cast wallet address --private-key "$key")
symbol=$(cast format-bytes32-string SPY)
node_interface=0x00000000000000000000000000000000000000C8
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$key" "$@" > /dev/null; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }
l2_gas() {
  cast call --rpc-url "$rpc" $node_interface "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" \
    "$cover" false "$1" --from "$me" | sed -n 1,2p | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}

[ "$(read_ "$cover" "engine()(address)")" = "$margin" ] || fail "the cover is not over the registry's margin program"
[ "$(read_ "$cover" "band()(address)")" = "$band" ] || fail "the cover is not over the registry's band"
feed=$(jq -r .chainlink.SPY_USD "$registry")
[ "$(read_ "$cover" "feed(bytes32)(address)" "$symbol")" = "$feed" ] ||
  fail "the cover does not settle SPY on the registry's SPY/USD feed"
engine_gap=$(read_ "$margin" "weekendGap(bytes32)(uint32,uint32)" "$symbol" | head -n 1)
read -r closes ends <<<"$(cast call --rpc-url "$rpc" "$cover" "sales()(uint64,uint64)" | cut -d' ' -f1 | tr '\n' ' ')"
read -r gap move <<<"$(cast call --rpc-url "$rpc" "$cover" "pricingGap(bytes32)(uint256,uint256)" "$symbol" |
  cut -d' ' -f1 | tr '\n' ' ')"
if [ "$ends" != 0 ]; then
  bounded=$((move * 2 > engine_gap * 5500 / 10000 ? move * 2 : engine_gap * 5500 / 10000))
  [ "$move" != 18446744073709551615 ] && [ "$gap" = $((bounded < engine_gap * 8 ? bounded : engine_gap * 8)) ] ||
    fail "SPY's gap $gap and week $move during the sales, over a feed with a round a day"
else
  [ "$move" = 0 ] && [ "$gap" = "$engine_gap" ] || fail "SPY's gap $gap and week $move outside the sales"
fi
last=$(cast call --rpc-url "$rpc" "$feed" "latestRoundData()(uint80,int256,uint256,uint256,uint80)" | sed -n 2p | cut -d' ' -f1)
mid=$(cast call --rpc-url "$rpc" "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$symbol" |
  sed -n 3p | cut -d' ' -f1)
stale=$(python3 -c "print(-(-($last - $mid) * 10**12 // $last) if $mid < $last else 0)")
deductible=$(read_ "$cover" "minDeductible(bytes32)(uint256)" "$symbol")
[ "$deductible" = $(((gap * 183288 + stale + 99999999) / 100000000)) ] ||
  fail "SPY's smallest deductible $deductible at gap $gap and a reference $stale picounits above the band's centre"
notional=10000000000
limit=$((deductible + 1000))
premium=$(read_ "$cover" "quote(bytes32,uint256,uint256,uint256)(uint256)" "$symbol" "$notional" "$deductible" "$limit")
expected=$(python3 -c "
g, n, x, l = $gap, $notional, $deductible * 10**8 - $stale, $limit * 10**8 - $stale
b, u = g * 95837, g * 183288
print(-(-n * 4 * 50224 * b * b * (l - x) * 15000 // (10**22 * (2 * b + x - u) * (2 * b + l - u))))")
[ "$premium" = "$expected" ] || fail "the premium of SPY cover is $premium, not $expected"
deposit=10000000000
send "$usdg" "mint(address,uint256)" "$me" $((deposit + premium))
send "$usdg" "approve(address,uint256)" "$cover" $((deposit + premium))
send "$cover" "deposit(uint256,address)" "$deposit" "$me"
[ "$(read_ "$cover" "held()(uint256)")" = "$deposit" ] || fail "the writer's deposit is not the cover's count"
if [ "$ends" != 0 ]; then
  measure=$(cast calldata "measure(bytes32)" "$symbol")
  measured=$(l2_gas "$measure")
  send "$cover" "$measure"
  [ "$(cast call --rpc-url "$rpc" "$cover" "pricingGap(bytes32)(uint256,uint256)" "$symbol" | cut -d' ' -f1 | tr '\n' ' ')" = \
    "$gap $move " ] || fail "the measured week is not the one SPY was priced at"
  calldata=$(cast calldata "buy(bytes32,uint256,uint256,uint256,uint256,address)" "$symbol" "$notional" "$deductible" \
    "$limit" "$premium" "$me")
  gas=$(l2_gas "$calldata")
  send "$cover" "$calldata"
  [ "$(read_ "$cover" "reserved()(uint256)")" = $((notional * 1000 / 10000)) ] || fail "the cover did not reserve the layer"
  [ "$(read_ "$cover" "premiums()(uint256)")" = "$premium" ] || fail "the cover did not hold the premium"
  [ "$(read_ "$usdg" "balanceOf(address)(uint256)" "$cover")" = $((deposit + premium)) ] || fail "the cover's USDG"
  [ "$(read_ "$cover" "outstandingIn(uint64)(uint256)" "$closes")" = 1 ] || fail "the cover did not join the series of $closes"
  outcome="the week measured in $measured L2 gas; sold 10,000 USDG of SPY cover from $deductible to $limit bps for $premium, sales ending at $ends, in $gas L2 gas"
else
  if refused=$(cast call --rpc-url "$rpc" "$cover" "buy(bytes32,uint256,uint256,uint256,uint256,address)" "$symbol" \
    "$notional" "$deductible" "$limit" "$premium" "$me" --from "$me" 2>&1); then
    fail "a purchase outside the sales went through: $refused"
  fi
  grep -q 0x0671dd5e <<<"$refused" || fail "a purchase outside the sales did not revert with SalesClosed: $refused"
  gas=$(l2_gas "$(cast calldata "record()")")
  send "$cover" "record()"
  send "$cover" "redeem(uint256,address,address)" "$(read_ "$cover" "balanceOf(address)(uint256)" "$me")" "$me" "$me"
  [ "$(read_ "$cover" "held()(uint256)")" = 0 ] || fail "the writer did not take its USDG back"
  [ "$closes" = 0 ] || fail "the sales name the closure of $closes while closed"
  outcome="no cover sold outside the sales, recorded in $gas L2 gas, and the writer took its 10,000 USDG back"
fi
echo "gap cover $cover: SPY's weekend gap $engine_gap, priced at $gap on a week of $move, smallest deductible $deductible bps, premium" \
  "$premium on 10,000 USDG; $outcome"
echo "PASS"
