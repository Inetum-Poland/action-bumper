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

  Describe 'curl_check'
    It 'checks curl'
      When call curl_check
      The status should be success
    End

    It 'checks curl with error'
      curl() {
        false
      }

      exit() {
        false
      }

      When call curl_check
      The status should be failure
      The output should eq '::error ::curl is not installed.'
    End
  End

  Describe 'git_check'
    It 'checks git'
      When call git_check
      The status should be success
    End

    It 'checks git with error'
      git() {
        false
      }

      exit() {
        false
      }

      When call git_check
      The status should be failure
      The output should eq '::error ::git is not installed.'
    End
  End
End
