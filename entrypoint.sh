#!/usr/bin/env bash

# https://github.com/fsaintjacques/semver-tool/tree/master

set -ex
export PS4='+(${BASH_SOURCE}:${LINENO}): ${FUNCNAME[0]:+${FUNCNAME[0]}(): }'

# --- set safe directory ------------------------------------------------------

if [[ -n "${GITHUB_WORKSPACE}" ]]; then
  git config --global --add safe.directory "${GITHUB_WORKSPACE}" || exit
  cd "${GITHUB_WORKSPACE}" || exit
fi

# --- debug helpers -----------------------------------------------------------

# Initial debug.
init_debug() {
  # For debugging.
  if [[ -n "${DEBUG_GITHUB_EVENT_PATH}" ]]; then
    GITHUB_EVENT_PATH="${DEBUG_GITHUB_EVENT_PATH}"
    # shellcheck disable=SC1091
    source .input.env
  fi
}

# --- functions for helpers ---------------------------------------------------

# Prepare and post a status.
post_pre_status() {
  head_label="$(jq -r '.pull_request.head.label' < "${GITHUB_EVENT_PATH}" )"
  compare=""

  if [[ -n "${BUMPER_CURRENT_VERSION}" ]]; then
    compare="**Changes**:[${BUMPER_CURRENT_VERSION}...${head_label}](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/compare/${BUMPER_CURRENT_VERSION}...${head_label})"
  fi

  post_txt="🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ ${ACTION}<br>**Next version**: ${BUMPER_NEXT_VERSION}<br>${compare}"

  FROM_FORK=$(jq -r '.pull_request.head.repo.fork' < "${GITHUB_EVENT_PATH}")

  if [[ "${FROM_FORK}" == "true" ]]; then
    post_warning "${post_txt}"
  else
    post_comment "${post_txt}"
  fi

  echo "::notice ::${post_txt}"
}

# Prepare and post a status.
post_post_status() {
  compare=""

  if [[ -n "${BUMPER_CURRENT_VERSION}" ]]; then
    compare="**Changes**:[${BUMPER_CURRENT_VERSION}...${BUMPER_NEXT_VERSION}](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/compare/${BUMPER_CURRENT_VERSION}...${BUMPER_NEXT_VERSION})"
  fi

  post_txt="🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/actions/runs/${GITHUB_RUN_ID})<br>**New version**: [${BUMPER_NEXT_VERSION}](${GITHUB_SERVER_URL}/${GITHUB_REPOSITORY}/releases/tag/${BUMPER_NEXT_VERSION})<br>${compare}"

  post_comment "${post_txt}"

  echo "::notice ::${post_txt}"
}

# TODO! opened and labeled are producing two separate comments on the same PR.
# Post a comment.
post_comment() {
  body_text="$1"
  endpoint="${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/issues/${PR_NUMBER}/comments"
  update_endpoint="${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/issues/comments/"

  body="$(echo -e "${body_text}" | jq -ncR "{body: input}")"

  # check if the comment has been already posted
  comment_id=$(curl -s -H "Authorization: token ${INPUT_GITHUB_TOKEN}" "${endpoint}" | jq -r '.[] | select((.body | contains("action-bumper")) and (.user.login == "github-actions[bot]") and (.user.type == "Bot")) | .id' | sort -V | tail -1)

  output=

  if [[ -n "${comment_id}" ]]; then
    # comment already posted, update it
    output=$(curl -H "Authorization: token ${INPUT_GITHUB_TOKEN}" -X PATCH -d "${body}" "${update_endpoint}${comment_id}")
  else
    output=$(curl -H "Authorization: token ${INPUT_GITHUB_TOKEN}" -d "${body}" "${endpoint}")
  fi

  echo "::notice ::$(echo "${output}" | jq -r '.id')"
}

# Post a warning comment.
post_warning() {
  body_text=$(echo "$1" | sed -e ':a' -e 'N' -e '$!ba' -e 's/\n/%0A/g')
  echo "::warning ::${body_text}"
}

# --- functions for pr event --------------------------------------------------

# Get label name from the pull request.
setup_labels_from_pr_event() {
  jq -r '.pull_request.labels[].name' < "${GITHUB_EVENT_PATH}" | tr '\n' ' '
}

# Get number from the pull request.
setup_pr_number_from_pr_event() {
  jq -r '.pull_request.number' < "${GITHUB_EVENT_PATH}"
}

# Get title from the pull request.
setup_pr_title_from_pr_event() {
  jq -r '.pull_request.title' < "${GITHUB_EVENT_PATH}"
}

# --- functions for push event ------------------------------------------------

# Get list of pull requests.
list_pulls() {
  pulls_endpoint="${GITHUB_API_URL}/repos/${GITHUB_REPOSITORY}/pulls?state=closed&sort=updated&direction=desc"

  if [[ -n "${INPUT_GITHUB_TOKEN}" ]]; then
    curl -s -H "Authorization: token ${INPUT_GITHUB_TOKEN}" "${pulls_endpoint}"
  else
    echo "INPUT_GITHUB_TOKEN is not available. Subscequent GitHub API call may fail due to API limit." >&2
    curl -s "${pulls_endpoint}"
  fi
}

# Get labels from the pull request.
setup_labels_from_push_event() {
  pull_request="$(list_pulls | jq ".[] | select(.merge_commit_sha==\"${GITHUB_SHA}\")")"
  echo "${pull_request}" | jq '.labels | .[].name'
}

# Get number from the pull request.
setup_pr_number_from_push_event() {
  pull_request="$(list_pulls | jq ".[] | select(.merge_commit_sha==\"${GITHUB_SHA}\")")"
  echo "${pull_request}" | jq -r .number
}

# Get title from the pull request.
setup_pr_title_from_push_event() {
  pull_request="$(list_pulls | jq ".[] | select(.merge_commit_sha==\"${GITHUB_SHA}\")")"
  echo "${pull_request}" | jq -r .title
}

# --- helper functions --------------------------------------------------------

# Check if semver is installed.
semver_check() {
  # Check if semver is installed.
  if ! semver -v &> /dev/null; then
    echo "::error ::semver is not installed."
    exit 1
  fi
}

# Check if jq is installed.
jq_check() {
  # Check if jq is installed.
  if ! jq -V &> /dev/null; then
    echo "::error ::jq is not installed."
    exit 1
  fi
}

# Check if the repository is shallowed.
git_shallow_repo() {
  # check the repository is shallowed.
  # comes from https://stackoverflow.com/questions/37531605/how-to-test-if-git-repository-is-shallow
  # the repository is shallowed, so we need to fetch all history.
  # Fetch history as well because bump uses git history (git tag --merged).
  if "$(git rev-parse --is-shallow-repository)"; then
    git fetch --tags -f # Fetch existing tags before bump.
    git fetch --prune --unshallow
  fi
}

# Setup the necessary variables based on the GitHub event.
setup_vars() {
  if [[ "${ACTION}" =~ ^(labeled|unlabeled|synchronize|opened|reopened)$ ]]; then
    PR_NUMBER=$(setup_pr_number_from_pr_event)
    PR_TITLE=$(setup_pr_title_from_pr_event)
    BUMPER_LABELS=$(setup_labels_from_pr_event)
    # elif [[ $(jq -r '.ref' < "${GITHUB_EVENT_PATH}") =~ "refs/" ]]; then
  else
    PR_NUMBER=$(setup_pr_number_from_push_event)
    PR_TITLE=$(setup_pr_title_from_push_event)
    BUMPER_LABELS=$(setup_labels_from_push_event)
  fi

  if echo "${BUMPER_LABELS}" | grep "${INPUT_BUMP_NONE}" ; then
    BUMPER_BUMP_LEVEL="none"
  fi

  if echo "${BUMPER_LABELS}" | grep "${INPUT_BUMP_PATCH}" ; then
    BUMPER_BUMP_LEVEL="patch"
  fi

  if echo "${BUMPER_LABELS}" | grep "${INPUT_BUMP_MINOR}" ; then
    BUMPER_BUMP_LEVEL="minor"
  fi

  if echo "${BUMPER_LABELS}" | grep "${INPUT_BUMP_MAJOR}" ; then
    BUMPER_BUMP_LEVEL="major"
  fi
}

# A function that processes the current version to determine the next version and generate a tag message.
setup_git_tag() {
  BUMPER_CURRENT_VERSION="$(git tag | grep -E "v?[0-9]+\.[0-9]+\.[0-9]+.*" | sort -V | tail -1)"

  if [[ -z "${DEBUG_GITHUB_EVENT_PATH}" ]]; then
    echo "current_version=${BUMPER_CURRENT_VERSION}" >> "$GITHUB_OUTPUT"
  fi

  BUMPER_BUMP_LEVEL="${INPUT_DEFAULT_BUMP_LEVEL}"
}

# A function that processes the bump level and current version to determine the next version and generate a tag message.
bump_tag() {
  if [[ -z "${BUMPER_BUMP_LEVEL}" ]]; then
    if [[ "${INPUT_FAIL_IF_NO_BUMP}" == "true" ]]; then
      echo "::error ::PR fails as no bump label is found."
      exit 1
    fi

    echo "::notice ::PR with labels for bump not found. Do nothing."

    if [[ -z "${DEBUG_GITHUB_EVENT_PATH}" ]]; then
      echo "skip=true" >> "$GITHUB_OUTPUT"
    fi
    exit
  fi

  BUMPER_NEXT_VERSION="v$(semver bump "${BUMPER_BUMP_LEVEL}" "${BUMPER_CURRENT_VERSION}")" || true

  if [[ -z "${DEBUG_GITHUB_EVENT_PATH}" ]]; then
    echo "next_version=${NEXT_VERSION}" >> "$GITHUB_OUTPUT"
  fi

  BUMPER_TAG_MESSAGE="${BUMPER_NEXT_VERSION}: PR #${PR_NUMBER} - ${PR_TITLE}"

  if [[ -z "${DEBUG_GITHUB_EVENT_PATH}" ]]; then
    echo "message=${BUMPER_TAG_MESSAGE}" >> "$GITHUB_OUTPUT"
  fi
}

# Set next version tag in case existing tags not found.
check_missing_tags() {
  # Set next version tag in case existing tags not found.
  if [[ -z "${BUMPER_NEXT_VERSION}" && -z "$(git tag)" ]]; then
    case "${BUMP_LEVEL}" in
      major)
        BUMPER_NEXT_VERSION="v1.0.0"
        ;;
      minor)
        BUMPER_NEXT_VERSION="v0.1.0"
        ;;
      patch)
        BUMPER_NEXT_VERSION="v0.0.1"
        ;;
      *)
        BUMPER_NEXT_VERSION="v0.0.0"
        ;;
    esac
  fi
}

# Remove 'v' prefix from BUMPER_NEXT_VERSION if INPUT_INCLUDE_V is false.
remove_v_prefix() {
  # Remove 'v' prefix if variable is false casted from string
  if [[ "${INPUT_INCLUDE_V}" == "false" ]]; then
    BUMPER_NEXT_VERSION="${BUMPER_NEXT_VERSION/^v/}"
  fi
}

# Process tags based on different conditions and perform Git operations accordingly.
make_and_push_tag() {
  if [[ "${INPUT_DRY_RUN}" == "true" || -n "${DEBUG_GITHUB_EVENT_PATH}" ]]; then
    echo "BUMPER_NEXT_VERSION=${BUMPER_NEXT_VERSION}"
    echo "TAG_MESSAGE=${BUMPER_TAG_MESSAGE}"
    exit 0
  fi

  # Push the next tag.
  git tag -a "${BUMPER_NEXT_VERSION}" -m "${TAG_MESSAGE}"

  if [[ -n "${INPUT_GITHUB_TOKEN}" ]]; then
    git remote set-url origin "https://${GITHUB_ACTOR}:${INPUT_GITHUB_TOKEN}@github.com/${GITHUB_REPOSITORY}.git"
  fi

  git push origin "${BUMPER_NEXT_VERSION}"
}

# Set up Git config.
setup_git_config() {
  # Set up Git.
  if [[ "${INPUT_DRY_RUN}" == "true" || -n "${DEBUG_GITHUB_EVENT_PATH}" ]]; then
    true
  else
    git config user.name "${INPUT_TAG_AS_USER:-${GITHUB_ACTOR}}"
    git config user.email "${INPUT_TAG_AS_EMAIL:-${GITHUB_ACTOR}@users.noreply.github.com}"
  fi
}

# Semver update for tags.
bump_semver_tags() {
  PATCH="${BUMPER_CURRENT_VERSION}" # v1.2.3
  MINOR="${PATCH%.*}"               # v1.2
  MAJOR="${MINOR%.*}"               # v1
}

# Semver update for tags.
make_and_push_semver_tags() {
  if [[ "${INPUT_DRY_RUN}" == "true" || -n "${DEBUG_GITHUB_EVENT_PATH}" ]]; then
    echo "PATCH=${PATCH}"
    echo "MINOR=${MINOR}"
    echo "MAJOR=${MAJOR}"
    exit 0
  fi

  git tag -fa "${MINOR}" -m "${TAG_MESSAGE}"
  git tag -fa "${MAJOR}" -m "${TAG_MESSAGE}"

  git push --force origin "${MINOR}"
  git push --force origin "${MAJOR}"
}

# --- main --------------------------------------------------------------------

ACTION=
PR_NUMBER=
PR_TITLE=
BUMPER_LABELS=
BUMPER_CURRENT_VERSION=
BUMPER_BUMP_LEVEL=
BUMPER_NEXT_VERSION=
BUMPER_TAG_MESSAGE=
PATCH=
MINOR=
MAJOR=

main() {
  jq_check
  semver_check

  if [[ -n "${GITHUB_EVENT_PATH}" ]]; then
    cat "${GITHUB_EVENT_PATH}"
  fi
  ACTION=$(jq -r '.action' < "${GITHUB_EVENT_PATH}")

  init_debug
  git_shallow_repo
  setup_git_tag
  setup_git_config

  echo "::notice ::${ACTION}"

  if [[ $(jq -r '.ref' < "${GITHUB_EVENT_PATH}") =~ "refs/tags/" && ${INPUT_BUMP_SEMVER} == "true" ]]; then
    bump_semver_tags
    remove_v_prefix
    make_and_push_semver_tags
  else
    setup_vars
    bump_tag
    check_missing_tags
    remove_v_prefix

    if [[ "${ACTION}" =~ ^(labeled|unlabeled|synchronize|opened|reopened)$ ]]; then
      post_pre_status
    else
      make_and_push_tag
      post_post_status
    fi
  fi

  exit 0
}

main

# TEST
