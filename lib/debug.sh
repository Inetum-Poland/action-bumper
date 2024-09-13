#!/usr/bin/env bash

init_debug() {
  if [[ (-n "${DEBUG:-}" && "${DEBUG}" == "true") || (-n "${ACTIONS_STEP_DEBUG:-}" && "${ACTIONS_STEP_DEBUG}" == "true") ]]; then
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

execute_or_debug() {
  if [[ "${INPUT_DRY_RUN:-}" == "true" ]]; then
    echo "> ${1}" 1>&2;
  else
    ${1}
  fi
}
