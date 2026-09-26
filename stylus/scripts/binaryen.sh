#!/usr/bin/env bash
# Usage: binaryen.sh VERSION DIR
# Installs Binaryen VERSION into DIR from its GitHub release, checked against the SHA-256 pinned below,
# unless DIR/bin/wasm-opt already reports VERSION. cargo-stylus runs that wasm-opt for every program whose
# Stylus.toml has a [wasm-opt] table, and refuses any other version.
set -euo pipefail

version=$1 dir=$2
if "$dir/bin/wasm-opt" --version 2>/dev/null | grep -q "^wasm-opt version $version "; then
  exit 0
fi

case "$version-$(uname -s)-$(uname -m)" in
  133-Darwin-arm64) platform=arm64-macos sha=ad66da82ac13f163e424b1643f16c6dfcccc98b5966296b43e52d3cab04f84a8 ;;
  133-Darwin-x86_64) platform=x86_64-macos sha=13a9b90be775c6389ce3d1f879cb8627bea56708ba8c122983941d53a8199b95 ;;
  133-Linux-x86_64) platform=x86_64-linux sha=2dc9c7813f5375db93d96ead4b78222fcc3e2677bbb832297af4797782a37489 ;;
  133-Linux-aarch64) platform=aarch64-linux sha=89c07ea56faf38d0fbecf36ca8ec0721756716185f265b568e133d427f299bf8 ;;
  *) echo "No pinned Binaryen $version release for $(uname -s) $(uname -m)." >&2; exit 1 ;;
esac

tarball=binaryen-version_$version-$platform.tar.gz
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
curl -fsSL --proto '=https' --tlsv1.2 -o "$tmp/$tarball" \
  "https://github.com/WebAssembly/binaryen/releases/download/version_$version/$tarball"
echo "$sha  $tmp/$tarball" | shasum -a 256 -c - > /dev/null
tar -xzf "$tmp/$tarball" -C "$tmp"
mkdir -p "$(dirname "$dir")"
rm -rf "$dir"
mv "$tmp/binaryen-version_$version" "$dir"
"$dir/bin/wasm-opt" --version
