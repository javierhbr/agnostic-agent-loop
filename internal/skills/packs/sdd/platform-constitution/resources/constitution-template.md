# Constitution Template

Full template for authoring a platform Constitution with all 6 sections, worked examples per section, and gate enforcement rules. Referenced by the main platform-constitution skill.

---

## Full Constitution Template

File location: `.specify/memory/constitution.md`

```markdown
# Platform Constitution v<VERSION>

Last updated: <YYYY-MM-DD>
Owner: <Platform Architect Name>

---

## 1. UX Rules

All public-facing experiences must adhere to:

### Accessibility
- WCAG 2.1 Level AA minimum standard (verified by automated scan)
- Keyboard navigation for all interactive elements
- Alt text for all images (no generic "image" labels)
- Color contrast ratio ≥ 4.5:1 for text
- Screen reader testing on 2+ screen readers per release

Gate enforcement: Gate 4 (NFR) requires accessibility checklist signed off by UX team.

### Design System
- All UI components from [Design System Name/Link]
- No custom one-off components without Design System approval
- Design review required before engineering starts
- Consistent spacing, typography, color palette per design system

Gate enforcement: Gate 4 requires design review sign-off.

### Mobile Experience
- Mobile-first responsive design (test at 320px, 768px, 1024px viewports)
- Touch targets minimum 44x44 pixels
- No hover-only interactions
- Performance target: LCP < 2.5s on 4G, mobile

Gate enforcement: Gate 4 requires mobile performance test results.

---

## 2. Security & PII Policy

### PII Handling (Mandatory)

All PII fields (email, phone, SSN, address, credit card, IP address) must:
- Be encrypted at rest using [encryption standard, e.g., AES-256]
- Be masked in all logs (e.g., email shown as "u****@example.com")
- Be never logged in raw form, anywhere
- Have masking enforced at API boundary (before any logging/tracing)
- Have retention policy documented (e.g., "deleted after 90 days")

### Cross-Service PII Access (Forbidden)

- PII fields CANNOT be accessed directly by other services
- Use contracts (events/APIs) to share PII-adjacent data only (e.g., user_id, not email)
- If another service needs to contact a user, use dedicated contact service via contract

### Authentication & Sessions

- All auth flows MUST support MFA (multi-factor authentication)
- Session TTL: max 24 hours for web clients, max 7 days for mobile clients
- All sessions MUST have logout/invalidation capability
- Password reset MUST invalidate all existing sessions
- Compromised token MUST be revocable immediately

### API Security

- All endpoints MUST require authentication (except /health, /status)
- HTTPS required for all traffic
- Rate limiting: [X requests/second] per user
- API key rotation: max [Y days]
- CORS: whitelist specific origins (never use *)

Gate enforcement: Gate 4 (NFR) requires security checklist + PII masking confirmation.

---

## 3. Observability Standards

### Logging (Required)

All services MUST emit structured logs in JSON format with minimum fields:

```json
{
  "timestamp": "2025-02-28T14:23:45Z",
  "level": "INFO",
  "service": "payment-processor",
  "request_id": "req-abc123def456",
  "message": "payment.authorized",
  "user_id": "user-789",
  "amount": 99.99
}
```

**Minimum fields:**
- `timestamp` (ISO 8601)
- `level` (DEBUG, INFO, WARN, ERROR, CRITICAL)
- `service` (service name)
- `request_id` (correlation ID across services)
- `message` (log message)

**PII masking rules:**
- Email: "u****@example.com" or "****@example.com"
- Phone: "***-***-1234"
- Credit card: "****-****-****-4242"

**Log aggregation:**
- Send to [log aggregation service: e.g., Datadog, ELK, CloudWatch]
- Retention: [X days]
- Searchable by: service, request_id, user_id, level

Gate enforcement: Gate 4 (NFR) requires logging declaration with format, fields, masking rules.

### Metrics (Required)

All services MUST emit Prometheus-format metrics:

**Standard metrics per endpoint:**
- Latency histogram: `<service>_<endpoint>_duration_ms` (buckets: 50, 100, 200, 500, 1000, 5000)
- Request counter: `<service>_<endpoint>_requests_total` (labels: method, status)
- Error rate: `<service>_<endpoint>_errors_total` (labels: error_type)
- Throughput: `<service>_<endpoint>_throughput_rps`

**Example:**
```
payment_processor_checkout_duration_ms{quantile="0.95"} 250
payment_processor_checkout_requests_total{method="POST", status="200"} 12450
payment_processor_checkout_errors_total{error_type="timeout"} 12
```

**Metric aggregation:**
- Send to [metrics service: e.g., Prometheus, Datadog, CloudWatch]
- Retention: [X days]
- Dashboards: [list critical dashboards]

Gate enforcement: Gate 4 (NFR) requires metric names, dimensions, aggregation destination.

### Tracing (Required)

All services MUST emit distributed traces:

**Trace span per request:**
- Span name: `<service>.<operation>` (e.g., `payment.authorize`)
- Parent span: request entry point
- Child spans: external calls (database, cache, other services)
- Trace ID: correlation across all services (from request header)

**Span attributes:**
- `service`: service name
- `operation`: operation name
- `duration`: milliseconds
- `status`: success | error
- `error.message`: if error
- Business attributes: `user_id`, `order_id`, `amount` (non-PII)

**Trace aggregation:**
- Send to [tracing service: e.g., Jaeger, Datadog, Lightstep]
- Retention: [X days]
- Sampling rate: [% of requests traced]

Gate enforcement: Gate 4 (NFR) requires span definitions and sampling strategy.

### Alerts (Required)

Critical alerts that MUST trigger paging:

- **p95 latency > [X ms]** — Response time degradation indicates service issue
- **Error rate > [Y%]** — Errors spike above baseline
- **Throughput drop > [Z%]** — Request volume drops (possible service down)
- **Circuit breaker open** — Dependency failure detected
- **PII exposure** — Any PII logged in plain text
- **Authentication failure** — Spike in auth failures

Gate enforcement: Gate 4 (NFR) requires alert rules + escalation policy.

---

## 4. Performance Baselines

### Customer-Facing Endpoints

All user-visible operations MUST meet these targets:

- **p95 latency: < 300ms** (strict — user perceives slowness at >300ms)
- **p99 latency: < 1000ms** (acceptable, 99% of requests complete within 1s)
- **Error rate: < 0.1%** (1 error per 1000 requests max)
- **Availability: 99.9%** (4.3 hours downtime per month max)

### Backend/Internal Endpoints

Internal service-to-service calls:

- **p95 latency: < 2 seconds** (internal communication can be slower)
- **p99 latency: < 5 seconds**
- **Error rate: < 0.5%**
- **Availability: 99%** (7.2 hours downtime per month)

### Caching Strategy

- Cache TTL: [minimum X minutes, maximum Y hours] per resource type
- Cache invalidation: [strategy: TTL expiry | event-driven | manual]
- Cache hit rate target: [X% for critical paths]
- Fallback to source: [required if cache miss or error]

### Database Query Limits

- Query timeout: [X seconds] per query
- Index requirements: [list critical indexes]
- Connection pool: [max N connections per service]
- Transaction timeout: [Y seconds] max per transaction

### Load Testing Requirements

- Test at 2x peak traffic volume before launch
- Measure: latency, throughput, error rate, CPU, memory
- Pass criteria: All performance targets met at 2x peak

Gate enforcement: Gate 4 (NFR) requires performance targets + load test results.

---

## 5. Domain Governance

### Domain Ownership

Each domain has a single owner responsible for:
- Contract definitions (events, APIs, versions)
- Invariants (immutable business rules)
- Consistency across all services in the domain

| Domain | Owner | Contact | Invariants |
|--------|-------|---------|-----------|
| Payment | Alice Chen (alice@company.com) | Slack: #payments | Amount > 0, Currency valid, Idempotent |
| User | Bob Singh (bob@company.com) | Slack: #users | Email unique, Active XOR Deleted |
| Order | Carol Lee (carol@company.com) | Slack: #orders | Status progression, Total = sum of items |
| Inventory | David Park (david@company.com) | Slack: #inventory | Quantity ≥ 0, Reserved ≤ Available |

### Invariants (Immutable Business Rules)

Each domain has invariants that CANNOT be violated:

**Payment domain:**
- Amount > 0 (no negative or zero payments)
- Currency must be valid (ISO 4217 code)
- Idempotency key ensures retry-safe processing
- Once authorized, payment cannot be reversed without refund operation

**User domain:**
- Email address is unique (no two users same email)
- User is either Active XOR Deleted (never both, never neither)
- Password hash is never stored in plain text
- User can have multiple sessions but only one primary session

**Order domain:**
- Status progression is only: DRAFT → CONFIRMED → SHIPPED → DELIVERED (no skips, no backwards)
- Order total = sum of item totals (always in sync)
- Each order has exactly one owner (user_id)
- Refund amount ≤ original order total

### Cross-Domain Communication Rules

- Services CANNOT access another domain's database directly
- Communication ONLY via versioned contracts (events or REST APIs)
- Domain owner approves all new contracts or contract versions
- Data ownership: if data belongs to Domain X, only Domain X modifies it

Gate enforcement: Gate 2 (Domain Validity) enforces invariants and cross-domain rules.

---

## 6. Contract Versioning

### Versioning Scheme

**NATS Events:**
- Topic format: `ServiceName-EventType-vX` (e.g., `Payment-Authorized-v1`, `Payment-Authorized-v2`)
- Semantic versioning: vX (major version only, no minor/patch)
- v1 vs v2 = breaking schema change (fields removed, type changed, required field added)

**HTTP APIs:**
- Path format: `/api/v1/resource` or `/api/v2/resource`
- v1 vs v2 = breaking change (field removed, type changed, endpoint removed)
- Deprecated version must be supported for 90 days minimum

**GraphQL:**
- Full version in `schema.graphql` file header: `# GraphQL Schema v2.0`
- Breaking changes increment major version

### Breaking Change Process

When a contract MUST change (field removed, type changed, endpoint removed):

1. **Dual-publish phase (30 days):**
   - Publish old format AND new format simultaneously
   - Old format: marked as deprecated with migration guide
   - New format: has new field/endpoint

   Example: Order event emits BOTH `total` (old) AND `totalAmount` (new)

2. **Consumer migration (30 days):**
   - All downstream services update to consume new format
   - Monitoring: check old format still being sent

3. **Deprecation phase (30 days):**
   - Send deprecation warnings when old format used
   - Metrics: track usage of old format
   - Notify teams: "old format sunset in 30 days"

4. **Removal (after 90 days):**
   - Remove old format completely
   - Only new format supported

### Compatibility Checks

Before any contract change:
- Query Integration MCP: "Who consumes this contract?"
- Notify all consumers before breaking change
- Get approval from all consumer teams
- Create contract change ADR if complexity is high

Gate enforcement: Gate 3 (Integration Safety) enforces contract versioning and dual-publish strategy.

---

## Definition of Done (DoD)

No change is "done" until ALL these items are complete:

- [ ] All acceptance criteria pass with observable evidence (verified in verify.md)
- [ ] All 5 gates pass (Context, Domain, Integration, NFR, Ready)
- [ ] Code reviewed and approved by 2+ engineers
- [ ] All tests pass (unit + integration): `make test` exits 0
- [ ] Performance targets met (load test if touching critical path)
- [ ] Logging/metrics/tracing implemented and tested in staging
- [ ] Security review passed (if touching auth/PII)
- [ ] Documentation updated (README, API docs, runbooks)
- [ ] Monitoring/alerts configured and tested
- [ ] Rollback strategy proven (feature flag tested in staging)
- [ ] Changelog entry added
- [ ] Release notes drafted for customers (if applicable)

Gate enforcement: Gate 5 (Ready-to-Implement) requires DoD checklist before merge.

---

## Amendment Process

Constitution changes are significant and require approval:

1. **Draft amendment** — Propose change with rationale
2. **Get sign-off** — Platform Architect + affected Domain Owners must approve
3. **Create PR** — Update this file with change
4. **After merge** — Increment version number (v1.0 → v1.1)
5. **Cascade effect** — All Approved/Implementing specs must recheck Gates 1 & 4
6. **ADR if conflict** — If change contradicts existing specs, require exception ADR

---

## Example: Complete Constitution Section

```markdown
### Logging (Specific Example)

All services MUST emit JSON logs with these fields:

**Payload example:**
{
  "timestamp": "2025-02-28T14:23:45.123Z",
  "level": "INFO",
  "service": "order-service",
  "request_id": "req-12345abcde",
  "message": "order.created",
  "user_id": "usr-789",
  "order_id": "ord-456",
  "amount_cents": 9999,
  "currency": "USD"
}

**Masking rules:**
- Email addresses: "u****@example.com"
- Phone numbers: "***-***-1234"
- Credit cards: "****-****-****-4242"
- SSN: "***-**-1234"

**Destination:** Datadog (logs.datadoghq.com)
**Retention:** 30 days
**Searchable by:** request_id, user_id, service, level
**Alert if:** ERROR or CRITICAL logs appear (immediate slack notification)
```
```

This is your platform's source of truth for compliance, performance, and security policies.
