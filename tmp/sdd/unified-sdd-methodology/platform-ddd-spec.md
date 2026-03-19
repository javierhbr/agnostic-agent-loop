# Platform DDD Spec

This document describes how three concepts from Domain-Driven Design are applied
at the platform level to make ownership, impact, and shared language explicit.

No DDD vocabulary is exposed to the team. The three concepts are translated into
plain artifacts that slot into the existing platform phase without adding new
methodology overhead.

---

## Why three concepts, not the full DDD model

DDD solves many problems. Most of them live at the implementation level inside
a single component. The platform level has three specific problems that DDD
addresses well and that are not solved by the existing methodology alone:

```text
Problem                          Solved by
─────────────────────────────────────────────────────────────────────
"Who owns this? Who decides?"    Component ownership boundaries
"Who else is affected?"          Dependency map (impact tiers)
"Do we mean the same thing?"     Shared glossary
```

Everything else in DDD (aggregates, repositories, domain events, anti-corruption
layers, open host service) lives inside individual components. Those belong in
a component's `design.md`, not in the platform spec.

---

## The three concepts

### Concept 1: Component ownership boundary

**What it is in DDD:** Bounded Context — a named boundary inside which a model
has a clear, consistent meaning and one team is responsible for it.

**What it means for the team:** Every component repo has a clear list of what
it owns and what it does not own. When a request arrives, the ownership boundary
answers "which component is the right home for this change?" without a team
discussion every time.

**Why it matters for the platform spec:** Without explicit ownership, two things
happen. Specs bleed across components — Product writes a requirement that spans
three services without realising it. And teams duplicate work — two components
each try to own validation for the same business rule. The ownership boundary
prevents both.

**Where it lives:** `component-ownership.md` in the platform repo.

**What it looks like:**

```text
  ┌─────────────────────────────────────────────────────────────┐
  │  component-ownership.md                                     │
  │                                                             │
  │  profile-service                                            │
  │    owns:     customer profile data                          │
  │              email update flow                              │
  │              profile update event                           │
  │    does NOT own: authentication, sessions, email delivery   │
  │                                                             │
  │  auth-service                                               │
  │    owns:     login sessions                                 │
  │              duplicate email enforcement                    │
  │    does NOT own: profile data shape, event publishing       │
  │                                                             │
  │  notification-service                                       │
  │    owns:     message delivery to customers                  │
  │    does NOT own: when to send, profile state, identity      │
  └─────────────────────────────────────────────────────────────┘
```

**How it connects to the methodology:**

- **Platform phase:** Architect writes `component-ownership.md` once.
  Updated only when a new component is added or an ownership boundary shifts.
- **Assess phase:** Team Lead reads it to answer "which component owns this
  change?" before opening a change package. Prevents the wrong epic owner.
- **Specify phase:** Product reads it to know what is in and out of scope
  per component. Prevents specs that span ownership lines.
- **Constitution rule added:** Every proposal MUST identify which component
  owns the change. If ownership is unclear, resolve it before writing the spec.

---

### Concept 2: Dependency map

**What it is in DDD:** Context Map — a document that describes the relationships
between bounded contexts and the integration patterns between them.

**What it means for the team:** A plain three-tier list that answers "if I
change this component, who else is affected and how urgently?" The three tiers
replace the manual "let's think about who might break" conversation in every
Assess meeting.

**Why it matters for the platform spec:** Impact assessment in Assess is
currently a judgment call. Different team leads make different calls on the same
type of change. The dependency map makes it deterministic — the tier of a
relationship tells you exactly what coordination is required.

**The three tiers (plain language, no DDD names):**

```text
  TIER 1  Must change together
  ─────────────────────────────────────────────────────────────
  Both components share the same rule or data contract.
  If one changes, the other MUST change in the same release.
  Missing this is a production incident.

  Example: profile-service and auth-service both enforce
  the email uniqueness rule. If profile-service changes
  the rule, auth-service must change it the same way.

  TIER 2  Watch for breakage
  ─────────────────────────────────────────────────────────────
  One component produces something the other consumes.
  The producer can change independently, but the consumer
  may break silently if the producer changes the shape.
  Requires monitoring and a follow-up check after deploy.

  Example: profile-service publishes a profile update event.
  notification-service consumes it. If profile-service
  changes the event fields, notification-service may fail
  silently on the next message.

  TIER 3  Adapts independently
  ─────────────────────────────────────────────────────────────
  One component reads from another with no negotiation.
  It adapts on its own schedule without blocking the producer.
  No coordination needed unless the consumer breaks.

  Example: audit-service reads profile data for logging.
  It adapts when it can. No joint release required.
```

**Where it lives:** `dependency-map.md` in the platform repo, alongside
`component-ownership.md`.

**What it looks like:**

```text
  ┌──────────────────────────────────────────────────────────────────┐
  │  dependency-map.md                                               │
  │                                                                  │
  │  TIER 1 — Must change together                                   │
  │  ┌──────────────────┐           ┌──────────────────┐            │
  │  │  profile-service │ ◄──────── │  auth-service    │            │
  │  │                  │ ────────► │                  │            │
  │  └──────────────────┘           └──────────────────┘            │
  │  Shared rule: email uniqueness enforcement                       │
  │  Trigger: any change to email validation logic                   │
  │  Contract: contracts.customer-profile.v2                         │
  │                                                                  │
  │  TIER 2 — Watch for breakage                                     │
  │  ┌──────────────────┐           ┌──────────────────────┐        │
  │  │  profile-service │ ────────► │  notification-service│        │
  │  └──────────────────┘  event    └──────────────────────┘        │
  │  Producer: profile-service (publishes profile-updated event)     │
  │  Consumer: notification-service (reads event fields)             │
  │  Trigger: any change to contracts.customer-profile.v2            │
  │                                                                  │
  │  TIER 3 — Adapts independently                                   │
  │  ┌──────────────────┐           ┌──────────────────┐            │
  │  │  profile-service │ ────────► │  audit-service   │            │
  │  └──────────────────┘  reads    └──────────────────┘            │
  │  Reason: audit reads but never negotiates                        │
  │  Trigger: no action required unless audit breaks                 │
  └──────────────────────────────────────────────────────────────────┘
```

**How it connects to the methodology:**

- **Platform phase:** Architect writes `dependency-map.md` once. Updated when
  a new integration is added, a relationship tier changes, or an incident
  reveals a hidden dependency.
- **Assess phase — Step 2.1:** Team Lead reads the dependency map to determine
  impact tier for the change. The tier drives the JIRA structure:
  - Tier 1 → open component epics for all involved components immediately
  - Tier 2 → open the primary epic, add a watch note for tier 2 consumers
  - Tier 3 → no additional epics required
- **Assess phase — Step 2.2:** Impact tiers are written into `platform-ref.yaml`
  as a structured field so every downstream phase can see the impact without
  re-reading the dependency map.
- **Plan phase:** `design.md` explicitly references tier 1 dependencies as
  hard constraints. Tier 2 dependencies appear as rollout risks.
- **Deliver phase:** PR descriptions note whether any tier 1 or tier 2
  components have been verified as part of this slice.

**What the impact field in `platform-ref.yaml` looks like after Assess:**

```yaml
impact:
  must_change_together: [auth-service]          # tier 1
  watch_for_breakage:   [notification-service]  # tier 2
  adapts_independently: [audit-service]         # tier 3
```

---

### Concept 3: Shared glossary

**What it is in DDD:** Ubiquitous Language — a shared vocabulary that is used
consistently in code, specs, conversations, and documents within a bounded
context.

**What it means for the team:** A short plain-language dictionary. Every term
that appears in a proposal, delta spec, or acceptance criterion must be defined
here. If Product and Engineering use the same word to mean different things,
the spec is broken before implementation starts.

**Why it matters for the platform spec:** Spec ambiguity almost always comes
from undefined or overloaded terms. "Validation" means something different to
Product ("does the form look right?") and to Engineering ("format + uniqueness
+ business rule checks before write"). Without a shared glossary, both write
specs using the same word and discover the mismatch during review or in
production.

**Where it lives:** `glossary.md` in the platform repo.

**What it looks like:**

```text
  ┌──────────────────────────────────────────────────────────────────┐
  │  glossary.md                                                     │
  │                                                                  │
  │  email update                                                    │
  │    A customer request to replace the email address on an         │
  │    existing account. Not the same as account creation,           │
  │    password reset, or email verification.                        │
  │                                                                  │
  │  validation                                                      │
  │    Checking that an input meets format and business rules        │
  │    before saving anything. Happens before persistence.           │
  │    Includes: format check, uniqueness check, business rules.     │
  │    Does not include: UI-level field highlighting.                │
  │                                                                  │
  │  profile event                                                   │
  │    The notification emitted after a successful profile write.    │
  │    Shape defined by contracts.customer-profile.v2.               │
  │    Consumed by notification-service and audit-service.           │
  │                                                                  │
  │  duplicate email                                                 │
  │    An email address already registered to a different account.  │
  │    Must be rejected with EMAIL_ALREADY_EXISTS.                   │
  │    Never silently accepted or queued.                            │
  └──────────────────────────────────────────────────────────────────┘
```

**How it connects to the methodology:**

- **Platform phase:** Architect seeds `glossary.md` from the constitution work
  and the brownfield review. Many terms will already exist implicitly in
  existing documentation — the glossary makes them explicit.
- **Specify phase — Step 3.1:** Product reads `glossary.md` before writing
  `proposal.md`. Terms used in the proposal MUST come from the glossary. If a
  term is missing, add it before finalising the proposal.
- **Specify phase — Step 3.2:** Delta specs use glossary terms in requirement
  language. "The feature MUST validate the email update" is unambiguous when
  "email update" and "validation" are both defined in the glossary.
- **Constitution rule added:** Any term used in a proposal or delta spec that
  is not in the glossary MUST be added to the glossary before the spec is
  approved.
- **Maintenance:** Glossary grows incrementally. New terms are added when a
  spec review surfaces a misunderstanding. Old terms are updated when a
  retro reveals drift between what was written and what was built.

---

## How the three artifacts fit into the platform repo

```text
  PLATFORM REPO  (canonical shared truth)
  ├── platform-baseline.md         (already exists)
  ├── component-ownership.md       [new — one-time, updated rarely]
  │     one entry per component
  │     owns / does NOT own
  │
  ├── dependency-map.md            [new — three tiers]
  │     tier 1: must change together
  │     tier 2: watch for breakage
  │     tier 3: adapts independently
  │
  ├── glossary.md                  [new — grows incrementally]
  │     one entry per shared term
  │     definition + what it is NOT
  │
  ├── capabilities/
  │     customer-identity.md       extended: add owns + impact_surface fields
  └── contracts/
        customer-profile.v2.md     unchanged
```

---

## How the three artifacts are used per phase

```text
  ┌───────────────────────────────────────────────────────────────────────────┐
  │  PLATFORM PHASE  (one time, owned by Architect)                           │
  │                                                                           │
  │  Write component-ownership.md                                             │
  │    "what does each component own and not own?"                            │
  │                                                                           │
  │  Write dependency-map.md                                                  │
  │    "who must change together, who to watch, who adapts alone?"            │
  │                                                                           │
  │  Seed glossary.md                                                         │
  │    "what do our shared terms mean precisely?"                             │
  └───────────────────────────────────────────────────────────────────────────┘
              │
              ▼
  ┌───────────────────────────────────────────────────────────────────────────┐
  │  ASSESS PHASE  (every change, owned by Team Lead)                         │
  │                                                                           │
  │  Read component-ownership.md                                              │
  │    → answer: which component owns this request?                           │
  │    → output: correct epic owner in jira-traceability.yaml                 │
  │                                                                           │
  │  Read dependency-map.md                                                   │
  │    → answer: which tier are the affected relationships?                   │
  │    → output: impact field in platform-ref.yaml                            │
  │                                                                           │
  │  Tier 1 → open additional component epics now                            │
  │  Tier 2 → note watch items in alignment_notes                             │
  │  Tier 3 → no action                                                       │
  └───────────────────────────────────────────────────────────────────────────┘
              │
              ▼
  ┌───────────────────────────────────────────────────────────────────────────┐
  │  SPECIFY PHASE  (every change, owned by Product)                          │
  │                                                                           │
  │  Read glossary.md before writing proposal.md                              │
  │    → use only defined terms in goals and acceptance criteria               │
  │    → add missing terms to glossary before approving the proposal           │
  │                                                                           │
  │  Read component-ownership.md before writing delta specs                   │
  │    → confirm scope stays within the owning component's boundary           │
  │    → flag any requirement that crosses an ownership line                  │
  └───────────────────────────────────────────────────────────────────────────┘
              │
              ▼
  ┌───────────────────────────────────────────────────────────────────────────┐
  │  PLAN PHASE  (every change, owned by Architect)                           │
  │                                                                           │
  │  Read platform-ref.yaml impact field                                      │
  │    → tier 1 impacts become hard constraints in design.md                  │
  │    → tier 2 impacts become rollout risks in design.md                     │
  │                                                                           │
  │  Read dependency-map.md for the affected relationships                    │
  │    → design.md must show how cross-component dependencies are handled     │
  └───────────────────────────────────────────────────────────────────────────┘
              │
              ▼
  ┌───────────────────────────────────────────────────────────────────────────┐
  │  DELIVER PHASE  (every slice, owned by Team Lead)                         │
  │                                                                           │
  │  PR description notes:                                                    │
  │    → which tier 1 dependencies were verified in this slice                │
  │    → which tier 2 consumers were checked after deploy                     │
  │                                                                           │
  │  Archive records:                                                         │
  │    → whether any ownership boundary or dependency tier changed            │
  │    → if yes, flag component-ownership.md or dependency-map.md for update  │
  └───────────────────────────────────────────────────────────────────────────┘
```

---

## The two constitution rules to add

These rules are added to the Speckit constitution at the end of the Platform
phase. They enforce the three artifacts without requiring anyone to remember to
use them.

```text
## Ownership and impact rules  (add to platform constitution)

Rule O-1: Ownership before scope
  Every proposal MUST identify which component owns the primary change.
  Use component-ownership.md to verify before writing proposal.md.
  If the ownership is unclear or spans multiple components, resolve it
  at Assess before Specify begins.

Rule O-2: Glossary terms only
  Every term used in a proposal, delta spec, or acceptance criterion
  MUST appear in glossary.md.
  If the term is missing, add it to the glossary before the proposal
  is approved.
  Do not invent synonyms for existing glossary terms inside a spec.

Rule O-3: Impact tiers drive JIRA structure
  The dependency map tier of the primary affected relationship MUST be
  recorded in platform-ref.yaml at Assess time.
  Tier 1 relationships MUST result in coordinated component epics.
  Tier 2 relationships MUST result in a watch note in alignment_notes.
  Tier 3 relationships require no additional action unless breakage occurs.
```

---

## The extended `platform-ref.yaml` fields

Two fields are added to the existing `platform-ref.yaml` structure to carry
the ownership and impact information from the platform artifacts into the
change package:

```yaml
# platform-ref.yaml  (additions only — all other fields unchanged)

ownership:
  primary_component: "profile-service"
  # verified against component-ownership.md

impact:
  must_change_together: [auth-service]          # tier 1 — open epics now
  watch_for_breakage:   [notification-service]  # tier 2 — monitor after deploy
  adapts_independently: [audit-service]         # tier 3 — no action needed
  # verified against dependency-map.md

alignment_notes:
  glossary_terms_used:
    - email update
    - validation
    - duplicate email
  # verified against glossary.md
```

---

## Concrete example: validated email updates

Using the running example to show all three artifacts working together.

### Ownership check at Assess

```text
  Request: "validate customer email updates"
              │
              ▼  read component-ownership.md
  ┌──────────────────────────────────────────┐
  │  profile-service                         │
  │    owns: email update flow ✓             │
  │    owns: profile update event ✓          │
  │                                          │
  │  auth-service                            │
  │    owns: duplicate email enforcement ✓   │
  │    does NOT own: email update flow       │
  └──────────────────────────────────────────┘

  Result: profile-service is the primary owner
          auth-service has a dependency — check tier
```

### Impact tier check at Assess

```text
  primary: profile-service
              │
              ▼  read dependency-map.md
  ┌────────────────────────────────────────────────────────┐
  │  profile-service ◄──────────► auth-service             │
  │    tier 1: must change together                        │
  │    trigger: email validation logic                     │
  │    result: open AUTH-234 alongside PROF-456            │
  │                                                        │
  │  profile-service ────────────► notification-service    │
  │    tier 2: watch for breakage                          │
  │    trigger: contracts.customer-profile.v2              │
  │    result: watch note in alignment_notes               │
  │                                                        │
  │  profile-service ────────────► audit-service           │
  │    tier 3: adapts independently                        │
  │    result: no action                                   │
  └────────────────────────────────────────────────────────┘

  platform-ref.yaml written:
    impact.must_change_together: [auth-service]
    impact.watch_for_breakage:   [notification-service]
    impact.adapts_independently: [audit-service]
```

### Glossary check at Specify

```text
  Product drafts proposal.md goal:
    "validate new email addresses before persistence"
              │
              ▼  check glossary.md
  ┌──────────────────────────────────────────────────────┐
  │  "email update"  — defined ✓                         │
  │    customer request to replace registered address    │
  │                                                      │
  │  "validation"    — defined ✓                         │
  │    format + uniqueness + business rules before write │
  │                                                      │
  │  "persistence"   — NOT in glossary                   │
  │    add before approving: "writing data to the store" │
  └──────────────────────────────────────────────────────┘

  Proposal approved only after glossary updated.
```

---

## Maintenance rules

These artifacts do not change per-change. They change per platform evolution.

```text
  Artifact                 Update when
  ──────────────────────────────────────────────────────────────
  component-ownership.md   New component added
                           Ownership boundary shifts between components
                           Incident reveals ownership was wrong

  dependency-map.md        New integration added between components
                           Tier changes (e.g. tier 2 becomes tier 1
                             after repeated incidents)
                           Incident reveals a hidden dependency

  glossary.md              Spec review reveals term ambiguity
                           Retro reveals mismatch between spec and build
                           New shared term introduced by a platform change

  platform-ref.yaml        Each change — ownership and impact populated
  (ownership + impact)     at Assess, carried through to archive
```

The Architect owns the first two. The whole team contributes to the glossary.
The Team Lead populates the `platform-ref.yaml` fields at Assess using the
first two as lookup documents.

---

## What the team learns — in three sentences

1. **Every component owns a clear set of things.** Before writing a spec, check
   which component owns the change and what it does not own.

2. **Some components must move together, others just need watching.**
   Check the dependency map to know which type this change is.

3. **Use the glossary, not your own words.** If the word is not in the glossary,
   add it. If it is there, use it exactly as written.
