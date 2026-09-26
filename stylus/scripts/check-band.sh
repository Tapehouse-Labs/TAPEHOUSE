#!/usr/bin/env bash
# Usage: check-band.sh RPC_URL BAND_ADDRESS DEPLOYMENTS_JSON
# Checks that a deployed band holds the configuration band-args.sh builds from DEPLOYMENTS_JSON, that
# each configured Chainlink feed describes its own asset, and that every asset band-args.sh leaves out is
# unconfigured.
set -euo pipefail

rpc=$1 band=$2 registry=$3
args=$("$(dirname "$0")/band-args.sh" "$registry")
read -r symbols feeds feed_ids <<<"$args"
IFS=, read -r -a symbols <<<"${symbols//[\[\]]/}"
IFS=, read -r -a feeds <<<"${feeds//[\[\]]/}"
IFS=, read -r -a feed_ids <<<"${feed_ids//[\[\]]/}"

for i in "${!symbols[@]}"; do
  asset=$(cast parse-bytes32-string "${symbols[$i]}")
  read -r feed feed_id <<<"$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32)" "${symbols[$i]}" | tr '\n' ' ')"
  [ "$feed" = "$(cast to-check-sum-address "${feeds[$i]}")" ] && [ "$feed_id" = "${feed_ids[$i]}" ] ||
    { echo "FAIL: band holds $feed $feed_id for $asset, $registry gives ${feeds[$i]} ${feed_ids[$i]}"; exit 1; }
  if [ "$feed" != 0x0000000000000000000000000000000000000000 ]; then
    description=$(cast call --rpc-url "$rpc" "$feed" "description()(string)")
    [[ "$description" == *"$asset / USD"* ]] || { echo "FAIL: the $asset feed $feed describes itself as $description"; exit 1; }
  fi
done
for asset in ${BAND_ASSETS:-NVDA TSLA AAPL MSFT GOOGL SPY}; do
  symbol=$(cast format-bytes32-string "$asset")
  case " ${symbols[*]} " in *" $symbol "*) continue ;; esac
  held=$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32)" "$symbol" | tr '\n' ' ')
  [ "$held" = "0x0000000000000000000000000000000000000000 0x0000000000000000000000000000000000000000000000000000000000000000 " ] ||
    { echo "FAIL: band configures $asset, which $registry leaves out: $held"; exit 1; }
done
echo "band $band matches $registry for ${#symbols[@]} assets."
