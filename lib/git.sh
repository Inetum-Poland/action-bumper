#!/usr/bin/env bash

make_and_push_tag() {
  exec_debug "git tag -a \"${BUMPER_NEXT_VERSION}\" -m \"${BUMPER_TAG_MESSAGE}\""
  exec_debug "git push origin \"${BUMPER_NEXT_VERSION}\""

  if [[ -n "${INPUT_ADD_LATEST}" && "${INPUT_ADD_LATEST}" == "true" ]]; then
    exec_debug "git tag -fa latest \"${BUMPER_NEXT_VERSION}^{commit}\" -m \"${BUMPER_TAG_MESSAGE}\""
    exec_debug "git push --force origin latest"
  fi
}

make_and_push_semver_tags() {
  exec_debug "git tag -fa \"${MINOR}\" \"${BUMPER_NEXT_VERSION}^{commit}\" -m \"${BUMPER_TAG_MESSAGE}\""
  exec_debug "git tag -fa \"${MAJOR}\" \"${BUMPER_NEXT_VERSION}^{commit}\" -m \"${BUMPER_TAG_MESSAGE}\""

  exec_debug "git push --force origin \"${MINOR}\""
  exec_debug "git push --force origin \"${MAJOR}\""
}
