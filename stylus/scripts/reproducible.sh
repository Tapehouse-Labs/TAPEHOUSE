#!/usr/bin/env bash
# Usage: reproducible.sh initcode
#        reproducible.sh verify RPC_URL DEPLOYMENT_TX CONTRACT
# Runs cargo-stylus in the image that `cargo stylus deploy` and `cargo stylus verify` build for
# reproducible builds, on the staged content of stylus/, mounted at /source as they mount the
# workspace. `initcode` writes stylus/target/reproducible/<contract>.initcode.hex for every program;
# `verify` succeeds only when cargo-stylus prints "Verification successful".
set -euo pipefail

case ${1:-} in
  initcode) [ $# -eq 1 ] ;;
  verify) [ $# -eq 4 ] ;;
  *) false ;;
esac || { echo "Usage: reproducible.sh initcode | verify RPC_URL DEPLOYMENT_TX CONTRACT" >&2; exit 2; }

stylus=$(cd "$(dirname "$0")/.." && pwd)
version=${CARGO_STYLUS_VERSION:?Set CARGO_STYLUS_VERSION, as the Makefile does}
toolchain=$(sed -n 's/^channel = "\(.*\)"$/\1/p' "$stylus/rust-toolchain.toml")
image="cargo-stylus-base-$version-toolchain-$toolchain"

docker image inspect "$image" > /dev/null 2>&1 || docker build --tag "$image" - <<DOCKERFILE
ARG BUILD_PLATFORM=linux/amd64
FROM --platform=\${BUILD_PLATFORM} offchainlabs/cargo-stylus-base:$version AS base
RUN rustup toolchain install $toolchain-x86_64-unknown-linux-gnu
RUN rustup default $toolchain-x86_64-unknown-linux-gnu
RUN rustup target add wasm32-unknown-unknown
RUN rustup component add rust-src --toolchain $toolchain-x86_64-unknown-linux-gnu
DOCKERFILE

source=$(mktemp -d)
run() { docker run --rm --network host --workdir /source --volume "$source:/source" "$image" "$@"; }
trap 'run rm -rf /source/target > /dev/null 2>&1 || true; rm -rf "$source"' EXIT
git -C "$stylus/.." archive "$(git -C "$stylus/.." write-tree)" stylus | tar -xf - -C "$source" --strip-components 1

if [ "$1" = initcode ]; then
  mkdir -p "$stylus/target/reproducible"
  for dir in "$source"/contracts/*/; do
    contract=$(basename "$dir")
    run sh -c "mkdir -p target && cargo stylus get-initcode --contract $contract --output target/$contract.initcode.hex"
    tr -d '\n' < "$source/target/$contract.initcode.hex" > "$stylus/target/reproducible/$contract.initcode.hex"
    hash=$(cast keccak "0x$(cat "$stylus/target/reproducible/$contract.initcode.hex")")
    echo "$contract: $hash"
  done
else
  run cargo stylus verify --no-verify --contract "$4" -e "$2" --deployment-tx "$3" | tee "$source/verify.log"
  grep -q '^Verification successful$' "$source/verify.log"
fi
