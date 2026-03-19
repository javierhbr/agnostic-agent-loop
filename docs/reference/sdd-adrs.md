# SDD Methodology ADR Index

Reference summary of all 14 Architecture Decision Records from the unified SDD methodology.

Source files: `tmp/sdd/unified-sdd-methodology/decisions/`

---

## Groupings

- **Foundation (ADR-001 to ADR-003):** Core workflow identity and accountability model.
- **Routing (ADR-004):** How changes are classified and routed.
- **Workflow Structure (ADR-005 to ADR-007):** Phase model, adoption sequence, and delivery gates.
- **Agent Model (ADR-008 to ADR-009):** Role-specific agents and brownfield platform handling.
- **Platform-Component (ADR-010 to ADR-011):** Canonical truth, repo boundaries, and spec ownership.
- **Infrastructure (ADR-012 to ADR-013):** Local MCP gateway implementation.
- **Ownership (ADR-014):** DDD-derived concepts for platform governance.

---

## Foundation

| ADR | Title | Decision | Key Consequence |
|-----|-------|----------|-----------------|
| ADR-001 | Federated Intake with Single Canonical Workflow | Allow work to start from platform, product, or component entry points but normalize to a single workflow. | Needs an intake router; one canonical execution unit per change. |
| ADR-002 | Change Package as Canonical Unit | One change package per approved change is the execution container. | Holds all artifacts: proposal, specs, design, tasks, PR refs, archive state. |
| ADR-003 | Stage-Based Role Accountability | Each phase has a clear owner with defined support roles. | Prevents diffused responsibility and unclear handoffs. |

## Routing

| ADR | Title | Decision | Key Consequence |
|-----|-------|----------|-----------------|
| ADR-004 | Separate Size from Impact in Routing | Size (small/medium/large) and impact (component-only/shared) are distinct routing dimensions. | Prevents undershooting shared changes and over-engineering small local ones. |

## Workflow Structure

| ADR | Title | Decision | Key Consequence |
|-----|-------|----------|-----------------|
| ADR-005 | Five-Phase v1 Model with Team Lead Owning Deliver | Use Platform, Assess, Specify, Plan, Deliver as the five phases. | Easier adoption; room to split Deliver later. |
| ADR-006 | Two-Iteration Adoption Plan | Iteration 1 covers Platform + Assess + Specify; Iteration 2 covers Plan + Deliver. | Learn to define work well before learning to deliver well. |
| ADR-007 | Make PR Creation and Review Explicit in Deliver | Pull requests are an explicit gate, not an implementation detail. | Makes code review a structured phase boundary. |

## Agent Model

| ADR | Title | Decision | Key Consequence |
|-----|-------|----------|-----------------|
| ADR-008 | Role-Specific Agents on Top of Unified Skill | Agents are role-aligned helpers; teams own approval. | Agents do not replace governance. |
| ADR-009 | Platform Contextualizer for Existing Platforms | Capture brownfield state before starting the formal Platform phase. | Makes platform baseline explicit instead of relying on implicit memory. |

## Platform-Component

| ADR | Title | Decision | Key Consequence |
|-----|-------|----------|-----------------|
| ADR-010 | Canonical Platform Truth with Versioned Component Alignment | One master platform repo is canonical; components pin version; JIRA tracks delivery. | Prevents drift; clear source of truth. |
| ADR-011 | Component Repos Use OpenSpec Only | No BMAD or Speckit inside component repos. | Prevents competing sources of intent. |

## Infrastructure

| ADR | Title | Decision | Key Consequence |
|-----|-------|----------|-----------------|
| ADR-012 | Local Platform MCP Gateway | Read-only MCP gateway for querying platform truth locally. | No server required; v1 is Go stdio. |
| ADR-013 | Go Stdio Platform MCP v1 | Implement platform MCP as a Go stdio server. | Simple implementation with no external dependencies. |

## Ownership

| ADR | Title | Decision | Key Consequence |
|-----|-------|----------|-----------------|
| ADR-014 | Three DDD-Derived Concepts for Platform Ownership and Impact | Adopt component ownership boundaries, dependency map, and shared glossary. | Makes Assess deterministic; prevents spec ambiguity. |

---

## Full Index (sorted by number)

| ADR | Title | Group |
|-----|-------|-------|
| ADR-001 | Federated Intake with Single Canonical Workflow | Foundation |
| ADR-002 | Change Package as Canonical Unit | Foundation |
| ADR-003 | Stage-Based Role Accountability | Foundation |
| ADR-004 | Separate Size from Impact in Routing | Routing |
| ADR-005 | Five-Phase v1 Model with Team Lead Owning Deliver | Workflow Structure |
| ADR-006 | Two-Iteration Adoption Plan | Workflow Structure |
| ADR-007 | Make PR Creation and Review Explicit in Deliver | Workflow Structure |
| ADR-008 | Role-Specific Agents on Top of Unified Skill | Agent Model |
| ADR-009 | Platform Contextualizer for Existing Platforms | Agent Model |
| ADR-010 | Canonical Platform Truth with Versioned Component Alignment | Platform-Component |
| ADR-011 | Component Repos Use OpenSpec Only | Platform-Component |
| ADR-012 | Local Platform MCP Gateway | Infrastructure |
| ADR-013 | Go Stdio Platform MCP v1 | Infrastructure |
| ADR-014 | Three DDD-Derived Concepts for Platform Ownership and Impact | Ownership |
