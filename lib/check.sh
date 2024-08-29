#!/usr/bin/env bash

semver_check() {
  if ! semver -v &> /dev/null; then
    echo "::error ::semver is not installed."
    exit 1
  fi
}

jq_check() {
  if ! jq -V &> /dev/null; then
    echo "::error ::jq is not installed."
    exit 1
  fi
}
