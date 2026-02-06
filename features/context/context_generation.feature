Feature: Context Generation
  As a developer working on a codebase
  I want to generate context files
  So that I have documentation for each module

  Scenario: Generate context for a directory
    Given a clean test environment
    And I have initialized a project
    And I have the following directory structure:
      | directory     |
      | src/api       |
      | src/models    |
    When I run context generation for "src/api"
    Then the command should succeed
    And a context file should exist at "src/api/context.md"

  Scenario: Context file contains expected sections
    Given a clean test environment
    And I have initialized a project
    And I have created directory "src/utils"
    And I have created file "src/utils/helper.go" with content "package utils"
    When I run context generation for "src/utils"
    Then the command should succeed
    And the context file should contain "# Context"
