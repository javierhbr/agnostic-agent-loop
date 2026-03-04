# ADR-001: Card Storage Strategy

**Status:** PENDING DECISION (blocks implementation)
**Date Created:** 2025-02-25
**Owner:** CTO + Security Team
**Scope:** Payment Processing Architecture

---

## Context

The Subscription Billing feature requires storing tokenized credit card information for automatic recurring charges. This storage decision has major implications for security, compliance, cost, and time-to-market.

**Source:** Discovery.md identified card vault as a blocking decision. Feature-spec gates validation requires resolution before developer handoff.

---

## Decision Required

**Question:** Where will ShopFlow store and manage tokenized credit card tokens?

---

## Options

### Option A: Self-Hosted Vault in Payments Service

**Approach:** Build a secure vault in the Payments Service to store tokenized card tokens.

**Pros:**
- Full control of card data
- No vendor lock-in
- Lower per-charge transaction fees (no Stripe markup)
- Customer card data stays within ShopFlow infrastructure

**Cons:**
- **PCI-DSS Compliance:** We become Level 1 merchant (highest burden)
  - Annual security audit required (~$50-100k)
  - Card rotation, expiry, deletion logic is our responsibility
  - Encryption, access control, monitoring all our liability
- **Security Risk:** Higher attack surface (we manage the vault)
- **Time-to-Market:** 4-6 weeks to build + audit
- **Maintenance Burden:** Ongoing updates for security patches, card networks changes

**Cost Estimate:** $30-50k initial + $20k/year maintenance

**Timeline:** 6-8 weeks (build + security audit)

---

### Option B: Stripe Vault (via Stripe Connect)

**Approach:** Stripe stores the tokenized card. We receive a token that we can re-use for charges.

**Pros:**
- **PCI-DSS:** Stripe handles compliance (they are Level 1 certified)
- **Proven:** Millions of merchants use Stripe vaults, battle-tested
- **Automatic Card Expiry:** Stripe updates expired cards automatically (if customer updates in Stripe)
- **Time-to-Market:** 2-3 weeks to integrate
- **Lower Liability:** Stripe bears PCI compliance responsibility
- **Security:** Stripe security team is world-class

**Cons:**
- Vendor Lock-in: If we ever want to switch payment processors, card migration is complex
- Higher Transaction Fees: Stripe charges 2.9% + $0.30 per charge (vs. 1-2% for self-hosted gateways)
- Customer Card Data on Stripe: We don't control the storage location (though Stripe's security is excellent)

**Cost Estimate:** 2.9% + $0.30 per transaction (~$0.40 per $10 subscription, scales with revenue)

**Timeline:** 2-3 weeks

---

### Option C: Hybrid (Stripe Tokenization, ShopFlow Storage)

**Approach:** Stripe tokenizes the card (they handle PCI), we store the token in our vault.

**Pros:**
- Stripe handles initial PCI burden
- We maintain token ownership
- Can switch processors later (tokens are ours)

**Cons:**
- **We still store payment data:** Token is still sensitive, still PCI scope
- **Compliance complexity:** Hybrid responsibility (Stripe's PCI + our token security)
- **Stripe tokenization cost:** Still pay Stripe transaction fee even though we don't use their vault
- **Time-to-Market:** 4-5 weeks (Stripe integration + our token vault)

**Cost Estimate:** Stripe fees + self-hosted vault costs

**Timeline:** 4-5 weeks

---

## Recommendation

**OPTION B: Stripe Vault**

### Rationale

1. **Time-to-Market is Critical:** 3-week timeline vs. 6-8 weeks means we capture Q1 merchant demand
2. **Risk Reduction:** PCI-DSS compliance is Stripe's responsibility, reducing our attack surface significantly
3. **Cost-Effective for Phase 1:** At 5,000 subscriptions × $120/year avg = $600k MRR, transaction fees are ~$18k/month (manageable)
4. **Phase 2 Optionality:** Once feature is proven and revenue justifies, Phase 2 can evaluate self-hosted vault
5. **Team Capacity:** Security + Infrastructure team is thin; Stripe reduces hiring needs

### Risks Mitigated

- **Vendor Lock-In Risk:** If cost becomes prohibitive in Phase 2, we can migrate tokens (Stripe provides export)
- **Card Data Risk:** Stripe's security is higher than we can achieve in Phase 1

### Decision Impact on Components

- **Payments Service:** Integrate with Stripe API for tokenization and charge processing
- **Billing Service:** Store Stripe token IDs (not card numbers) in subscription records
- **Observability:** Log every token creation and charge with Stripe token ID (not card data)

---

## Next Steps

1. **Decision:** CTO + Security approve this ADR (3-day timeline)
2. **Stripe Setup:** Finance provisioning Stripe account (2 days)
3. **Integrate:** Payments Service developer implements Stripe integration (3-4 days)
4. **Component Specs:** Re-validate component-spec-payments with Stripe choice
5. **Unblock:** Payments component-spec is no longer blocked; hand off to developer

---

## Sign-Off (Pending)

- **Decision Maker:** [CTO Name] — Pending approval
- **Security Review:** [Security Lead] — Pending review
- **Finance:** [Finance] — Pending budget approval

**Decision Date:** [To be filled on approval]

**Status Update:** When this ADR reaches `Approved`, the feature-spec and component-specs transition from `Draft` to `Approved`, and developer work can begin.
