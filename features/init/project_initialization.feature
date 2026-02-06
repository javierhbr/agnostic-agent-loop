Feature: Project Initialization
  As a developer starting a new project
  I want to initialize the agentic project structure
  So that I have the proper directories and files to begin work

  Scenario: Initialize a new project successfully
    Given a clean test environment
    When I run "init BeginnerProject"
    Then the command should succeed
    And the following directories should exist:
      | directory                   |
      | .agentic/tasks              |
      | .agentic/context            |
      | .agentic/spec               |
      | .agentic/agent-rules        |
    And the following files should exist:
      | file                              |
      | agnostic-agent.yaml               |
      | .agentic/tasks/backlog.yaml       |
      | .agentic/tasks/in-progress.yaml   |
      | .agentic/tasks/done.yaml          |

  Scenario: Initialize project with project structure verification
    Given a clean test environment
    When I initialize a project with name "TestProject"
    Then the command should succeed
    And the project structure should be created
    And git should be initialized
