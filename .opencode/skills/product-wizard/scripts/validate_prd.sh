#!/bin/bash

# PRD Validation Script
# Checks PRD completeness and quality

set -e

GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

VERBOSE=false
SECTIONS="all"

usage() {
    echo "Usage: $0 <prd_file.md> [options]"
    echo ""
    echo "Options:"
    echo "  --verbose           Show detailed suggestions"
    echo "  --sections <list>   Check specific sections (comma-separated)"
    echo "                      e.g., --sections user-stories,metrics"
    echo ""
    echo "Example:"
    echo "  $0 my_prd.md"
    echo "  $0 my_prd.md --verbose"
    echo "  $0 my_prd.md --sections user-stories,metrics"
    exit 1
}

if [ $# -lt 1 ]; then
    usage
fi

PRD_FILE="$1"
shift

while [[ $# -gt 0 ]]; do
    case $1 in
        --verbose) VERBOSE=true; shift ;;
        --sections) SECTIONS="$2"; shift 2 ;;
        *) echo -e "${RED}Unknown option: $1${NC}"; usage ;;
    esac
done

if [ ! -f "$PRD_FILE" ]; then
    echo -e "${RED}✗ Error: File not found: $PRD_FILE${NC}"
    exit 1
fi

echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║        PRD Validation Report           ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "File: ${BLUE}$PRD_FILE${NC}"
echo ""

ISSUES=0
WARNINGS=0
PASSED=0

check_section() {
    local name="$1" pattern="$2" required="$3"
    if grep -qi "$pattern" "$PRD_FILE"; then
        echo -e "${GREEN}✓${NC} $name found"
        ((PASSED++))
        return 0
    else
        if [ "$required" = "true" ]; then
            echo -e "${RED}✗${NC} $name missing (REQUIRED)"
            ((ISSUES++))
        else
            echo -e "${YELLOW}⚠${NC} $name missing (recommended)"
            ((WARNINGS++))
        fi
        return 1
    fi
}

check_content() {
    local name="$1" pattern="$2" msg="$3"
    if grep -q "$pattern" "$PRD_FILE"; then
        echo -e "${YELLOW}⚠${NC} $name: $msg"
        ((WARNINGS++))
    else
        echo -e "${GREEN}✓${NC} $name passed"
        ((PASSED++))
    fi
}

# === Required Sections ===
echo -e "${BLUE}━━━ Required Sections ━━━${NC}"
echo ""
check_section "Problem Statement" "problem statement\|## problem\|the problem" true
check_section "Goals & Objectives" "goals.*objectives\|## goals\|business goals" true
check_section "User Stories" "user stor\|as a.*i want" true
check_section "Success Metrics" "success.*metric\|## success\|north star\|KPI" true
check_section "Scope" "## scope\|in scope\|in-scope\|out of scope" true

echo ""
echo -e "${BLUE}━━━ Recommended Sections ━━━${NC}"
echo ""
check_section "Executive Summary" "executive summary" false
check_section "User Personas" "user persona\|## persona" false
check_section "Technical Specifications" "technical.*spec\|technical.*consider\|## architecture" false
check_section "Design & UX" "design.*ux\|## design\|wireframe\|accessibility" false
check_section "Timeline & Milestones" "timeline.*milestone\|## timeline\|## milestones\|## roadmap" false
check_section "Risks & Mitigation" "risk.*mitigation\|## risk" false
check_section "Non-Goals / Out of Scope" "non-goal\|out of scope\|out-of-scope" false
check_section "Why Now" "why now\|why.*matter.*now\|urgency" false

# === Content Quality ===
echo ""
echo -e "${BLUE}━━━ Content Quality ━━━${NC}"
echo ""
check_content "Placeholder text" "\[.*\]" "Contains [bracket] placeholders — fill in all values"
check_content "TBD markers" "TBD\|TODO" "Contains TBD/TODO — resolve open items"

# === User Story Validation ===
echo ""
echo -e "${BLUE}━━━ User Story Validation ━━━${NC}"
echo ""

STORY_COUNT=$(grep -ci "as a.*i want\|as a.*want to" "$PRD_FILE" || true)
if [ "$STORY_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✓${NC} Found $STORY_COUNT user stories"
    ((PASSED++))

    AC_COUNT=$(grep -ci "acceptance criteria\|given.*when.*then" "$PRD_FILE" || true)
    if [ "$AC_COUNT" -ge "$STORY_COUNT" ]; then
        echo -e "${GREEN}✓${NC} Stories have acceptance criteria"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠${NC} Some stories may be missing acceptance criteria"
        ((WARNINGS++))
    fi
else
    echo -e "${RED}✗${NC} No user stories found (expected: 'As a [user], I want...')"
    ((ISSUES++))
fi

# === Metrics Validation ===
echo ""
echo -e "${BLUE}━━━ Metrics Validation ━━━${NC}"
echo ""

if grep -qi "KPI\|metric\|measure\|north star" "$PRD_FILE"; then
    echo -e "${GREEN}✓${NC} Success metrics mentioned"
    ((PASSED++))

    if grep -q "[0-9]\+%" "$PRD_FILE" || grep -qP "\d+\s*(ms|seconds|days|users)" "$PRD_FILE"; then
        echo -e "${GREEN}✓${NC} Contains quantifiable targets"
        ((PASSED++))
    else
        echo -e "${YELLOW}⚠${NC} Add quantifiable targets (%, ms, numbers)"
        ((WARNINGS++))
    fi
else
    echo -e "${RED}✗${NC} Success metrics not defined"
    ((ISSUES++))
fi

# === Scope Validation ===
echo ""
echo -e "${BLUE}━━━ Scope Validation ━━━${NC}"
echo ""

if grep -qi "in scope\|in-scope\|phase 1\|MVP" "$PRD_FILE"; then
    echo -e "${GREEN}✓${NC} In-scope defined"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠${NC} In-scope not clearly defined"
    ((WARNINGS++))
fi

if grep -qi "out of scope\|out-of-scope\|non-goal\|not building" "$PRD_FILE"; then
    echo -e "${GREEN}✓${NC} Out-of-scope defined"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠${NC} Out-of-scope not defined (prevents scope creep)"
    ((WARNINGS++))
fi

# === Vague Language Check ===
echo ""
echo -e "${BLUE}━━━ Vague Language Check ━━━${NC}"
echo ""

VAGUE_WORDS=("should be fast" "easy to use" "intuitive" "user-friendly" "modern" "seamless" "performant" "scalable" "robust")
VAGUE_FOUND=0
for word in "${VAGUE_WORDS[@]}"; do
    if grep -qi "$word" "$PRD_FILE"; then
        echo -e "${YELLOW}⚠${NC} Vague language: \"$word\" — replace with concrete metric"
        ((WARNINGS++))
        ((VAGUE_FOUND++))
    fi
done
if [ "$VAGUE_FOUND" -eq 0 ]; then
    echo -e "${GREEN}✓${NC} No vague language detected"
    ((PASSED++))
fi

# === Document Stats ===
echo ""
echo -e "${BLUE}━━━ Document Stats ━━━${NC}"
echo ""

WORD_COUNT=$(wc -w < "$PRD_FILE")
echo "Word count: $WORD_COUNT"

if [ "$WORD_COUNT" -lt 300 ]; then
    echo -e "${YELLOW}⚠${NC} Very short (< 300 words) — may need more detail"
    ((WARNINGS++))
elif [ "$WORD_COUNT" -gt 5000 ]; then
    echo -e "${YELLOW}⚠${NC} Very long (> 5000 words) — consider splitting or using Lean format"
    ((WARNINGS++))
else
    echo -e "${GREEN}✓${NC} Document length appropriate"
    ((PASSED++))
fi

# === Summary ===
echo ""
echo -e "${BLUE}╔════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║           Validation Summary           ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════╝${NC}"
echo ""
echo -e "Checks passed:  ${GREEN}$PASSED${NC}"
echo -e "Warnings:       ${YELLOW}$WARNINGS${NC}"
echo -e "Issues found:   ${RED}$ISSUES${NC}"
echo ""

if [ "$VERBOSE" = true ]; then
    echo -e "${BLUE}━━━ Recommendations ━━━${NC}"
    echo ""
    echo "1. Ensure all required sections are present"
    echo "2. Fill in all [placeholder] text"
    echo "3. Use 'As a [user], I want [action], So that [benefit]' format"
    echo "4. Include specific, measurable success metrics with targets"
    echo "5. Define what's in AND out of scope"
    echo "6. Add acceptance criteria for all user stories"
    echo "7. Replace vague language with concrete numbers"
    echo "8. Review with stakeholders before finalizing"
    echo ""
fi

if [ "$ISSUES" -gt 0 ]; then
    echo -e "${RED}❌ Validation failed — address $ISSUES critical issue(s)${NC}"
    exit 1
elif [ "$WARNINGS" -gt 0 ]; then
    echo -e "${YELLOW}⚠ Passed with $WARNINGS warning(s)${NC}"
    exit 0
else
    echo -e "${GREEN}✅ Validation passed — ready for stakeholder review${NC}"
    exit 0
fi
