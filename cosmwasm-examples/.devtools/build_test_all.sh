#!/bin/bash
set -o errexit -o nounset -o pipefail
command -v shellcheck > /dev/null && shellcheck "$0"

# Iterates over all example projects, builds and tests them.
# This script is for development purposes only. In the CI, each example project
# is configured separately.

export RUST_BACKTRACE=1

for example in ./*; do
  if [[ -d "$example" ]]; then
    echo "Building and testing $example ..."

    (
        cd "$example"
        cargo fmt
        cargo build --locked
        cargo unit-test --locked
        cargo wasm --locked
        cargo integration-test --locked
        cargo schema --locked

        if [[ -n "${REGEN_OPT:-}" ]]; then
          docker run --rm -v "$(pwd)":/code \
            --mount type=volume,source="$(basename "$(pwd)")_cache",target=/code/target \
            --mount type=volume,source=registry_cache,target=/usr/local/cargo/registry \
            cosmwasm/rust-optimizer:0.8.0
        fi
    )
  fi
done
