# Generic 3-Layer Context Instruction (Rules + Skills)

Purpose: keep agent context lean by separating routing, skill rules, and deep references while ensuring every rule and skill lives in exactly one place.

## Layer Definitions
- **Tier 1 — Routers (≤70 lines):** global rules, session start loaders, one-line links to skills. No examples or long text.
- **Tier 2 — Skills (60–120 lines):** per-skill intent, triggers, compact steps, skill-specific rules, and explicit links to resources.
- **Tier 3 — Resources (150–500 lines):** detailed procedures, examples, templates, FAQs. Split files once they near 500 lines; keep a single `resources/` folder, no nesting.

## Placement Rules
- Put universal rules in Tier 1; put skill-scoped rules in that skill’s Tier 2; put all lengthy guidance in Tier 3 and link to it from Tier 2.
- Never duplicate content upward. If a Tier 2 file grows, extract detail to Tier 3 immediately.
- Every Tier 2 file must point to its Tier 3 resources; every Tier 1 entry must point to its Tier 2 skill file.

## Load & Use
- Load Tier 1 at session start. When a skill triggers, load only its Tier 2. Pull Tier 3 sections on demand. Golden rule: load only what you need.
- If a response needs examples or templates, fetch the linked Tier 3 section instead of expanding Tier 2.

## Maintenance Cadence
- Monthly quick check: routers <70 lines; skills <130; resources <500; all Tier 2 links resolve; extract anything trending large.
- Quarterly deep dive: size audit, duplication pass, staleness review, broken-link check, and resource splitting where needed.

## Implementation Steps (new repo)
1) Create Tier 1 router file listing skills with one-line descriptions and paths.  
2) For each skill, write a 60–120 line Tier 2 file with triggers, steps, and links to `resources/`.  
3) Move long examples/templates into Tier 3 files under `resources/`.  
4) Add the size checks to CI or a monthly script; fail when limits are exceeded.  
5) Teach contributors: “router → skill → resource” and “extract when in doubt.”

## Quick Checklist
- Router updated?  
- Skill under 130 lines with clear triggers?  
- Resources exist for every linked detail and stay under 500 lines?  
- Broken-link scan clean?  
- Extraction performed as soon as growth appears?
