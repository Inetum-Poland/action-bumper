#!/usr/bin/env bash

set -e

for dir in test/*; do
  echo "${dir}:"
  ASSERT=
  DEBUG_GITHUB_EVENT_PATH=${dir} ./entrypoint.sh && ASSERT="OK" || ASSERT="NOK"

  echo -e "$(cat "${dir}"/assert.txt) == ${ASSERT}\n"
  [ "${ASSERT}" == "$(cat "${dir}"/assert.txt)" ]
done
