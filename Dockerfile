# Copyright (c) 2024 Inetum Poland.
FROM ubuntu:25.10

# Workdir
WORKDIR /opt

# Add user
RUN groupadd -g 1001 runtime && \
  useradd --create-home --no-log-init -u 1001 -g 1001 runtime

# Install dependencies
RUN apt update && apt install -y \
  bash \
  curl \
  git \
  gpg \
  jq \
  wget \
  xz-utils

RUN wget -O /usr/bin/semver \
  https://raw.githubusercontent.com/fsaintjacques/semver-tool/master/src/semver && \
  chmod +x /usr/bin/semver

# set the runtime user to a non-root user and the same user as used by the github runners for actions runs.
USER runtime

# Project files
COPY lib lib
COPY bumper.sh bumper.sh

# SHELL
SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

# Initial command
ENTRYPOINT ["/opt/bumper.sh"]
