#!/usr/bin/env bash
# Usage: deploy-band-feeds.sh RPC_URL DEPLOYMENTS_JSON SIGNER...
# Deploys one BandFeed, answering the low side, for each launch asset the registry's band configures,
# with the asset's Uniswap v3 USDG pool where the registry names one. The description says what is
# priced: the Robinhood Stock Token, or the share where the band has no token. Prints each asset and its
# feed's address. SIGNER is forge's wallet flags, such as `--account NAME --password-file FILE`. BAND_ASSETS
# overrides the launch assets.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
zero_address=0x0000000000000000000000000000000000000000
zero_id=0x0000000000000000000000000000000000000000000000000000000000000000
band=$(jq -er .tapehouse.Band "$registry") || { echo "$registry has no .tapehouse.Band" >&2; exit 1; }

for asset in ${BAND_ASSETS:-NVDA TSLA AAPL MSFT GOOGL SPY}; do
  symbol=$(cast format-bytes32-string "$asset")
  info=$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "$symbol")
  read -r feed feed_id index_id token <<<"$(tr '\n' ' ' <<<"$info")"
  [ "$feed" != $zero_address ] || [ "$feed_id" != $zero_id ] || [ "$index_id" != $zero_id ] || continue
  pool=$zero_address priced="$asset share"
  if [ "$token" != $zero_address ]; then
    priced="Robinhood $asset Stock Token"
    pool=$(jq -r --arg prefix "${asset}_USDG_" --arg zero $zero_address \
      '[.uniswapV3 // {} | to_entries[] | select(.key | startswith($prefix)) | .value][0] // $zero' "$registry")
  fi
  address=$(forge create --root "$root/contracts" src/BandFeed.sol:BandFeed --rpc-url "$rpc" --broadcast --json "$@" \
    --constructor-args "$band" "$symbol" 0 "$pool" "$asset / USD Tapehouse band, low side ($priced)" | jq -r .deployedTo)
  echo "$asset $address"
done
