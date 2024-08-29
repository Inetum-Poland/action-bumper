Describe 'lib/message.sh'
  Include lib/message.sh

  GITHUB_API_URL="https://api.github.com"
  GITHUB_REPOSITORY="inetum-poland/action-bumper"

  setup() {
    GITHUB_OUTPUT="/tmp/shellspec-${RANDOM}"
    touch ${GITHUB_OUTPUT}
  }

  cleanup() {
    rm -rf ${GITHUB_OUTPUT}
  }

  BeforeEach 'setup'
  AfterEach 'cleanup'

  Describe 'make_pr_status'
    BUMPER_CURRENT_VERSION="1.2.3"
    BUMPER_NEXT_VERSION="1.2.4"
    GITHUB_SERVER_URL="https://github.com"
    GITHUB_REPOSITORY="inetum-poland/action-bumper"
    GITHUB_RUN_ID="1"

    __get_head_label() {
      echo "feature/test"
    }

    It 'creates a comment content'
      INPUT_ADD_LATEST="false"
      ACTION="labeled"

      When call make_pr_status
      The status should be success
      The contents of file ${GITHUB_OUTPUT} should eq 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ labeled<br>**Next version**: 1.2.4<br>**Changes**: [1.2.3...feature/test](https://github.com/inetum-poland/action-bumper/compare/1.2.3...feature/test)'
    End

    It 'creates a comment content without latest'
      INPUT_ADD_LATEST="true"
      ACTION="labeled"

      When call make_pr_status
      The status should be success
      The contents of file ${GITHUB_OUTPUT} should eq 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ labeled<br>**Next version**: 1.2.4 / latest<br>**Changes**: [1.2.3...feature/test](https://github.com/inetum-poland/action-bumper/compare/1.2.3...feature/test)'
    End
  End

  Describe 'make_push_status'
    BUMPER_CURRENT_VERSION="1.2.3"
    BUMPER_NEXT_VERSION="1.2.4"
    GITHUB_SERVER_URL="https://github.com"
    GITHUB_REPOSITORY="inetum-poland/action-bumper"
    GITHUB_RUN_ID="1"

    It 'creates a comment content'
      INPUT_ADD_LATEST="false"

      When call make_push_status
      The status should be success
      The contents of file ${GITHUB_OUTPUT} should eq 'tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](https://github.com/inetum-poland/action-bumper/actions/runs/1)<br>**New version**: [1.2.4](https://github.com/inetum-poland/action-bumper/releases/tag/1.2.4)<br>**Changes**: [1.2.3...1.2.4](https://github.com/inetum-poland/action-bumper/compare/1.2.3...1.2.4)'
    End

    It 'creates a comment content with the latest'
      INPUT_ADD_LATEST="true"


      When call make_push_status
      The status should be success
      The contents of file ${GITHUB_OUTPUT} should eq 'tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](https://github.com/inetum-poland/action-bumper/actions/runs/1)<br>**New version**: [1.2.4 / latest](https://github.com/inetum-poland/action-bumper/releases/tag/1.2.4)<br>**Changes**: [1.2.3...1.2.4](https://github.com/inetum-poland/action-bumper/compare/1.2.3...1.2.4)'
    End
  End

  Describe 'make_merge_semver_status'
    PATCH="1.2.3"
    MINOR="1.2"
    MAJOR="1"

    It 'passes the info to the next step'
      When call make_merge_semver_status
      The contents of file ${GITHUB_OUTPUT} should eq 'tag_status=New patch: 1.2.3<br>New minor: 1.2<br>New major: 1'
      The status should be success
    End
  End
End
