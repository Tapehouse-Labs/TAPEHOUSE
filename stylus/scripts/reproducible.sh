#!/usr/bin/env bash
# Usage: reproducible.sh deploy RPC_URL CONTRACT SIGNER -- CONSTRUCTOR_ARGS...
#        reproducible.sh verify RPC_URL DEPLOYMENT_TX CONTRACT
# Runs cargo-stylus in the image that `cargo stylus deploy` and `cargo stylus verify` build for
# reproducible builds, on the staged content of stylus/, mounted at /source as they mount the
# workspace. `deploy` deploys, activates and constructs the program through StylusDeployer, in as
# many code fragments as it needs, and prints the deployment transaction and the program's address.
# SIGNER is `--private-key-path FILE`, or `--account NAME --password-file FILE` for a Foundry keystore;
# the files are mounted read-only.
# `verify` succeeds only when cargo-stylus prints "Verification successful".
set -euo pipefail

usage="Usage: reproducible.sh deploy RPC_URL CONTRACT SIGNER -- CONSTRUCTOR_ARGS... | verify RPC_URL DEPLOYMENT_TX CONTRACT"
case ${1:-} in
  deploy) [ $# -ge 5 ] ;;
  verify) [ $# -eq 4 ] ;;
  *) false ;;
esac || { echo "$usage" >&2; exit 2; }

stylus=$(cd "$(dirname "$0")/.." && pwd)
version=${CARGO_STYLUS_VERSION:?Set CARGO_STYLUS_VERSION, as the Makefile does}
toolchain=$(sed -n 's/^channel = "\(.*\)"$/\1/p' "$stylus/rust-toolchain.toml")
binaryen=$(awk '/^\[/ { table = $0 } table == "[wasm-opt]" && $1 == "version" { gsub(/"/, "", $3); print $3 }' "$stylus/Stylus.toml")
image="cargo-stylus-base-$version-toolchain-$toolchain${binaryen:+-binaryen-$binaryen}"
layer=
if [ -n "$binaryen" ]; then
  tarball=binaryen-version_$binaryen-x86_64-linux.tar.gz
  release=https://github.com/WebAssembly/binaryen/releases/download/version_$binaryen
  sha=$(sed -n "s/^ *$binaryen-Linux-x86_64) platform=x86_64-linux sha=\([0-9a-f]*\) .*/\1/p" "$stylus/scripts/binaryen.sh")
  [ -n "$sha" ] || { echo "No pinned SHA-256 for Binaryen $binaryen in binaryen.sh" >&2; exit 1; }
  layer="RUN cd /tmp && curl -fsSL --proto '=https' --tlsv1.2 -O $release/$tarball \\
  && curl -fsSL --proto '=https' --tlsv1.2 -O $release/$tarball.sha256 && sha256sum -c $tarball.sha256 \\
  && echo \"$sha  $tarball\" | sha256sum -c && tar -xzf $tarball -C /opt && rm $tarball $tarball.sha256
ENV PATH=\"/opt/binaryen-version_$binaryen/bin:\${PATH}\""
fi

docker image inspect "$image" > /dev/null 2>&1 || docker build --tag "$image" - <<DOCKERFILE
ARG BUILD_PLATFORM=linux/amd64
FROM --platform=\${BUILD_PLATFORM} offchainlabs/cargo-stylus-base:$version AS base
RUN rustup toolchain install $toolchain-x86_64-unknown-linux-gnu
RUN rustup default $toolchain-x86_64-unknown-linux-gnu
RUN rustup target add wasm32-unknown-unknown
RUN rustup component add rust-src --toolchain $toolchain-x86_64-unknown-linux-gnu
$layer
DOCKERFILE

signer=() keys=()
if [ "$1" = deploy ]; then
  case $4 in
    --private-key-path)
      keys=(--volume "$(cd "$(dirname "$5")" && pwd)/$(basename "$5"):/keys/key:ro")
      signer=(--private-key-path /keys/key) next=6 ;;
    --account)
      [ "${6:-}" = --password-file ] || { echo "$usage" >&2; exit 2; }
      keys=(--volume "$(cd "${FOUNDRY_KEYSTORES:-$HOME/.foundry/keystores}" && pwd)/$5:/keys/keystore:ro" --volume "$(cd "$(dirname "$7")" && pwd)/$(basename "$7"):/keys/password:ro")
      signer=(--keystore-path /keys/keystore --keystore-password-path /keys/password) next=8 ;;
    *) echo "$usage" >&2; exit 2 ;;
  esac
  [ "${!next:-}" = -- ] || { echo "$usage" >&2; exit 2; }
fi

source=$(mktemp -d)
run() {
  docker run --rm --network host --workdir /source --volume "$source:/source" ${keys[@]+"${keys[@]}"} "$image" "$@"
}
trap 'run rm -rf /source/target > /dev/null 2>&1 || true; rm -rf "$source"' EXIT
git -C "$stylus/.." archive "$(git -C "$stylus/.." write-tree)" stylus | tar -xf - -C "$source" --strip-components 1

if [ "$1" = deploy ]; then
  rpc=$2 contract=$3
  shift "$next"
  run cargo stylus deploy --no-verify --contract "$contract" -e "$rpc" "${signer[@]}" --constructor-args "$@" |
    perl -pe 's/\e\[[0-9;]*m//g' | tee "$source/deploy.log" >&2
  tx=$(grep -o 'deployment tx hash: 0x[0-9a-f]\{64\}' "$source/deploy.log" | grep -o '0x.*')
  address=$(grep -o 'deployed code at address: 0x[0-9a-f]\{40\}' "$source/deploy.log" | grep -o '0x.*')
  address=$(cast to-check-sum-address "$address")
  echo "$tx $address"
else
  run cargo stylus verify --no-verify --contract "$4" -e "$2" --deployment-tx "$3" | tee "$source/verify.log"
  grep -q '^Verification successful$' "$source/verify.log"
fi
