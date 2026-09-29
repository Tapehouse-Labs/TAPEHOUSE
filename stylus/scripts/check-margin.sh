#!/usr/bin/env bash
# Usage: check-margin.sh RPC_URL MARGIN_ADDRESS DEPLOYMENTS_JSON
# Checks that a deployed margin program holds the assets, floors, depth ceilings, market, pools, ETH/USD feed,
# band and owner margin-args.sh builds from DEPLOYMENTS_JSON, that each pool trades the file's Stock Token for
# its asset against its USDG or WETH, and that the file's band, where it names one, prices every one of those
# assets. The depth ceilings are the file's initial depths. Until the first update the parameters themselves
# must be the file's initial values; after it they must sit at or above their floors and at or below their
# ceilings, since the owner moves them.
set -euo pipefail

rpc=$1 margin=$2 registry=$3
args=$("$(dirname "$0")/margin-args.sh" "$registry")
read -r assets correlation_floors correlations market usdg weth eth_usd band owner <<<"$args"
symbols=() volatility_floors=() volatilities=() gap_floors=() gaps=() depths=() pools=() tokens=()
entries=${assets#[(}
entries=${entries%)]}
while IFS=, read -r symbol volatility_floor volatility gap_floor gap selling buying pool token; do
  symbols+=("$symbol") volatility_floors+=("$volatility_floor") volatilities+=("$volatility")
  gap_floors+=("$gap_floor") gaps+=("$gap") depths+=("$selling" "$buying") pools+=("$pool") tokens+=("$token")
done <<<"${entries//),(/$'\n'}"
held=$(cast call --rpc-url "$rpc" "$margin" "owner()(address)")
[ "$held" = "$(cast to-check-sum-address "$owner")" ] || { echo "FAIL: margin's owner is $held, $registry gives $owner"; exit 1; }
held=$(cast call --rpc-url "$rpc" "$margin" "assets()(bytes32[])" | tr -d " ")
expected="[$(IFS=,; echo "${symbols[*]}")]"
[ "$held" = "$expected" ] || { echo "FAIL: margin holds the assets $held, $registry gives $expected"; exit 1; }
held=$(cast call --rpc-url "$rpc" "$margin" "market()(bytes32)")
[ "$held" = "$market" ] || { echo "FAIL: margin's market is $held, $registry gives $market"; exit 1; }
held=$(cast call --rpc-url "$rpc" "$margin" "ethUsdFeed()(address)")
[ "$held" = "$(cast to-check-sum-address "$eth_usd")" ] || { echo "FAIL: margin's ETH/USD feed is $held, $registry gives $eth_usd"; exit 1; }
updated=$(cast call --rpc-url "$rpc" "$margin" "lastUpdate()(uint64)")
band=${band#0x0000000000000000000000000000000000000000}
zero_address=0x0000000000000000000000000000000000000000
zero_id=0x0000000000000000000000000000000000000000000000000000000000000000
IFS=, read -r -a correlation_floors <<<"${correlation_floors//[\[\]]/}"
IFS=, read -r -a correlations <<<"${correlations//[\[\]]/}"

holds() {
  local value=$1 floor=$2 expected_value=$3 expected_floor=$4
  [ "$floor" = "$expected_floor" ] || return 1
  if [ "$updated" = 0 ]; then [ "$value" = "$expected_value" ]; else [ "$value" -ge "$floor" ]; fi
}
under() {
  local value=$1 ceiling=$2 expected_value=$3 expected_ceiling=$4
  [ "$ceiling" = "$expected_ceiling" ] || return 1
  if [ "$updated" = 0 ]; then [ "$value" = "$expected_value" ]; else [ "$value" -le "$ceiling" ]; fi
}
k=0
for i in "${!symbols[@]}"; do
  asset=$(cast parse-bytes32-string "${symbols[$i]}")
  if [ -n "$band" ]; then
    priced=$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "${symbols[$i]}") ||
      { echo "FAIL: band $band does not answer asset() for $asset"; exit 1; }
    read -r feed feed_id index_id _ <<<"$(tr '\n' ' ' <<<"$priced")"
    [ "$feed" != "$zero_address" ] || [ "$feed_id" != "$zero_id" ] || [ "$index_id" != "$zero_id" ] ||
      { echo "FAIL: band $band does not price $asset"; exit 1; }
  fi
  read -r value floor <<<"$(cast call --rpc-url "$rpc" "$margin" "volatility(bytes32)(uint32,uint32)" "${symbols[$i]}" | sed 's/ \[[^]]*\]//' | tr '\n' ' ')"
  holds "$value" "$floor" "${volatilities[$i]}" "${volatility_floors[$i]}" ||
    { echo "FAIL: margin holds volatility $value, floor $floor for $asset; $registry gives ${volatilities[$i]}, floor ${volatility_floors[$i]}"; exit 1; }
  read -r value floor <<<"$(cast call --rpc-url "$rpc" "$margin" "weekendGap(bytes32)(uint32,uint32)" "${symbols[$i]}" | sed 's/ \[[^]]*\]//' | tr '\n' ' ')"
  holds "$value" "$floor" "${gaps[$i]}" "${gap_floors[$i]}" ||
    { echo "FAIL: margin holds weekend gap $value, floor $floor for $asset; $registry gives ${gaps[$i]}, floor ${gap_floors[$i]}"; exit 1; }
  read -r selling buying selling_ceiling buying_ceiling <<<"$(cast call --rpc-url "$rpc" "$margin" "depth(bytes32)(uint32,uint32,uint32,uint32)" "${symbols[$i]}" | sed 's/ \[[^]]*\]//' | tr '\n' ' ')"
  under "$selling" "$selling_ceiling" "${depths[$((2 * i))]}" "${depths[$((2 * i))]}" &&
    under "$buying" "$buying_ceiling" "${depths[$((2 * i + 1))]}" "${depths[$((2 * i + 1))]}" ||
    { echo "FAIL: margin holds depths $selling, $buying, ceilings $selling_ceiling, $buying_ceiling for $asset; $registry gives ${depths[$((2 * i))]}, ${depths[$((2 * i + 1))]}, ceilings the same"; exit 1; }
  held=$(cast call --rpc-url "$rpc" "$margin" "pool(bytes32)(address)" "${symbols[$i]}")
  [ "$held" = "$(cast to-check-sum-address "${pools[$i]}")" ] ||
    { echo "FAIL: margin's pool for $asset is $held, $registry gives ${pools[$i]}"; exit 1; }
  if [ "$held" != 0x0000000000000000000000000000000000000000 ]; then
    pair=$(for side in token0 token1; do cast call --rpc-url "$rpc" "$held" "$side()(address)"; done | paste -sd' ' -)
    stock=$(cast to-check-sum-address "${tokens[$i]}")
    case $pair in
      "$stock $(cast to-check-sum-address "$usdg")" | "$(cast to-check-sum-address "$usdg") $stock" | \
        "$stock $(cast to-check-sum-address "$weth")" | "$(cast to-check-sum-address "$weth") $stock") ;;
      *) echo "FAIL: margin's pool for $asset trades $pair, not $registry's Stock Token $stock against USDG or WETH"; exit 1 ;;
    esac
  fi
  for ((j = i + 1; j < ${#symbols[@]}; j++)); do
    other=$(cast parse-bytes32-string "${symbols[$j]}")
    read -r value floor <<<"$(cast call --rpc-url "$rpc" "$margin" "correlation(bytes32,bytes32)(uint16,uint16)" "${symbols[$i]}" "${symbols[$j]}" | sed 's/ \[[^]]*\]//' | tr '\n' ' ')"
    holds "$value" "$floor" "${correlations[$k]}" "${correlation_floors[$k]}" ||
      { echo "FAIL: margin holds correlation $value, floor $floor for $asset/$other; $registry gives ${correlations[$k]}, floor ${correlation_floors[$k]}"; exit 1; }
    k=$((k + 1))
  done
done
held=$(cast call --rpc-url "$rpc" "$margin" "band()(address)")
[ "$held" = "$(cast to-check-sum-address "${band:-0x0000000000000000000000000000000000000000}")" ] ||
  { echo "FAIL: margin's band is $held, $registry gives ${band:-none}"; exit 1; }
echo "margin $margin matches $registry for ${#symbols[@]} assets."
