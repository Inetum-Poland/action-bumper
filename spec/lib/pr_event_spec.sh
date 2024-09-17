Describe 'lib/pr_event.sh'
  Include lib/pr_event.sh

  GITHUB_EVENT_PATH="./spec/lib/pr_event.json"
  __PR_EVENT="${GITHUB_EVENT_PATH}"

  Describe 'setup_labels_from_pr_event'
    It 'returns labels'

      When call setup_labels_from_pr_event "${__PR_EVENT}"
      The output should eq 'bumper:major wiki '
      The status should be success
    End
  End

  Describe 'setup_pr_number_from_pr_event'
    It 'returns number'
      When call setup_pr_number_from_pr_event "${__PR_EVENT}"
      The output should eq '1'
      The status should be success
    End
  End

  Describe 'setup_pr_title_from_pr_event'
    It 'returns title'
      When call setup_pr_title_from_pr_event "${__PR_EVENT}"
      The output should eq 'feat(test): test the features'
      The status should be success
    End
  End
End
