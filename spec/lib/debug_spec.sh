Describe 'lib/debug.sh'
  Include lib/debug.sh

  Describe 'execute_or_debug'
    It 'runs wet'
      INPUT_DRY_RUN="false"

      When call execute_or_debug "echo 'hello'"
      The output should equal "'hello'"
      The status should be success
    End

    It 'runs dry'
      INPUT_DRY_RUN="true"

      When call execute_or_debug "echo 'hello'"
      The stderr should equal "> echo 'hello'"
      The status should be success
    End
  End

  Describe 'init_debug'
    It 'does debug'
      DEBUG=
      ACTIONS_STEP_DEBUG=
      DEBUG_GITHUB_EVENT_PATH=spec/action_bumper/opened_event_bumper_auto

      When call init_debug
      The variable 'GITHUB_EVENT_PATH' should be defined
      The status should be success
    End
  End
End
