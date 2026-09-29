#!/usr/bin/env bash
# Usage: devnode-deploy.sh RPC_URL PRIVATE_KEY
# Deploys stub Chainlink aggregators for NVDA, TSLA and SPY and stub Stock Tokens for NVDA and SPY, writes
# them to stylus/target/devnode-registry.json with Anvil's second test account as the halt signer and the
# deploying account as the owner, and deploys the band and margin programs configured from that file through
# StylusDeployer. Writes each program's address to stylus/target/devnode-<program>.
set -euo pipefail

rpc=$1 key=$2
root=$(cd "$(dirname "$0")/../.." && pwd)
target=$root/stylus/target
mkdir -p "$target"

stub() {
  forge create --root "$root/contracts" test/devnode/StubAggregator.sol:StubAggregator --rpc-url "$rpc" \
    --private-key "$key" --broadcast --json --constructor-args 8 "$1" "$now" "$2" | jq -r .deployedTo
}

token() {
  forge create --root "$root/contracts" test/devnode/StubStockToken.sol:StubStockToken --rpc-url "$rpc" \
    --private-key "$key" --broadcast --json --constructor-args "$1" | jq -r .deployedTo
}

now=$(cast block --rpc-url "$rpc" latest -f timestamp)
jq -n --arg nvda "$(stub 22900000000 "NVDA / USD")" --arg tsla "$(stub 37800000000 "TSLA / USD")" \
  --arg spy "$(stub 77232802713 "SPY / USD")" --arg nvda_token "$(token 1000775159164630595)" \
  --arg spy_token "$(token 1001717991187472003)" --arg owner "$(cast wallet address --private-key "$key")" \
  '{chainId: 412346, chainlink: {NVDA_USD: $nvda, TSLA_USD: $tsla, SPY_USD: $spy}, tokens: {NVDA: $nvda_token, SPY: $spy_token},
    tapehouse: {HaltSigner: "0x70997970C51812dc3A010C7d01b50e0d17dc79C8", Owner: $owner}}' \
  > "$target/devnode-registry.json"

deploy() {
  local program=$1 words args logged tx
  words=$(BAND_ASSETS="NVDA TSLA SPY" "$root/stylus/scripts/$program-args.sh" "$target/devnode-registry.json")
  read -r -a args <<<"$words"
  (cd "$root/stylus/contracts/$program" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
    --constructor-args "${args[@]}") > "$target/devnode-$program.log"
  cast to-check-sum-address "$(grep 'deployed code at address' "$target/devnode-$program.log" | grep -o '0x[0-9a-f]\{40\}')" \
    > "$target/devnode-$program"
  tx=$(grep 'deployment tx hash' "$target/devnode-$program.log" | grep -o '0x[0-9a-f]\{64\}')
  logged=$(cast receipt --rpc-url "$rpc" "$tx" --json | jq -r --arg topic "$(cast keccak 'ContractDeployed(address)')" \
    '.logs[] | select(.address == "0xcecba2f1dc234f70dd89f2041029807f8d03a990" and .topics[0] == $topic) | .data')
  [ "$(cast to-check-sum-address "0x${logged: -40}")" = "$(cat "$target/devnode-$program")" ] ||
    { echo "StylusDeployer's ContractDeployed log names 0x${logged: -40}, not $(cat "$target/devnode-$program")" >&2; exit 1; }
  echo "$program deployed on the dev node at $(cat "$target/devnode-$program")."
}

deploy band
deploy margin
