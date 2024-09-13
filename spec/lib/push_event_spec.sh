Describe 'lib/push_event.sh'
  Include lib/push_event.sh


  Describe 'list_pulls'
    GITHUB_API_URL="https://api.github.com"
    GITHUB_REPOSITORY="inetum-poland/action-bumper"

    It 'returns pulls from file'
      DEBUG_GITHUB_EVENT_PATH="./spec/lib/push_event"

      When call list_pulls
      The output should eq '[{"number": 1, "title":"feat(gha): align the gh actions before publish","merge_commit_sha":"54fa23aef40b58c8f22350c830f7a89dad0121bc","labels":[{"name":"feature"},{"name":"bumper:patch"},{"name":"doc"}]},{"number": 2, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"0f9485fc7a911a6ed4acbc0d35c62fc844d83c56","labels":[{"name":"feature"},{"name":"bumper:patch"}]},{"number": 3, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"05d0d3ff25ee0a40dac1aba4d48492cf72c4da5d","labels":[{"name":"feature"},{"name":"bumper:patch"}]}]'
      The status should be success
    End

    It 'returns pulls from curl with token'
      DEBUG_GITHUB_EVENT_PATH=
      INPUT_GITHUB_TOKEN="XXX"

      curl() {
        echo '[{"number": 1, "title":"feat(gha): align the gh actions before publish","merge_commit_sha":"54fa23aef40b58c8f22350c830f7a89dad0121bc","labels":[{"name":"feature"},{"name":"bumper:patch"},{"name":"doc"}]},{"number": 2, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"0f9485fc7a911a6ed4acbc0d35c62fc844d83c56","labels":[{"name":"feature"},{"name":"bumper:patch"}]},{"number": 3, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"05d0d3ff25ee0a40dac1aba4d48492cf72c4da5d","labels":[{"name":"feature"},{"name":"bumper:patch"}]}]'
      }

      When call list_pulls
      The output should eq '[{"number": 1, "title":"feat(gha): align the gh actions before publish","merge_commit_sha":"54fa23aef40b58c8f22350c830f7a89dad0121bc","labels":[{"name":"feature"},{"name":"bumper:patch"},{"name":"doc"}]},{"number": 2, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"0f9485fc7a911a6ed4acbc0d35c62fc844d83c56","labels":[{"name":"feature"},{"name":"bumper:patch"}]},{"number": 3, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"05d0d3ff25ee0a40dac1aba4d48492cf72c4da5d","labels":[{"name":"feature"},{"name":"bumper:patch"}]}]'
      The status should be success
    End

    It 'returns pulls from curl without token'
      DEBUG_GITHUB_EVENT_PATH=
      INPUT_GITHUB_TOKEN=

      curl() {
        echo '[{"number": 1, "title":"feat(gha): align the gh actions before publish","merge_commit_sha":"54fa23aef40b58c8f22350c830f7a89dad0121bc","labels":[{"name":"feature"},{"name":"bumper:patch"},{"name":"doc"}]},{"number": 2, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"0f9485fc7a911a6ed4acbc0d35c62fc844d83c56","labels":[{"name":"feature"},{"name":"bumper:patch"}]},{"number": 3, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"05d0d3ff25ee0a40dac1aba4d48492cf72c4da5d","labels":[{"name":"feature"},{"name":"bumper:patch"}]}]'
      }

      When call list_pulls
      The output should eq '[{"number": 1, "title":"feat(gha): align the gh actions before publish","merge_commit_sha":"54fa23aef40b58c8f22350c830f7a89dad0121bc","labels":[{"name":"feature"},{"name":"bumper:patch"},{"name":"doc"}]},{"number": 2, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"0f9485fc7a911a6ed4acbc0d35c62fc844d83c56","labels":[{"name":"feature"},{"name":"bumper:patch"}]},{"number": 3, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"05d0d3ff25ee0a40dac1aba4d48492cf72c4da5d","labels":[{"name":"feature"},{"name":"bumper:patch"}]}]'
      The status should be success
      The error should eq 'INPUT_GITHUB_TOKEN is not available. Subscequent GitHub API call may fail due to API limit.'
    End
  End

  Describe 'uses list_pulls'
    DEBUG_GITHUB_EVENT_PATH=
    GITHUB_SHA="54fa23aef40b58c8f22350c830f7a89dad0121bc"

    list_pulls() {
      echo '[{"number": 1, "title":"feat(gha): align the gh actions before publish","merge_commit_sha":"54fa23aef40b58c8f22350c830f7a89dad0121bc","labels":[{"name":"feature"},{"name":"bumper:patch"},{"name":"doc"}]},{"number": 2, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"0f9485fc7a911a6ed4acbc0d35c62fc844d83c56","labels":[{"name":"feature"},{"name":"bumper:patch"}]},{"number": 3, "title":"chore(tf): action-bumper merge test","merge_commit_sha":"05d0d3ff25ee0a40dac1aba4d48492cf72c4da5d","labels":[{"name":"feature"},{"name":"bumper:patch"}]}]'
    }

    Describe 'setup_labels_from_push_event'
      It 'returns labels'
        When call setup_labels_from_push_event
        The line 1 of output should eq "\"feature\""
        The line 2 of output should eq "\"bumper:patch\""
        The line 3 of output should eq "\"doc\""
        The status should be success
      End
    End

    Describe 'setup_pr_number_from_push_event'
      It 'returns number'
        When call setup_pr_number_from_push_event
        The output should eq '1'
        The status should be success
      End
    End

    Describe 'setup_pr_title_from_push_event'
      It 'returns title'
        When call setup_pr_title_from_push_event
        The output should eq 'feat(gha): align the gh actions before publish'
        The status should be success
      End
    End
  End
End
