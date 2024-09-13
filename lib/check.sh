#!/usr/bin/env bash

tool_check() {
  if ! $2 &> /dev/null; then
    echo "::error ::$1 is not installed."
    exit 1
  fi
}

semver_check() {
  tool_check "semver" "semver -v"
}

jq_check() {
  tool_check "jq" "jq -V"
}

curl_check() {
  tool_check "curl" "curl --version"
}

git_check() {
  tool_check "git" "git --version"
}
