#!/usr/bin/env bash

setup_labels_from_pr_event() {
  jq -r '.pull_request.labels[].name' < "${1}" | tr '\n' ' '
}

setup_pr_number_from_pr_event() {
  jq -r '.pull_request.number' < "${1}"
}

setup_pr_title_from_pr_event() {
  jq -r '.pull_request.title' < "${1}"
}
