Describe 'lib/message.sh'
  Include lib/message.sh

  GITHUB_API_URL="https://api.github.com"
  GITHUB_REPOSITORY="inetum-poland/action-bumper"

  setup() {
    GITHUB_OUTPUT="/tmp/shellspec-${RANDOM}"
    bash -c "touch ${GITHUB_OUTPUT}"
  }

  cleanup() {
    bash -c "rm -rf ${GITHUB_OUTPUT}"
  }

  BeforeEach 'setup'
  AfterEach 'cleanup'

  Describe 'make_pr_status'
    BUMPER_CURRENT_VERSION="1.2.3"
    BUMPER_NEXT_VERSION="1.2.4"
    PATCH="1.2.4"
    MINOR="1.2"
    MAJOR="1"
    GITHUB_SERVER_URL="https://github.com"
    GITHUB_REPOSITORY="inetum-poland/action-bumper"
    GITHUB_RUN_ID="1"

    __get_head_label() {
      echo "feature/test"
    }

    Parameters
      false false 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ labeled<br>**Next version**: 1.2.4<br>**Changes**: [1.2.3...feature/test](https://github.com/inetum-poland/action-bumper/compare/1.2.3...feature/test)'
      false true 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ labeled<br>**Next version**: 1.2.4 / latest<br>**Changes**: [1.2.3...feature/test](https://github.com/inetum-poland/action-bumper/compare/1.2.3...feature/test)'
      true false 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ labeled<br>**Next version**: 1.2.4 / 1.2 / 1<br>**Changes**: [1.2.3...feature/test](https://github.com/inetum-poland/action-bumper/compare/1.2.3...feature/test)'
      true true 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ labeled<br>**Next version**: 1.2.4 / 1.2 / 1 / latest<br>**Changes**: [1.2.3...feature/test](https://github.com/inetum-poland/action-bumper/compare/1.2.3...feature/test)'
    End
    It "creates a pr status content with semver: ${1}, latest: ${2}"
      INPUT_BUMP_SEMVER="${1}"
      INPUT_ADD_LATEST="${2}"
      ACTION="labeled"

      When call make_pr_status
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should eq "${3}"
    End
  End

  Describe 'make_push_status'
    BUMPER_CURRENT_VERSION="1.2.3"
    BUMPER_NEXT_VERSION="1.2.4"
    PATCH="1.2.4"
    MINOR="1.2"
    MAJOR="1"
    GITHUB_SERVER_URL="https://github.com"
    GITHUB_REPOSITORY="inetum-poland/action-bumper"
    GITHUB_RUN_ID="1"

    Parameters
      false false 'tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](https://github.com/inetum-poland/action-bumper/actions/runs/1)<br>**New version**: [1.2.4](https://github.com/inetum-poland/action-bumper/releases/tag/1.2.4)<br>**Changes**: [1.2.3...1.2.4](https://github.com/inetum-poland/action-bumper/compare/1.2.3...1.2.4)'
      false true 'tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](https://github.com/inetum-poland/action-bumper/actions/runs/1)<br>**New version**: [1.2.4](https://github.com/inetum-poland/action-bumper/releases/tag/1.2.4) / latest<br>**Changes**: [1.2.3...1.2.4](https://github.com/inetum-poland/action-bumper/compare/1.2.3...1.2.4)'
      true false 'tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](https://github.com/inetum-poland/action-bumper/actions/runs/1)<br>**New version**: [1.2.4](https://github.com/inetum-poland/action-bumper/releases/tag/1.2.4) / 1.2 / 1<br>**Changes**: [1.2.3...1.2.4](https://github.com/inetum-poland/action-bumper/compare/1.2.3...1.2.4)'
      true true 'tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](https://github.com/inetum-poland/action-bumper/actions/runs/1)<br>**New version**: [1.2.4](https://github.com/inetum-poland/action-bumper/releases/tag/1.2.4) / 1.2 / 1 / latest<br>**Changes**: [1.2.3...1.2.4](https://github.com/inetum-poland/action-bumper/compare/1.2.3...1.2.4)'
    End
    It "creates a push status content with semver: ${1}, latest: ${2}"
      INPUT_BUMP_SEMVER="${1}"
      INPUT_ADD_LATEST="${2}"

      When call make_push_status
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should eq "${3}"
    End
  End
End
