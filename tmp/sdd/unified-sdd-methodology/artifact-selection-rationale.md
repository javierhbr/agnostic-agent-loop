# Artifact Selection Rationale

This document explains why each artifact is chosen from each tool at each phase
of the unified SDD methodology. For each phase it answers:

- which artifact is used
- which tool produces it
- why that artifact from that tool — not a different one

---

## The core boundary rule

Before reading the per-phase rationale, this rule must be clear:

```
Platform-side work (cross-cutting):  Speckit + OpenSpec + BMAD  (combined)
Component repos (local work):         OpenSpec ONLY
```

Every artifact selection below follows this rule.

---

## Phase 1: Platform

**Owner:** Architect | **Support:** Product, Team Lead

### Artifact: Constitution / platform principles
**Tool:** Speckit — `constitution` command + `rules/constitution-rules.md`

**Why this artifact from this tool:**

Speckit is designed to turn vague values into explicit, testable principles.
No other tool does this. BMAD can describe what a system looks like, but it
does not produce governance-level rules. OpenSpec can encode rules in config,
but it needs the rules to exist before it can store them.

The constitution must be produced first because every artifact in every later
phase traces back to it. If the quality bar, security expectations, and testing
standards are not explicit here, they will be invented inconsistently in
individual change packages.

Speckit's constitution workflow enforces that principles are:
- stated as specific behaviors, not aspirational phrases
- linked to a testable outcome
- short enough to be followed

A constitution from a different tool would likely produce a list of values
rather than enforced rules.

---

### Artifact: Project config / durable context
**Tool:** OpenSpec — `config.yaml` + `rules/project-config-template.yaml`

**Why this artifact from this tool:**

OpenSpec is designed to separate durable context from change-specific detail.
This is exactly what platform config is. Platform config should not change with
every feature. It should encode the stable rules, versioning conventions, and
artifact expectations that every change inherits.

BMAD could describe context narratively, but it does not produce a reusable,
structured config file. Speckit produces principles, not configuration.

OpenSpec's `config.yaml` is used here because:
- it survives multiple change cycles without modification
- it is machine-readable and agent-readable
- it becomes the reference that component repos inherit from
- it separates what is stable from what changes per request

If context were written in a narrative BMAD planning doc, teams would duplicate
it in every change package and it would drift.

---

### Artifact: Role and context framing
**Tool:** BMAD — Architect agent + brownfield/project-context playbook

**Why this artifact from this tool:**

BMAD is used third in Platform because its job here is to frame the operating
model for all later routing decisions. BMAD's role-based context model ensures
that the platform baseline explicitly records who owns what, how decisions
escalate, and how the team expects future changes to flow.

This is not an artifact in the traditional sense. It is the framing that ensures
the Assess phase can classify work correctly. Without BMAD's project-context
contribution, teams lack a shared model for what counts as Quick Flow versus
architecture-heavy work.

Speckit and OpenSpec do not produce this role and routing framing. That is
BMAD's contribution.

---

### Artifact: Component ownership boundary files
**Tool:** Architect (human-authored) — `ownership/component-ownership-<name>.md`

**Why this artifact at this phase:**

Without explicit ownership records, two problems recur in every Assess step:
scope bleeds across components because nobody wrote down where one component's
responsibility ends, and impact assessment relies on informal memory that varies
by Team Lead.

The ownership boundary file records what each component owns, what it explicitly
does NOT own, and which contracts it publishes or consumes. It is written once
during Platform and updated only when the platform structure changes.

This is not produced by BMAD, OpenSpec, or Speckit because it is not a
change-package artifact — it is a durable platform truth artifact. Its home is
the platform repo alongside the platform baseline and shared contracts.

In Assess, the Team Lead reads it to confirm the correct epic owner before any
JIRA issue is created. This makes ownership classification a lookup, not a
judgment call.

See: [platform-ddd-spec.md](platform-ddd-spec.md) — Concept 1

---

### Artifact: Dependency map
**Tool:** Architect (human-authored) — `ownership/dependency-map.md`

**Why this artifact at this phase:**

Impact assessment in Assess has historically been the most variable step in the
methodology. Different Team Leads classify the same type of change differently,
leading to missing JIRA epics, uncoordinated releases, and silent breakage in
consuming components.

The dependency map replaces that variable judgment with a deterministic lookup.
It records which components depend on which other components, what the
dependency is (contract, event, data), and which impact tier applies:

- `must_change_together` — both must ship in the same release
- `watch_for_breakage` — consumer may break; monitor after deploy
- `adapts_independently` — no forced coordination

The tier directly drives the JIRA structure at Assess time. This connection
makes the JIRA decision mechanical: read the tier, apply the rule.

This is not produced by BMAD, OpenSpec, or Speckit because it is a platform
governance artifact, not a change-package artifact. It must exist before any
change starts and it must not be recreated per change.

See: [platform-ddd-spec.md](platform-ddd-spec.md) — Concept 2

---

### Artifact: Shared glossary
**Tool:** Architect (human-authored, team-maintained) — `ownership/glossary.md`

**Why this artifact at this phase:**

Spec ambiguity almost always originates from undefined or overloaded terms.
Product and Engineering may use the same word ("validation", "customer",
"profile") with different meanings. The mismatch surfaces in review cycles or
in production, not during Specify.

The shared glossary is seeded during Platform — when shared contracts and
capabilities are defined — because that is when terms first acquire precise
meaning. Once seeded, the glossary grows incrementally: a term is added whenever
a spec review surfaces ambiguity.

The critical rule: every term used in a proposal, delta spec, or acceptance
criterion must appear in the glossary. This constraint belongs at Platform
because it must be in place before the first Specify phase can run.

This is not produced by Speckit (which produces governance rules) or OpenSpec
(which manages change packages). The glossary is a durable platform artifact
with a different lifecycle — it lives in the platform repo, not inside any
component change package.

See: [platform-ddd-spec.md](platform-ddd-spec.md) — Concept 3

---

## Phase 2: Assess

**Owner:** Team Lead | **Support:** Product, Architect, Engineering Manager

### Artifact: Change package (scope + classification + path selection)
**Tool:** BMAD — routing agent + `rules/track-selection-rules.md`

**Why this artifact from this tool:**

BMAD is the only tool with explicit rules for classifying work by:
- greenfield vs brownfield
- size (small / medium / large)
- planning depth required (Quick Flow / PRD / architecture-heavy)

No other tool answers the question "how much workflow does this change need?"
OpenSpec manages what artifacts to produce once the path is chosen. Speckit
clarifies ambiguity but does not classify work. BMAD's track selection is the
routing decision that prevents teams from spending architecture effort on small
changes and from undershooting on large shared ones.

The output is not a document. It is the classification decision that determines
which artifacts come next and who owns them.

---

### Artifact: Change package scaffold + `platform-ref.yaml` + `jira-traceability.yaml`
**Tool:** OpenSpec — `/opsx:explore` or `/opsx:propose` + `rules/artifact-rules.md`

**Why this artifact from this tool:**

OpenSpec's change package is the canonical execution unit. Once BMAD has
classified the request, OpenSpec gives it a home. The change package becomes
the container that holds every downstream artifact: proposal, specs, design,
tasks, PR links, and archive state.

`platform-ref.yaml` is an OpenSpec convention. It pins the platform version and
platform refs the change is accountable to. Without this file, there is no
machine-readable record of which platform contracts the change respects.

`jira-traceability.yaml` is also an OpenSpec convention. It records the JIRA
issue chain: platform issue → component epic → stories. Without it, JIRA and
the change package exist independently and drift.

BMAD produces planning artifacts, not change containers. Speckit produces
governance artifacts, not traceability files. OpenSpec is the right tool for
this because the change package must survive through Deliver and Archive.

---

### Artifact: Clarify pass (optional)
**Tool:** Speckit — `clarify` command

**Why this artifact from this tool, and why it is optional:**

Speckit's clarify command is the right tool when the request is too vague to
classify safely. It surfaces hidden assumptions, missing actors, and undefined
failure behavior.

This artifact is optional in Assess because most requests can be classified
from BMAD's routing rules alone. It is only needed when ambiguity would cause
a wrong path selection or a change package that cannot be scoped.

If clarify is used here, it is lightweight. A deep clarify pass belongs in
Specify, not Assess.

---

## Phase 3: Specify

**Owner:** Product | **Support:** Team Lead, Architect, Developers

### Artifact: `proposal.md`
**Tool:** OpenSpec — `artifact-rules.md`

**Why this artifact from this tool:**

`proposal.md` is the first artifact inside the component change package. It
defines the problem statement, goals, non-goals, and acceptance summary in
plain language.

OpenSpec is the only tool used inside component repos. The proposal must live
inside the change package because everything downstream — specs, design, tasks,
PR — traces back to it. A proposal written in a BMAD doc or Speckit spec would
not be inside the change package and would not survive the archive step.

OpenSpec's artifact structure enforces that the proposal separates goals from
non-goals, which prevents the change from expanding during execution.

---

### Artifact: Delta specs — `ADDED`, `MODIFIED`, `REMOVED` sections
**Tool:** OpenSpec — `rules/artifact-rules.md`

**Why this artifact from this tool:**

Delta specs capture the behavioral change in explicit, testable language. The
`ADDED / MODIFIED / REMOVED` structure forces teams to say exactly what
changes, not just what the new behavior is. This makes the spec comparable to
current reality and reviewable before any implementation begins.

This is an OpenSpec-native artifact format. Neither BMAD nor Speckit produce
delta specs. BMAD produces planning and architecture descriptions. Speckit
produces governance rules and clarification outputs.

The delta format is critical because:
- it shows what changes, not just what exists
- it is explicit enough for acceptance test generation
- it stays inside the change package and links to platform refs

---

### Artifact: Clarify pass (platform-side only, optional)
**Tool:** Speckit — `clarify` command

**Why this artifact from this tool, and why it is platform-side only:**

When a shared platform change is involved, Speckit may be used upstream to
expose ambiguity before the component package is approved. This is the last
point where BMAD and Speckit are permitted in the flow.

Once the work enters the component repo, only OpenSpec is used. This is not
a preference — it is an architectural boundary rule. Mixing BMAD or Speckit
artifacts into the local component change package creates two competing
sources of intent.

---

## Phase 4: Plan

**Owner:** Architect | **Support:** Team Lead, Developers, Product

### Artifact: `design.md`
**Tool:** OpenSpec — `rules/artifact-rules.md`

**Why this artifact from this tool:**

`design.md` is the technical execution model. It records architecture decisions,
service boundaries, data flow, dependencies, and rollout constraints. It lives
inside the change package because implementation must trace back to it.

OpenSpec is the only tool used inside component repos. BMAD's architecture
playbook is still available for platform-level shared design, but once the work
enters the component repo, all design artifacts belong to OpenSpec.

The key reason OpenSpec is used for `design.md` and not BMAD's architecture
output is traceability. An OpenSpec `design.md` inside the change package is
directly connected to the `proposal.md`, delta specs, and `tasks.md` that
precede and follow it. A BMAD architecture doc would live outside the change
package and require manual synchronization.

---

### Artifact: `tasks.md`
**Tool:** OpenSpec — `rules/artifact-rules.md`

**Why this artifact from this tool:**

`tasks.md` is the implementation work order. Each task maps to a story key
and has an explicit validation step. Tasks are inside the change package
because they are the units of delivery — their status, story links, PR links,
and verification notes are tracked through the same OpenSpec structure.

Speckit also has a `tasks` command, but Speckit tasks are not inside a change
package. They are standalone planning artifacts. OpenSpec tasks are the right
choice here because:
- they are inside the change package
- they link to JIRA stories
- they carry validation criteria
- they feed directly into the Deliver phase's CLI tracking

---

### Artifact: ADRs when needed
**Tool:** OpenSpec — inside the change package, or standalone in the platform repo

**Why this artifact from this tool:**

ADRs document significant technical decisions. When a decision is
component-local, it belongs inside the change package as part of the design
artifacts. When it affects shared platform truth, it belongs in the platform
repo.

The methodology does not assign ADRs to BMAD or Speckit. BMAD's architecture
playbook surfaces the need for an ADR, but the ADR itself is an OpenSpec
artifact in the component flow.

---

## Phase 5: Deliver

**Owner:** Team Lead | **Support:** Developers, QA, Architect, Product

### Artifact: Task status updates and PR traceability
**Tool:** OpenSpec — task updates, `artifact-rules.md`

**Why this artifact from this tool:**

Every task in `tasks.md` must be updated as work progresses. The same change
package that held the proposal, specs, and design now tracks execution state.
This is intentional. Archive is only possible when the change package reflects
what was actually delivered.

No other tool manages this. BMAD does not track task state. Speckit does not
track PRs. OpenSpec's change package is the only structure that carries the
full artifact chain from proposal through archive.

---

### Artifact: PR description
**Tool:** OpenSpec — `artifact-rules.md`

**Why this artifact from this tool:**

A PR description produced under OpenSpec conventions references the change
package, affected tasks, story keys, and validation performed. This is not a
freeform commit message — it is a traceable delivery artifact.

The reason this matters is audit and archive. When the change is archived, the
PR description becomes part of the delivery history. A PR description not
connected to the change package would be an orphan with no traceability.

---

### Artifact: Verification evidence
**Tool:** OpenSpec — inside the change package

**Why this artifact from this tool:**

Verification evidence (test results, log samples, rollout readiness notes) is
recorded inside the change package before archive. This keeps the evidence
co-located with the requirements it proves.

The methodology requires verification evidence before deploy. OpenSpec's change
package is the only structure that enforces this co-location.

---

### Artifact: Archive
**Tool:** OpenSpec — `/openspec-archive` command

**Why this artifact from this tool:**

Archive is the final promotion of the new truth. It merges the local delta
specs into the component's main spec directory, closes the change package, and
records the delivery as history.

This is an OpenSpec-only operation. BMAD and Speckit do not have an archive
step. The archive step is what transforms a change package from in-progress
work into the component's new canonical truth.

Skipping archive is treated as an incomplete delivery. The methodology treats
archive as a required closure step, not optional cleanup.

---

## Summary table

| Phase | Artifact | Tool | Why this tool |
|---|---|---|---|
| Platform | Constitution / principles | Speckit | Only tool that produces explicit, testable governance rules |
| Platform | Project config / durable context | OpenSpec | Only tool that separates durable context from change-specific detail in a reusable config structure |
| Platform | Role and context framing | BMAD | Only tool with a role-based context and routing model |
| Platform | Component ownership boundary files | Architect (human-authored) | Durable platform truth artifact — not a change-package artifact; makes ownership classification a lookup in every future Assess step |
| Platform | Dependency map | Architect (human-authored) | Durable platform truth artifact — three impact tiers make JIRA structure decisions deterministic; must exist before Assess runs |
| Platform | Shared glossary | Architect + team (human-authored) | Durable platform truth artifact — prevents spec ambiguity by defining shared terms with "what it is NOT" clauses; seeded at Platform, grown during Specify |
| Assess | Classification + path selection | BMAD | Only tool with explicit track-selection rules (size, type, depth) |
| Assess | Change package scaffold + `platform-ref.yaml` + `jira-traceability.yaml` | OpenSpec | Only tool that produces the canonical execution container with versioned platform alignment and JIRA traceability |
| Assess | Clarify pass (optional) | Speckit | Only tool with a structured ambiguity-surfacing step |
| Specify | `proposal.md` | OpenSpec | Must live inside the change package; only OpenSpec manages the change package |
| Specify | Delta specs | OpenSpec | OpenSpec-native format; only tool that captures behavioral change as ADDED / MODIFIED / REMOVED |
| Specify | Clarify pass (platform-side, optional) | Speckit | Last point where ambiguity can be resolved before component repo work begins |
| Plan | `design.md` | OpenSpec | Must trace back to proposal and forward to tasks inside the same change package |
| Plan | `tasks.md` | OpenSpec | Must carry story keys, validation criteria, and feed into Deliver CLI tracking |
| Plan | ADRs | OpenSpec | Significant decisions must stay inside or alongside the change package |
| Deliver | Task status + PR traceability | OpenSpec | Only tool that tracks execution state inside the change package |
| Deliver | PR description | OpenSpec | Must reference change package, tasks, and stories for archive traceability |
| Deliver | Verification evidence | OpenSpec | Must be co-located with the requirements it proves inside the change package |
| Deliver | Archive | OpenSpec | Archive is an OpenSpec-only operation that closes the change package and promotes new truth |
