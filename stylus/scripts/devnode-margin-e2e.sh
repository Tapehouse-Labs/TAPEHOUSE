#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-margin-e2e.sh RPC_URL PRIVATE_KEY MARGIN_ADDRESS DEPLOYMENTS_JSON
# Checks the margin program deployed by devnode-deploy.sh for NVDA, TSLA and SPY: its configuration against the
# registry and the dev node's band, its scenario set and a requirement against the references in testdata, the
# bit a pool without history sets, the current requirement against the band's session and without a band, every
# refused update, an accepted one and the day that must pass before the next, and a constructor that refuses a
# matrix that is not positive definite. Measures the scenario set by lattice size, the requirement of three and
# of six assets, the current requirement of three and of one position among six, and an update of every
# volatility on programs with the six launch assets and with eight. Appends its L2 gas to
# stylus/target/devnode-gas.txt.
set -euo pipefail

rpc=$1 key=$2 margin=$3 registry=$4
root=$(cd "$(dirname "$0")/../.." && pwd)
node_interface=0x00000000000000000000000000000000000000C8
other=0x70997970C51812dc3A010C7d01b50e0d17dc79C8
nvda=$(cast format-bytes32-string NVDA)
spy=$(cast format-bytes32-string SPY)
gas_log=$root/stylus/target/devnode-gas.txt
owner=$(cast wallet address --private-key "$key")
update="setParameters(uint32[],uint16[],uint32[],uint32[])"
requirement="requirement(int256[],uint256[],uint64,bool)"
gaps="[118601,135345,54840]"
depths="[3561773,1377156,164067,166576,320861,877772]"

fail() { echo "FAIL: $*"; exit 1; }
record() {
  echo "$1 $2" >> "$gas_log"
  [[ $2 =~ ^[0-9]+$ ]] || { echo "FAIL: no L2 gas for $1" >&2; exit 1; }
  echo "$2"
}
l2_gas() {
  cast call --rpc-url "$rpc" --from "$owner" $node_interface \
    "gasEstimateComponents(address,bool,bytes)(uint64,uint64,uint256,uint256)" "${2:-$margin}" false "$1" |
    head -n 2 | cut -d' ' -f1 | { read -r total; read -r l1; echo $((total - l1)); }
}
measure() { record "$1" "$(l2_gas "$(cast calldata "$2" "${@:3}")")"; }
expect_refusal() {
  local out
  if out=$(cast call --rpc-url "$rpc" --from "$1" "$margin" "$update" "$2" "$3" "$4" "$5" 2>&1); then
    fail "expected $6, got $out"
  fi
  grep -q "$(cast sig "$6")" <<<"$out" || fail "expected $6, got $out"
}

BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry"
if out=$(BAND_ASSETS="NVDA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry"); then
  fail "check-margin.sh accepted a program that differs from the registry"
fi
grep -q '^FAIL: margin holds the assets' <<<"$out" || fail "check-margin.sh refused a narrower registry for another reason: $out"
jq --arg band "$other" '.tapehouse.Band = $band' "$registry" > "$registry.nothing"
if out=$(BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry.nothing" 2>&1); then
  fail "check-margin.sh accepted a band that answers nothing"
fi
grep -q '^FAIL: band' <<<"$out" || fail "check-margin.sh refused a band that answers nothing for another reason: $out"
jq 'del(.tapehouse.Band)' "$registry" > "$registry.unbanded"
if out=$(BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry.unbanded" 2>&1); then
  fail "check-margin.sh accepted a registry without the program's band"
fi
grep -q "^FAIL: margin's band is" <<<"$out" || fail "check-margin.sh refused a registry without the band for another reason: $out"
jq --arg token "$other" '.tokens.NVDA = $token' "$registry" > "$registry.token"
if out=$(BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry.token" 2>&1); then
  fail "check-margin.sh accepted a registry whose Stock Token NVDA's pool does not trade"
fi
grep -q "^FAIL: margin's pool for NVDA" <<<"$out" || fail "check-margin.sh refused another Stock Token for another reason: $out"
vols="[31352,37436,11335]" corrs="[4637,7203,6232]"
expect_refusal "$other" "$vols" "$corrs" "$gaps" "$depths" "OwnableUnauthorizedAccount(address)"
expect_refusal "$owner" "[47029,37436,11335]" "$corrs" "$gaps" "$depths" "VolatilityStepTooLarge(bytes32,uint32,uint32)"
expect_refusal "$owner" "[31351,37436,11335]" "$corrs" "$gaps" "$depths" "InvalidVolatility(bytes32,uint32,uint32)"
expect_refusal "$owner" "$vols" "[4637,9990,5178]" "$gaps" "$depths" "NotPositiveDefinite()"
expect_refusal "$owner" "$vols" "$corrs" "[177902,135345,54840]" "$depths" "GapStepTooLarge(bytes32,uint32,uint32)"
expect_refusal "$owner" "$vols" "$corrs" "$gaps" "[3561774,1377156,164067,166576,320861,877772]" \
  "InvalidDepth(bytes32,uint32,uint32)"
expect_refusal "$owner" "[31352,37436]" "$corrs" "$gaps" "$depths" "LengthMismatch()"

vectors=$root/stylus/contracts/margin/testdata/requirement-vectors.json
list() { jq -r "$1"' | map(tostring) | "[" + join(",") + "]"' "$vectors"; }
quantities=$(list .devnode.quantities)
prices=$(list .devnode.prices)
expected="$(jq -r '.devnode.requirement' "$vectors") $(jq -r '.devnode.missing' "$vectors")"
held=$(cast call --rpc-url "$rpc" "$margin" "$requirement(uint256,uint8)" "$quantities" "$prices" 172800 false | sed 's/ \[[^]]*\]//' | tr '\n' ' ')
[ "$held" = "$expected " ] || fail "the requirement on chain is $held, the reference $expected"
nvda_only=$(list '.devnode.quantities | [.[0], "0", "0"]')
pool=$(jq -r '.uniswapV3.NVDA_USDG_500' "$registry")
live=$(cast call --rpc-url "$rpc" "$margin" "$requirement(uint256,uint8)" "$nvda_only" "$prices" 172800 false | sed 's/ \[[^]]*\]//' | tr '\n' ' ')
[ "${live#* }" = "0 " ] || fail "NVDA's pool did not count: $live"
cast send --rpc-url "$rpc" --private-key "$key" "$pool" "setYoung(bool)" true > /dev/null
young=$(cast call --rpc-url "$rpc" "$margin" "$requirement(uint256,uint8)" "$nvda_only" "$prices" 172800 false | sed 's/ \[[^]]*\]//' | tr '\n' ' ')
cast send --rpc-url "$rpc" --private-key "$key" "$pool" "setYoung(bool)" false > /dev/null
less() { [ ${#1} -lt ${#2} ] || { [ ${#1} -eq ${#2} ] && [[ $1 < $2 ]]; }; }
[ "${young#* }" = "1 " ] && [ "${young%% *}" = "${live%% *}" ] ||
  fail "a pool without 30 minutes of history did not set NVDA's bit, or changed its charge: $young, with it $live"
spy_only=$(list '.devnode.quantities | ["0", "0", .[2]]')
feed=$(jq -r '.chainlink.ETH_USD' "$registry")
now=$(cast block --rpc-url "$rpc" latest -f timestamp)
live=$(cast call --rpc-url "$rpc" "$margin" "$requirement(uint256,uint8)" "$spy_only" "$prices" 172800 false | sed 's/ \[[^]]*\]//' | tr '\n' ' ')
cast send --rpc-url "$rpc" --private-key "$key" "$feed" "setRound(int256,uint256)" 268330550000 $((now - 86461)) > /dev/null
stale=$(cast call --rpc-url "$rpc" "$margin" "$requirement(uint256,uint8)" "$spy_only" "$prices" 172800 false | sed 's/ \[[^]]*\]//' | tr '\n' ' ')
cast send --rpc-url "$rpc" --private-key "$key" "$feed" "setRound(int256,uint256)" 268330550000 "$now" > /dev/null
[ "${live#* }" = "0 " ] && [ "${stale#* }" = "4 " ] && ! less "${live%% *}" "${stale%% *}" ||
  fail "a stale ETH/USD answer did not set SPY's bit and drop its pool's discount: $stale, with it $live"
band=$(jq -r '.tapehouse.Band' "$registry")
read -r open_code nyse nyse_next change_ms boundary_ms <<<"$(cast call --rpc-url "$rpc" "$band" "session()(uint8,uint8,uint8,uint64,uint64)" | sed 's/ \[[^]]*\]//' | tr '\n' ' ')"
now_ms=$(( $(cast block --rpc-url "$rpc" latest -f timestamp) * 1000 ))
requirement_at() { cast call --rpc-url "$rpc" "$margin" "$requirement(uint256,uint8)" "$quantities" "$prices" "$1" "$2" | head -n 1 | cut -d' ' -f1; }
held=$(cast call --rpc-url "$rpc" "$margin" "currentRequirement(int256[],uint256[])(uint256,uint8,uint8)" "$quantities" "$prices" | sed 's/ \[[^]]*\]//' | tr '\n' ' ')
expected_current() {
  python3 -c '
import json, sys
open_, closed, instant, code, boundary, now = map(int, sys.argv[1:7])
values = [(1 if q * p >= 0 else -1) * (abs(q * p) // 10**8) for q, p in zip(json.loads(sys.argv[7]), json.loads(sys.argv[8]))]
moves = [150_000, 150_000, 80_000 if values[2] > 0 else 60_000]
addon = instant - max(abs(e) * m // 10**6 for e, m in zip(values, moves))
ramp = 25_200_000
buffered = -(-open_ * 5 // 4)
capped = -(-sum(abs(e) for e in values) * 10_000 // 50_000) + addon
across = max(buffered, closed, capped)
if code != 2:
    print(across, code, "cap" if capped > max(buffered, closed) else "model")
elif boundary != 0 and boundary <= now:
    print(across, 1)
elif boundary != 0 and boundary - now < ramp:
    elapsed = ramp - (boundary - now)
    print(buffered + (across - buffered) * elapsed // ramp, 3)
else:
    print(buffered, 2)
' "$@"
}
expected=$(expected_current "$(requirement_at 172800 false)" "$(requirement_at 172800 true)" "$(requirement_at 0 false)" \
  "$open_code" "$boundary_ms" "$now_ms" "$quantities" "$prices")
expected=${expected% cap}
expected=${expected% model}
read -r value _ code <<<"$held"
[ "$(cast call --rpc-url "$rpc" "$margin" "weekendLeverage()(uint32)" | cut -d' ' -f1)" = 50000 ] ||
  fail "the weekend leverage cap is not 5×"
[ "$value $code" = "$expected" ] ||
  fail "the current requirement is $value in regime $code; from the band's session $open_code $nyse $nyse_next $change_ms $boundary_ms it is $expected"
later=$((now_ms / 1000 + 7200))
at_later() {
  local signature=$1
  shift
  cast abi-decode "$signature" "$(cast rpc --rpc-url "$rpc" eth_call \
    "{\"to\":\"$margin\",\"data\":\"$(cast calldata "${signature%%)(*})" "$@")\"}" latest '{}' \
    "{\"time\":\"$(printf '0x%x' "$later")\"}" | tr -d '"')" | sed 's/ \[[^]]*\]//' | tr '\n' ' '
}
hedge=$(python3 -c 'import json, sys; q, p = json.loads(sys.argv[1]), json.loads(sys.argv[2]); print(f"[{q[0]},0,{-(q[0] * p[0] // p[2])}]")' \
  "$quantities" "$prices")
hedge_at() { at_later "$requirement(uint256,uint8)" "$hedge" "$prices" "$1" "$2" | cut -d' ' -f1; }
read -r value _ code <<<"$(at_later "currentRequirement(int256[],uint256[])(uint256,uint8,uint8)" "$hedge" "$prices")"
[ "$value $code cap" = "$(expected_current "$(hedge_at 172800 false)" "$(hedge_at 172800 true)" "$(hedge_at 0 false)" 0 0 \
  "$now_ms" "$hedge" "$prices")" ] ||
  fail "long NVDA against short SPY in an unknown session needs $value in regime $code, not a fifth of its gross plus its add-on"
current_gas=$(measure "margin.currentRequirement(3)" "currentRequirement(int256[],uint256[])" "$quantities" "$prices")
estimate_at() {
  cast rpc --rpc-url "$rpc" eth_estimateGas "{\"to\":\"$margin\",\"data\":\"$1\"}" latest '{}' \
    "{\"time\":\"$(printf '0x%x' "$2")\"}" | tr -d '"'
}
call_data=$(cast calldata "currentRequirement(int256[],uint256[])" "$quantities" "$prices")
spread=$(($(estimate_at "$call_data" $((now_ms / 1000 + 7200))) - $(estimate_at "$call_data" $((now_ms / 1000)))))
[ "${spread#-}" -le 200 ] ||
  fail "currentRequirement's gas moves by $spread from regime $code to an unknown session, beyond the gas gate's noise"
echo "margin current requirement: regime $code from the band's session, at most 5× across a closure; L2 gas $current_gas, the same within ${spread#-} in an unknown session"
echo "margin requirement: the requirement on chain is the reference requirement; a pool without history or a stale ether price sets its bit; L2 gas \
$(measure "margin.requirement(3)" "$requirement" "$quantities" "$prices" 172800 false) for three assets"

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
gas=$(record "margin.setParameters(3)" "$(l2_gas "$(cast calldata "$update" "[40000,37436,11335]" "[4637,7203,7000]" "$gaps" "$depths")")")
receipt=$(cast send --rpc-url "$rpc" --private-key "$key" "$margin" "$update" "[40000,37436,11335]" "[4637,7203,7000]" "$gaps" "$depths" --json)
[ "$(jq -r .status <<<"$receipt")" = 0x1 ] || fail "an update within every bound reverted"
[ "$(jq '.logs | length' <<<"$receipt")" = 2 ] || fail "the update did not log exactly the two values that changed"
[ "$(cast call --rpc-url "$rpc" "$margin" "volatility(bytes32)(uint32,uint32)" "$nvda" | sed 's/ \[[^]]*\]//' | tr '\n' ' ')" = "40000 31352 " ] ||
  fail "the update did not store NVDA's volatility"
BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/check-margin.sh" "$rpc" "$margin" "$registry" > /dev/null
expect_refusal "$owner" "[40000,37436,11335]" "[4637,7203,7000]" "$gaps" "$depths" "UpdateTooSoon(uint64)"
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
column() {
  local entries=${1#[(} out="" field
  entries=${entries%)]}
  while IFS=, read -r -a field; do
    for k in "${@:2}"; do out+=",${field[$k]}"; done
  done <<<"${entries//),(/$'\n'}"
  echo "[${out#,}]"
}
six=$(deploy_margin "${args[@]}" | address_of)
six_gas=$(record "margin.setParameters(6)" "$(l2_gas "$(cast calldata "$update" "$(up_a_fifth "$(column "${args[0]}" 2)")" \
  "${args[2]}" "$(column "${args[0]}" 4)" "$(column "${args[0]}" 5 6)")" "$six")")
zero=0x0000000000000000000000000000000000000000
assets=$(for i in 0 1 2 3 4 5 6 7; do echo "($(cast format-bytes32-string "A$i"),1,10000,1,10000,1000000,1000000,$zero,$zero)"; done |
  paste -sd, -)
ones="[$(for _ in $(seq 28); do echo 1; done | paste -sd, -)]"
eights="[10000,10000,10000,10000,10000,10000,10000,10000]"
sixteen="[$(for _ in $(seq 16); do echo 1000000; done | paste -sd, -)]"
eight=$(deploy_margin "[$assets]" "$ones" "$ones" 0x0000000000000000000000000000000000000000000000000000000000000000 \
  $zero $zero $zero $zero "$owner" | address_of)
eight_gas=$(record "margin.setParameters(8)" \
  "$(l2_gas "$(cast calldata "$update" "[12000,12000,12000,12000,12000,12000,12000,12000]" "$ones" "$eights" "$sixteen")" "$eight")")
echo "margin update of every volatility: L2 gas $six_gas for the six launch assets, $eight_gas for eight"
read -r value missing code <<<"$(cast call --rpc-url "$rpc" "$eight" "currentRequirement(int256[],uint256[])(uint256,uint8,uint8)" \
  "[0,0,0,0,0,0,0,0]" "$eights" | tr '\n' ' ')"
[ "$value $missing $code" = "0 0 0" ] ||
  fail "a margin program without a band gave $value, $missing in regime $code, not an unknown session"
echo "margin scenarios for the six launch assets: L2 gas $(record "margin.scenarioDigest(256,launch)" \
  "$(l2_gas "$(cast calldata "scenarioDigest(uint16,uint64)" 256 172800)" "$six")")"
launch_quantities=$(list '[.vectors[] | select(.label == "100k each of the six")][0].quantities')
echo "margin requirement of the six launch assets: L2 gas $(record "margin.requirement(6)" \
  "$(l2_gas "$(cast calldata "$requirement" "$launch_quantities" "$(list .prices)" 172800 false)" "$six")")"
isolated=$(list '[.vectors[] | select(.label == "100k each of the six")][0].quantities | [.[0]] + [range(5) | 0]')
echo "margin current requirement of one position among the six launch assets: L2 gas $(record "margin.currentRequirement(1 of 6)" \
  "$(l2_gas "$(cast calldata "currentRequirement(int256[],uint256[])" "$isolated" "$(list .prices)")" "$six")")"

words=$(BAND_ASSETS="NVDA TSLA SPY" "$(dirname "$0")/margin-args.sh" "$registry")
read -r -a args <<<"$words"
if out=$(deploy_margin "${args[0]}" "[1000,1000,1000]" "[9000,9000,1000]" "${args[@]:3}"); then
  fail "a correlation matrix that is not positive definite was accepted"
fi
refused_with() {
  grep -q "$(cast sig "ContractInitializationError(address,bytes)")" <<<"$2" && grep -q "$(cast sig "$1" | cut -c3-)" <<<"$2" ||
    fail "expected $1 inside ContractInitializationError, got: $(tail -n 1 <<<"$2")"
}
refused_with "NotPositiveDefinite()" "$out"
if out=$(deploy_margin "${args[@]:0:4}" $zero $zero "${args[@]:6}"); then
  fail "a pool against neither USDG nor WETH was accepted"
fi
refused_with "InvalidPool(bytes32,address)" "$out"
if out=$(deploy_margin "${args[@]:0:6}" $zero "${args[@]:7}"); then
  fail "a WETH pool without an ETH/USD feed was accepted"
fi
refused_with "InvalidFeed(address)" "$out"
unmeasured=$(grep -v ' [0-9][0-9]*$' "$gas_log" || true)
[ -z "$unmeasured" ] || fail "no L2 gas for: $unmeasured"
echo "PASS"
