# ── Variables ──────────────────────────────────────────────────────
BINARY_NAME  := agentic-agent
CMD_PKG      := ./cmd/agentic-agent
BUILD_DIR    := build
VERSION      := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT       := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE   := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS      := -ldflags "-X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.BuildDate=$(BUILD_DATE)"

.PHONY: help build build-all install run run-example clean \
        test test-verbose test-functional test-integration test-bdd test-bdd-verbose test-all \
        test-coverage coverage-html coverage-func coverage-summary coverage-all clean-coverage \
        lint vet fmt check \
        deps deps-tidy deps-update \
        bdd-init

# ── Help ──────────────────────────────────────────────────────────
help:
	@echo ""
	@echo "Usage: make <target>"
	@echo ""
	@echo "Build:"
	@echo "  build              Build the binary to $(BUILD_DIR)/$(BINARY_NAME)"
	@echo "  build-all          Cross-compile for darwin/linux (amd64 + arm64)"
	@echo "  install            Install to GOPATH/bin"
	@echo "  clean              Remove build artifacts and example binaries"
	@echo ""
	@echo "Run:"
	@echo "  run                Run the CLI (usage: make run ARGS='task list')"
	@echo "  run-example        Copy binary to example dir (usage: make run-example EXAMPLE=tdd)"
	@echo ""
	@echo "Test:"
	@echo "  test               Run all tests"
	@echo "  test-verbose       Run all tests with verbose output"
	@echo "  test-functional    Run functional CLI tests"
	@echo "  test-integration   Run integration tests"
	@echo "  test-bdd           Run BDD feature tests"
	@echo "  test-bdd-verbose   Run BDD tests with verbose output"
	@echo "  test-all           Run all test types (unit + functional + integration + BDD)"
	@echo ""
	@echo "Coverage:"
	@echo "  test-coverage      Run tests with basic coverage"
	@echo "  coverage-html      Generate HTML coverage report and open in browser"
	@echo "  coverage-func      Show coverage by function"
	@echo "  coverage-summary   Show coverage summary by package"
	@echo "  coverage-all       Run all tests with merged coverage"
	@echo "  clean-coverage     Remove coverage files"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint               Run golangci-lint"
	@echo "  vet                Run go vet"
	@echo "  fmt                Format code with gofmt"
	@echo "  check              Run vet + fmt check + test (pre-commit)"
	@echo ""
	@echo "Dependencies:"
	@echo "  deps               Download dependencies"
	@echo "  deps-tidy          Tidy go.mod and go.sum"
	@echo "  deps-update        Update all dependencies"
	@echo ""
	@echo "Setup:"
	@echo "  bdd-init           Install godog and setup BDD infrastructure"
	@echo ""

# ── Build ─────────────────────────────────────────────────────────

build:
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	@go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PKG)
	@echo "Built $(BUILD_DIR)/$(BINARY_NAME)"

build-all:
	@echo "Building $(BINARY_NAME) $(VERSION) for all platforms..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin  GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64  $(CMD_PKG)
	@GOOS=darwin  GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64  $(CMD_PKG)
	@GOOS=linux   GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64   $(CMD_PKG)
	@GOOS=linux   GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64   $(CMD_PKG)
	@echo "Binaries:"
	@ls -lh $(BUILD_DIR)/$(BINARY_NAME)-*

install:
	@echo "Installing $(BINARY_NAME)..."
	@go install $(LDFLAGS) $(CMD_PKG)
	@echo "Installed to $$(go env GOPATH)/bin/$(BINARY_NAME)"

# ── Run ───────────────────────────────────────────────────────────

run: build
	@$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

run-example: build
	@if [ -z "$(EXAMPLE)" ]; then \
		echo "Usage: make run-example EXAMPLE=<name>"; \
		echo ""; \
		echo "Available examples:"; \
		ls -d examples/*/  | sed 's|examples/||;s|/||' | grep -v foo | while read e; do echo "  $$e"; done; \
		exit 1; \
	fi
	@cp $(BUILD_DIR)/$(BINARY_NAME) examples/$(EXAMPLE)/$(BINARY_NAME)
	@echo "Binary copied to examples/$(EXAMPLE)/$(BINARY_NAME)"
	@echo "  cd examples/$(EXAMPLE) && ./$(BINARY_NAME)"

# ── Clean ─────────────────────────────────────────────────────────

clean: clean-coverage
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)/
	@rm -f examples/*/$(BINARY_NAME)
	@echo "Clean"

# ── Test ──────────────────────────────────────────────────────────

test:
	@echo "Running all tests..."
	@go test ./...

test-verbose:
	@echo "Running all tests with verbose output..."
	@go test -v ./...

test-functional:
	@echo "Running functional CLI tests..."
	@go test -v ./test/functional

test-integration:
	@echo "Running integration tests..."
	@go test -v ./test/integration

test-bdd:
	@echo "Running BDD feature tests..."
	@go test ./test/bdd -v

test-bdd-verbose:
	@echo "Running BDD tests with verbose output..."
	@go test ./test/bdd -v -tags=godog

test-all: test test-functional test-integration test-bdd
	@echo "All tests completed!"

# ── Coverage ──────────────────────────────────────────────────────

test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -cover

coverage-html:
	@echo "Generating HTML coverage report..."
	@mkdir -p $(BUILD_DIR)/coverage
	@go test ./... -coverprofile=$(BUILD_DIR)/coverage/coverage.out -covermode=count
	@go tool cover -html=$(BUILD_DIR)/coverage/coverage.out -o $(BUILD_DIR)/coverage/coverage.html
	@echo "Coverage report generated: $(BUILD_DIR)/coverage/coverage.html"
	@echo "Opening coverage report in browser..."
	@which open > /dev/null && open $(BUILD_DIR)/coverage/coverage.html || \
	 which xdg-open > /dev/null && xdg-open $(BUILD_DIR)/coverage/coverage.html || \
	 echo "Please open $(BUILD_DIR)/coverage/coverage.html manually"

coverage-func:
	@echo "Generating coverage by function..."
	@mkdir -p $(BUILD_DIR)/coverage
	@go test ./... -coverprofile=$(BUILD_DIR)/coverage/coverage.out -covermode=count
	@go tool cover -func=$(BUILD_DIR)/coverage/coverage.out

coverage-summary:
	@echo "Generating coverage summary..."
	@mkdir -p $(BUILD_DIR)/coverage
	@go test ./... -coverprofile=$(BUILD_DIR)/coverage/coverage.out -covermode=count > /dev/null 2>&1
	@echo ""
	@echo "=== Test Coverage Summary ==="
	@echo ""
	@go tool cover -func=$(BUILD_DIR)/coverage/coverage.out | grep total | awk '{print "Total Coverage: " $$3}'
	@echo ""
	@echo "=== Coverage by Package ==="
	@go test ./... -coverprofile=$(BUILD_DIR)/coverage/coverage.out -covermode=count 2>&1 | grep coverage: | sort -t: -k2 -rn
	@echo ""
	@echo "Detailed report available at: $(BUILD_DIR)/coverage/coverage.html"
	@echo "Run 'make coverage-html' to generate and open the HTML report"

clean-coverage:
	@rm -rf $(BUILD_DIR)/coverage/

coverage-all:
	@echo "Running all tests with coverage..."
	@mkdir -p $(BUILD_DIR)/coverage
	@echo "Running unit tests..."
	@go test ./internal/... ./pkg/... ./cmd/... -coverprofile=$(BUILD_DIR)/coverage/unit-coverage.out -covermode=count 2>&1 | grep -v "no test files"
	@echo "Running functional tests..."
	@go test ./test/functional -coverprofile=$(BUILD_DIR)/coverage/functional-coverage.out -covermode=count
	@echo "Running BDD tests..."
	@go test ./test/bdd -coverprofile=$(BUILD_DIR)/coverage/bdd-coverage.out -covermode=count
	@echo "Merging coverage reports..."
	@echo "mode: count" > $(BUILD_DIR)/coverage/merged-coverage.out
	@tail -q -n +2 $(BUILD_DIR)/coverage/unit-coverage.out $(BUILD_DIR)/coverage/functional-coverage.out $(BUILD_DIR)/coverage/bdd-coverage.out 2>/dev/null >> $(BUILD_DIR)/coverage/merged-coverage.out || true
	@go tool cover -html=$(BUILD_DIR)/coverage/merged-coverage.out -o $(BUILD_DIR)/coverage/merged-coverage.html
	@echo ""
	@echo "=== Coverage Summary ==="
	@go tool cover -func=$(BUILD_DIR)/coverage/merged-coverage.out | grep total | awk '{print "Total Coverage: " $$3}'
	@echo ""
	@echo "Individual reports:"
	@echo "  - Unit tests:       $(BUILD_DIR)/coverage/unit-coverage.html"
	@echo "  - Functional tests: $(BUILD_DIR)/coverage/functional-coverage.html"
	@echo "  - BDD tests:        $(BUILD_DIR)/coverage/bdd-coverage.html"
	@echo "  - Merged report:    $(BUILD_DIR)/coverage/merged-coverage.html"
	@echo ""
	@go tool cover -html=$(BUILD_DIR)/coverage/unit-coverage.out -o $(BUILD_DIR)/coverage/unit-coverage.html 2>/dev/null || true
	@go tool cover -html=$(BUILD_DIR)/coverage/functional-coverage.out -o $(BUILD_DIR)/coverage/functional-coverage.html 2>/dev/null || true
	@go tool cover -html=$(BUILD_DIR)/coverage/bdd-coverage.out -o $(BUILD_DIR)/coverage/bdd-coverage.html 2>/dev/null || true

# ── Code Quality ──────────────────────────────────────────────────

lint:
	@echo "Running golangci-lint..."
	@golangci-lint run ./...

vet:
	@echo "Running go vet..."
	@go vet ./...

fmt:
	@echo "Formatting code..."
	@gofmt -s -w .

check: vet
	@echo "Checking formatting..."
	@test -z "$$(gofmt -l .)" || (echo "Files need formatting:"; gofmt -l .; exit 1)
	@echo "Running tests..."
	@go test ./...
	@echo "All checks passed"

# ── Dependencies ──────────────────────────────────────────────────

deps:
	@echo "Downloading dependencies..."
	@go mod download

deps-tidy:
	@echo "Tidying dependencies..."
	@go mod tidy

deps-update:
	@echo "Updating dependencies..."
	@go get -u ./...
	@go mod tidy

# ── Setup ─────────────────────────────────────────────────────────

bdd-init:
	@echo "Installing godog..."
	@go get github.com/cucumber/godog/cmd/godog@latest
	@go mod tidy
	@echo "Creating feature directories..."
	@mkdir -p features/{init,tasks,context,validation,workflows}
	@mkdir -p test/bdd/{steps,helpers}
	@echo "BDD infrastructure initialized!"
