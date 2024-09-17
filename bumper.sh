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

# https://stackoverflow.com/questions/59895/how-do-i-get-the-directory-where-a-bash-script-is-located-from-within-the-script
SCRIPT_FOLDER=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )

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

  if [[ "${ACTION}" =~ ^(labeled|unlabeled|synchronize|opened|reopened)$ ]]; then
    bump_tag
    remove_v_prefix
    bump_semver_tags
    make_pr_status
  else
    bump_tag
    remove_v_prefix
    bump_semver_tags
    setup_git_config
    make_and_push_tag
    make_push_status
  fi
}

action_bumper
