#!/usr/bin/env bash

# KCOV_EXCL_START
__get_git_tag_from_api() {
  curl --fail -s -H "Authorization: token ${INPUT_GITHUB_TOKEN}" "${1}" | jq -r '[.[] | select(.name != "latest")] | .[0].name'
}
# KCOV_EXCL_STOP

__get_action() {
  jq -r '.action' < "${GITHUB_EVENT_PATH}"
}

__get_git_tag_from_file() {
  jq -r '[.[] | select(.name != "latest")] | .[0].name' < "${DEBUG_GITHUB_EVENT_PATH}/tags.json"
}

__append_github_output() {
  echo "$1" >> "$GITHUB_OUTPUT"
}

__setup_bump_level() {
  BUMPER_BUMP_LEVEL="${INPUT_DEFAULT_BUMP_LEVEL}"

  if echo "${BUMPER_LABELS}" | grep -q "${INPUT_BUMP_NONE}" ; then
    BUMPER_BUMP_LEVEL="none"
  fi

  if echo "${BUMPER_LABELS}" | grep -q "${INPUT_BUMP_PATCH}" ; then
    BUMPER_BUMP_LEVEL="patch"
  fi

  if echo "${BUMPER_LABELS}" | grep -q "${INPUT_BUMP_MINOR}" ; then
    BUMPER_BUMP_LEVEL="minor"
  fi

  if echo "${BUMPER_LABELS}" | grep -q "${INPUT_BUMP_MAJOR}" ; then
    BUMPER_BUMP_LEVEL="major"
  fi
}

__setup_vars_from_pr_event() {
  PR_NUMBER=$(setup_pr_number_from_pr_event)
  PR_TITLE=$(setup_pr_title_from_pr_event)
  BUMPER_LABELS=$(setup_labels_from_pr_event)
}

__setup_vars_from_push_event() {
  PR_NUMBER=$(setup_pr_number_from_push_event)
  PR_TITLE=$(setup_pr_title_from_push_event)
  BUMPER_LABELS=$(setup_labels_from_push_event)
}

setup_vars() {
  ACTION=$(__get_action)

  if [[ "${ACTION}" =~ ^(labeled|unlabeled|synchronize|opened|reopened)$ ]]; then
    __setup_vars_from_pr_event
  else
    __setup_vars_from_push_event
  fi

  __setup_bump_level
}

setup_git_tag() {
  if [[ -n "${DEBUG_GITHUB_EVENT_PATH:-}" ]]; then
    # shellcheck disable=SC2002
    BUMPER_CURRENT_VERSION=$(__get_git_tag_from_file)
  else
    endpoint="${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/tags?sort=updated&direction=desc"
    BUMPER_CURRENT_VERSION=$(__get_git_tag_from_api "${endpoint}")
  fi

  if [[ "${BUMPER_CURRENT_VERSION}" == "null" ]]; then
    unset BUMPER_CURRENT_VERSION
  else
    __append_github_output "current_version=${BUMPER_CURRENT_VERSION}"
  fi
}

bump_tag() {
  if [[ -z "${BUMPER_BUMP_LEVEL}" ]]; then
    if [[ "${INPUT_FAIL_IF_NO_BUMP}" == "true" ]]; then
      echo "::error ::Job failed as no bump label is found."
      exit 1
    else
      echo "::notice ::Job skiped as no bump label is found. Do nothing."

      __append_github_output "skip=true"

      exit 0
    fi
  elif [[ "${BUMPER_BUMP_LEVEL}" == "none" ]]; then
    echo "::notice ::Job skiped as bump level is 'none'. Do nothing."

    __append_github_output "skip=true"

    exit 0
  else
    if [[ -z "${BUMPER_CURRENT_VERSION:-}" || "${BUMPER_CURRENT_VERSION:-}" == "null" ]]; then
      check_missing_tags
    else
      BUMPER_NEXT_VERSION="v$(semver bump "${BUMPER_BUMP_LEVEL}" "${BUMPER_CURRENT_VERSION}")" || exit 1
    fi

    __append_github_output "next_version=${BUMPER_NEXT_VERSION}"

    BUMPER_TAG_MESSAGE="${BUMPER_NEXT_VERSION}: PR #${PR_NUMBER} - ${PR_TITLE}"

    __append_github_output "message=${BUMPER_TAG_MESSAGE}"
  fi
}

check_missing_tags() {
  case "${BUMPER_BUMP_LEVEL}" in
    major)
      BUMPER_NEXT_VERSION="v1.0.0"
      ;;
    minor)
      BUMPER_NEXT_VERSION="v0.1.0"
      ;;
    patch)
      BUMPER_NEXT_VERSION="v0.0.1"
      ;;
    *)
      BUMPER_NEXT_VERSION="v0.0.0"
      ;;
  esac
}

remove_v_prefix() {
  if [[ "${INPUT_INCLUDE_V}" == "false" ]]; then
    BUMPER_NEXT_VERSION="${BUMPER_NEXT_VERSION//v}"
  fi
}

setup_git_config() {
  execute_or_debug "git config user.name \"${INPUT_TAG_AS_USER:-${GITHUB_ACTOR}}\""
  execute_or_debug "git config user.email \"${INPUT_TAG_AS_EMAIL:-${GITHUB_ACTOR}@users.noreply.github.com}\""

  if [[ -n "${INPUT_GITHUB_TOKEN}" ]]; then
    execute_or_debug "git remote set-url origin \"https://${GITHUB_ACTOR}:${INPUT_GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git\""
  fi
}

# shellcheck disable=SC2034
bump_semver_tags() {
  PATCH="${BUMPER_CURRENT_VERSION}" # v1.2.3
  MINOR="${PATCH%.*}"               # v1.2
  MAJOR="${MINOR%.*}"               # v1
}
