#!/usr/bin/env bash

# KCOV_EXCL_START
__get_pulls_with_token() {
  curl --fail -s -H "Authorization: token ${INPUT_GITHUB_TOKEN}" "${1}"
}

__get_pulls_without_token() {
  curl --fail -s "${1}"
}
# KCOV_EXCL_STOP

list_pulls() {
  if [[ -n "${DEBUG_GITHUB_EVENT_PATH}" ]]; then
    cat "${DEBUG_GITHUB_EVENT_PATH}/pull_request.json"
  else
    pulls_endpoint="${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/pulls?state=closed&sort=updated&direction=desc"
    if [[ -n "${INPUT_GITHUB_TOKEN}" ]]; then
      __get_pulls_with_token "${pulls_endpoint}"
    else
      echo "INPUT_GITHUB_TOKEN is not available. Subscequent GitHub API call may fail due to API limit." >&2
      __get_pulls_without_token "${pulls_endpoint}"
    fi
  fi
}

setup_labels_from_push_event() {
  pull_request="$(list_pulls | jq ".[] | select(.merge_commit_sha==\"${GITHUB_SHA}\")")"
  echo "${pull_request}" | jq '.labels | .[].name'
}

setup_pr_number_from_push_event() {
  pull_request="$(list_pulls | jq ".[] | select(.merge_commit_sha==\"${GITHUB_SHA}\")")"
  echo "${pull_request}" | jq -r .number
}

setup_pr_title_from_push_event() {
  pull_request="$(list_pulls | jq ".[] | select(.merge_commit_sha==\"${GITHUB_SHA}\")")"
  echo "${pull_request}" | jq -r .title
}
