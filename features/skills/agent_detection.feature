Feature: Agent Detection
  As a developer using the agentic-agent CLI
  I want the tool to detect which AI coding agent is active
  So that it configures itself correctly for my agent tool

  The detection follows a strict priority order:
  1. Explicit --agent flag
  2. Environment variables (AGENTIC_AGENT first, then agent-specific)
  3. Filesystem heuristics (marker files/directories in the project root)

  # --- Flag Detection (highest priority) ---

  Scenario: Explicit flag overrides environment variable
    Given the environment variable "AGENTIC_AGENT" is set to "gemini"
    When I detect the agent with flag "claude-code"
    Then the detected agent name should be "claude-code"
    And the detected agent source should be "flag"

  Scenario: Explicit flag overrides filesystem markers
    Given a project root with directory ".cursor"
    When I detect the agent with flag "windsurf"
    Then the detected agent name should be "windsurf"
    And the detected agent source should be "flag"

  Scenario: Flag accepts any arbitrary agent name
    When I detect the agent with flag "my-custom-agent"
    Then the detected agent name should be "my-custom-agent"
    And the detected agent source should be "flag"

  # --- Environment Variable Detection ---

  Scenario: AGENTIC_AGENT takes priority over agent-specific env vars
    Given the environment variable "AGENTIC_AGENT" is set to "gemini"
    And the environment variable "CLAUDE" is set to "1"
    And the environment variable "CURSOR_SESSION" is set to "abc"
    When I detect the agent without a flag
    Then the detected agent name should be "gemini"
    And the detected agent source should be "env"

  Scenario Outline: Agent-specific environment variables
    Given the environment variable "<env_var>" is set to "<value>"
    When I detect the agent without a flag
    Then the detected agent name should be "<expected_agent>"
    And the detected agent source should be "env"

    Examples:
      | env_var          | value   | expected_agent |
      | CLAUDE           | 1       | claude-code    |
      | CLAUDE_CODE      | 1       | claude-code    |
      | CURSOR_SESSION   | abc123  | cursor         |
      | GEMINI_CLI       | 1       | gemini         |
      | WINDSURF_SESSION | sess-42 | windsurf       |
      | CODEX_SANDBOX    | 1       | codex          |

  Scenario: Environment variables override filesystem markers
    Given a project root with directory ".cursor"
    And the environment variable "GEMINI_CLI" is set to "1"
    When I detect the agent without a flag
    Then the detected agent name should be "gemini"
    And the detected agent source should be "env"

  # --- Filesystem Heuristic Detection ---

  Scenario Outline: Filesystem markers detect the correct agent
    Given a project root with <marker_type> "<marker>"
    When I detect the agent without a flag
    Then the detected agent name should be "<expected_agent>"
    And the detected agent source should be "filesystem"

    Examples:
      | marker_type | marker    | expected_agent |
      | directory   | .claude   | claude-code    |
      | file        | CLAUDE.md | claude-code    |
      | directory   | .cursor   | cursor         |
      | directory   | .gemini   | gemini         |
      | directory   | .windsurf | windsurf       |
      | directory   | .codex    | codex          |
      | directory   | .agent    | antigravity    |

  # --- Unknown / Fallback ---

  Scenario: No agent detected returns unknown
    Given an empty project root
    When I detect the agent without a flag
    Then the detected agent name should be ""
    And the detected agent source should be "unknown"

  # --- DetectAllAgents ---

  Scenario: Detect all agents from filesystem markers
    Given a project root with directory ".claude"
    And a project root with directory ".cursor"
    And a project root with directory ".gemini"
    And a project root with directory ".windsurf"
    When I detect all agents
    Then the detected agents should include "claude-code"
    And the detected agents should include "cursor"
    And the detected agents should include "gemini"
    And the detected agents should include "windsurf"

  Scenario: Duplicate filesystem markers produce no duplicate agents
    Given a project root with directory ".claude"
    And a project root with file "CLAUDE.md"
    When I detect all agents
    Then the detected agents should include "claude-code"
    And the detected agent count should be 1

  Scenario: Empty project detects no agents
    Given an empty project root
    When I detect all agents
    Then the detected agent count should be 0
