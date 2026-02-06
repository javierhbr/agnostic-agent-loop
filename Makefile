.PHONY: test test-verbose test-coverage coverage-html coverage-func coverage-summary clean-coverage test-functional help

# Help target
help:
	@echo "Available targets:"
	@echo "  test              - Run all tests"
	@echo "  test-verbose      - Run all tests with verbose output"
	@echo "  test-functional   - Run functional CLI tests"
	@echo "  test-coverage     - Run tests and generate coverage report"
	@echo "  coverage-html     - Generate HTML coverage report and open in browser"
	@echo "  coverage-func     - Show coverage by function"
	@echo "  coverage-summary  - Show coverage summary by package"
	@echo "  clean-coverage    - Clean coverage files"

# Run all tests
test:
	@echo "Running all tests..."
	@go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running all tests with verbose output..."
	@go test -v ./...

# Run functional CLI tests
test-functional:
	@echo "Running functional CLI tests..."
	@go test -v ./tests/functional

# Run tests and generate basic coverage report
test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -cover

# Generate detailed HTML coverage report
coverage-html:
	@echo "Generating HTML coverage report..."
	@mkdir -p coverage
	@go test ./... -coverprofile=coverage/coverage.out -covermode=count
	@go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Coverage report generated: coverage/coverage.html"
	@echo "Opening coverage report in browser..."
	@which open > /dev/null && open coverage/coverage.html || \
	 which xdg-open > /dev/null && xdg-open coverage/coverage.html || \
	 echo "Please open coverage/coverage.html manually"

# Show coverage by function
coverage-func:
	@echo "Generating coverage by function..."
	@mkdir -p coverage
	@go test ./... -coverprofile=coverage/coverage.out -covermode=count
	@go tool cover -func=coverage/coverage.out

# Show coverage summary by package
coverage-summary:
	@echo "Generating coverage summary..."
	@mkdir -p coverage
	@go test ./... -coverprofile=coverage/coverage.out -covermode=count > /dev/null 2>&1
	@echo ""
	@echo "=== Test Coverage Summary ==="
	@echo ""
	@go tool cover -func=coverage/coverage.out | grep total | awk '{print "Total Coverage: " $$3}'
	@echo ""
	@echo "=== Coverage by Package ==="
	@go test ./... -coverprofile=coverage/coverage.out -covermode=count 2>&1 | grep coverage: | sort -t: -k2 -rn
	@echo ""
	@echo "Detailed report available at: coverage/coverage.html"
	@echo "Run 'make coverage-html' to generate and open the HTML report"

# Clean coverage files
clean-coverage:
	@echo "Cleaning coverage files..."
	@rm -rf coverage/
	@echo "Coverage files cleaned"
