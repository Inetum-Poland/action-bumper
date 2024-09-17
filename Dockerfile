FROM ubuntu:24.10

# ENVS
# ENV MISE_DATA_DIR=/opt/mise
# ENV MISE_CACHE_DIR=/opt/mise/cache
# ENV PATH="/opt/mise/shims:${PATH}"

# Workdir
WORKDIR /opt/bumper

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

# MISE
# RUN install -dm 755 /etc/apt/keyrings && \
#   bash -c 'wget -qO - https://mise.jdx.dev/gpg-key.pub | gpg --dearmor | tee /etc/apt/keyrings/mise-archive-keyring.gpg 1> /dev/null' && \
#   bash -c 'echo "deb [signed-by=/etc/apt/keyrings/mise-archive-keyring.gpg arch=amd64] https://mise.jdx.dev/deb stable main" | tee /etc/apt/sources.list.d/mise.list' && \
#   apt update && apt install -y mise

# install mise folder with permissions
# RUN install -dm 775 -o runtime -g runtime /opt/mise

# set the runtime user to a non-root user and the same user as used by the github runners for actions runs.
USER runtime

# COPY --chown=runtime:runtime .mise.toml .mise.toml
# RUN mise install -y --verbose

# Project files
COPY lib lib
COPY bumper.sh bumper.sh

# SHELL
SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

# Initial command
ENTRYPOINT ["/opt/bumper/bumper.sh"]
