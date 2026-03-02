---
name: api-documentation-generator
description: Generate comprehensive API documentation for REST, GraphQL, and WebSocket APIs
---

# skill:api-docs

## Does exactly this

Generate user-friendly, production-ready documentation for any API type — REST, GraphQL, or WebSocket — covering endpoints, authentication, errors, and working code examples in multiple languages.

---

## When to use

- Documenting a new API for the first time
- Updating existing API documentation after changes
- API lacks clear or complete documentation
- Onboarding new developers to your API
- Preparing documentation for external or public consumers
- Creating an OpenAPI/Swagger specification

---

## Steps — in order, no skipping

1. **Analyze** — Scan the codebase for routes, HTTP methods, parameters, response shapes, auth patterns, and error handling conventions.

2. **Generate** — Produce per-endpoint docs: method + URL, auth requirements, request schema, response schema, all error codes.

3. **Add** — Write multi-language code examples (cURL, JavaScript, Python) for every endpoint.

4. **Document** — Compile error handling section: all codes, message formats, troubleshooting steps, and common scenarios.

5. **Create** — Produce deliverables: Postman collection, OpenAPI/Swagger YAML, and interactive examples where applicable.

---

## Output

A complete documentation set:
- Per-endpoint reference pages (method, auth, request/response schema, error codes, code examples)
- Authentication guide (token acquisition, usage, refresh flow)
- Error code reference with troubleshooting
- OpenAPI/Swagger YAML snippet
- Postman collection JSON

---

## Done when

- Every endpoint has method, auth requirement, full request/response schema documented
- At least cURL + one language example exists per endpoint
- All error codes listed with descriptions and example payloads
- Authentication setup steps are complete and testable
- OpenAPI snippet or Postman collection is included

---

## Best practices

**Do:** Be consistent across all endpoints — same format, same depth.
**Do:** Include examples with realistic data, not placeholder "foo"/"bar" values.
**Do:** Document every error code and what it means for the caller.
**Do:** Version your API paths (/api/v1/) and note it in the introduction.
**Do:** Test every code example before including it.
**Do:** Add rate limit details and pagination patterns where they apply.
**Do:** Link related endpoints so consumers can discover adjacent functionality.

**Don't:** Skip error cases — users need to know what can go wrong.
**Don't:** Use vague descriptions like "Gets data" — be specific.
**Don't:** Let documentation drift from code — regenerate or review on every API change.
**Don't:** Omit response headers that carry semantic meaning.

---

## Recommended documentation structure

1. **Introduction** — purpose, base URL, API version, support contact
2. **Authentication** — how to authenticate, token lifecycle, security notes
3. **Quick Start** — minimal working example, common use-case walkthrough
4. **Endpoints** — organized by resource, full details per endpoint
5. **Data Models** — schema definitions, field types, validation rules
6. **Error Handling** — error code reference, response format, troubleshooting
7. **Rate Limiting** — limits, relevant headers, handling 429 responses
8. **Changelog** — version history, breaking changes, deprecation notices
9. **SDKs and Tools** — client libraries, Postman collection, OpenAPI spec

---

## If you need more detail

→ `resources/api-docs-examples.md` — Full examples (REST, GraphQL, auth flow), expanded best practices, section body copy, common pitfalls, OpenAPI snippet, Postman JSON

---

## Related skills

- `@doc-coauthoring` — collaborative documentation writing
- `@copywriting` — clear, user-friendly endpoint descriptions
- `@test-driven-development` — ensuring API behavior matches documented contracts

---

## External resources

- https://swagger.io/specification/
- https://restfulapi.net/
- https://graphql.org/learn/
- https://www.apiguide.com/
- https://learning.postman.com/docs/
