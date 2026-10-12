#!/usr/bin/env bash
# SPDX-License-Identifier: MIT OR Apache-2.0
# Usage: clone-stylus.sh RPC_URL CONTRACT SOURCE --account NAME --password-file FILE -- CONSTRUCTOR_ARGS...
# Deploys a new instance of a Stylus program that is already active, with its own constructor arguments, and
# prints the deployment transaction and the instance's address. SOURCE is a deployed instance of CONTRACT, `band`
# or `margin`. An activation belongs to the program's code hash, so an instance with the same root code, which
# points at the same code fragments, runs without being activated again and works while a chain pauses new
# activations. The instance runs the code SOURCE runs: a change to the program needs `deploy-stylus`.
# The instance is created and constructed through StylusDeployer, at its canonical address, and the deployment
# is simulated first. CONSTRUCTOR_ARGS are those band-args.sh and margin-args.sh print.
set -euo pipefail

usage="Usage: clone-stylus.sh RPC_URL CONTRACT SOURCE --account NAME --password-file FILE -- CONSTRUCTOR_ARGS..."
[ $# -ge 8 ] && [ "$4" = --account ] && [ "$6" = --password-file ] && [ "$8" = -- ] || { echo "$usage" >&2; exit 2; }
rpc=$1 contract=$2 source=$3 account=$5 password=$7
shift 8

stylus_deployer=0xcEcba2F1DC234f70Dd89F2041029807F8D03A990
asset='(bytes32,uint32,uint32,uint32,uint32,uint32,uint32,address,address)'
case $contract in
  band) types='bytes32[],address[],bytes32[],bytes32[],address[],address,address,address,bool' ;;
  margin) types="${asset}[],uint16[],uint16[],bytes32,address,address,address,address,address" ;;
  *) echo "No program $contract to clone. $usage" >&2; exit 2 ;;
esac

root=$(cast code "$source" --rpc-url "$rpc")
[[ $root =~ ^0xeff00[0-9a-f]+$ ]] || { echo "$source is not a Stylus program" >&2; exit 1; }
length=$(printf '%04x' $(((${#root} - 2) / 2)))
initcode="0x61${length}80600c6000396000f3${root#0x}"
constructor="0x5585258d$(cast abi-encode "f($types)" "$@" | cut -c3-)"
salt=$(cast keccak "$contract-$(date +%s)-$RANDOM")
sender=$(cast wallet address --account "$account" --password-file "$password")
deploy="deploy(bytes,bytes,uint256,bytes32)(address)"

address=$(cast call "$stylus_deployer" "$deploy" "$initcode" "$constructor" 0 "$salt" --from "$sender" --rpc-url "$rpc")
tx=$(cast send "$stylus_deployer" "$deploy" "$initcode" "$constructor" 0 "$salt" --account "$account" \
  --password-file "$password" --rpc-url "$rpc" --json | jq -er 'select(.status == "0x1") | .transactionHash')
[ "$(cast codehash "$address" --rpc-url "$rpc")" = "$(cast codehash "$source" --rpc-url "$rpc")" ] ||
  { echo "$address does not have the code of $source" >&2; exit 1; }
echo "$tx $(cast to-check-sum-address "$address")"
