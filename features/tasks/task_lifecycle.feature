Feature: Task Lifecycle Management
  As a developer using the agentic-agent CLI
  I want to move tasks through their lifecycle stages
  So that I can track work progress from backlog to completion

  Scenario: Successfully claim a task from backlog
    Given a clean test environment
    And I have initialized a project
    When I create a task with title "Implement authentication"
    And I claim the task
    Then the command should succeed
    And the task should be in "in-progress" state
    And the task should not exist in backlog

  Scenario: Complete an in-progress task
    Given a clean test environment
    And I have initialized a project
    When I create a task with title "Fix bug in login"
    And I claim the task
    And I complete the task
    Then the command should succeed
    And the task should be in "done" state
    And the task should not exist in "in-progress" state

  Scenario: Multiple tasks in different states
    Given a clean test environment
    And I have initialized a project
    When I create a task with title "First task"
    When I create a task with title "Second task"
    And I list all tasks
    Then I should see 2 tasks in backlog
    And I should see 0 tasks in progress
