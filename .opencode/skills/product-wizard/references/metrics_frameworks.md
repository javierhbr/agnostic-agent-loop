# Product Metrics Frameworks

Guide to choosing and implementing metrics for PRD success criteria.

## Choosing a Framework

| Framework | Best For | Stage |
|-----------|---------|-------|
| **AARRR** | Growth-focused, startups | Early/Growth |
| **HEART** | UX quality, feature launches | Per feature |
| **North Star** | Company alignment, focus | Ongoing |
| **OKRs** | Goal setting, team alignment | Quarterly |

## AARRR (Pirate Metrics)

Customer lifecycle: **Acquisition → Activation → Retention → Revenue → Referral**

### Acquisition — How do users find you?
- Website traffic, app store impressions, CTR, CPA, traffic sources
- Example: 10,000 monthly visitors, CAC < $50

### Activation — Do users have a great first experience?
- Sign-up completion rate, time to "aha moment", onboarding completion
- Example: 60% sign-up completion, 40% reach aha moment in first session
- Aha moments: Slack (send first message), Dropbox (upload first file), Airbnb (first booking)

### Retention — Do users come back?
- DAU/WAU/MAU, retention curves (Day 1/7/30), churn rate, session frequency
- Example: 40% Day 7 retention, < 5% monthly churn

### Revenue — How do you monetize?
- MRR, ARPU, LTV, free-to-paid conversion, upsell rate
- Example: 5% free-to-paid, LTV:CAC > 3:1
- Formulas: LTV = ARPU × Avg Lifespan, Churn Rate = Lost ÷ Total × 100

### Referral — Do users tell others?
- K-factor, NPS, referral rate, social shares
- Example: 15% referral rate, NPS > 50, K-factor > 1

## HEART Framework (Google)

UX quality: **Happiness + Engagement + Adoption + Retention + Task Success**

| Dimension | Metrics | Example Target |
|-----------|---------|---------------|
| **Happiness** | NPS, CSAT, ratings | CSAT > 4.5/5 |
| **Engagement** | Session duration, actions/session | 2+ sessions/day |
| **Adoption** | Feature adoption rate, time to first use | 40% within 30 days |
| **Retention** | DAU/MAU, week-over-week active | 70% weekly active |
| **Task Success** | Completion rate, error rate, time to complete | 95% completion |

### HEART Template for PRDs

| Dimension | Goals | Signals | Metrics |
|-----------|-------|---------|---------|
| Happiness | Users love it | Positive feedback | NPS > 40 |
| Engagement | Frequent interaction | Daily active usage | 60% DAU/MAU |
| Adoption | Most users try it | Feature activation | 70% adoption |
| Retention | Users keep returning | Weekly return rate | 50% W1 retention |
| Task Success | Users complete goals | Low error rate | 95% success |

## North Star Metric

Single metric capturing core value delivery to customers.

### Characteristics
- Reflects value delivery to customers
- Leading indicator of revenue
- Actionable by the team
- Understandable by everyone

### Examples

| Company | North Star | Why |
|---------|-----------|-----|
| Airbnb | Nights booked | Core value: successful stays |
| Spotify | Time spent listening | Core value: music enjoyment |
| Slack | Messages sent by teams | Core value: communication |
| Netflix | Hours watched | Core value: entertainment |

### Finding Yours
1. Define your value proposition — what core value do you deliver?
2. Identify the metric — what measurement captures that value?
3. Validate — does it correlate with revenue? Can teams influence it?
4. Break down into a metric tree (contributing metrics)

## OKRs (Objectives & Key Results)

### Structure
- **Objective**: Qualitative, inspirational, time-bound goal
- **Key Results**: 3–5 quantitative, measurable outcomes per objective

### Example: Growth OKR
**Objective**: Become the go-to platform for small business invoicing

**Key Results:**
1. Increase monthly active businesses from 10,000 to 25,000
2. Achieve 40% month-over-month retention
3. Reach NPS of 50+
4. Generate $500K MRR

### Example: Feature Launch OKR
**Objective**: Successfully launch team collaboration features

**Key Results:**
1. 60% of active users try features within 30 days
2. 25% become weekly active collaborators
3. Features drive 15% increase in paid conversions
4. CSAT score of 4.2/5 for collaboration features

## Engagement Deep Dive

### DAU/MAU Stickiness Benchmarks
- **Excellent** (> 50%): Messaging apps
- **Good** (20–50%): Social media
- **Average** (10–20%): Utilities

### Product-Market Fit Signals
- Sean Ellis Test: 40%+ would be "very disappointed" without your product
- Retention: 40%+ month 1 retention
- NPS: Score > 50
- Growth: 10%+ MoM organic growth
- LTV:CAC: Ratio > 3:1

## Anti-Patterns

| Mistake | Problem | Fix |
|---------|---------|-----|
| Too many metrics | Focus on nothing | 3–5 key metrics per initiative |
| Vanity metrics | Doesn't drive decisions | Focus on active users, engagement |
| Lagging only | Rear-view mirror | Balance with leading indicators |
| No targets | Tracking without goals | Set specific, time-bound targets |
| Not segmenting | Hides patterns | Segment by user type, cohort |

## PRD Metrics Template

```markdown
## Success Metrics

### North Star Metric
**Metric**: [Single most important metric]
**Current**: [Baseline] | **Target**: [Goal by launch + X months]
**Why**: [Why this measures success]

### Supporting Metrics
| Stage | Metric | Current | Target |
|-------|--------|---------|--------|
| Acquisition | [Name] | [X] | [Y] |
| Activation | [Name] | [X] | [Y] |
| Retention | [Name] | [X] | [Y] |

### Counter-Metrics
- [Ensure support tickets don't increase > 10%]

### Measurement Plan
- **Dashboard**: [Link]
- **Review Cadence**: [Weekly/bi-weekly]
- **Owner**: [Name]
```
