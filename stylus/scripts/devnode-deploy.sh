#!/usr/bin/env bash
# Usage: devnode-deploy.sh RPC_URL PRIVATE_KEY
# Deploys stub Chainlink aggregators for NVDA and TSLA, writes them to stylus/target/devnode-registry.json,
# and deploys the band program configured from that file through StylusDeployer. Writes the program's
# address to stylus/target/devnode-band.
set -euo pipefail

rpc=$1 key=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
target=$root/stylus/target
mkdir -p "$target"

stub() {
  forge create --root "$root/contracts" test/devnode/StubAggregator.sol:StubAggregator --rpc-url "$rpc" \
    --private-key "$key" --broadcast --json --constructor-args 8 "$1" "$now" "$2" | jq -r .deployedTo
}

now=$(cast block --rpc-url "$rpc" latest -f timestamp)
jq -n --arg nvda "$(stub 22900000000 "NVDA / USD")" --arg tsla "$(stub 37800000000 "TSLA / USD")" \
  '{chainId: 412346, chainlink: {NVDA_USD: $nvda, TSLA_USD: $tsla}}' > "$target/devnode-registry.json"

args=$(BAND_ASSETS="NVDA TSLA" "$root/stylus/scripts/band-args.sh" "$target/devnode-registry.json")
read -r symbols feeds feed_ids <<<"$args"
(cd "$root/stylus/contracts/band" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
  --constructor-args "$symbols" "$feeds" "$feed_ids") > "$target/devnode-band.log"
cast to-check-sum-address "$(grep 'deployed code at address' "$target/devnode-band.log" | grep -o '0x[0-9a-f]\{40\}')" \
  > "$target/devnode-band"
tx=$(grep 'deployment tx hash' "$target/devnode-band.log" | grep -o '0x[0-9a-f]\{64\}')
logged=$(cast receipt --rpc-url "$rpc" "$tx" --json | jq -r --arg topic "$(cast keccak 'ContractDeployed(address)')" \
  '.logs[] | select(.address == "0xcecba2f1dc234f70dd89f2041029807f8d03a990" and .topics[0] == $topic) | .data')
[ "$(cast to-check-sum-address "0x${logged: -40}")" = "$(cat "$target/devnode-band")" ] ||
  { echo "StylusDeployer's ContractDeployed log names 0x${logged: -40}, not $(cat "$target/devnode-band")" >&2; exit 1; }
echo "band deployed on the dev node at $(cat "$target/devnode-band")."
