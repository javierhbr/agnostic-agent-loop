# Initiative: Subscription Billing Platform

**Initiative ID:** SBL-2025-Q1

## Problem

ShopFlow currently supports only one-time purchases. Merchants report that 40% of support requests are for subscription/recurring payment features. This prevents us from capturing the recurring revenue (SaaS) market segment, which has 3x higher lifetime customer value than one-time purchases.

We will launch subscription billing to enable merchants to offer subscriptions, reducing support churn and opening a new revenue channel. **Success metric:** 25% of new merchants enable subscriptions within 60 days of launch.

## Goals

- Enable merchants to create and manage subscription plans (monthly, quarterly, annual)
- Process recurring payments automatically on schedule
- Handle common subscription events: upgrades, downgrades, cancellations, renewals
- Reduce support requests for "how do I do recurring payments?" from 40% to < 5%
- Increase merchant retention for accounts with subscriptions by 30%

## Non-Goals

- We are NOT implementing dunning (automatic retry on card decline) — this is Phase 2
- We are NOT adding multi-currency support yet — English/USD only
- We are NOT building merchant-facing subscription analytics dashboard — that's Phase 2
- We are NOT changing the existing one-time payment flow

## Success Criteria

- Given a merchant logs in, When they access the products section, Then they can create a subscription plan with name, price, billing interval
- Given a subscription plan exists, When a customer purchases it, Then a subscription record is created and billing is scheduled automatically
- Given a subscription exists, When the billing date arrives, Then an automatic charge is attempted
- Given a subscription was active, When a customer cancels, Then no future charges occur and the subscription transitions to cancelled state
- Given the feature launches, When 60 days have passed, Then at least 25% of new merchants have created at least one subscription plan

## Affected Components

- **Billing Service** (new) — manages subscription lifecycle: create, update, cancel, billing schedule
- **Payments Service** (modify) — tokenized card storage, recurring charge processing
- **User Service** (minimal change) — subscription notification webhooks
- **Email Service** (modify) — send subscription confirmation, renewal reminder, cancellation confirmation emails
- **Platform API** (modify) — new `/subscriptions/` endpoints

## Timeline

- Risk Assessment + Workflow Selection: 2 days
- Analyst Discovery (if Full workflow): 3 days
- Architecture & Component Specs: 5 days
- Parallel Development (Billing + Payments): 8-10 days
- Verification & Testing: 3-4 days
- **Target Go-Live:** 3 weeks from start

## Stakeholders

- Merchants: need recurring payment capability
- Finance: needs new revenue stream tracking
- Support: needs to reduce subscription-related tickets
- Platform Team: owns payment infrastructure stability
