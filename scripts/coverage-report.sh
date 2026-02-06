#!/bin/bash

# Test Coverage Report Generator
# Generates comprehensive test coverage reports for the Agentic Agent project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COVERAGE_DIR="coverage"
COVERAGE_FILE="${COVERAGE_DIR}/coverage.out"
COVERAGE_HTML="${COVERAGE_DIR}/coverage.html"
COVERAGE_JSON="${COVERAGE_DIR}/coverage.json"
COVERAGE_THRESHOLD=50 # Minimum acceptable coverage percentage

# Create coverage directory
mkdir -p "${COVERAGE_DIR}"

echo -e "${BLUE}================================================${NC}"
echo -e "${BLUE}  Agentic Agent - Test Coverage Report${NC}"
echo -e "${BLUE}================================================${NC}"
echo ""

# Run tests with coverage
echo -e "${YELLOW}Running tests with coverage...${NC}"
if go test ./... -coverprofile="${COVERAGE_FILE}" -covermode=count -json > "${COVERAGE_JSON}" 2>&1; then
    echo -e "${GREEN}✓ Tests completed successfully${NC}"
else
    echo -e "${RED}✗ Tests failed${NC}"
    exit 1
fi
echo ""

# Generate coverage summary
echo -e "${BLUE}=== Test Coverage Summary ===${NC}"
echo ""

# Extract total coverage
TOTAL_COVERAGE=$(go tool cover -func="${COVERAGE_FILE}" | grep total | awk '{print $3}' | sed 's/%//')

# Display total coverage with color (using awk for float comparison)
if awk -v cov="$TOTAL_COVERAGE" -v thresh="$COVERAGE_THRESHOLD" 'BEGIN {exit !(cov >= thresh)}'; then
    echo -e "${GREEN}Total Coverage: ${TOTAL_COVERAGE}%${NC}"
else
    echo -e "${RED}Total Coverage: ${TOTAL_COVERAGE}% (below threshold of ${COVERAGE_THRESHOLD}%)${NC}"
fi
echo ""

# Coverage by package
echo -e "${BLUE}=== Coverage by Package ===${NC}"
echo ""
go test ./... -coverprofile="${COVERAGE_FILE}" -covermode=count 2>&1 | \
    grep coverage: | \
    grep -v "coverage: 0.0%" | \
    sort -t: -k2 -rn | \
    while IFS= read -r line; do
        coverage=$(echo "$line" | grep -o '[0-9.]*%' | head -1 | sed 's/%//')
        if awk -v cov="$coverage" 'BEGIN {exit !(cov >= 70)}'; then
            echo -e "${GREEN}${line}${NC}"
        elif awk -v cov="$coverage" 'BEGIN {exit !(cov >= 40)}'; then
            echo -e "${YELLOW}${line}${NC}"
        else
            echo -e "${RED}${line}${NC}"
        fi
    done
echo ""

# Packages with no tests
echo -e "${BLUE}=== Packages Without Tests ===${NC}"
echo ""
NO_TESTS=$(go test ./... -coverprofile="${COVERAGE_FILE}" -covermode=count 2>&1 | grep "coverage: 0.0%" | wc -l)
if [ "$NO_TESTS" -gt 0 ]; then
    go test ./... -coverprofile="${COVERAGE_FILE}" -covermode=count 2>&1 | \
        grep "coverage: 0.0%" | \
        sed 's/\tcoverage: 0.0% of statements//' | \
        while IFS= read -r package; do
            echo -e "${RED}  • ${package}${NC}"
        done
    echo ""
    echo -e "${YELLOW}Total packages without tests: ${NO_TESTS}${NC}"
else
    echo -e "${GREEN}All packages have tests!${NC}"
fi
echo ""

# Top uncovered packages
echo -e "${BLUE}=== Top 5 Packages Needing Attention ===${NC}"
echo ""
go tool cover -func="${COVERAGE_FILE}" | \
    grep -v "^mode:" | \
    awk '{coverage=$3; sub(/%/, "", coverage); if (coverage+0 < 100 && coverage+0 > 0) print coverage "% " $1}' | \
    sort -n | \
    head -5 | \
    while IFS= read -r line; do
        echo -e "${YELLOW}  • ${line}${NC}"
    done
echo ""

# Generate HTML report
echo -e "${YELLOW}Generating HTML coverage report...${NC}"
go tool cover -html="${COVERAGE_FILE}" -o "${COVERAGE_HTML}"
echo -e "${GREEN}✓ HTML report generated: ${COVERAGE_HTML}${NC}"
echo ""

# Summary
echo -e "${BLUE}=== Summary ===${NC}"
echo ""
echo -e "Coverage Profile: ${COVERAGE_FILE}"
echo -e "HTML Report:      ${COVERAGE_HTML}"
echo -e "JSON Output:      ${COVERAGE_JSON}"
echo ""

# Coverage badge suggestion
BADGE_COLOR="red"
if awk -v cov="$TOTAL_COVERAGE" 'BEGIN {exit !(cov >= 80)}'; then
    BADGE_COLOR="brightgreen"
elif awk -v cov="$TOTAL_COVERAGE" 'BEGIN {exit !(cov >= 60)}'; then
    BADGE_COLOR="yellow"
elif awk -v cov="$TOTAL_COVERAGE" 'BEGIN {exit !(cov >= 40)}'; then
    BADGE_COLOR="orange"
fi

echo -e "${BLUE}Coverage Badge:${NC}"
echo "[![Coverage](https://img.shields.io/badge/coverage-${TOTAL_COVERAGE}%25-${BADGE_COLOR})](./coverage/coverage.html)"
echo ""

# Open HTML report
if [[ "$1" == "--open" ]]; then
    echo -e "${YELLOW}Opening HTML report in browser...${NC}"
    if command -v open &> /dev/null; then
        open "${COVERAGE_HTML}"
    elif command -v xdg-open &> /dev/null; then
        xdg-open "${COVERAGE_HTML}"
    else
        echo -e "${RED}Cannot open browser automatically. Please open ${COVERAGE_HTML} manually.${NC}"
    fi
fi

echo -e "${BLUE}================================================${NC}"
echo -e "${GREEN}Coverage report complete!${NC}"
echo -e "${BLUE}================================================${NC}"

# Exit with error if below threshold
if awk -v cov="$TOTAL_COVERAGE" -v thresh="$COVERAGE_THRESHOLD" 'BEGIN {exit !(cov < thresh)}'; then
    exit 1
fi
