Feature: Agent Pack Installation
  As a developer using the agentic-agent CLI
  I want the agentic-helper agent to be installed automatically
  So that Claude Code gets a pre-configured assistant for CLI workflows

  Background:
    Given a clean test environment
    And I have initialized a project

  Scenario: Installing agentic-helper for claude-code creates agent file
    When I install the "agentic-helper" pack for "claude-code"
    Then the command should succeed
    And the file ".claude/skills/agentic-helper.md" should exist
    And the file ".claude/skills/agentic-helper.md" should contain "name: agentic-helper"
    And the file ".claude/skills/agentic-helper.md" should contain "tools: Read, Write, Edit, Bash"
    And the file ".claude/skills/agentic-helper/SKILL.md" should exist

  Scenario: Agent file contains the workflow decision tree
    When I install the "agentic-helper" pack for "claude-code"
    Then the file ".claude/skills/agentic-helper.md" should contain "Workflow Decision Tree"
    And the file ".claude/skills/agentic-helper.md" should contain "Full SDD"
    And the file ".claude/skills/agentic-helper.md" should contain "openspec init"

  Scenario: Agent file contains CLI command examples
    When I install the "agentic-helper" pack for "claude-code"
    Then the file ".claude/skills/agentic-helper.md" should contain "task claim"
    And the file ".claude/skills/agentic-helper.md" should contain "context generate"
    And the file ".claude/skills/agentic-helper.md" should contain "agentic-agent validate"

  Scenario: Installing agentic-helper for cursor falls back to skill directory
    When I install the "agentic-helper" pack for "cursor"
    Then the command should succeed
    And the file ".cursor/skills/agentic-helper/SKILL.md" should exist

  Scenario: agentic-helper is installed during skills ensure
    When I run skills ensure for "claude-code"
    Then the file ".claude/skills/agentic-helper.md" should exist
    And the file ".claude/skills/agentic-helper.md" should not be empty

  Scenario: agentic-helper agent file uses sonnet model
    When I install the "agentic-helper" pack for "claude-code"
    Then the file ".claude/skills/agentic-helper.md" should contain "model: sonnet"

  Scenario: agentic-helper agent file enables project memory
    When I install the "agentic-helper" pack for "claude-code"
    Then the file ".claude/skills/agentic-helper.md" should contain "memory: project"

  Scenario: agentic-helper skill file explains when to use it
    When I install the "agentic-helper" pack for "claude-code"
    Then the file ".claude/skills/agentic-helper/SKILL.md" should contain "Use this skill when"
    And the file ".claude/skills/agentic-helper/SKILL.md" should contain "agentic-agent CLI"

  Scenario: agentic-helper agent file contains error recovery guidance
    When I install the "agentic-helper" pack for "claude-code"
    Then the file ".claude/skills/agentic-helper.md" should contain "Error Recovery"
    And the file ".claude/skills/agentic-helper.md" should contain "Red Flags"

  Scenario: agentic-helper is marked as mandatory
    When I run the installed agentic-agent CLI
    Then "agentic-helper" should be in the mandatory packs list
