FROM alpine:3.19.1

RUN apk --no-cache add git git-lfs jq grep curl bash
RUN wget -O /usr/local/bin/semver \
https://raw.githubusercontent.com/fsaintjacques/semver-tool/master/src/semver && \
chmod +x /usr/local/bin/semver
COPY entrypoint.sh /entrypoint.sh

SHELL ["/bin/bash", "-eo", "pipefail", "-c"]

# set the runtime user to a non-root user and the same user as used by the github runners for actions runs.
USER 1001

ENTRYPOINT ["/entrypoint.sh"]
