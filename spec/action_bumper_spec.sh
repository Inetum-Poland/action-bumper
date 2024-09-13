# shellcheck disable=SC2148

Describe 'action_bumper.sh'
  setup() {
    GITHUB_OUTPUT="/tmp/shellspec-${RANDOM}"
    bash -c "touch ${GITHUB_OUTPUT}"
  }

  cleanup() {
    bash -c "rm -rf ${GITHUB_OUTPUT}"
  }

  BeforeEach 'setup'
  AfterEach 'cleanup'

  Describe 'opened_event_bumper_auto'
    Include spec/action_bumper/opened_event_bumper_auto/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should include 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ opened<br>**Next version**: v0.9.2<br>**Changes**: [v0.9.1...feature/test](https://github.com/inetum-poland/action-bumper/compare/v0.9.1...feature/test)'
    End
  End

  Describe 'opened_event_bumper_major'
    Include spec/action_bumper/opened_event_bumper_major/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should include 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ opened<br>**Next version**: v1.0.0<br>**Changes**: [v0.9.1...feature/test](https://github.com/inetum-poland/action-bumper/compare/v0.9.1...feature/test)'
    End
  End

  Describe 'opened_event_bumper_minor'
    Include spec/action_bumper/opened_event_bumper_minor/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should include 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ opened<br>**Next version**: v0.10.0<br>**Changes**: [v0.9.1...feature/test](https://github.com/inetum-poland/action-bumper/compare/v0.9.1...feature/test)'
    End
  End

  Describe 'opened_event_bumper_major_latest'
    Include spec/action_bumper/opened_event_bumper_major_latest/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should include 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ opened<br>**Next version**: v1.0.0 / latest<br>**Changes**: [v0.9.1...feature/test](https://github.com/inetum-poland/action-bumper/compare/v0.9.1...feature/test)'
    End
  End

  Describe 'opened_event_bumper_none'
    Include spec/action_bumper/opened_event_bumper_none/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The stdout should eq "::notice ::Job skiped as bump level is 'none'. Do nothing."
    End
  End

  Describe 'opened_event_without_tags_without_labels'
    Include spec/action_bumper/opened_event_without_tags_without_labels/.input.env

    It 'does not the magic'
      When run source action_bumper.sh
      The status should be failure
      The stdout should eq "::error ::Job failed as no bump label is found."
    End
  End

  Describe 'opened_event_without_tags_without_labels_allow'
    Include spec/action_bumper/opened_event_without_tags_without_labels_allow/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The stdout should eq '::notice ::Job skiped as no bump label is found. Do nothing.'
    End
  End

  Describe 'push_event_tags'
    Include spec/action_bumper/push_event_tags/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should include 'tag_status=New patch: v0.9.1<br>New minor: v0.9<br>New major: v0'
      The line 1 of stderr should eq '> git config user.name "github-actions[bot]"'
      The line 2 of stderr should eq '> git config user.email "github-actions[bot]@users.noreply.github.com"'
      The line 3 of stderr should eq '> git remote set-url origin "https://github-actions[bot]:XXX@github.com/inetum-poland/action-bumper.git"'
      The line 4 of stderr should eq '> git tag -fa "v0.9" "^{commit}" -m ""'
      The line 5 of stderr should eq '> git tag -fa "v0" "^{commit}" -m ""'
      The line 6 of stderr should eq '> git push --force origin "v0.9"'
      The line 7 of stderr should eq '> git push --force origin "v0"'
    End
  End

  Describe 'push_event_with_labels'
    Include spec/action_bumper/push_event_with_labels/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should include 'tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](https://github.com/inetum-poland/action-bumper/actions/runs/1)<br>**New version**: [v0.9.2](https://github.com/inetum-poland/action-bumper/releases/tag/v0.9.2)<br>**Changes**: [v0.9.1...v0.9.2](https://github.com/inetum-poland/action-bumper/compare/v0.9.1...v0.9.2)'
      The line 1 of stderr should eq '> git config user.name "github-actions[bot]"'
      The line 2 of stderr should eq '> git config user.email "github-actions[bot]@users.noreply.github.com"'
      The line 3 of stderr should eq '> git remote set-url origin "https://github-actions[bot]:XXX@github.com/inetum-poland/action-bumper.git"'
      The line 4 of stderr should eq '> git tag -a "v0.9.2" -m "v0.9.2: PR #null - feat(gha): align the gh actions before publish"'
      The line 5 of stderr should eq '> git push origin "v0.9.2"'

    End
  End

  Describe 'push_event_with_labels_without_v'
    Include spec/action_bumper/push_event_with_labels_without_v/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should include 'tag_status=🚀 [[bumper]](https://github.com/inetum-poland/action-bumper) [Bumped!](https://github.com/inetum-poland/action-bumper/actions/runs/1)<br>**New version**: [0.9.2](https://github.com/inetum-poland/action-bumper/releases/tag/0.9.2)<br>**Changes**: [v0.9.1...0.9.2](https://github.com/inetum-poland/action-bumper/compare/v0.9.1...0.9.2)'
      The line 1 of stderr should eq '> git config user.name "github-actions[bot]"'
      The line 2 of stderr should eq '> git config user.email "github-actions[bot]@users.noreply.github.com"'
      The line 3 of stderr should eq '> git remote set-url origin "https://github-actions[bot]:XXX@github.com/inetum-poland/action-bumper.git"'
      The line 4 of stderr should eq '> git tag -a "0.9.2" -m "v0.9.2: PR #null - feat(gha): align the gh actions before publish"'
      The line 5 of stderr should eq '> git push origin "0.9.2"'
    End
  End

  Describe 'push_event_without_labels'
    Include spec/action_bumper/push_event_without_labels/.input.env

    It 'does not the magic'
      When run source action_bumper.sh
      The status should be failure
      The stdout should eq "::error ::Job failed as no bump label is found."
    End
  End

  Describe 'synchronize_event'
    Include spec/action_bumper/synchronize_event/.input.env

    It 'does the magic'
      When run source action_bumper.sh
      The status should be success
      The contents of file "${GITHUB_OUTPUT}" should include 'tag_status=🏷️ [[bumper]](https://github.com/inetum-poland/action-bumper) @ synchronize<br>**Next version**: v0.9.2<br>**Changes**: [v0.9.1...feature/test](https://github.com/inetum-poland/action-bumper/compare/v0.9.1...feature/test)'
    End
  End
End
