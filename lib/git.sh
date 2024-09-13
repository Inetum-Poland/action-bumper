#!/usr/bin/env bash

make_and_push_tag() {
  execute_or_debug "git tag -a \"${BUMPER_NEXT_VERSION}\" -m \"${BUMPER_TAG_MESSAGE}\""
  execute_or_debug "git push origin \"${BUMPER_NEXT_VERSION}\""

  if [[ -n "${INPUT_ADD_LATEST}" && "${INPUT_ADD_LATEST}" == "true" ]]; then
    execute_or_debug "git tag -fa latest \"${BUMPER_NEXT_VERSION}^{commit}\" -m \"${BUMPER_TAG_MESSAGE}\""
    execute_or_debug "git push --force origin latest"
  fi
}

make_and_push_semver_tags() {
  execute_or_debug "git tag -fa \"${MINOR}\" \"${BUMPER_NEXT_VERSION}^{commit}\" -m \"${BUMPER_TAG_MESSAGE}\""
  execute_or_debug "git tag -fa \"${MAJOR}\" \"${BUMPER_NEXT_VERSION}^{commit}\" -m \"${BUMPER_TAG_MESSAGE}\""

  execute_or_debug "git push --force origin \"${MINOR}\""
  execute_or_debug "git push --force origin \"${MAJOR}\""
}
