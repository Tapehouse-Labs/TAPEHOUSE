#!/usr/bin/env bash
# Usage: devnode-margin-e2e.sh RPC_URL PRIVATE_KEY MARGIN_ADDRESS DEPLOYMENTS_JSON
# Checks the margin program deployed by devnode-deploy.sh for NVDA, TSLA and SPY: its configuration against the
# registry and the dev node's band, its scenario set against the reference in testdata, every refused update, an
# accepted one and the day that must pass before the next, and a constructor that refuses a matrix that is not
# positive definite. Measures the scenario set by lattice size, and an update of every volatility on programs
# with the six launch assets and with eight. Appends its L2 gas to stylus/target/devnode-gas.txt.
set -euo pipefail

rpc=$1 key=$2 margin=$3 registry=$4
root=$(cd "$(dirname "$0")/../.." && pwd)
node_interface=0x00000000000000000000000000000000000000C8
other=0x70997970C51812dc3A010C7d01b50e0d17dc79C8
nvda=$(cast format-bytes32-string NVDA)
spy=$(cast format-bytes32-string SPY)
gas_log=$root/stylus/target/devnode-gas.txt
owner=$(cast wallet address --private-key "$key")
update="setParameters(uint32[],uint16[],uint32[])"
gaps="[118601,135345,54840]"

fail() { echo "FAIL: $*"; exit 1; }
record() { echo "$1 $2" >> "$gas_log"; echo "$2"; }
l2_gas() {
  cast call --rpc-url "$rpc" --from "$owner" $node_interface \
    "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" "${2:-$margin}" false "$1" |
    head -n 2 | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}
measure() { record "$1" "$(l2_gas "$(cast calldata "$2" "${@:3}")")"; }
expect_refusal() {
  local out
  if out=$(cast call --rpc-url "$rpc" --from "$1" "$margin" "$update" "$2" "$3" "$4" 2>&1); then
    fail "expected $5, got $out"
  fi
  grep -q "$(cast sig "$5")" <<<"$out" || fail "expected $5, got $out"
}

jq --arg band "$(cat "$root/stylus/target/devnode-band")" '.tapehouse.Band = $band' "$registry" > "$registry.band"
BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry.band"
if out=$(BAND_ASSETS="NVDA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry"); then
  fail "check-margin.sh accepted a program that differs from the registry"
fi
grep -q '^FAIL: margin holds the assets' <<<"$out" || fail "check-margin.sh refused a narrower registry for another reason: $out"
jq --arg band "$other" '.tapehouse.Band = $band' "$registry" > "$registry.nothing"
if out=$(BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry.nothing" 2>&1); then
  fail "check-margin.sh accepted a band that answers nothing"
fi
grep -q '^FAIL: band' <<<"$out" || fail "check-margin.sh refused a band that answers nothing for another reason: $out"
expect_refusal "$other" "[31352,37436,11335]" "[4637,7203,6232]" "$gaps" "OwnableUnauthorizedAccount(address)"
expect_refusal "$owner" "[47029,37436,11335]" "[4637,7203,6232]" "$gaps" "VolatilityStepTooLarge(bytes32,uint32,uint32)"
expect_refusal "$owner" "[31351,37436,11335]" "[4637,7203,6232]" "$gaps" "InvalidVolatility(bytes32,uint32,uint32)"
expect_refusal "$owner" "[31352,37436,11335]" "[4637,9990,5178]" "$gaps" "NotPositiveDefinite()"
expect_refusal "$owner" "[31352,37436,11335]" "[4637,7203,6232]" "[177902,135345,54840]" "GapStepTooLarge(bytes32,uint32,uint32)"
expect_refusal "$owner" "[31352,37436]" "[4637,7203,6232]" "$gaps" "LengthMismatch()"

reference=$(jq -r '.digests[] | select(.config == "dev node" and .horizon == 172800 and .size == 256) | .digest' \
  "$root/stylus/contracts/margin/testdata/scenario-vectors.json")
[ "$(cast call --rpc-url "$rpc" "$margin" "scenarioDigest(uint16,uint64)(bytes32)" 256 172800)" = "$reference" ] ||
  fail "the scenario set on chain is not the reference set"
sizes=""
for size in 32 64 128 256; do
  sizes+=" $size $(measure "margin.scenarioDigest($size)" "scenarioDigest(uint16,uint64)" "$size" 172800),"
done
echo "margin scenarios: the set on chain is the reference set; L2 gas by lattice size${sizes%,}; one scenario $(measure \
  "margin.scenario(0)" "scenario(uint16,uint64)" 0 172800)"

echo "margin views: L2 gas assets $(measure "margin.assets()" "assets()"), volatility $(measure "margin.volatility(NVDA)" \
  "volatility(bytes32)" "$nvda"), correlation $(measure "margin.correlation(NVDA,SPY)" "correlation(bytes32,bytes32)" "$nvda" "$spy")"
gas=$(record "margin.setParameters(3)" "$(l2_gas "$(cast calldata "$update" "[40000,37436,11335]" "[4637,7203,7000]" "$gaps")")")
receipt=$(cast send --rpc-url "$rpc" --private-key "$key" "$margin" "$update" "[40000,37436,11335]" "[4637,7203,7000]" "$gaps" --json)
[ "$(jq -r .status <<<"$receipt")" = 0x1 ] || fail "an update within every bound reverted"
[ "$(jq '.logs | length' <<<"$receipt")" = 2 ] || fail "the update did not log exactly the two values that changed"
[ "$(cast call --rpc-url "$rpc" "$margin" "volatility(bytes32)(uint32,uint32)" "$nvda" | sed 's/ \[[^]]*\]//' | tr '\n' ' ')" = "40000 31352 " ] ||
  fail "the update did not store NVDA's volatility"
BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry" > /dev/null
expect_refusal "$owner" "[40000,37436,11335]" "[4637,7203,7000]" "$gaps" "UpdateTooSoon(uint64)"
echo "margin update of 3 assets: L2 gas $gas; the next one must wait a day"

deploy_margin() {
  (cd "$root/stylus/contracts/margin" && cargo stylus deploy --no-verify -e "$rpc" --private-key "$key" \
    --constructor-args "$@" 2>&1)
}
address_of() { grep 'deployed code at address' | grep -o '0x[0-9a-fA-F]\{40\}'; }
up_a_fifth() { tr -d '[]' <<<"$1" | tr ',' '\n' | awk '{ printf "%s%d", (NR > 1 ? "," : "["), $1 * 6 / 5 } END { print "]" }'; }
jq '.chainlink += {AAPL_USD: .chainlink.NVDA_USD, MSFT_USD: .chainlink.NVDA_USD, GOOGL_USD: .chainlink.NVDA_USD}' \
  "$registry" > "$registry.launch"
words=$("$(dirname "$0")/margin-args.sh" "$registry.launch")
read -r -a args <<<"$words"
six=$(deploy_margin "${args[@]}" | address_of)
six_gas=$(record "margin.setParameters(6)" \
  "$(l2_gas "$(cast calldata "$update" "$(up_a_fifth "${args[2]}")" "${args[4]}" "${args[6]}")" "$six")")
symbols=$(for i in 0 1 2 3 4 5 6 7; do cast format-bytes32-string "A$i"; done | paste -sd, -)
ones="[$(for _ in $(seq 28); do echo 1; done | paste -sd, -)]"
eights="[10000,10000,10000,10000,10000,10000,10000,10000]"
eight=$(deploy_margin "[$symbols]" "[1,1,1,1,1,1,1,1]" "$eights" "$ones" "$ones" "[1,1,1,1,1,1,1,1]" "$eights" \
  0x0000000000000000000000000000000000000000000000000000000000000000 "$owner" | address_of)
eight_gas=$(record "margin.setParameters(8)" \
  "$(l2_gas "$(cast calldata "$update" "[12000,12000,12000,12000,12000,12000,12000,12000]" "$ones" "$eights")" "$eight")")
echo "margin update of every volatility: L2 gas $six_gas for the six launch assets, $eight_gas for eight"
echo "margin scenarios for the six launch assets: L2 gas $(record "margin.scenarioDigest(256,launch)" \
  "$(l2_gas "$(cast calldata "scenarioDigest(uint16,uint64)" 256 172800)" "$six")")"

words=$(BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/margin-args.sh" "$registry")
read -r -a args <<<"$words"
if out=$(deploy_margin "${args[0]}" "${args[1]}" "${args[2]}" "[1000,1000,1000]" "[9000,9000,1000]" "${args[@]:5}"); then
  fail "a correlation matrix that is not positive definite was accepted"
fi
if ! grep -q 0x88d8f57d <<<"$out" || ! grep -q "$(cast sig "NotPositiveDefinite()" | cut -c3-)" <<<"$out"; then
  fail "expected NotPositiveDefinite inside ContractInitializationError, got: $(tail -n 1 <<<"$out")"
fi
echo "PASS"
