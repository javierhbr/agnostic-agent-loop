APP_NAME := platform-truth-mcp
PKG := ./cmd/platform-truth-mcp
BIN_DIR ?= $(CURDIR)/bin
BUILD_BIN := $(BIN_DIR)/$(APP_NAME)
PREFIX ?= $(HOME)/.local
INSTALL_DIR ?= $(PREFIX)/bin
CONFIG ?= ./examples/demo-platform-mcp-config.yaml
COMPONENT_REPO ?= ../unified-sdd-methodology/example/component-repo
RUNTIME_DIR ?= $(CURDIR)/.runtime
PID_FILE ?= $(RUNTIME_DIR)/$(APP_NAME).pid
CHILD_PID_FILE ?= $(RUNTIME_DIR)/$(APP_NAME).child.pid
REQUEST_PIPE ?= $(RUNTIME_DIR)/$(APP_NAME).requests.pipe
LOG_FILE ?= $(RUNTIME_DIR)/$(APP_NAME).log
HARNESS := ./scripts/dev-harness.sh
GO ?= go
GOCACHE ?= $(CURDIR)/.gocache

.PHONY: help setup fmt test build install run-demo start status stop smoke-test clean uninstall check-config check-component

help:
	@echo "Targets:"
	@echo "  make setup      Create local build/cache directories"
	@echo "  make fmt        Run gofmt on the module"
	@echo "  make test       Run go test with a local GOCACHE"
	@echo "  make build      Build the MCP server into ./bin"
	@echo "  make install    Install the binary into $(INSTALL_DIR)"
	@echo "  make run-demo   Run the server with the demo config"
	@echo "  make start      Start the local stdio harness for manual MCP testing"
	@echo "  make status     Show whether the local stdio harness is running"
	@echo "  make stop       Stop the local stdio harness"
	@echo "  make smoke-test Run initialize + alignment validation against the example component"
	@echo "  make clean      Remove local build artifacts"
	@echo "  make uninstall  Remove the installed binary from $(INSTALL_DIR)"

setup:
	@mkdir -p $(BIN_DIR) $(GOCACHE) $(INSTALL_DIR) $(RUNTIME_DIR)
	@echo "Prepared $(BIN_DIR), $(GOCACHE), $(INSTALL_DIR), and $(RUNTIME_DIR)"

check-config:
	@test -f "$(CONFIG)" || (echo "Config not found: $(CONFIG)"; \
		echo "Use CONFIG=/absolute/path/to/platform-mcp-config.yaml or the demo config at ./examples/demo-platform-mcp-config.yaml"; \
		exit 1)

check-component:
	@test -d "$(COMPONENT_REPO)" || (echo "Component repo not found: $(COMPONENT_REPO)"; \
		echo "Set COMPONENT_REPO=/absolute/path/to/your-component-repo"; \
		exit 1)

fmt:
	@gofmt -w .

test: setup
	@GOCACHE=$(GOCACHE) $(GO) test ./...

build: setup
	@GOCACHE=$(GOCACHE) $(GO) build -o $(BUILD_BIN) $(PKG)
	@echo "Built $(BUILD_BIN)"

install: build
	@mkdir -p $(INSTALL_DIR)
	@cp $(BUILD_BIN) $(INSTALL_DIR)/$(APP_NAME)
	@chmod 755 $(INSTALL_DIR)/$(APP_NAME)
	@echo "Installed $(INSTALL_DIR)/$(APP_NAME)"

run-demo: build check-config
	@$(BUILD_BIN) serve --config $(CONFIG)

start: build check-config
	@APP_NAME="$(APP_NAME)" \
	BUILD_BIN="$(BUILD_BIN)" \
	CONFIG="$(CONFIG)" \
	RUNTIME_DIR="$(RUNTIME_DIR)" \
	PID_FILE="$(PID_FILE)" \
	CHILD_PID_FILE="$(CHILD_PID_FILE)" \
	REQUEST_PIPE="$(REQUEST_PIPE)" \
	LOG_FILE="$(LOG_FILE)" \
	"$(HARNESS)" start

status:
	@APP_NAME="$(APP_NAME)" \
	BUILD_BIN="$(BUILD_BIN)" \
	CONFIG="$(CONFIG)" \
	RUNTIME_DIR="$(RUNTIME_DIR)" \
	PID_FILE="$(PID_FILE)" \
	CHILD_PID_FILE="$(CHILD_PID_FILE)" \
	REQUEST_PIPE="$(REQUEST_PIPE)" \
	LOG_FILE="$(LOG_FILE)" \
	"$(HARNESS)" status

stop:
	@APP_NAME="$(APP_NAME)" \
	BUILD_BIN="$(BUILD_BIN)" \
	CONFIG="$(CONFIG)" \
	RUNTIME_DIR="$(RUNTIME_DIR)" \
	PID_FILE="$(PID_FILE)" \
	CHILD_PID_FILE="$(CHILD_PID_FILE)" \
	REQUEST_PIPE="$(REQUEST_PIPE)" \
	LOG_FILE="$(LOG_FILE)" \
	"$(HARNESS)" stop

smoke-test: build check-config check-component
	@printf '%s\n%s\n%s\n' \
	  '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"make","version":"dev"}}}' \
	  '{"jsonrpc":"2.0","method":"notifications/initialized","params":{}}' \
	  '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"validate_component_alignment","arguments":{"componentRepoPath":"$(COMPONENT_REPO)"}}}' \
	  | $(BUILD_BIN) serve --config $(CONFIG)

clean:
	@rm -rf $(BIN_DIR) $(GOCACHE) $(RUNTIME_DIR)
	@echo "Removed local build artifacts"

uninstall:
	@rm -f $(INSTALL_DIR)/$(APP_NAME)
	@echo "Removed $(INSTALL_DIR)/$(APP_NAME)"
