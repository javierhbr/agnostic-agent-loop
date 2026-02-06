Feature: Error Handling
  As a CLI user
  I want meaningful error messages
  So that I can fix issues quickly

  Scenario: Claim nonexistent task
    Given a clean test environment
    And I have initialized a project
    When I try to claim task "TASK-NONEXISTENT"
    Then the command should fail
    And the error message should contain "not found"

  Scenario: Decompose task adds more subtasks
    Given a clean test environment
    And I have initialized a project
    When I create a task with title "Task to Decompose"
    And I decompose the task into the following subtasks:
      | subtask    |
      | Subtask 1  |
    And I decompose the task into the following subtasks:
      | subtask    |
      | Subtask 2  |
    Then the command should succeed
    And the task should have 2 subtasks
