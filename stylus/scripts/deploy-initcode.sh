#!/usr/bin/env bash
# Usage: deploy-initcode.sh RPC_URL INITCODE_FILE INIT_DATA SIGNER_FLAGS...
# Deploys, activates and constructs a Stylus program in one transaction through the canonical
# StylusDeployer. Sends twice the activation data fee (StylusDeployer refunds the excess), or nothing
# when the code is already activated. Prints the transaction hash and the program address from the
# ContractDeployed log.
set -euo pipefail

rpc=$1 initcode=0x$(tr -d '\n' < "$2") init_data=$3
shift 3
stylus_deployer=0xcEcba2F1DC234f70Dd89F2041029807F8D03A990
arb_wasm=0x0000000000000000000000000000000000000071
contract_deployed=0x8ffcdc15a283d706d38281f500270d8b5a656918f555de0913d7455e3e6bc1bf
probe=0x000000000000000000000000000000000000bEEF
funded=0x000000000000000000000000000000000000cAFE
runtime=0x${initcode:88}

code=$(cast code --rpc-url "$rpc" $stylus_deployer) || { echo "Cannot read $stylus_deployer on $rpc" >&2; exit 1; }
[ "$(cast keccak "$code")" = 0x85c3998f7541c47e69d221966c2725e197d243190f367eefab6eda7c90c3994a ] ||
  { echo "No canonical StylusDeployer at $stylus_deployer on $rpc" >&2; exit 1; }

if cast call --rpc-url "$rpc" $arb_wasm "codehashVersion(bytes32)(uint16)" "$(cast keccak "$runtime")" > /dev/null 2>&1; then
  value=0
else
  fee=$(cast call --rpc-url "$rpc" $arb_wasm "activateProgram(address)(uint16,uint256)" $probe --value 1ether \
    --from $funded --override-balance $funded:10000000000000000000 --override-code "$probe:$runtime" | sed -n 2p | cut -d' ' -f1)
  value=$((fee * 2))
fi

receipt=$(cast send --rpc-url "$rpc" "$@" --value "$value" --json $stylus_deployer \
  "deploy(bytes,bytes,uint256,bytes32)" "$initcode" "$init_data" 0 0x0000000000000000000000000000000000000000000000000000000000000000) ||
  { echo "$receipt" >&2; exit 1; }
tx=$(jq -r .transactionHash <<<"$receipt")
[ "$(jq -r .status <<<"$receipt")" = 0x1 ] || { echo "deploy reverted: $tx" >&2; exit 1; }
data=$(jq -r --arg from $stylus_deployer --arg topic $contract_deployed \
  '.logs[] | select((.address | ascii_downcase) == ($from | ascii_downcase) and .topics[0] == $topic) | .data' <<<"$receipt")
address=$(cast to-check-sum-address "0x${data:26}")
echo "$tx $address"
