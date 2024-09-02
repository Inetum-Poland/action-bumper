#!/usr/bin/env bash

source lib/helpers.sh

__get_head_label() {
  jq -r '.pull_request.head.label' < "${GITHUB_EVENT_PATH}"
}

make_pr_status() {
  head_label="$(__get_head_label)"
  compare=""

  if [[ -n "${BUMPER_CURRENT_VERSION:-}" ]]; then
    compare="**Changes**: [${BUMPER_CURRENT_VERSION}...${head_label}](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/compare/${BUMPER_CURRENT_VERSION}...${head_label})"
  fi

  if [[ -n "${INPUT_ADD_LATEST}" && "${INPUT_ADD_LATEST}" == "true" ]]; then
    LATEST=" / latest"
    true
  fi

  __append_github_output "tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ ${ACTION}<br>**Next version**: ${BUMPER_NEXT_VERSION}${LATEST:-}<br>${compare}"
}

make_push_status() {
  compare=""

  if [[ -n "${BUMPER_CURRENT_VERSION}" ]]; then
    compare="**Changes**: [${BUMPER_CURRENT_VERSION}...${BUMPER_NEXT_VERSION}](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/compare/${BUMPER_CURRENT_VERSION}...${BUMPER_NEXT_VERSION})"
  fi

  if [[ -n "${INPUT_ADD_LATEST}" && "${INPUT_ADD_LATEST}" == "true" ]]; then
    LATEST=" / latest"
  fi

  __append_github_output "tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID})<br>**New version**: [${BUMPER_NEXT_VERSION}${LATEST:-}](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/releases/tag/${BUMPER_NEXT_VERSION})<br>${compare}"
}

make_merge_semver_status(){
  __append_github_output "tag_status=New patch: ${PATCH}<br>New minor: ${MINOR}<br>New major: ${MAJOR}"
}
