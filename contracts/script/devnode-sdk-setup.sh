#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: devnode-sdk-setup.sh RPC_URL PRIVATE_KEY DEPLOYMENTS_JSON
# Prepares the dev node for the SDKs' examples: the stub router and quoter trade SPY at its band's centre, the
# account lends 1 SPY through the margin accounts, so SPY's lending vault has a token to lend the shorts, and holds
# 1,000 USDG of margin, 6 NVDA and 3 SPY for the examples' mints of the basket PAIR, and 8,000 USDG for the
# three SDKs' and the MCP server's examples to write gap cover with and buy it. Checks that the registry names SPY's
# band feed, the quoter over the stub router, and the basket PAIR. Each run adds to what the last left.
set -euo pipefail

rpc=$1 key=$2 registry=$3
accounts=$(jq -r .tapehouse.MarginAccounts "$registry")
router=$(jq -r .uniswapV3.SwapRouter02 "$registry")
spy=$(jq -r .tokens.SPY "$registry")
nvda=$(jq -r .tokens.NVDA "$registry")
me=$(cast wallet address --private-key "$key")
symbol=$(cast format-bytes32-string SPY)
one=1000000000000000000
fail() { echo "FAIL: $*"; exit 1; }
send() { cast send --rpc-url "$rpc" --private-key "$key" "$@" > /dev/null; }

jq -e '.bandFeeds.SPY' "$registry" > /dev/null || fail "the registry has no band feed for SPY"
jq -e '.tapehouse.Baskets.PAIR' "$registry" > /dev/null || fail "the registry has no basket PAIR"
[ "$(cast call --rpc-url "$rpc" "$(jq -r .uniswapV3.QuoterV2 "$registry")" "factory()(address)")" = "$router" ] ||
  fail "the registry's quoter does not quote the stub router"
band=$(jq -r .tapehouse.Band "$registry")
mid=$(cast call --rpc-url "$rpc" "$band" "quote(bytes32)(uint8,uint8,uint64,uint64,uint64,uint128)" "$symbol" |
  sed -n 3p | cut -d' ' -f1)
send "$router" "setPrice(address,uint256)" "$spy" "$mid"
send "$spy" "mint(address,uint256)" "$me" "$one"
send "$spy" "approve(address,uint256)" "$accounts" "$one"
send "$accounts" "deposit(bytes32,address,uint256,address)" "$symbol" "$spy" "$one" "$me"
send "$accounts" "lend(bytes32,address,uint256,address)" "$symbol" "$spy" "$one" "$me"
send "$(jq -r .tokens.USDG "$registry")" "mint(address,uint256)" "$me" 9000000000
send "$nvda" "mint(address,uint256)" "$me" "$((6 * one))"
send "$spy" "mint(address,uint256)" "$me" "$((3 * one))"
echo "the SDKs' examples short SPY at its band's centre, $mid, against $(jq -r .bandFeeds.SPY "$registry")"
