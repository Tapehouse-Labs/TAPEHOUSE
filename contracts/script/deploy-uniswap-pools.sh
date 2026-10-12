#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: POOL_USDG=<USDG units> [POOL_FEE=500|3000|10000] [BAND_ASSETS="TSLA AMZN"] deploy-uniswap-pools.sh RPC_URL \
#   DEPLOYMENTS_JSON SIGNER...
# For a chain without Uniswap v3, such as Robinhood Chain's testnet: deploys Uniswap v3's factory, a single-pool swap
# router and a quoter over it, records them under .uniswapV3.Factory, .SwapRouter02 and .QuoterV2, and for each asset the
# registry's band gives a Stock Token, creates and seeds its pool against USDG at the price of the asset's band feed, at
# POOL_FEE in hundredths of a basis point (3000, 0.3%, unless it says another), and records it under
# .uniswapV3.<ASSET>_USDG_<POOL_FEE>. A short sells only at the band's low edge or above, so a pool whose fee is as wide
# as the band, as 0.3% is, cannot take one: a chain whose band is narrow wants 500. A pool holds POOL_USDG USDG units from
# the signer, in the ticks 11% either side of its price, which is where the margin engine reads its depth, with the Stock
# Token those ticks take. The signer wraps the Stock Token it holds into the registry's wrapper if it holds less than
# that. The margin engine reads a pool with more than 1,800 observations and 30 minutes of history, so its observations
# are raised here and a pool serves the engine half an hour after it is created. A pool the registry names already is left
# as it is. SIGNER is cast's wallet flags, such as `--account NAME --password-file FILE`.
set -euo pipefail

rpc=$1 registry=$2
shift 2
root=$(cd "$(dirname "$0")/../.." && pwd)
project=$root/contracts/uniswap-testnet
: "${POOL_USDG:?POOL_USDG, the USDG units each pool holds}"
field() { jq -er "$1" "$registry" || { echo "$registry has no $1" >&2; exit 1; }; }
record() {
  jq --arg value "$2" "$1 = \$value" "$registry" > "$registry.tmp"
  mv "$registry.tmp" "$registry"
}
# A public RPC answers from nodes a block or two apart, and a transaction sent right after another can meet a nonce that
# is already used: it is sent again, with the nonce the chain has by then.
retry() {
  local output
  for _ in 1 2 3 4 5 6; do
    if output=$("$@" 2>&1); then
      printf '%s' "$output"
      return 0
    fi
    grep -q "nonce too low" <<<"$output" || { printf '%s\n' "$output" >&2; return 1; }
    sleep 3
  done
  printf '%s\n' "$output" >&2
  return 1
}
send() { retry cast send --rpc-url "$rpc" "$@" "${signer_flags[@]}" > /dev/null; }
signer_flags=("$@")
signer=$(cast wallet address "$@")
band=$(field .tapehouse.Band)
usdg=$(field .tokens.USDG)
zero=0x0000000000000000000000000000000000000000
fee=${POOL_FEE:-3000}
case $fee in 500) spacing=10 ;; 3000) spacing=60 ;; 10000) spacing=200 ;; *) echo "POOL_FEE must be 500, 3000 or 10000" >&2; exit 1 ;; esac

out=$(mktemp -d)
trap 'rm -rf "$out"' EXIT
forge build --root "$project" --skip test --out "$out" > /dev/null
artifact() { jq -er .bytecode.object "$out/$1.sol/$1.json"; }
create() {
  local code
  code=$(artifact "$1")${2:+$(cast abi-encode "constructor(address)" "$2" | cut -c3-)}
  retry cast send --rpc-url "$rpc" "${signer_flags[@]}" --json --create "$code" | jq -er 'select(.status == "0x1") | .contractAddress'
}

factory=$(jq -r '.uniswapV3.Factory // empty' "$registry")
if [ -z "$factory" ]; then
  factory=$(cast to-check-sum-address "$(create UniswapV3Factory)")
  record .uniswapV3.Factory "$factory"
  echo "Factory $factory"
fi
for entry in "SwapRouter02 UniswapV3Router" "QuoterV2 UniswapV3Quoter"; do
  read -r key contract <<<"$entry"
  [ -z "$(jq -r --arg key "$key" '.uniswapV3[$key] // empty' "$registry")" ] || continue
  address=$(cast to-check-sum-address "$(create "$contract" "$factory")")
  record ".uniswapV3.$key" "$address"
  echo "$key $address"
done
seeder=$(cast to-check-sum-address "$(create UniswapV3Seeder)")

for asset in ${BAND_ASSETS:-TSLA AMZN}; do
  name=${asset}_USDG_$fee
  [ -z "$(jq -r --arg name "$name" '.uniswapV3[$name] // empty' "$registry")" ] || { echo "$name: already in the registry"; continue; }
  read -r _ _ _ token <<<"$(cast call --rpc-url "$rpc" "$band" "asset(bytes32)(address,bytes32,bytes32,address)" "$(cast format-bytes32-string "$asset")" | tr '\n' ' ')"
  [ "$token" != $zero ] || { echo "$asset: no Stock Token in the band, skipped" >&2; continue; }
  read -r price _ <<<"$(cast call --rpc-url "$rpc" "$(field ".bandFeeds.$asset")" "latestRoundData()(uint80,int256,uint256,uint256,uint80)" | sed -n 2p)"
  read -r tick sqrt lower upper liquidity stock_amount <<<"$(python3 - "$token" "$usdg" "$price" "$POOL_USDG" "$(cast call --rpc-url "$rpc" "$token" "decimals()(uint8)")" "$spacing" <<'EOF'
import sys
from math import isqrt, log, floor, ceil
stock, usdg, price, pool_usdg, stock_dec, spacing = sys.argv[1], sys.argv[2], int(sys.argv[3]), int(sys.argv[4]), int(sys.argv[5]), int(sys.argv[6])
usdg_dec, half = 6, 1100
stock_is_token0 = int(stock, 16) < int(usdg, 16)
# price of token0 in token1, raw units; the band feed answers USD with 8 decimals
num, den = (price * 10**usdg_dec, 10**8 * 10**stock_dec) if stock_is_token0 else (10**8 * 10**stock_dec, price * 10**usdg_dec)
sqrt = isqrt(num * 2**192 // den)
tick = floor(log(num / den) / log(1.0001))
lower = floor((tick - half) / spacing) * spacing
upper = ceil((tick + half) / spacing) * spacing
rp = sqrt / 2**96
ra, rb = 1.0001 ** (lower / 2), 1.0001 ** (upper / 2)
if stock_is_token0:
    liquidity = int(pool_usdg / (rp - ra) * 0.98)
    stock_amount = liquidity * (1 / rp - 1 / rb)
else:
    liquidity = int(pool_usdg / (1 / rp - 1 / rb) * 0.98)
    stock_amount = liquidity * (rp - ra)
print(tick, sqrt, lower, upper, liquidity, int(stock_amount * 1.05) + 1)
EOF
)"
  pair=("$token" "$usdg")
  [ "$(tr '[:upper:]' '[:lower:]' <<<"$token")" \< "$(tr '[:upper:]' '[:lower:]' <<<"$usdg")" ] || pair=("$usdg" "$token")
  pool=$(cast call --rpc-url "$rpc" "$factory" "getPool(address,address,uint24)(address)" "${pair[0]}" "${pair[1]}" "$fee")
  if [ "$pool" = $zero ]; then
    send "$factory" "createPool(address,address,uint24)" "${pair[0]}" "${pair[1]}" "$fee"
    pool=$(cast call --rpc-url "$rpc" "$factory" "getPool(address,address,uint24)(address)" "${pair[0]}" "${pair[1]}" "$fee")
  fi
  slot0=$(cast call --rpc-url "$rpc" "$pool" "slot0()(uint160,int24,uint16,uint16,uint16,uint8,bool)")
  [ "$(sed -n 1p <<<"$slot0" | cut -d' ' -f1)" != 0 ] || send "$pool" "initialize(uint160)" "$sqrt"
  next_now=$(cast call --rpc-url "$rpc" "$pool" "slot0()(uint160,int24,uint16,uint16,uint16,uint8,bool)" | sed -n 5p | cut -d' ' -f1)
  for next in 500 1000 1500 1801; do
    [ "$next_now" -ge "$next" ] || send "$pool" "increaseObservationCardinalityNext(uint16)" "$next"
  done
  if [ "$(cast call --rpc-url "$rpc" "$pool" "liquidity()(uint128)" | cut -d' ' -f1)" = 0 ]; then
    held=$(cast call --rpc-url "$rpc" "$token" "balanceOf(address)(uint256)" "$signer" | cut -d' ' -f1)
    if [ "$held" -lt "$stock_amount" ]; then
      stock=$(cast call --rpc-url "$rpc" "$token" "underlying()(address)")
      send "$stock" "approve(address,uint256)" "$token" "$stock_amount"
      send "$token" "depositFor(address,uint256)" "$signer" "$((stock_amount - held))"
    fi
    send "$usdg" "approve(address,uint256)" "$seeder" "$POOL_USDG"
    send "$token" "approve(address,uint256)" "$seeder" "$stock_amount"
    send "$seeder" "seed(address,int24,int24,uint128)" "$pool" "$lower" "$upper" "$liquidity"
  fi
  record ".uniswapV3.$name" "$(cast to-check-sum-address "$pool")"
  echo "$name $(cast to-check-sum-address "$pool"), tick $tick in [$lower, $upper], liquidity $liquidity"
done
