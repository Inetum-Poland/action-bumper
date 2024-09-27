#!/usr/bin/env bash
#
# Copyright (c) 2024 Inetum Poland.

source "${SCRIPT_FOLDER}/lib/helpers.sh"

__get_head_label() {
  jq -r '.pull_request.head.label' < "${GITHUB_EVENT_PATH}"
}

make_pr_status() {
  head_label="$(__get_head_label)"
  COMPARE=""
  ADDITIONAL_INFO=""

  if [[ -n "${BUMPER_CURRENT_VERSION:-}" ]]; then
    COMPARE="**Changes**: [${BUMPER_CURRENT_VERSION:-}...${head_label}](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/compare/${BUMPER_CURRENT_VERSION:-}...${head_label})"
  fi

  if [[ -n "${INPUT_BUMP_SEMVER}" && "${INPUT_BUMP_SEMVER}" == "true" ]]; then
    ADDITIONAL_INFO="${ADDITIONAL_INFO} / ${MINOR} / ${MAJOR}"
  fi

  if [[ -n "${INPUT_BUMP_LATEST}" && "${INPUT_BUMP_LATEST}" == "true" ]]; then
    ADDITIONAL_INFO="${ADDITIONAL_INFO} / latest"
  fi

  __append_github_output "tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ ${ACTION}<br>**Next version**: ${BUMPER_NEXT_VERSION}${ADDITIONAL_INFO:-}<br>${COMPARE}"
}

make_push_status() {
  COMPARE=""
  ADDITIONAL_INFO=""

  if [[ -n "${BUMPER_CURRENT_VERSION}" ]]; then
    COMPARE="**Changes**: [${BUMPER_CURRENT_VERSION}...${BUMPER_NEXT_VERSION}](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/compare/${BUMPER_CURRENT_VERSION}...${BUMPER_NEXT_VERSION})"
  fi

  if [[ -n "${INPUT_BUMP_SEMVER}" && "${INPUT_BUMP_SEMVER}" == "true" ]]; then
    ADDITIONAL_INFO="${ADDITIONAL_INFO} / ${MINOR} / ${MAJOR}"
  fi

  if [[ -n "${INPUT_BUMP_LATEST}" && "${INPUT_BUMP_LATEST}" == "true" ]]; then
    ADDITIONAL_INFO="${ADDITIONAL_INFO} / latest"
  fi

  __append_github_output "tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID})<br>**New version**: [${BUMPER_NEXT_VERSION}${SEMVER:-}](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/releases/tag/${BUMPER_NEXT_VERSION})${ADDITIONAL_INFO:-}<br>${COMPARE}"
}
