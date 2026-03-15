# Platform Truth MCP Server

This is the v1 local read-only MCP server for the unified SDD methodology.

It is designed for developer machines that keep a local clone of the platform
repository and want fast local query and validation without hosted
infrastructure.

## Design goals

- Go single binary
- stdio transport
- newline-delimited JSON-RPC messages
- read-only only
- no external Go dependencies
- version-aware validation based on `platform-ref.yaml`

## Implemented tools

- `get_platform_version`
- `list_platform_refs`
- `get_platform_ref`
- `get_jira_mapping`
- `validate_component_alignment`
- `validate_component_jira_chain`
- `detect_platform_drift_from_pinned_version`

## Quick start

From `/Users/javierbenavides/others/dev/poc/random-poc/custom-skills/sdd/platform-truth-mcp-server`:

```bash
make setup
make test
make build
make install
make start
make status
make stop
```

## Make targets

Useful targets:

- `make setup` prepares `./bin`, `./.gocache`, and the local install directory
- `make fmt` runs `gofmt -w .`
- `make test` runs `go test ./...` with a local `GOCACHE`
- `make build` writes the binary to `./bin/platform-truth-mcp`
- `make install` copies the binary to `$HOME/.local/bin/platform-truth-mcp` by default
- `make run-demo` starts the server with the demo config
- `make start` starts a local stdio harness for manual MCP testing
- `make status` shows whether the local stdio harness is running
- `make stop` stops the local stdio harness
- `make smoke-test` runs `initialize` and `validate_component_alignment` against the sample component repo
- `make uninstall` removes the installed binary from the install directory

Override install paths when needed:

```bash
make install PREFIX=/tmp/platform-mcp
make install INSTALL_DIR=/usr/local/bin
```

The `start` target is a local dev harness, not a network daemon. It keeps the
stdio server alive behind a request pipe so you can inspect liveness and logs on
a developer machine. Normal MCP client integration should still launch the
binary directly.

Runtime files:

- request pipe: `./.runtime/platform-truth-mcp.requests.pipe`
- log file: `./.runtime/platform-truth-mcp.log`
- supervisor pid: `./.runtime/platform-truth-mcp.pid`

If you override `CONFIG` or `COMPONENT_REPO`, use real local paths instead of
the placeholder `/absolute/path/...` examples.

## How to test the server

Use three levels of testing.

### 1. Unit tests and build

```bash
make test
make build
```

This checks:

- the Go tests pass
- the binary builds into `./bin/platform-truth-mcp`

### 2. End-to-end smoke test with the bundled example

```bash
make smoke-test
```

This uses:

- config: `./examples/demo-platform-mcp-config.yaml`
- component repo: `../unified-sdd-methodology/example/component-repo`

Expected result:

- one `initialize` response
- one `validate_component_alignment` response
- `"valid": true` in the returned JSON payload

### 3. Start the local stdio harness

```bash
make start
make status
make stop
```

Use this when you want the server running locally while you send manual JSON-RPC
requests or connect a local MCP client.

## Manual MCP testing

With the harness running, inspect the log in one terminal:

```bash
tail -f ./.runtime/platform-truth-mcp.log
```

Send requests into the request pipe from another terminal:

```bash
printf '%s\n' \
'{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"manual-test","version":"1.0"}}}' \
> ./.runtime/platform-truth-mcp.requests.pipe

printf '%s\n' \
'{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}' \
> ./.runtime/platform-truth-mcp.requests.pipe

printf '%s\n' \
'{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' \
> ./.runtime/platform-truth-mcp.requests.pipe
```

Try a real tool call:

```bash
printf '%s\n' \
'{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"validate_component_alignment","arguments":{"componentRepoPath":"../unified-sdd-methodology/example/component-repo"}}}' \
> ./.runtime/platform-truth-mcp.requests.pipe
```

Read the responses in `./.runtime/platform-truth-mcp.log`.

## Test with your own platform repo and component repo

Run the smoke test with real local paths:

```bash
make smoke-test \
  CONFIG=/real/path/to/platform-mcp-config.yaml \
  COMPONENT_REPO=/real/path/to/your-component-repo
```

Your component repo should contain:

- `platform-ref.yaml`
- `jira-traceability.yaml`

If the config path is wrong, `make smoke-test` now fails with a direct message
that tells you which path is missing.

## Test with the platform template in this repo

Use this config:

- `./examples/platform-template-mcp-config.yaml`

It points to:

- platform repo: `/Users/javierbenavides/others/dev/poc/random-poc/custom-skills/sdd/unified-sdd-methodology/templates/platform-template`

Important constraint:

- the platform template is still a scaffold
- `platform-baseline.md` and the component boilerplate contain placeholders
- use `sample-platform-baseline.md` for version discovery and ref scanning
- expect meaningful alignment results only after you replace placeholders in a real copied repo

### Platform-only tests against the template repo

Build and start the server against the platform template:

```bash
make build
make start CONFIG=./examples/platform-template-mcp-config.yaml
make status CONFIG=./examples/platform-template-mcp-config.yaml
```

Initialize the session:

```bash
printf '%s\n' \
'{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"platform-template-test","version":"1.0"}}}' \
> ./.runtime/platform-truth-mcp.requests.pipe

printf '%s\n' \
'{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}' \
> ./.runtime/platform-truth-mcp.requests.pipe
```

List tools:

```bash
printf '%s\n' \
'{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}' \
> ./.runtime/platform-truth-mcp.requests.pipe
```

Get the latest platform version from the sample platform baseline:

```bash
printf '%s\n' \
'{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"get_platform_version","arguments":{"mode":"latest"}}}' \
> ./.runtime/platform-truth-mcp.requests.pipe
```

List discoverable refs:

```bash
printf '%s\n' \
'{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"list_platform_refs","arguments":{}}}' \
> ./.runtime/platform-truth-mcp.requests.pipe
```

Fetch one specific ref from the template:

```bash
printf '%s\n' \
'{"jsonrpc":"2.0","id":5,"method":"tools/call","params":{"name":"get_platform_ref","arguments":{"id":"contracts.customer-profile.v2"}}}' \
> ./.runtime/platform-truth-mcp.requests.pipe
```

Inspect the responses:

```bash
tail -f ./.runtime/platform-truth-mcp.log
```

Stop the harness when done:

```bash
make stop CONFIG=./examples/platform-template-mcp-config.yaml
```

### Template smoke test with the component boilerplate

This is useful to verify wiring, not full business alignment:

```bash
make smoke-test \
  CONFIG=./examples/platform-template-mcp-config.yaml \
  COMPONENT_REPO=../unified-sdd-methodology/templates/component-boilerplate
```

Because the boilerplate still contains placeholders, treat this as a structure
test. For a real alignment test, copy the template, replace placeholders, and
run the same command against the filled component repo.

## Run without Make

```bash
go run ./cmd/platform-truth-mcp serve --config ./examples/demo-platform-mcp-config.yaml
```

## Example config

- `examples/demo-platform-mcp-config.yaml`
- `examples/platform-template-mcp-config.yaml`

## Notes

- the server reads platform truth from a local platform clone
- the server reads component alignment from `platform-ref.yaml` and
  `jira-traceability.yaml`
- the server defaults to pinned-version behavior for component validation
