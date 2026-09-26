#!/usr/bin/env bash
# Usage: deploy-stylus-deployer.sh RPC_URL PRIVATE_KEY
# Deploys the deterministic-deployment proxy and StylusDeployer at their canonical addresses on a
# Nitro dev node whose chain owner is PRIVATE_KEY. Skips whatever is already deployed.
set -euo pipefail

rpc=$1 key=$2
create2=0x4e59b44847b379578588920cA78FbF26c0B4956C
create2_signer=0x3fab184622dc19b6109349b94811493bf2a45362
create2_tx=0xf8a58085174876e800830186a08080b853604580600e600039806000f350fe7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe03601600081602082378035828234f58015156039578182fd5b8082525050506014600cf31ba02222222222222222222222222222222222222222222222222222222222222222a02222222222222222222222222222222222222222222222222222222222222222
stylus_deployer=0xcEcba2F1DC234f70Dd89F2041029807F8D03A990
stylus_deployer_input=$(cat "$(dirname "$0")/stylus-deployer.hex")
arb_owner=0x0000000000000000000000000000000000000070
arb_gas_info=0x000000000000000000000000000000000000006C

if [ "$(cast codesize --rpc-url "$rpc" $create2)" = 0 ]; then
  l1_price=$(cast call --rpc-url "$rpc" $arb_gas_info "getL1BaseFeeEstimate()(uint256)" | cut -d' ' -f1)
  trap 'cast send --rpc-url "$rpc" --private-key "$key" $arb_owner "setL1PricePerUnit(uint256)" "$l1_price" >/dev/null' EXIT
  cast send --rpc-url "$rpc" --private-key "$key" $arb_owner "setL1PricePerUnit(uint256)" 0 >/dev/null
  cast send --rpc-url "$rpc" --private-key "$key" --value 0.01ether $create2_signer >/dev/null
  cast publish --rpc-url "$rpc" $create2_tx >/dev/null
fi

if [ "$(cast codesize --rpc-url "$rpc" $stylus_deployer)" = 0 ]; then
  [ "$(cast keccak "$stylus_deployer_input")" = 0x958e594c8d313d3b04bca21f39e9d334c2f8821ba9a5af9faafbfe6087a2558f ] ||
    { echo "stylus-deployer.hex is not the canonical salt and initcode"; exit 1; }
  cast send --rpc-url "$rpc" --private-key "$key" $create2 "$stylus_deployer_input" >/dev/null
fi

[ "$(cast keccak "$(cast code --rpc-url "$rpc" $stylus_deployer)")" = 0x85c3998f7541c47e69d221966c2725e197d243190f367eefab6eda7c90c3994a ] ||
  { echo "No canonical StylusDeployer at $stylus_deployer"; exit 1; }
echo "StylusDeployer ready at $stylus_deployer."
