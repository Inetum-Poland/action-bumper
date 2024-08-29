#!/usr/bin/env bash

setup_labels_from_pr_event() {
  jq -r '.pull_request.labels[].name' < "${GITHUB_EVENT_PATH}" | tr '\n' ' '
}

setup_pr_number_from_pr_event() {
  jq -r '.pull_request.number' < "${GITHUB_EVENT_PATH}"
}

setup_pr_title_from_pr_event() {
  jq -r '.pull_request.title' < "${GITHUB_EVENT_PATH}"
}
