#!/usr/bin/env bash
# Usage: margin-args.sh DEPLOYMENTS_JSON
# Prints the margin program's constructor arguments as nine shell words: symbols, volatility floors,
# volatilities, correlation floors, correlations, weekend-gap floors, weekend gaps, the market and the owner.
# The assets are those band-args.sh configures for the same file, in its order, and BAND_ASSETS narrows them
# as it does there. Each asset's parameters come from stylus/contracts/margin/parameters.json: an asset or a
# pair missing from it is an error. Correlations are listed for each pair (i, j), i < j, row by row. The
# market is the file's .market where the band prices it, and zero, the equal-weighted portfolio, where it
# does not. The owner is the file's .tapehouse.Owner.
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
  jq -er --arg kind "$3" --arg key "$1" --arg field "$2" '.[$kind][$key][$field]' "$parameters" ||
    { echo "${parameters#"$root"/} has no $3 for $1" >&2; exit 1; }
}
volatility_floors="" volatilities="" correlation_floors="" correlations="" gap_floors="" gaps=""
market=0x0000000000000000000000000000000000000000000000000000000000000000
for ((i = 0; i < ${#names[@]}; i++)); do
  volatility_floors+=",$(value "${names[i]}" floor volatility)"
  volatilities+=",$(value "${names[i]}" initial volatility)"
  gap_floors+=",$(value "${names[i]}" floor weekendGap)"
  gaps+=",$(value "${names[i]}" initial weekendGap)"
  [ "${names[i]}" != "$(jq -r .market "$parameters")" ] || market=$(cast format-bytes32-string "${names[i]}")
  for ((j = i + 1; j < ${#names[@]}; j++)); do
    pair="${names[i]}/${names[j]}"
    jq -e --arg pair "$pair" '.correlation | has($pair)' "$parameters" > /dev/null || pair="${names[j]}/${names[i]}"
    correlation_floors+=",$(value "$pair" floor correlation)"
    correlations+=",$(value "$pair" initial correlation)"
  done
done
echo "${band[0]}" "[${volatility_floors#,}]" "[${volatilities#,}]" "[${correlation_floors#,}]" "[${correlations#,}]" \
  "[${gap_floors#,}]" "[${gaps#,}]" "$market" "$owner"
