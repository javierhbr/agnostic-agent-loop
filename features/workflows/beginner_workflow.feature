Feature: Beginner Workflow - First Project Setup
  As a developer new to agentic-agent
  I want to follow a step-by-step workflow
  So that I can learn the basics of task management

  @tutorial @beginner
  Scenario: Complete beginner workflow from CLI tutorial
    # Step 1: Initialize project
    Given a clean test environment
    When I initialize a project with name "BeginnerProject"
    Then the command should succeed
    And the project structure should be created

    # Step 2: Create sample task
    When I create a task with title "My First Task"
    Then a task should be created successfully
    And the task should appear in the backlog

    # Step 3: List tasks
    When I list all tasks
    Then I should see 1 task in backlog
    And I should see 0 tasks in progress

    # Step 4: Claim the task
    When I claim the task
    Then the command should succeed
    And the task should move to in-progress
    And the task status should be "in-progress"

    # Step 5: Complete the task
    When I complete the task
    Then the command should succeed
    And the task should move to done
    And the task status should be "done"

    # Step 6: Run validation
    When I run validation
    Then the validation should complete successfully
