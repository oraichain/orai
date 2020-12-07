#!/bin/bash
set -o errexit -o nounset -o pipefail
command -v shellcheck > /dev/null && shellcheck "$0"

# Iterates over all example projects and checks them. This updates lockfiles
# and provides a quick sanity check.
# This script is for development purposes only. In the CI, each example project
# is configured separately.

for example in ./*; do
  if [[ -d "$example" ]]; then
    echo "Checking $example ..."

    (
        cd "$example"
        cargo fmt
        cargo check --tests
    )
  fi
done
