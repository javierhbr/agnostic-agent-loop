# Phase, Role, and Skill Worked Examples

The detailed worked examples now live under:

- `example/README.md`

Use that directory for:

- phase-by-phase examples
- role and skill examples
- entry-point examples
- platform/component interaction examples
- concrete artifact examples

Main entry points:

- `example/README.md`
- `example/01-platform-phase.md`
- `example/02-assess-phase.md`
- `example/03-specify-phase.md`
- `example/04-plan-phase.md`
- `example/05-deliver-phase.md`
- `example/06-entry-point-examples.md`
- `example/07-platform-component-interaction.md`

## DDD concepts — platform ownership, impact, and glossary

Three DDD-derived concepts are integrated into the methodology to make
ownership classification, impact assessment, and terminology consistent
across all phases. They add three durable artifacts to the platform repo.

Reference documents:

- `platform-ddd-spec.md` — full rationale, usage per phase, worked example,
  and maintenance rules for the three concepts
- `decisions/ADR-014-three-concept-ddd-for-platform-ownership-and-impact.md` —
  decision record with options considered and consequences

Templates:

- `templates/platform-template/ownership/component-ownership-template.md` —
  one file per component; records owns / does NOT own / published contracts /
  consumed contracts
- `templates/platform-template/ownership/dependency-map-template.md` —
  one platform file; records tier 1 / tier 2 / tier 3 relationships with
  their JIRA implications
- `templates/platform-template/ownership/glossary-template.md` —
  one platform file; records shared terms with plain definitions and
  "what it is NOT" clauses

Phase integration summary:

| Phase | Uses these artifacts | How |
|---|---|---|
| Platform | All three | Architect writes them once as platform truth |
| Assess | component-ownership + dependency-map | Team Lead reads to confirm owner and populate `platform-ref.yaml` impact tiers |
| Specify | glossary | Product reads before writing `proposal.md`; all terms must be defined |
| Plan | platform-ref.yaml impact field | Architect reads tiers; tier 1 → hard constraint, tier 2 → rollout risk |
| Deliver | platform-ref.yaml impact field | PR notes tier 1 verification; archive records tier changes |
