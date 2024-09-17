#!/usr/bin/env bash

__get_pulls() {
  if [[ -n "${INPUT_GITHUB_TOKEN}" ]]; then
    curl --fail -s -H "Authorization: token ${INPUT_GITHUB_TOKEN}" "${1}"
  else
    echo "INPUT_GITHUB_TOKEN is not available. Subscequent GitHub API call may fail due to API limit." >&2
    curl --fail -s "${1}"
  fi
}

list_pulls() {
  if [[ -n "${DEBUG_GITHUB_EVENT_PATH:-}" ]]; then
    jq -c -a "." "${DEBUG_GITHUB_EVENT_PATH}/pull_request.json"
  else
    __get_pulls "${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/pulls?state=closed&sort=updated&direction=desc"
  fi
}

setup_labels_from_push_event() {
  pull_request="$(echo -n "${1}" | jq ".[] | select(.merge_commit_sha==\"${GITHUB_SHA}\")")"
  echo "${pull_request}" | jq '.labels | .[].name'
}

setup_pr_number_from_push_event() {
  pull_request="$(echo -n "${1}" | jq ".[] | select(.merge_commit_sha==\"${GITHUB_SHA}\")")"
  echo "${pull_request}" | jq -r .number
}

setup_pr_title_from_push_event() {
  pull_request="$(echo -n "${1}" | jq ".[] | select(.merge_commit_sha==\"${GITHUB_SHA}\")")"
  echo "${pull_request}" | jq -r .title
}
