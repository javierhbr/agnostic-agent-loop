# Product Requirements Document Template

Adapt sections based on scope. Remove sections not relevant; add domain-specific ones.

## Document Header

| Field | Value |
|-------|-------|
| **Product/Feature Name** | [Name] |
| **Status** | [Draft / In Review / Approved] |
| **Author** | [Your Name] |
| **Stakeholders** | [List key stakeholders] |
| **Date Created** | [YYYY-MM-DD] |
| **Last Updated** | [YYYY-MM-DD] |
| **Version** | [1.0] |

---

## 1. Executive Summary

**One-liner**: [Single sentence describing the product/feature]

**Overview**: [2–3 paragraph summary of what, why, and expected impact]

| Quick Fact | Value |
|-----------|-------|
| Target Users | [Primary user segment] |
| Problem Solved | [Core problem] |
| Key Metric | [Primary success metric] |
| Target Launch | [Date or Quarter] |

---

## 2. Problem Statement

### The Problem
[Clearly articulate the problem. What pain point exists today?]

### Current State
[How users currently handle this, including workarounds]

### Impact
- **User Impact**: [How this affects users. Quantify: "Users spend 30 min daily on workarounds"]
- **Business Impact**: [How this affects business. Include metrics: "Costs $X in support monthly"]

### Why Now?
[Explain urgency or strategic importance of solving this now versus later]

---

## 3. Goals & Objectives

### Business Goals
1. **[Goal 1]**: [Description and expected impact]
2. **[Goal 2]**: [Description and expected impact]
3. **[Goal 3]**: [Description and expected impact]

### User Goals
1. **[Goal 1]**: [What users want to achieve]
2. **[Goal 2]**: [What users want to achieve]

### Non-Goals
- [What we're explicitly NOT trying to achieve]

---

## 4. User Personas

### Primary Persona: [Name/Type]

| Attribute | Detail |
|-----------|--------|
| Role/Title | [Role] |
| Tech Savviness | [Low/Medium/High] |

**Behaviors**: [Key behavior patterns]

**Needs & Motivations**: [What they need to accomplish, what drives decisions]

**Pain Points**: [Current frustrations]

**Quote**: "[Verbatim user quote capturing their perspective]"

### Secondary Persona: [Name/Type]
[Repeat as needed]

---

## 5. User Stories & Requirements

### Epic: [Epic Name]

#### Must-Have (P0)

**US-001: [Feature Name]**

> As a [user type],
> I want to [perform action],
> So that [achieve benefit/value].

**Acceptance Criteria:**
- [ ] Given [context], when [action], then [expected outcome]
- [ ] Given [context], when [action], then [expected outcome]
- [ ] Edge case: [Specific scenario]

**Priority**: Must Have (P0)
**Effort**: [T-shirt size: XS/S/M/L/XL]
**Dependencies**: [List any]

#### Should-Have (P1)
[Stories using same format]

#### Nice-to-Have (P2)
[Stories using same format]

### Functional Requirements

| Req ID | Description | Priority | Status |
|--------|-------------|----------|--------|
| FR-001 | [Specific, measurable requirement] | Must Have | Open |
| FR-002 | [Specific, measurable requirement] | Should Have | Open |

### Non-Functional Requirements

| Req ID | Category | Description | Target |
|--------|----------|-------------|--------|
| NFR-001 | Performance | Page load time | < 2 seconds |
| NFR-002 | Availability | Uptime SLA | 99.9% |
| NFR-003 | Security | Data encryption | AES-256 |
| NFR-004 | Accessibility | WCAG compliance | Level AA |

---

## 6. Success Metrics

### North Star Metric

| Attribute | Value |
|-----------|-------|
| **Metric** | [Single most important metric] |
| **Current** | [Baseline value] |
| **Target** | [Goal by launch + X months] |
| **Why** | [Why this measures success] |

### Supporting Metrics

| Stage | Metric | Current | Target | Timeframe |
|-------|--------|---------|--------|-----------|
| Acquisition | [Metric] | [Value] | [Value] | [When] |
| Activation | [Metric] | [Value] | [Value] | [When] |
| Retention | [Metric] | [Value] | [Value] | [When] |
| Revenue | [Metric] | [Value] | [Value] | [When] |

### Counter-Metrics
[Metrics to ensure we're not sacrificing other areas, e.g., "Support tickets don't increase > 10%"]

### Measurement Plan
- **Dashboard**: [Link]
- **Review Cadence**: [Weekly/bi-weekly]
- **Owner**: [Name]

---

## 7. Scope

### In Scope

**Phase 1 (MVP):**
- [Feature/capability 1]
- [Feature/capability 2]

**Phase 2 (Post-MVP):**
- [Feature/capability 1]

### Out of Scope
- [Item 1 and why it's excluded]
- [Item 2 and why it's excluded]

### Future Considerations
- [Enhancement 1]
- [Enhancement 2]

---

## 8. AI System Requirements (If Applicable)

### Model & Tool Requirements
- **Models**: [Which LLMs/ML models]
- **Tools/APIs**: [External services needed]
- **Prompt Strategy**: [Key prompt design approaches]

### Evaluation Strategy
- **Benchmark Test Set**: [Number of test cases, expected pass rate]
- **Automated Evals**: [What runs continuously]
- **Human Review**: [Where human-in-the-loop is required]

### Guardrails
- **Safety**: [Content filtering, harmful output prevention]
- **Hallucination Prevention**: [Grounding, citation requirements]
- **Failure Modes**: [What happens when the model fails]

---

## 9. Technical Specifications

### Architecture Overview
[High-level data flow and component interaction]

### Integration Points
- **APIs**: [External/internal APIs]
- **Databases**: [Storage systems]
- **Auth**: [Authentication/authorization method]

### Performance Requirements
| Metric | Target |
|--------|--------|
| Response Time (p95) | < [X]ms |
| Throughput | [X] req/s |
| Concurrent Users | [X] |

### Security & Privacy
- **Encryption**: [At rest and in transit]
- **Compliance**: [GDPR, HIPAA, SOC 2, etc.]
- **Data Handling**: [PII handling, retention, deletion]

### Scalability
- **Expected Load**: [Users, requests, data volume]
- **Scaling Strategy**: [Horizontal/vertical, auto-scaling]

---

## 10. Design & UX Requirements (If Applicable)

### User Flow
1. [Step 1]
2. [Step 2]
3. [Final state]

### Design Assets
- [Link to Figma/Sketch files]
- [Link to design system]

### Accessibility
- WCAG 2.1 Level AA compliance
- Keyboard navigation support
- Screen reader compatibility
- Color contrast ratios (4.5:1 for text)

### Responsive Breakpoints
- Mobile: 320px – 767px
- Tablet: 768px – 1023px
- Desktop: 1024px+

---

## 11. Timeline & Milestones

| Phase | Deliverables | Owner | Start | End |
|-------|-------------|-------|-------|-----|
| Discovery | Requirements finalized | PM/Design | [Date] | [Date] |
| Design | High-fidelity mockups | Design | [Date] | [Date] |
| Development | Implementation | Engineering | [Date] | [Date] |
| QA | Testing complete | QA | [Date] | [Date] |
| Launch | Production release | Engineering | [Date] | [Date] |
| Post-Launch | Monitoring, iteration | PM | [Date] | [Date] |

---

## 12. Risks & Mitigation

| Risk | Impact | Probability | Mitigation | Owner |
|------|--------|-------------|------------|-------|
| [Risk 1] | High | Medium | [Strategy] | [Name] |
| [Risk 2] | Medium | High | [Strategy] | [Name] |

### Contingency Plans
- If [scenario]: Action plan [steps], Decision maker [who], Trigger [what indicates this]

---

## 13. Dependencies & Assumptions

### Dependencies
- **Internal**: [Dependency 1], [Dependency 2]
- **External**: [Third-party API approval], [Partner integration]

### Assumptions
- [Assumption 1: e.g., "Users have app version 2.0+"]
- [Assumption 2: e.g., "Budget approved for $X infrastructure"]

---

## 14. Open Questions

| # | Question | Context | Options | Owner | Deadline |
|---|----------|---------|---------|-------|----------|
| 1 | [Question] | [Why it matters] | [Options] | [Who decides] | [When] |
| 2 | [Question] | [Why it matters] | [Options] | [Who decides] | [When] |

---

## 15. Stakeholder Sign-Off

| Stakeholder | Role | Status | Approved | Date |
|-------------|------|--------|----------|------|
| [Name] | Product Lead | ⏳ Pending | ☐ | - |
| [Name] | Engineering Lead | ⏳ Pending | ☐ | - |
| [Name] | Design Lead | ⏳ Pending | ☐ | - |

---

## Appendix

### References
- [User research findings link]
- [Competitive analysis link]
- [Technical design doc link]

### Glossary
- **[Term 1]**: [Definition]
- **[Term 2]**: [Definition]

### Change Log

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | [YYYY-MM-DD] | [Name] | Initial draft |
