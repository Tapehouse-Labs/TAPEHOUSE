#!/usr/bin/env bash
# Usage: margin-args.sh DEPLOYMENTS_JSON
# Prints the margin program's constructor arguments as nine shell words: the assets, correlation floors,
# correlations, the market, USDG, WETH, the ETH/USD feed, the band and the owner. Each asset is a tuple of its
# symbol, volatility floor, volatility, weekend-gap floor, weekend gap, selling depth, buying depth, pool and
# Stock Token.
# The assets are those band-args.sh configures for the same file, in its order, and BAND_ASSETS narrows them
# as it does there. Each asset's parameters come from stylus/contracts/margin/parameters.json: an asset or a
# pair missing from it, or a value that is not a positive integer, is an error. Correlations are listed for
# each pair (i, j), i < j, row by row. The market is the file's .market where the band prices it, and zero,
# the equal-weighted portfolio, where it does not. Each asset's pool is the file's .uniswapV3 entry that
# parameters.json's .pool names, and zero where the file has none; an asset missing from .pool is an error.
# Each asset's Stock Token, and USDG, WETH and the ETH/USD feed, are the file's .tokens.<asset>, .tokens.USDG,
# .tokens.WETH and .chainlink.ETH_USD, zero where absent. The band is the file's .tapehouse.Band, zero where
# absent, and the owner its .tapehouse.Owner.
set -euo pipefail

registry=$1
root=$(cd "$(dirname "$0")/../.." && pwd)
parameters=$root/stylus/contracts/margin/parameters.json
words=$("$root/stylus/scripts/band-args.sh" "$registry")
read -r -a band <<<"$words"
owner=${band[6]}
names=()
for symbol in $(tr -d '[]' <<<"${band[0]}" | tr ',' ' '); do
  names+=("$(cast parse-bytes32-string "$symbol")")
done

value() {
  jq -er --arg kind "$3" --arg key "$1" --arg field "$2" \
    '.[$kind][$key][$field] | select(type == "number" and . == floor and . > 0)' "$parameters" ||
    { echo "${parameters#"$root"/} has no positive integer $3 $2 for $1" >&2; exit 1; }
}
market_symbol=$(jq -er '.market | strings' "$parameters") || { echo "${parameters#"$root"/} names no market" >&2; exit 1; }
assets="" correlation_floors="" correlations=""
zero_address=0x0000000000000000000000000000000000000000
market=0x0000000000000000000000000000000000000000000000000000000000000000
for ((i = 0; i < ${#names[@]}; i++)); do
  volatility_floor=$(value "${names[i]}" floor volatility)
  volatility=$(value "${names[i]}" initial volatility)
  gap_floor=$(value "${names[i]}" floor weekendGap)
  gap=$(value "${names[i]}" initial weekendGap)
  depths=$(jq -er --arg key "${names[i]}" '.depth[$key].initial
    | select(type == "array" and length == 2 and all(type == "number" and . == floor and . > 0)) | join(",")' \
    "$parameters") || { echo "${parameters#"$root"/} has no two positive integer depths for ${names[i]}" >&2; exit 1; }
  pool=$(jq -er --arg key "${names[i]}" '.pool[$key] | strings' "$parameters") ||
    { echo "${parameters#"$root"/} has no pool for ${names[i]}" >&2; exit 1; }
  pool=$(jq -r --arg name "$pool" --arg zero $zero_address '.uniswapV3[$name] // $zero' "$registry")
  token=$(jq -r --arg name "${names[i]}" --arg zero $zero_address '.tokens[$name] // $zero' "$registry")
  symbol=$(cast format-bytes32-string "${names[i]}")
  assets+=",($symbol,$volatility_floor,$volatility,$gap_floor,$gap,$depths,$pool,$token)"
  [ "${names[i]}" != "$market_symbol" ] || market=$(cast format-bytes32-string "${names[i]}")
  for ((j = i + 1; j < ${#names[@]}; j++)); do
    pair="${names[i]}/${names[j]}"
    jq -e --arg pair "$pair" '.correlation | has($pair)' "$parameters" > /dev/null || pair="${names[j]}/${names[i]}"
    correlation_floors+=",$(value "$pair" floor correlation)"
    correlations+=",$(value "$pair" initial correlation)"
  done
done
optional() { jq -r --arg zero $zero_address "$1 // \$zero" "$registry"; }
echo "[${assets#,}]" "[${correlation_floors#,}]" "[${correlations#,}]" "$market" "$(optional .tokens.USDG)" "$(optional .tokens.WETH)" "$(optional .chainlink.ETH_USD)" "$(optional .tapehouse.Band)" "$owner"
