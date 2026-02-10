#!/bin/bash
# Copyright (c) 2024 Inetum Poland.
# Entrypoint script that selects between Bash and Go implementations

set -eo pipefail

# Check if Go implementation is requested
if [[ "${BUMPER_USE_GO:-false}" == "true" ]]; then
  echo "Using Go implementation"
  exec /opt/bumper-go "$@"
else
  # Default to Bash implementation
  exec /opt/bumper.sh "$@"
fi
