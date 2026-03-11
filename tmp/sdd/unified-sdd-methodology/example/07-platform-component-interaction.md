# Example: Platform and Component Interaction

This document shows how component teams interact with the platform specification while keeping one source of truth for shared rules.

## Core rule

- platform repo = canonical shared truth
- component repo = local implementation truth
- JIRA = tracking and coordination truth
- once the work enters the component repo, use `OpenSpec` only

## Example interaction model

```text
[Platform master repo]
  capability: customer-identity
  contract: customer-profile.v2
  version: 2026.03
        |
        v
[profile-service repo]
  platform-ref.yaml
  proposal.md
  design.md
  tasks.md
        |
        v
[PROF-456 epic]
  PROF-789 -> PR -> verify
```

## Step 1: Platform publishes truth

The platform team defines:

- platform version `2026.03`
- capability ref `capabilities.customer-identity`
- contract ref `contracts.customer-profile.v2`
- JIRA convention: one platform issue -> many component epics -> stories

## Step 2: Component team pins the platform version

The component team creates `platform-ref.yaml` in the component repo.

Example prompt:

- `Team Lead`: "Using the OpenSpec skill, create platform-ref.yaml for profile-service and record the platform version, capability refs, and contract refs that constrain this change."

Expected output:

- aligned `platform-ref.yaml`

## Step 3: Component team writes local specs

The component team uses OpenSpec only.

Example prompts:

- `Product`: "Using the OpenSpec skill, write the local proposal and acceptance criteria for profile-service, keeping the shared platform refs visible."
- `Architect`: "Using the OpenSpec skill, draft the local component design so `profile-service` stays aligned to platform version `2026.03` and `contracts.customer-profile.v2`."
- `Developer`: "Using the OpenSpec skill, refine the local component spec and tasks so they stay consistent with platform version `2026.03` and `contracts.customer-profile.v2`."

Expected outputs:

- local proposal
- local delta specs
- local design and tasks

## Step 4: JIRA links the delivery chain

The component team updates `jira-traceability.yaml` and keeps the issue chain aligned.

Example prompts:

- `Team Lead`: "Using the OpenSpec skill, map the local tasks for profile-service to story keys under PROF-456 and keep the issue chain aligned to PLAT-123."
- `Developer`: "Using the OpenSpec skill, update the current story and PR links in jira-traceability.yaml as the slice moves through delivery."

Expected outputs:

- story mapping
- PR mapping
- final traceability chain

## Step 5: Shared changes update both sides

If shared truth changes, both repositories move.

```text
[Platform issue PLAT-123]
        |
        +--> [Platform change package]
        |
        +--> [Profile component change]
        |
        +--> [Auth component change]
```

Example prompts:

- `Architect`: "Using the OpenSpec skill, review the local component package and determine whether the proposed behavior changes shared contract behavior and therefore requires a linked platform change."
- `Product`: "Using the OpenSpec skill, update the shared and local acceptance expectations so they stay consistent across the platform issue and component epic."

Expected outputs:

- linked platform and component change packages
- consistent shared and local specs

## What good looks like

Teams are applying the model correctly when:

- component repos do not copy editable platform truth
- every component change knows its platform version
- local specs always show which platform refs constrain them
- stories, PRs, and archive records stay linked to the same platform issue chain
