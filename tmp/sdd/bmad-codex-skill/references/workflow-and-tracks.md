# Workflow and track reference

BMAD’s workflow map says the system builds context progressively across **4 distinct phases**, with each phase producing documents that inform the next. It also points users to `bmad-help` when they are unsure what to do next. citeturn1view2

## Track selection

### Quick Flow
Quick Spec Flow is the streamlined BMAD track. The docs say it is for bug fixes, small features, rapid prototyping, and quick enhancements. It goes straight to a context-aware technical specification instead of requiring Product Brief → PRD → Architecture. citeturn1view6

Use Quick Flow when:
- single bug fix or small enhancement
- small feature with clear scope
- rapid prototyping
- brownfield work on an existing codebase
- the requester already knows what they want to build

### Full BMAD
Use the fuller BMAD path when the docs’ “Use BMad Method or Enterprise when” conditions apply:
- new products or major features
- stakeholder alignment is needed
- multi-team coordination exists
- architecture depth is required citeturn1view6

## Brownfield sizing guidance
The “How to Add a Feature to an Existing Project” guide gives a practical scale rule:
- small (1–5 stories): Quick Flow with tech-spec
- medium (5–15 stories): BMAD with PRD
- large (15+ stories): full BMAD with architecture citeturn1view7

## Typical artifact chain
For larger work, a useful BMAD-aligned chain is:
1. Product brief
2. PRD
3. Architecture
4. Project context
5. Sprint planning / create-story
6. Dev story
7. Code review

This chain is an implementation-friendly synthesis of the workflow map plus the brownfield feature guide. citeturn1view2turn1view7
