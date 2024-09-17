#!/usr/bin/env bash
#
# Copyright (c) 2024 Inetum Poland.

setup_git_config() {
  exec_debug "git config user.name \"${INPUT_BUMP_TAG_AS_USER:-${GITHUB_ACTOR}}\""
  exec_debug "git config user.email \"${INPUT_BUMP_TAG_AS_EMAIL:-${GITHUB_ACTOR}@users.noreply.github.com}\""

  if [[ -n "${INPUT_GITHUB_TOKEN}" ]]; then
    exec_debug "git remote set-url origin \"https://${GITHUB_ACTOR}:${INPUT_GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git\""
  fi
}

make_and_push_tag() {
  exec_debug "git tag -a \"${BUMPER_NEXT_VERSION}\" -m \"${BUMPER_TAG_MESSAGE}\""
  exec_debug "git push origin \"${BUMPER_NEXT_VERSION}\""

  if [[ -n "${INPUT_BUMP_LATEST}" && "${INPUT_BUMP_LATEST}" == "true" ]]; then
    exec_debug "git tag -fa latest \"${BUMPER_NEXT_VERSION}^{commit}\" -m \"${BUMPER_TAG_MESSAGE}\""
    exec_debug "git push --force origin latest"
  fi

  if [[ "${INPUT_BUMP_SEMVER}" == "true" ]]; then
    exec_debug "git tag -fa \"${MINOR}\" \"${BUMPER_NEXT_VERSION}^{commit}\" -m \"${BUMPER_TAG_MESSAGE}\""
    exec_debug "git tag -fa \"${MAJOR}\" \"${BUMPER_NEXT_VERSION}^{commit}\" -m \"${BUMPER_TAG_MESSAGE}\""

    exec_debug "git push --force origin \"${MINOR}\""
    exec_debug "git push --force origin \"${MAJOR}\""
  fi
}
