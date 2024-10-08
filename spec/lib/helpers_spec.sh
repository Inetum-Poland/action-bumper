Describe 'lib/helpers.sh'
  Include lib/helpers.sh
  Include lib/debug.sh

  INETUM_POLAND_ACTION_BUMPER_DEBUG="true"

  Describe 'bump_semver_tags'
    BUMPER_NEXT_VERSION="1.2.3"

    Parameters
      true 1.2.3 1.2 1
      false '' '' ''
    End
    It "bumps semver tags semver: ${1}"
      INPUT_BUMP_SEMVER="${1}"
      MAJOR=
      MINOR=
      PATCH=

      When call bump_semver_tags
      The status should be success
      The variable 'PATCH' should eq "${2}"
      The variable 'MINOR' should eq "${3}"
      The variable 'MAJOR' should eq "${4}"
    End
  End

  Describe 'remove_v_prefix'
    BUMPER_NEXT_VERSION="v1.2.3"

    It 'removes v prefix'
      INPUT_BUMP_INCLUDE_V="true"

      When call remove_v_prefix
      The status should be success
      The variable 'BUMPER_NEXT_VERSION' should eq 'v1.2.3'
    End

    It 'does not remove v prefix'
      INPUT_BUMP_INCLUDE_V="false"

      When call remove_v_prefix
      The status should be success
      The variable 'BUMPER_NEXT_VERSION' should eq '1.2.3'
    End
  End

  Describe 'check_missing_tags'
    Parameters
      "major" "v1.0.0"
      "minor" "v0.1.0"
      "patch" "v0.0.1"
      "none" "v0.0.0"
    End


    It "sets version if tags not found and level is '$1'"
      BUMPER_BUMP_LEVEL="$1"

      When call check_missing_tags
      The status should be success
      The variable 'BUMPER_NEXT_VERSION' should eq "$2"
    End
  End

  Describe 'bump_tag'
    DEBUG_GITHUB_EVENT_PATH=

    exit() {
      true
    }

    __append_github_output() {
      true
    }

    It 'skips bump tag'
      INPUT_BUMP_FAIL_IF_NO_LEVEL="false"
      BUMPER_BUMP_LEVEL=

      When call bump_tag
      The status should be success
      The output should eq '::notice ::Job skiped as no bump label is found. Do nothing.'
    End

    It 'fails while no bump label is found'
      INPUT_BUMP_FAIL_IF_NO_LEVEL="true"
      BUMPER_BUMP_LEVEL=

      When call bump_tag
      The status should be success
      The output should eq '::error ::Job failed as no bump label is found.'
    End

    It 'skips bump tag while a "none" bump label is found'
      INPUT_BUMP_FAIL_IF_NO_LEVEL=
      BUMPER_BUMP_LEVEL="none"

      When call bump_tag
      The status should be success
      The output should eq "::notice ::Job skiped as bump level is 'none'. Do nothing."
    End

    It 'checks missing tags while bumping'
      INPUT_BUMP_FAIL_IF_NO_LEVEL=
      BUMPER_BUMP_LEVEL="patch"
      BUMPER_CURRENT_VERSION=
      PR_NUMBER="1"
      PR_TITLE="test"

      When call bump_tag
      The status should be success
      The variable 'BUMPER_NEXT_VERSION' should eq 'v0.0.1'
      The variable 'BUMPER_TAG_MESSAGE' should eq "v0.0.1: PR #1 - test"
    End

    Describe 'checks all options'
      Parameters
        "patch" "v1.2.3" "v1.2.4"
        "patch" "v2.0.0" "v2.0.1"
        "minor" "v1.2.3" "v1.3.0"
        "minor" "v2.0.0" "v2.1.0"
        "major" "v1.2.3" "v2.0.0"
        "major" "v2.0.0" "v3.0.0"
      End

      It "bumps tag version $1 $2 -> $3"
        INPUT_BUMP_FAIL_IF_NO_LEVEL=
        BUMPER_BUMP_LEVEL="$1"
        BUMPER_CURRENT_VERSION="$2"
        PR_NUMBER="1"
        PR_TITLE="test"

        When call bump_tag
        The status should be success
        The variable 'BUMPER_NEXT_VERSION' should eq "$3"
        The variable 'BUMPER_TAG_MESSAGE' should eq "$3: PR #1 - test"
        The variable 'BUMPER_CURRENT_VERSION' should eq "$2"
      End
    End
  End

  Describe 'setup_git_tag'
    DEBUG_GITHUB_EVENT_PATH=
    GITHUB_API_URL="https://api.github.com"
    GITHUB_REPOSITORY="inetum-poland/action-bumper"

    __append_github_output() {
      true
    }


    It 'returns git tag from api'
      INPUT_GITHUB_TOKEN=

      __get_git_tag_from_api() {
        echo 'v1.2.3'
      }

      When call setup_git_tag
      The status should be success
      The variable 'BUMPER_CURRENT_VERSION' should eq 'v1.2.3'
    End

    It 'returns git tag from file'
      INPUT_GITHUB_TOKEN=
      DEBUG_GITHUB_EVENT_PATH="test"

      __get_git_tag_from_file() {
        echo 'v1.2.3'
      }

      When call setup_git_tag
      The status should be success
      The variable 'BUMPER_CURRENT_VERSION' should eq 'v1.2.3'
    End

    It 'returns nothing'
      INPUT_GITHUB_TOKEN=
      DEBUG_GITHUB_EVENT_PATH="test"

      __get_git_tag_from_file() {
        echo 'null'
      }

      When call setup_git_tag
      The status should be success
      The variable 'BUMPER_CURRENT_VERSION' should be undefined
    End
  End

  Describe 'setup_vars'
    # Include lib/pr_event.sh
    # Include lib/push_event.sh

    DEBUG_GITHUB_EVENT_PATH="spec/lib/push_event"
    GITHUB_API_URL="https://api.github.com"
    GITHUB_EVENT_PATH="test"
    GITHUB_REPOSITORY="inetum-poland/action-bumper"
    GITHUB_SHA="54fa23aef40b58c8f22350c830f7a89dad0121bc"
    INPUT_BUMP_MAJOR="bumper:major"
    INPUT_BUMP_MINOR="bumper:minor"
    INPUT_BUMP_NONE="bumper:none"
    INPUT_BUMP_PATCH="bumper:patch"
    INPUT_BUMP_DEFAULT_LEVEL="patch"
    INPUT_GITHUB_TOKEN=

    Context 'pr_event'
      __get_action() {
        echo 'labeled'
      }

      setup_pr_number_from_pr_event() {
        echo '1'
      }

      setup_pr_title_from_pr_event() {
        echo 'test'
      }

      setup_labels_from_pr_event() {
        echo 'feature'
      }

      It 'returns defaults'
        When call setup_vars
        The status should be success
        The variable 'BUMPER_BUMP_LEVEL' should eq 'patch'
      End
    End

    Context 'push_event'
      __get_action() {
        echo 'push'
      }

      setup_pr_number_from_push_event() {
        echo '1'
      }

      setup_pr_title_from_push_event() {
        echo 'test'
      }

      setup_labels_from_push_event() {
        echo 'feature'
      }

      It 'returns defaults'
        When call setup_vars
        The status should be success
        The variable 'BUMPER_BUMP_LEVEL' should eq 'patch'
      End
    End

    Context 'labels'
      __get_action() {
        echo 'labeled'
      }

      setup_pr_number_from_pr_event() {
        echo '1'
      }

      setup_pr_title_from_pr_event() {
        echo 'test'
      }


      It 'returns none'
        setup_labels_from_pr_event() {
          echo 'bumper:none feature'
        }

        When call setup_vars
        The status should be success
        The variable 'BUMPER_BUMP_LEVEL' should eq 'none'
      End

      It 'returns patch'
        setup_labels_from_pr_event() {
          echo 'bumper:patch feature'
        }

        When call setup_vars
        The status should be success
        The variable 'BUMPER_BUMP_LEVEL' should eq 'patch'
      End

      It 'returns minor'
        setup_labels_from_pr_event() {
          echo 'bumper:minor feature'
        }

        When call setup_vars
        The status should be success
        The variable 'BUMPER_BUMP_LEVEL' should eq 'minor'
      End

      It 'returns major'
        setup_labels_from_pr_event() {
          echo 'bumper:major feature'
        }

        When call setup_vars
        The status should be success
        The variable 'BUMPER_BUMP_LEVEL' should eq 'major'
      End
    End
  End
End
