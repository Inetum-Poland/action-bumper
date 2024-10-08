#!/usr/bin/env bash
#
# Copyright (c) 2024 Inetum Poland.

init_debug() {
  if [[ -n "${DEBUG_GITHUB_EVENT_PATH:-}" ]]; then
    # shellcheck disable=SC2034
    GITHUB_EVENT_PATH="${DEBUG_GITHUB_EVENT_PATH}/data.json"
    # shellcheck disable=SC1091
    source "${DEBUG_GITHUB_EVENT_PATH}/.input.env"
  fi
}

exec_debug() {
  # -n; True if the length of string is non-zero.
  if [[ (-n "${INETUM_POLAND_ACTION_BUMPER_DEBUG:-}" && "${INETUM_POLAND_ACTION_BUMPER_DEBUG}" == "true") ]]; then
    echo "> ${1}" 2>&1;
  else
    bash -c "${1}"
  fi
}
