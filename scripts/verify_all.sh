#!/bin/bash
set -e

echo "Building agentic-agent..."
go build -o agentic-agent ./cmd/agentic-agent

TEST_DIR="test_workspace_$(date +%s)"
mkdir "$TEST_DIR"
cp agentic-agent "$TEST_DIR/"
cd "$TEST_DIR"

echo "1. Initializing Project..."
./agentic-agent init --name "Test Project"

echo "2. Creating Task..."
OUTPUT=$(./agentic-agent task create --title "Test Task")
TASK_ID=$(echo "$OUTPUT" | grep -o 'TASK-[0-9]*')
echo "Created Task ID: $TASK_ID"

echo "3. Claiming Task..."
./agentic-agent task claim "$TASK_ID"

echo "4. Creating Mock Source..."
mkdir -p src/testpkg
echo "package testpkg" > src/testpkg/main.go

echo "5. Generating Context..."
./agentic-agent context generate src/testpkg

echo "6. Validating..."
./agentic-agent validate

echo "7. Generating Skills..."
./agentic-agent skills generate --tool claude-code
if [ ! -f CLAUDE.md ]; then
    echo "Error: CLAUDE.md not generated"
    exit 1
fi

echo "8. Running Orchestrator..."
./agentic-agent run --task "$TASK_ID"

echo "9. Token Status..."
./agentic-agent token status

echo "✅ End-to-End Verification Passed!"
cd ..
rm -rf "$TEST_DIR"
rm -f agentic-agent
