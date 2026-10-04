#!/bin/sh
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: verify-deployment.sh <indexer-base-url>
# Environment:
#   TAPEHOUSE_RPC_URL  the chain's RPC endpoint (required)
#   MAX_LAG            most blocks the indexer's head may trail the chain's (default 20)
#   SERVICE_URLS       space-separated URLs of the other services; each must answer with a status below 500
set -u

base=${1:-}
[ -n "$base" ] || { echo "usage: $0 <indexer-base-url>" >&2; exit 2; }
base=${base%/}
rpc=${TAPEHOUSE_RPC_URL:-}
[ -n "$rpc" ] || { echo "TAPEHOUSE_RPC_URL is required" >&2; exit 2; }
max_lag=${MAX_LAG:-20}
failed=0

fail() { echo "FAIL $*"; failed=1; }
ok() { echo "ok   $*"; }

status=$(curl -fsS --max-time 15 "$base/v1/status") || { fail "the API does not answer at $base/v1/status"; exit 1; }
ok "the API answers"

head=$(printf '%s' "$status" | sed -n 's/.*"head":\([0-9][0-9]*\).*/\1/p')
[ -n "$head" ] || { fail "the status carries no head: $status"; exit 1; }

reply=$(curl -fsS --max-time 15 -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"eth_blockNumber","params":[]}' "$rpc") || { fail "the RPC does not answer"; exit 1; }
hex=$(printf '%s' "$reply" | sed -n 's/.*"result":"0x\([0-9a-fA-F][0-9a-fA-F]*\)".*/\1/p')
[ -n "$hex" ] || { fail "the RPC returned no block number"; exit 1; }
chain=$((16#$hex))

lag=$((chain - head))
if [ "$lag" -le "$max_lag" ]; then
  ok "indexer head $head, chain head $chain, lag $lag within $max_lag"
else
  fail "indexer head $head trails chain head $chain by $lag blocks, more than $max_lag"
fi

for url in ${SERVICE_URLS:-}; do
  code=$(curl -sS -o /dev/null -w '%{http_code}' --max-time 15 "$url") || code=000
  if [ "$code" != 000 ] && [ "$code" -lt 500 ]; then
    ok "$url answers $code"
  else
    fail "$url answers $code"
  fi
done

exit "$failed"
