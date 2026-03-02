# api-docs — Detailed Examples & Reference

This file is the detailed companion to `../SKILL.md`. It contains full worked examples, expanded best practice explanations, complete section body copy, common pitfalls, and tool format snippets. Reference this when the slim SKILL.md does not provide enough context for a specific situation.

---

## Example 1: REST API Endpoint Documentation

The output below is the exact format to produce for a POST endpoint. Replicate this structure for every endpoint in the API surface.

### Create User

Creates a new user account.

**Endpoint:** `POST /api/v1/users`

**Authentication:** Required (Bearer token)

**Request Body:**
```json
{
  "email": "user@example.com",      // Required: Valid email address
  "password": "SecurePass123!",     // Required: Min 8 chars, 1 uppercase, 1 number
  "name": "John Doe",               // Required: 2-50 characters
  "role": "user"                    // Optional: "user" or "admin" (default: "user")
}
```

**Success Response (201 Created):**
```json
{
  "id": "usr_1234567890",
  "email": "user@example.com",
  "name": "John Doe",
  "role": "user",
  "createdAt": "2026-01-20T10:30:00Z",
  "emailVerified": false
}
```

**Error Responses:**

- `400 Bad Request` — Invalid input data
  ```json
  {
    "error": "VALIDATION_ERROR",
    "message": "Invalid email format",
    "field": "email"
  }
  ```

- `409 Conflict` — Email already exists
  ```json
  {
    "error": "EMAIL_EXISTS",
    "message": "An account with this email already exists"
  }
  ```

- `401 Unauthorized` — Missing or invalid authentication token

**Example Request (cURL):**
```bash
curl -X POST https://api.example.com/api/v1/users \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "SecurePass123!",
    "name": "John Doe"
  }'
```

**Example Request (JavaScript):**
```javascript
const response = await fetch('https://api.example.com/api/v1/users', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    email: 'user@example.com',
    password: 'SecurePass123!',
    name: 'John Doe'
  })
});

const user = await response.json();
console.log(user);
```

**Example Request (Python):**
```python
import requests

response = requests.post(
    'https://api.example.com/api/v1/users',
    headers={
        'Authorization': f'Bearer {token}',
        'Content-Type': 'application/json'
    },
    json={
        'email': 'user@example.com',
        'password': 'SecurePass123!',
        'name': 'John Doe'
    }
)

user = response.json()
print(user)
```

---

## Example 2: GraphQL API Documentation

Use this format for GraphQL queries and mutations. Note that error structure differs from REST — GraphQL always returns 200 OK and puts errors in a top-level `errors` array.

### User Query

Fetch user information by ID.

**Query:**
```graphql
query GetUser($id: ID!) {
  user(id: $id) {
    id
    email
    name
    role
    createdAt
    posts {
      id
      title
      publishedAt
    }
  }
}
```

**Variables:**
```json
{
  "id": "usr_1234567890"
}
```

**Response:**
```json
{
  "data": {
    "user": {
      "id": "usr_1234567890",
      "email": "user@example.com",
      "name": "John Doe",
      "role": "user",
      "createdAt": "2026-01-20T10:30:00Z",
      "posts": [
        {
          "id": "post_123",
          "title": "My First Post",
          "publishedAt": "2026-01-21T14:00:00Z"
        }
      ]
    }
  }
}
```

**Errors:**
```json
{
  "errors": [
    {
      "message": "User not found",
      "extensions": {
        "code": "USER_NOT_FOUND",
        "userId": "usr_1234567890"
      }
    }
  ]
}
```

---

## Example 3: Authentication Flow Documentation

Document authentication as a first-class section, not an afterthought. Include token acquisition, usage, and refresh. Every auth step needs a working request + response pair.

### Getting a Token

**Endpoint:** `POST /api/v1/auth/login`

**Request:**
```json
{
  "email": "user@example.com",
  "password": "your-password"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expiresIn": 3600,
  "refreshToken": "refresh_token_here"
}
```

### Using the Token

Include the token in the Authorization header for every protected request:

```
Authorization: Bearer YOUR_TOKEN
```

### Token Expiration

Tokens expire after 1 hour. Use the refresh token to obtain a new access token without requiring the user to log in again.

**Endpoint:** `POST /api/v1/auth/refresh`

**Request:**
```json
{
  "refreshToken": "refresh_token_here"
}
```

---

## Expanded Best Practices

These are the full explanations for each bullet in SKILL.md. Read these when a 1-line rule is not enough context to apply it correctly.

### Be Consistent

Apply the same structure to every endpoint: same heading order, same field naming in examples, same way of indicating required vs optional. Readers scan docs, they do not read them linearly. If section order shifts between endpoints they will miss critical details.

### Include Examples with Realistic Data

Replace placeholder strings ("foo", "bar", "string", "123") with data that mirrors production shapes. A user ID like `usr_1234567890` immediately tells the reader the format. An email like `user@example.com` shows validation intent. Realistic data also helps readers catch bugs — if your example shows a date as `"2026-01-20"` but the API actually returns `"2026-01-20T10:30:00Z"`, the mismatch surfaces before someone files a bug.

### Document Every Error Code

Every HTTP status code and application-level error code the API can return must appear in the docs with: what triggers it, what the response body looks like, and what the caller should do about it. Error documentation is often the most-read section for developers integrating your API.

### Version API Paths

Include the version segment (/api/v1/) in every URL shown. State the current version prominently in the Introduction section. When a breaking change is shipped, document it in the Changelog and leave the old version endpoints documented (even if deprecated) so integrators know what changed and when.

### Test Every Code Example

Run each cURL, JavaScript, and Python snippet against a real or sandbox environment before publishing. Broken examples are worse than no examples — they destroy trust and generate support load. If a sandbox is not available, add a clear note indicating the endpoint requires a real account.

### Add Rate Limit and Pagination Details

Document the X-RateLimit-Limit, X-RateLimit-Remaining, and X-RateLimit-Reset headers. Show a 429 response example. For paginated endpoints, document all pagination parameters (page, limit, cursor, offset) and show the envelope structure including total count, next/previous links, or cursor values.

### Link Related Endpoints

At the bottom of each endpoint section, add a "See also" or "Related" list pointing to endpoints that are commonly chained. For example, Create User should link to Get User, Update User, and the authentication guide. This reduces context-switching and helps readers discover the full surface area organically.

---

## Documentation Section Body Copy

This section provides full body copy for each of the nine recommended sections. Use these as starting templates and fill in API-specific details.

### 1. Introduction

```markdown
## Introduction

[API Name] provides a [REST|GraphQL|WebSocket] interface for [core capability].
All requests are made to the base URL:

**Base URL:** `https://api.example.com`

**Current Version:** v1

**API Style:** REST over HTTPS. All request and response bodies are JSON unless
otherwise noted.

**Support:** api-support@example.com | [Status Page](https://status.example.com)
```

### 2. Authentication

```markdown
## Authentication

[API Name] uses Bearer token authentication. Include your token in the
Authorization header of every request:

```
Authorization: Bearer YOUR_TOKEN
```

Tokens are obtained via `POST /api/v1/auth/login`. They expire after 1 hour.
Use `POST /api/v1/auth/refresh` to obtain a new token without re-authenticating.

See [Authentication Flow](#authentication-flow) for complete request/response details.
```

### 3. Quick Start

```markdown
## Quick Start

The fastest way to make your first request:

1. Obtain a token (see [Authentication](#authentication))
2. Call an endpoint:

```bash
curl -X GET https://api.example.com/api/v1/users/me \
  -H "Authorization: Bearer YOUR_TOKEN"
```

3. Parse the JSON response.

**Common use case:** [Describe the most typical integration scenario in 2-3 sentences.]
```

### 4. Endpoints

```markdown
## Endpoints

Endpoints are organized by resource. Each section covers a single resource type
and lists all available operations.

- [Users](#users) — Create, read, update, and delete user accounts
- [Posts](#posts) — Manage published and draft content
- [Auth](#auth) — Token management

[Expand with actual resource list]
```

### 5. Data Models

```markdown
## Data Models

### User

| Field         | Type     | Required | Description                          |
|---------------|----------|----------|--------------------------------------|
| id            | string   | Yes      | Unique identifier, format: usr_*     |
| email         | string   | Yes      | Valid RFC 5321 email address         |
| name          | string   | Yes      | Display name, 2-50 characters        |
| role          | string   | Yes      | One of: user, admin                  |
| createdAt     | string   | Yes      | ISO 8601 timestamp                   |
| emailVerified | boolean  | Yes      | Whether email has been confirmed     |
```

### 6. Error Handling

```markdown
## Error Handling

All errors follow a consistent response format:

```json
{
  "error": "ERROR_CODE",
  "message": "Human-readable description",
  "field": "fieldName"   // Present only for validation errors
}
```

**Standard HTTP Status Codes:**

| Code | Meaning              | When it occurs                              |
|------|----------------------|---------------------------------------------|
| 400  | Bad Request          | Invalid input, missing required fields      |
| 401  | Unauthorized         | Missing, expired, or invalid token          |
| 403  | Forbidden            | Valid token but insufficient permissions    |
| 404  | Not Found            | Resource does not exist                     |
| 409  | Conflict             | Duplicate resource (e.g. email exists)      |
| 422  | Unprocessable Entity | Input format valid but semantically invalid |
| 429  | Too Many Requests    | Rate limit exceeded                         |
| 500  | Internal Server Error| Unexpected server-side failure              |
```

### 7. Rate Limiting

```markdown
## Rate Limiting

The API enforces rate limits to ensure fair use. Limits are applied per API key.

**Default limits:** 1000 requests per hour.

**Rate limit headers returned with every response:**

| Header                  | Description                                |
|-------------------------|--------------------------------------------|
| X-RateLimit-Limit       | Maximum requests allowed in the window     |
| X-RateLimit-Remaining   | Requests remaining in current window       |
| X-RateLimit-Reset       | Unix timestamp when the window resets      |

When you exceed the limit, the API returns `429 Too Many Requests`. Wait until
X-RateLimit-Reset before retrying. Use exponential backoff for automated clients.
```

### 8. Changelog

```markdown
## Changelog

### v1.1.0 — 2026-02-01
- Added `role` field to User response
- Added `POST /api/v1/auth/refresh` endpoint

### v1.0.0 — 2026-01-01
- Initial public release

**Deprecation policy:** Deprecated endpoints are supported for 6 months after
the deprecation notice. Breaking changes are introduced only in major versions.
```

### 9. SDKs and Tools

```markdown
## SDKs and Tools

**Official client libraries:**
- JavaScript/TypeScript: `npm install @example/api-client`
- Python: `pip install example-api`

**Postman Collection:** [Download](https://api.example.com/postman-collection.json)

**OpenAPI Specification:** [Download YAML](https://api.example.com/openapi.yaml)
| [View in Swagger UI](https://api.example.com/docs)
```

---

## Common Pitfalls

### Pitfall 1: Documentation Gets Out of Sync

**Problem:** Over time the codebase evolves but the documentation is not updated in the same PRs. The docs begin to describe a different API than what is actually deployed.

**Symptoms:**
- Code examples return errors or unexpected shapes
- Parameters listed in docs are rejected or silently ignored
- New fields appear in responses without documentation
- Removed endpoints still appear as available

**Solution:**
- Generate documentation from code annotations (JSDoc, OpenAPI decorators, godoc comments) so docs update automatically when code changes
- Use tooling like Swagger UI, Redoc, or Stoplight that renders live from an OpenAPI YAML/JSON spec committed to the repository
- Add a CI gate that validates the OpenAPI spec against actual API responses using tools like Dredd or Schemathesis
- Treat documentation updates as a required part of any PR that changes API behaviour — block merge if docs are not updated

---

### Pitfall 2: Missing Error Documentation

**Problem:** Only the happy path is documented. Error responses are undocumented or described only with status codes and no body examples.

**Symptoms:**
- Integrators open support tickets asking what a specific error code means
- Client applications display raw error JSON to end users because they did not know the error format in advance
- Integrators write defensive catch-all error handlers instead of targeted ones because they cannot enumerate possible failures

**Solution:**
- For every endpoint, enumerate every HTTP status code it can return
- For each error, provide: the error code constant (e.g. VALIDATION_ERROR), the English message pattern, which fields might be present (e.g. `field` for 400s), and a complete example JSON response
- Include a troubleshooting guide with the top 5 errors and their fixes
- Test error documentation the same way you test happy-path docs

---

### Pitfall 3: Examples Don't Work

**Problem:** Code examples in the documentation are untested and broken — wrong URLs, outdated request shapes, or missing required headers.

**Symptoms:**
- New integrators cannot complete the Quick Start guide
- Stack Overflow questions and GitHub issues reference documentation examples that fail
- High drop-off rate at the first code example in onboarding

**Solution:**
- Run every cURL, JavaScript, and Python example against a real or sandbox environment before publishing
- Store examples as runnable test fixtures in the repository and execute them in CI
- Provide a sandbox base URL (e.g. https://sandbox.api.example.com) so examples can be run safely without affecting production data
- If a sandbox is not available, add a prominent note on every example: "Replace YOUR_TOKEN and base URL with your credentials before running"

---

### Pitfall 4: Unclear Parameter Requirements

**Problem:** Parameters are listed without enough metadata — the reader cannot tell which are required, what format is expected, or what values are valid.

**Symptoms:**
- Integrators send invalid requests frequently, generating 400 errors
- Validation error messages reference field names users have never seen in the docs
- Optional parameters are treated as required (or vice versa), breaking integrations

**Solution:**
- For every parameter, document: name, type, whether required or optional, default value (if optional), allowed values or format, and an example value
- Use a table format for request body schemas — it is easier to scan than prose
- Mark required fields explicitly in JSON examples (inline comments or a legend)
- Show validation rules: min/max length, regex patterns, numeric ranges, enum values

---

## OpenAPI / Swagger Snippet

Use this as the starting skeleton for a full OpenAPI 3.0 specification. Expand `components/schemas` with your actual data model definitions.

```yaml
openapi: 3.0.0
info:
  title: My API
  version: 1.0.0
  description: Brief description of what this API does
  contact:
    email: api-support@example.com

servers:
  - url: https://api.example.com
    description: Production
  - url: https://sandbox.api.example.com
    description: Sandbox

paths:
  /api/v1/users:
    post:
      summary: Create a new user
      operationId: createUser
      security:
        - bearerAuth: []
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/CreateUserRequest'
      responses:
        '201':
          description: User created successfully
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '400':
          description: Invalid input
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'
        '409':
          description: Email already exists
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/ErrorResponse'

components:
  securitySchemes:
    bearerAuth:
      type: http
      scheme: bearer
      bearerFormat: JWT

  schemas:
    CreateUserRequest:
      type: object
      required: [email, password, name]
      properties:
        email:
          type: string
          format: email
          example: user@example.com
        password:
          type: string
          minLength: 8
          example: SecurePass123!
        name:
          type: string
          minLength: 2
          maxLength: 50
          example: John Doe
        role:
          type: string
          enum: [user, admin]
          default: user

    User:
      type: object
      properties:
        id:
          type: string
          example: usr_1234567890
        email:
          type: string
          example: user@example.com
        name:
          type: string
          example: John Doe
        role:
          type: string
          example: user
        createdAt:
          type: string
          format: date-time
          example: '2026-01-20T10:30:00Z'
        emailVerified:
          type: boolean
          example: false

    ErrorResponse:
      type: object
      properties:
        error:
          type: string
          example: VALIDATION_ERROR
        message:
          type: string
          example: Invalid email format
        field:
          type: string
          example: email
```

---

## Postman Collection JSON Snippet

Use this as the starting template for a Postman collection. Add one item per endpoint, grouped by resource. Use `{{baseUrl}}` and `{{token}}` environment variables so the collection works across environments without edits.

```json
{
  "info": {
    "name": "My API",
    "description": "Complete API collection for My API v1",
    "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json"
  },
  "variable": [
    {
      "key": "baseUrl",
      "value": "https://api.example.com",
      "type": "string"
    },
    {
      "key": "token",
      "value": "",
      "type": "string"
    }
  ],
  "item": [
    {
      "name": "Auth",
      "item": [
        {
          "name": "Login",
          "request": {
            "method": "POST",
            "url": "{{baseUrl}}/api/v1/auth/login",
            "header": [
              {
                "key": "Content-Type",
                "value": "application/json"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"user@example.com\",\n  \"password\": \"SecurePass123!\"\n}"
            }
          }
        }
      ]
    },
    {
      "name": "Users",
      "item": [
        {
          "name": "Create User",
          "request": {
            "method": "POST",
            "url": "{{baseUrl}}/api/v1/users",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{token}}"
              },
              {
                "key": "Content-Type",
                "value": "application/json"
              }
            ],
            "body": {
              "mode": "raw",
              "raw": "{\n  \"email\": \"user@example.com\",\n  \"password\": \"SecurePass123!\",\n  \"name\": \"John Doe\"\n}"
            }
          }
        },
        {
          "name": "Get Current User",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}/api/v1/users/me",
            "header": [
              {
                "key": "Authorization",
                "value": "Bearer {{token}}"
              }
            ]
          }
        }
      ]
    }
  ]
}
```

---

*This file is the reference companion to `../SKILL.md`. The slim SKILL.md is the agent's operating instruction. This file is the human (and agent) reference for when more detail is needed.*
