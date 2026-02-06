Feature: Intermediate Workflow - Feature with Full Metadata
  As a developer building features
  I want to create tasks with complete metadata
  So that I have clear requirements and traceability

  @tutorial @intermediate
  Scenario: Create and complete feature with acceptance criteria
    # Step 1: Initialize project with metadata
    Given a clean test environment
    When I initialize a project with name "BlogAPI"
    Then the command should succeed
    And the project structure should be created

    # Step 2: Create task with full metadata
    When I create a task with title "Implement Blog Post API"
    And I set task description to "Create REST API endpoints for blog post CRUD operations"
    And I add the following acceptance criteria:
      | criterion                          |
      | GET /posts returns all posts       |
      | POST /posts creates a new post     |
      | PUT /posts/:id updates a post      |
      | DELETE /posts/:id deletes a post   |
      | All endpoints have tests           |
    And I add the following outputs:
      | output                      |
      | src/api/posts.go           |
      | src/models/post.go         |
      | tests/api/posts_test.go    |
    Then the task should have 5 acceptance criteria

    # Step 3: Claim and work on the task
    When I claim the task
    Then the command should succeed

    # Step 4: Create output files
    When I create the following files:
      | file                       | content                        |
      | src/api/posts.go          | package api // Post handlers   |
      | src/models/post.go        | package models // Post model   |
      | tests/api/posts_test.go   | package api_test // Tests      |

    # Step 5: Complete the task
    When I complete the task
    Then the command should succeed
    And the task should be in "done" state
    And the task should preserve all acceptance criteria
    And the task should preserve all outputs
