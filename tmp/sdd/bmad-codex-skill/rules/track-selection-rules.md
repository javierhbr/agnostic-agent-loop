# BMAD track selection rules

Use these routing rules before generating any artifact.

## 1. Classify project type
- New project / greenfield
- Existing project / brownfield

## 2. Classify size
- Small: 1–5 stories
- Medium: 5–15 stories
- Large: 15+ stories

## 3. Pick track
- Small or bugfix → Quick Flow / tech-spec
- Medium → PRD-first
- Large or architecture-heavy → PRD + architecture

## 4. Escalate to architecture when any of these are true
- new service or subsystem
- major data model changes
- external integration changes
- security, reliability, or operational implications
- multiple teams or domains affected

## 5. Brownfield-specific rules
- conform to repo conventions unless there is explicit reason not to
- identify integration points early
- generate or refresh project-context when the repo’s norms are unclear
