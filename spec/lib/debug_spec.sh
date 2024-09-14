Describe 'lib/debug.sh'
  Include lib/debug.sh

  Describe 'exec_debug'
    It 'runs wet'
      ACTIONS_STEP_DEBUG="false"

      When call exec_debug "echo 'hello'"
      The output should equal "hello"
      The status should be success
    End

    It 'runs dry'
      When call exec_debug "echo 'hello'"
      The output should equal "> echo 'hello'"
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
