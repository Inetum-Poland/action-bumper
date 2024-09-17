#!/usr/bin/env bash

# https://github.com/fsaintjacques/semver-tool/tree/master

set -eu

# KCOV_EXCL_START
# -n; True if the length of string is non-zero.
if [[ (-n "${INETUM_POLAND_ACTION_BUMPER_TRACE:-}" && "${INETUM_POLAND_ACTION_BUMPER_TRACE}" == "true") && "${SHELLSPEC:-}" != "true" ]]; then
  set -x
  export PS4='+(${BASH_SOURCE}:${LINENO}): ${FUNCNAME[0]:+${FUNCNAME[0]}(): }'
fi
# KCOV_EXCL_STOP

if [[ "${SHELLSPEC:-}" != "true" ]]; then
  SCRIPT_FOLDER="$(dirname "$(readlink -f "$0")")"
else
  SCRIPT_FOLDER="."
fi

source "${SCRIPT_FOLDER}/lib/debug.sh"
source "${SCRIPT_FOLDER}/lib/helpers.sh"
source "${SCRIPT_FOLDER}/lib/message.sh"
source "${SCRIPT_FOLDER}/lib/pr_event.sh"
source "${SCRIPT_FOLDER}/lib/push_event.sh"
source "${SCRIPT_FOLDER}/lib/check.sh"
source "${SCRIPT_FOLDER}/lib/git.sh"

# KCOV_EXCL_START
if [[ -n "${GITHUB_WORKSPACE:-}" ]]; then
  git config --global --add safe.directory "${GITHUB_WORKSPACE}" || exit
  cd "${GITHUB_WORKSPACE}" || exit
fi
# KCOV_EXCL_STOP

action_bumper() {
  ACTION=
  PR_NUMBER=
  PR_TITLE=
  BUMPER_LABELS=
  BUMPER_CURRENT_VERSION=
  BUMPER_BUMP_LEVEL=
  BUMPER_NEXT_VERSION=
  BUMPER_TAG_MESSAGE=
  PATCH=
  MINOR=
  MAJOR=

  init_debug

  jq_check
  semver_check
  curl_check
  git_check

  setup_git_tag
  setup_vars

  if [[ $(jq -r '.ref' < "${GITHUB_EVENT_PATH}") =~ "refs/tags/" && ${INPUT_BUMP_SEMVER} == "true" ]]; then
    bump_semver_tags
    remove_v_prefix
    setup_git_config
    make_and_push_semver_tags
    make_merge_semver_status
  elif [[ "${ACTION}" =~ ^(labeled|unlabeled|synchronize|opened|reopened)$ ]]; then
    bump_tag
    remove_v_prefix
    make_pr_status
  else
    bump_tag
    remove_v_prefix
    setup_git_config
    make_and_push_tag
    make_push_status
  fi
}

action_bumper
