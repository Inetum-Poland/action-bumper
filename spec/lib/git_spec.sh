Describe 'lib/git.sh'
  Include lib/git.sh
  Include lib/debug.sh

  INETUM_POLAND_ACTION_BUMPER_DEBUG="true"
  GITHUB_REPOSITORY="inetum-poland/action-bumper"
  BUMPER_TAG_MESSAGE="1.2.3: PR #1 - test"
  BUMPER_NEXT_VERSION="1.2.3"

  git() {
    echo "git ${@}"
  }

  Describe 'setup_git_config'
    GITHUB_ACTOR="github-actions[bot]"

    git() {
      echo "git ${@}"
    }

    It 'sets up git config'
      INPUT_GITHUB_TOKEN=

      When call setup_git_config
      The status should be success
      The line 1 of output should eq '> git config user.name "github-actions[bot]"'
      The line 2 of output should eq '> git config user.email "github-actions[bot]@users.noreply.github.com"'
    End

    It 'sets up git config with token'
      INPUT_GITHUB_TOKEN="XXX"
      GITHUB_ACTOR="github-actions[bot]"
      GITHUB_REPOSITORY="inetum-poland/action-bumper"

      When call setup_git_config
      The status should be success
      The line 1 of output should eq '> git config user.name "github-actions[bot]"'
      The line 2 of output should eq '> git config user.email "github-actions[bot]@users.noreply.github.com"'
      The line 3 of output should eq '> git remote set-url origin "https://github-actions[bot]:XXX@github.com/inetum-poland/action-bumper.git"'
    End
  End

  Describe 'make_and_push_tag'
    GITHUB_ACTOR="github-actions[bot]"
    INPUT_GITHUB_TOKEN="XXX"

    It 'pushes tag'
      INPUT_ADD_LATEST="false"

      When call make_and_push_tag
      The status should be success
      The line 1 of output should eq '> git tag -a "1.2.3" -m "1.2.3: PR #1 - test"'
      The line 2 of output should eq '> git push origin "1.2.3"'
    End

    It 'pushes tag without latest'
      INPUT_ADD_LATEST="true"

      When call make_and_push_tag
      The status should be success
      The line 1 of output should eq '> git tag -a "1.2.3" -m "1.2.3: PR #1 - test"'
      The line 2 of output should eq '> git push origin "1.2.3"'
      The line 3 of output should eq '> git tag -fa latest "1.2.3^{commit}" -m "1.2.3: PR #1 - test"'
    End
  End

  Describe 'make_and_push_semver_tags'
    MINOR="1.2"
    MAJOR="1"

    It 'pushes tag'
      When call make_and_push_semver_tags
      The status should be success
      The line 1 of output should eq '> git tag -fa "1.2" "1.2.3^{commit}" -m "1.2.3: PR #1 - test"'
      The line 2 of output should eq '> git tag -fa "1" "1.2.3^{commit}" -m "1.2.3: PR #1 - test"'
      The line 3 of output should eq '> git push --force origin "1.2"'
      The line 4 of output should eq '> git push --force origin "1"'
    End
  End
End
