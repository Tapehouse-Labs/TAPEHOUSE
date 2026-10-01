#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-reopening-auction-e2e.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Checks the reopening auction the registry names on the dev node: it is the liquidator's auction, over it; any other
# caller is refused the liquidator's hold, settle and collect and the backstop's buyRemainder; and its phase follows the
# band's session as it stands: outside the night before a regular open that follows a closure, nothing is enrolled or
# committed.
set -euo pipefail

rpc=$1 key=$2 registry=$3
liquidator=$(jq -r .tapehouse.Liquidator "$registry")
auction=$(jq -r .tapehouse.ReopeningAuction "$registry")
backstop=$(jq -r .tapehouse.GapBackstop "$registry")
band=$(jq -r .tapehouse.Band "$registry")
me=$(cast wallet address --private-key "$key")
nvda=$(cast format-bytes32-string NVDA)
fail() { echo "FAIL: $*"; exit 1; }
read_() { cast call --rpc-url "$rpc" "$@" | cut -d' ' -f1; }
refused() {
  local out
  out=$(cast call --rpc-url "$rpc" --from "$me" "$@" 2>&1) && fail "$2 went through"
  grep -q "$1" <<<"$out" || fail "$2 reverted otherwise: $out"
}

[ "$(read_ "$liquidator" "auction()(address)")" = "$auction" ] || fail "the liquidator's auction is not the registry's"
[ "$(read_ "$auction" "liquidator()(address)")" = "$liquidator" ] || fail "the auction's liquidator is not the registry's"
refused "$(cast sig "NotAuction(address)")" "$liquidator" "hold(address,bytes32,uint64)" "$me" "$nvda" 1
refused "$(cast sig "NotAuction(address)")" "$liquidator" "settle(address,bytes32,address,uint256,uint256)" "$me" "$nvda" "$me" 1 1
refused "$(cast sig "NotAuction(address)")" "$liquidator" "collect(uint256)" 1
refused "$(cast sig "NotAuction(address)")" "$backstop" "buyRemainder(address,bytes32,address,uint256,uint256)" "$me" "$nvda" "$me" 1 1
read -r state nyse next _ boundary <<<"$(cast call --rpc-url "$rpc" "$band" "session()(uint8,uint8,uint8,uint64,uint64)" | cut -d' ' -f1 | tr '\n' ' ')"
read -r open revealing <<<"$(cast call --rpc-url "$rpc" "$auction" "phase()(uint64,bool)" | cut -d' ' -f1 | tr '\n' ' ')"
if [ "$state" = 2 ] && [ "$nyse" = 3 ] && [ "$boundary" = 0 ] && [ "$next" != 0 ]; then
  [ "$open" != 0 ] || fail "the night before a regular open has no phase"
else
  [ "$open" = 0 ] && [ "$revealing" = false ] || fail "session $state $nyse $next has phase $open $revealing"
  refused "$(cast sig "WrongPhase()")" "$auction" "enroll(address,bytes32,bytes32)" "$me" 0x0000000000000000000000000000000000000000000000000000000000000000 "$nvda"
  refused "$(cast sig "WrongPhase()")" "$auction" "commit(bytes32,bytes32,uint256)" "$nvda" 0x0000000000000000000000000000000000000000000000000000000000000001 100000000
fi
echo "reopening auction $auction: the liquidator's, refusing other callers; session $state $nyse $next, phase $open $revealing"
echo "PASS"
