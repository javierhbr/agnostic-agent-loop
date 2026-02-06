Feature: Gemini Skill Generation
  As a developer using the agentic-agent CLI
  I want to generate Gemini-compatible configuration files
  So that I can use Gemini CLI as my AI coding assistant with the framework

  Background:
    Given a clean test environment
    And I have initialized a project

  Scenario: Generate Gemini rules file
    When I generate skills for "gemini"
    Then the command should succeed
    And the file ".gemini/GEMINI.md" should exist
    And the file ".gemini/GEMINI.md" should contain "Agnostic Agent Rules"
    And the file ".gemini/GEMINI.md" should contain "Gemini-Specific Rules"
    And the file ".gemini/GEMINI.md" should contain "agentic-agent task claim"

  Scenario: Generate Gemini slash command files
    When I generate Gemini slash commands
    Then the command should succeed
    And the file ".gemini/commands/prd/gen.toml" should exist
    And the file ".gemini/commands/ralph/convert.toml" should exist
    And the file ".gemini/commands/prd/gen.toml" should contain "Product Requirements Document"
    And the file ".gemini/commands/ralph/convert.toml" should contain "Convert a PRD to YAML"

  Scenario: Gemini rules include base rules
    Given I have custom base rules:
      | rule                                |
      | Always write tests before code      |
      | Keep functions under 20 lines       |
    When I generate skills for "gemini"
    Then the command should succeed
    And the file ".gemini/GEMINI.md" should contain "Always write tests before code"
    And the file ".gemini/GEMINI.md" should contain "Keep functions under 20 lines"

  Scenario: Gemini appears in skill registry
    Then the skill registry should contain "gemini"
    And the skill registry should contain "claude-code"
    And the skill registry should contain "cursor"

  Scenario: Generate all skills includes Gemini
    When I generate skills for all tools
    Then the command should succeed
    And the file ".gemini/GEMINI.md" should exist
    And the file "CLAUDE.md" should exist

  Scenario: Detect drift in Gemini rules file
    When I generate skills for "gemini"
    And I modify the file ".gemini/GEMINI.md" with "tampered content"
    And I check for skill drift
    Then drift should be detected for ".gemini/GEMINI.md"

  Scenario: Detect missing Gemini rules file as drift
    When I check for skill drift
    Then drift should be detected for ".gemini/GEMINI.md"

  Scenario: Gemini directory structure is created automatically
    When I generate skills for "gemini"
    Then the following directories should exist:
      | directory             |
      | .gemini               |

  Scenario: Gemini slash commands use configured PRD path
    When I generate Gemini slash commands with PRD path "docs/requirements/"
    Then the command should succeed
    And the file ".gemini/commands/prd/gen.toml" should contain "docs/requirements/"
