Feature: Advanced Workflow - Task Decomposition
  As a developer working on large features
  I want to decompose tasks into subtasks
  So that I can manage complexity and track progress

  @tutorial @advanced
  Scenario: Decompose large feature into subtasks
    # Step 1: Initialize project
    Given a clean test environment
    When I initialize a project with name "AdvancedProject"
    Then the command should succeed
    And the project structure should be created

    # Step 2: Create a large feature task
    When I create a task with title "Implement Complete User Management System"
    Then a task should be created successfully
    And the task should appear in the backlog

    # Step 3: Decompose into subtasks
    When I decompose the task into the following subtasks:
      | subtask                                    |
      | Create user model and database schema      |
      | Implement user registration endpoint       |
      | Implement user authentication              |
      | Add password reset functionality           |
      | Write comprehensive tests                  |
    Then the command should succeed
    And the task should have 5 subtasks
    And the task should remain in backlog

    # Step 4: Verify subtasks were created
    When I list all tasks
    Then I should see 1 task in backlog
