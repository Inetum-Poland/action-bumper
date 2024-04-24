FROM alpine:3.19.1

SHELL ["/bin/ash", "-eo", "pipefail", "-c"]

RUN apk --no-cache add git git-lfs jq grep curl bash
RUN wget -O /usr/local/bin/semver \
  https://raw.githubusercontent.com/fsaintjacques/semver-tool/master/src/semver
COPY entrypoint.sh /entrypoint.sh

# set the runtime user to a non-root user and the same user as used by the github runners for actions runs.
USER 1001

ENTRYPOINT ["/entrypoint.sh"]
