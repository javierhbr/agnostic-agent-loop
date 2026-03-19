# Platform vs Component Boundary

The unified SDD methodology draws a hard line between platform-level work and
component-level work. This boundary determines which tools are allowed, which
artifacts are produced, and how ownership flows from shared intent to local
implementation. Understanding this boundary is essential before working with
change packages, federated intake, or the ownership artifacts.

---

## The Core Boundary Rule

Platform-side work is cross-cutting. It frames shared context, defines
contracts, sets governance rules, and reasons about multi-team impact. At this
level, three SDD tools collaborate:

- **BMAD** -- context framing, routing depth, quality rules
- **OpenSpec** -- proposals, specs, design, tasks
- **Speckit** -- shared planning discipline, cross-team coordination

Component repositories have a different job. They pin platform alignment,
define local behavior, produce a local design and task breakdown, and execute
the change. That job maps directly to OpenSpec, so the rule is:

> Inside a component repository, use OpenSpec only. No Speckit. No BMAD.

The rationale is documented in ADR-011. Allowing all three tools inside a
component repo creates competing sources of intent. Teams stop knowing which
skill owns the local change package, artifacts drift between styles, and
onboarding becomes harder. Restricting to OpenSpec gives every component repo
one simple, repeatable workflow.

---

## The Change Package as Canonical Unit

Every approved change produces exactly one change package (ADR-002). The change
package is the standard container for planning, execution, validation, and
status tracking. It is small enough to execute but large enough to hold
everything a team needs.

A change package contains:

- `proposal.md` -- what the change is and why it matters
- delta specs -- precise modifications to existing capability or contract specs
- `design.md` -- architectural decisions, constraints, component mapping
- `tasks.md` -- implementation breakdown with acceptance criteria
- PR references -- links to implementation pull requests
- archive state -- final status, learnings, and closure record

Change packages live in `.sdd-spec/changes/<slug>/` within the repository
where the work executes. Each package links upward to its parent initiative
(when one exists) and outward to affected capability specs, contracts, and
ADRs.

---

## Federated Intake

Work does not have to start from one place. ADR-001 establishes a federated
intake model: new work can enter the methodology from any of three entry
points.

- **Platform entry** -- an architect identifies a cross-cutting initiative
- **Product entry** -- a product manager submits a requirement
- **Component entry** -- a team proposes a local change

All three entry points normalize into the same canonical workflow before
planning and implementation begin. The intake router classifies each request
by scope, risk, and affected teams, then routes it into the standard five-phase
pipeline. Entry artifacts remain valid records of origin, but none of them
become the sole source of truth. The change package does.

This model preserves team speed for small local changes while maintaining
governance and traceability for cross-cutting work.

---

## Three DDD Ownership Artifacts

ADR-014 introduces three artifacts derived from Domain-Driven Design concepts.
They use plain-English names so teams adopt them without DDD training. All
three live in the platform master repository under `ownership/`.

### Component Ownership Boundaries

**File:** `ownership/component-ownership-<name>.md` (one per component)

Each component gets a file that records three things:

1. What the component owns -- its authoritative capabilities and data
2. What the component explicitly does NOT own -- prevents scope drift
3. Which contracts it publishes and which it consumes

Ownership boundary files are created during the Platform phase and updated only
through a platform change package. During Assess, the Team Lead reads the
relevant ownership files to determine whether a change stays within one
component or crosses an ownership boundary. That classification drives the
rest of the workflow.

### Dependency Map

**File:** `ownership/dependency-map.md` (one per platform)

The dependency map records which component depends on which other component,
what the dependency is (contract, event, shared data), and which impact tier
applies. Three tiers replace informal judgment:

| Tier | Meaning | Action |
|------|---------|--------|
| `must_change_together` | Shared contract requires coordinated releases | Open coordinated epics immediately |
| `watch_for_breakage` | Breaking contract change will break the consumer | Add watch note in `platform-ref.yaml` |
| `adapts_independently` | Loose coupling; consumer absorbs change on own schedule | No additional coordination needed |

The dependency map is read during Assess to populate the impact fields in
`platform-ref.yaml`. Tier 1 dependencies become hard constraints in
`design.md` -- they cannot be planned around or deferred.

### Shared Glossary

**File:** `ownership/glossary.md` (one per platform)

The glossary records shared terms with precise definitions. Each term includes
a "what it is NOT" clause to prevent the common problem where two teams use the
same word with different meanings (for example, "customer" in auth-service vs
"customer" in profile-service).

Constitution rule O-2 requires that all terms in proposals and delta specs
appear in the glossary. If a term is missing, the team must add it before
Specify can proceed. This prevents terminology drift from silently introducing
misalignment between specs written by different teams.

---

## Platform Truth Model

The methodology distributes truth across three layers. Each layer owns a
specific kind of truth, and none of them replaces the others.

**Truth Layer (platform repository):**
Canonical shared truth. Capability specs, contract specs, ownership boundaries,
the dependency map, and the glossary live here. This is what the architecture
actually says.

**Tracking Layer (JIRA):**
Workflow and ownership truth. Epics, stories, status transitions, and
assignment records live here. This is who is doing what and where it stands.

**Execution Layer (git + pull requests):**
Implementation truth. Branches, commits, code reviews, CI results, and
deployment records live here. This is what was actually built and shipped.

The three layers connect through traceability links. `jira-traceability.yaml`
records the JIRA chain. `platform-ref.yaml` pins the platform version.
PR references in the change package close the loop between execution and
specification.

---

## The Handoff Pattern

The boundary between platform and component is crossed through a structured
handoff, not an informal hand-wave.

**Step 1: Platform Plan sets shared boundaries.**
The platform architect produces the cross-cutting plan, including ownership
boundaries, impact tiers, glossary terms, and coordination requirements. This
work uses all three tools (BMAD, OpenSpec, Speckit) and lives in the platform
repository.

**Step 2: Component repo turns boundaries into local OpenSpec artifacts.**
The component team creates a local change package using OpenSpec only. The
proposal references the platform plan. Delta specs describe local behavior
changes. The design pins constraints from the platform impact assessment.

**Step 3: `platform-ref.yaml` pins platform version and references.**
This file records which version of the platform truth the component is aligned
to, which ownership boundary applies, which impact tier was assessed, and which
glossary terms are in use. It makes alignment auditable without requiring
anyone to read the full platform repository.

**Step 4: `jira-traceability.yaml` records the JIRA chain.**
This file links the local change package to the JIRA epic chain derived from
the impact tier lookup. Tier 1 dependencies produce coordinated epics. Tier 2
dependencies produce watch notes. Tier 3 dependencies produce no forced
coordination.

The handoff is complete when the component change package has a valid
`platform-ref.yaml` and `jira-traceability.yaml`. From that point forward,
the component team works entirely within OpenSpec until the change is
delivered and archived.

---

## Summary

The platform-component boundary exists to prevent competing sources of intent
inside component repositories. Platform work uses the full tool suite to govern
shared truth. Component work uses OpenSpec only to execute local changes.
The change package is the canonical unit that crosses this boundary. Three
ownership artifacts (boundaries, dependency map, glossary) make the Assess
phase deterministic. And the handoff pattern -- mediated by `platform-ref.yaml`
and `jira-traceability.yaml` -- ensures traceability without coupling the
two levels of work together.
