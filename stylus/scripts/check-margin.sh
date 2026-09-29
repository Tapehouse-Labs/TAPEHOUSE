#!/usr/bin/env bash
# Usage: check-margin.sh RPC_URL MARGIN_ADDRESS DEPLOYMENTS_JSON
# Checks that a deployed margin program holds the assets, floors, market and owner margin-args.sh builds from
# DEPLOYMENTS_JSON, and that the file's band, where it names one, prices every one of those assets. Until
# the first update the parameters themselves must be the file's initial values; after it they must sit at
# or above their floors, since the owner moves them.
set -euo pipefail

rpc=$1 margin=$2 registry=$3
args=$("$(dirname "$0")/margin-args.sh" "$registry")
read -r symbols volatility_floors volatilities correlation_floors correlations gap_floors gaps market owner <<<"$args"
held=$(cast call --rpc-url "$rpc" "$margin" "owner()(address)")
[ "$held" = "$(cast to-check-sum-address "$owner")" ] || { echo "FAIL: margin's owner is $held, $registry gives $owner"; exit 1; }
held=$(cast call --rpc-url "$rpc" "$margin" "assets()(bytes32[])" | tr -d " ")
[ "$held" = "$symbols" ] || { echo "FAIL: margin holds the assets $held, $registry gives $symbols"; exit 1; }
held=$(cast call --rpc-url "$rpc" "$margin" "market()(bytes32)")
[ "$held" = "$market" ] || { echo "FAIL: margin's market is $held, $registry gives $market"; exit 1; }
updated=$(cast call --rpc-url "$rpc" "$margin" "lastUpdate()(uint64)")
band=$(jq -r '.tapehouse.Band // empty' "$registry")
zero_address=0x0000000000000000000000000000000000000000
zero_id=0x0000000000000000000000000000000000000000000000000000000000000000
IFS=, read -r -a symbols <<<"${symbols//[\[\]]/}"
IFS=, read -r -a volatility_floors <<<"${volatility_floors//[\[\]]/}"
IFS=, read -r -a volatilities <<<"${volatilities//[\[\]]/}"
IFS=, read -r -a correlation_floors <<<"${correlation_floors//[\[\]]/}"
IFS=, read -r -a correlations <<<"${correlations//[\[\]]/}"
IFS=, read -r -a gap_floors <<<"${gap_floors//[\[\]]/}"
IFS=, read -r -a gaps <<<"${gaps//[\[\]]/}"

holds() {
  local value=$1 floor=$2 expected_value=$3 expected_floor=$4
  [ "$floor" = "$expected_floor" ] || return 1
  if [ "$updated" = 0 ]; then [ "$value" = "$expected_value" ]; else [ "$value" -ge "$floor" ]; fi
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
  for ((j = i + 1; j < ${#symbols[@]}; j++)); do
    other=$(cast parse-bytes32-string "${symbols[$j]}")
    read -r value floor <<<"$(cast call --rpc-url "$rpc" "$margin" "correlation(bytes32,bytes32)(uint16,uint16)" "${symbols[$i]}" "${symbols[$j]}" | sed 's/ \[[^]]*\]//' | tr '\n' ' ')"
    holds "$value" "$floor" "${correlations[$k]}" "${correlation_floors[$k]}" ||
      { echo "FAIL: margin holds correlation $value, floor $floor for $asset/$other; $registry gives ${correlations[$k]}, floor ${correlation_floors[$k]}"; exit 1; }
    k=$((k + 1))
  done
done
echo "margin $margin matches $registry for ${#symbols[@]} assets."
