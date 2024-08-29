Describe 'lib/check.sh'
  Include lib/check.sh

  Describe 'semver_check'
    It 'checks semver'
      When call semver_check
      The status should be success
    End

    It 'checks semver with error'
      semver() {
        false
      }

      exit() {
        false
      }

      When call semver_check
      The status should be failure
      The output should eq '::error ::semver is not installed.'
    End
  End

  Describe 'jq_check'
    It 'checks jq'
      When call jq_check
      The status should be success
    End

    It 'checks jq with error'
      jq() {
        false
      }

      exit() {
        false
      }

      When call jq_check
      The status should be failure
      The output should eq '::error ::jq is not installed.'
    End
  End
End
