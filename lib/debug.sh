#!/usr/bin/env bash

init_debug() {
  if [[ -n "${ACTIONS_STEP_DEBUG:-}" && "${ACTIONS_STEP_DEBUG}" == "true" && -z "${SHELLSPEC:-}" ]]; then
    # KCOV_EXCL_START
    set -x
    export PS4='+(${BASH_SOURCE}:${LINENO}): ${FUNCNAME[0]:+${FUNCNAME[0]}(): }'
    # KCOV_EXCL_STOP
  fi

  if [[ -n "${DEBUG_GITHUB_EVENT_PATH:-}" ]]; then
    # shellcheck disable=SC2034
    GITHUB_EVENT_PATH="${DEBUG_GITHUB_EVENT_PATH}/data.json"
    # shellcheck disable=SC1091
    source "${DEBUG_GITHUB_EVENT_PATH}/.input.env"
  fi
}

exec_debug() {
  if [[ -n "${ACTIONS_STEP_DEBUG}" && "${ACTIONS_STEP_DEBUG}" == "true" ]]; then
    echo "> ${1}" 2>&1;
  else
    bash -c "${1}"
  fi
}
