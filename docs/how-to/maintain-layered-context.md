# Tier Maintenance Guide — Practical Workflows

This guide provides step-by-step procedures for maintaining the 3-tier layered context model in production.

---

## Quick Reference: Is Your Skill File the Right Size?

```
Under 60 lines   → Too small? Check if you're missing steps
60-80 lines      → Perfect (Tier 2 target)
80-130 lines     → Still OK, but approaching extraction threshold
130-200 lines    → ⚠️  RED FLAG — extract detail immediately
200+ lines       → 🚨 CRITICAL — this is a bloated Tier 2 file
```

---

## Monthly Maintenance Routine (15 minutes)

Run these commands at the start of each month to keep the system healthy.

### Step 1: Audit Router Sizes

```bash
# Check both routers
wc -l internal/skills/packs/SKILLS.md internal/skills/packs/sdd/SKILLS.md

# Expected output:
#    53 internal/skills/packs/SKILLS.md
#    59 internal/skills/packs/sdd/SKILLS.md
#   112 total
```

**If either router exceeds 100 lines:**
- Something was added inline instead of as a new skill
- Review the git diff to find what was added
- Extract the added content to a new resource file
- Update the router to just add the mapping line

### Step 2: Audit Skill File Sizes

```bash
# Find all SKILL.md files and their line counts
find internal/skills/packs -name "SKILL.md" -exec wc -l {} + | sort -n

# Example output:
#    54 internal/skills/packs/sdd/component-spec/SKILL.md
#    58 internal/skills/packs/api-docs/SKILL.md
#    59 internal/skills/packs/openspec/SKILL.md
#    ...
#   118 internal/skills/packs/sdd/hotfix/SKILL.md
# 1625 total
```

**Interpretation:**
- Most files are 54-120 lines ✅
- If any file > 150 lines, flag for extraction
- Average is ~85 lines (good)

**If you find files > 130 lines:**

```bash
# List them for review
find internal/skills/packs -name "SKILL.md" -exec wc -l {} + | awk '$1 > 130 {print}'
```

For each file > 130 lines, follow "Extracting Detail From a Skill File" section below.

### Step 3: Count Resource Files

```bash
# Should be around 15-20 resource files total
find internal/skills/packs -path "*/resources/*.md" | wc -l

# List them to see what exists
find internal/skills/packs -path "*/resources/*.md" | sort
```

### Step 4: Check for Broken Links

```bash
# Find all "→ resources/" references in SKILL.md files
grep -r "→ resources/" internal/skills/packs --include="SKILL.md"

# For each link found, verify the file exists
# Example: If you see "→ resources/phases.md"
ls internal/skills/packs/openspec/resources/phases.md
# Should exist and show file size
```

**If a link is broken:**
1. Find the skill file that references it
2. Determine what resource it should link to (or doesn't exist)
3. Either create the resource file or fix the link

### Step 5: Quick Health Report

Create a summary and save to `docs/TIER-MAINTENANCE.md`:

```markdown
# Tier Maintenance Report

**Date:** [YYYY-MM-DD]
**Status:** ✅ HEALTHY | ⚠️  NEEDS ATTENTION | 🚨 CRITICAL

## Routers
- SKILLS.md: [N] lines ✅
- sdd/SKILLS.md: [N] lines ✅

## Skills Summary
- Total skill files: [N]
- Average size: [N] lines
- Largest: [filename] ([N] lines)
- Files over 130 lines: 0 ✅

## Resources
- Total resource files: [N]
- Largest: [filename] ([N] lines)
- Broken links: 0 ✅

## Action Items (if any)
- [ ] [Item 1]
- [ ] [Item 2]

Next audit: [date in 1 month]
```

---

## Extracting Detail From a Skill File

**When to do this:** Your skill file exceeds 130 lines

**Time estimate:** 20-30 minutes

### Step 1: Identify What to Extract

Read the entire skill file and ask:

1. Are there 3+ code examples? → Extract to `resources/examples.md`
2. Are there 2+ templates? → Extract to `resources/templates.md`
3. Are there multi-sentence explanations? → Condense + extract to `resources/guide.md`
4. Are there edge cases or troubleshooting? → Extract to `resources/troubleshooting.md`
5. Are there many rules or conventions? → Extract to `resources/rules.md`

### Step 2: Create the Resource File

Create `[skill]/resources/[topic].md`:

```bash
# Navigate to the skill directory
cd internal/skills/packs/[skill-name]/

# Create resources directory if it doesn't exist
mkdir -p resources

# Create the resource file
touch resources/[topic].md
```

### Step 3: Copy Content to Resource File

Take the content from the skill file and move it to the resource file with clear section anchors.

**Example: Extracting examples from openspec/SKILL.md**

**Before (in SKILL.md):**
```markdown
## Steps

1. Read requirements

Example: A user says "I want a login form"
You respond: "What validation rules? What's the mobile UX?"

2. Create proposal

Example proposal structure:
- Goal: Users can log in with email/password
- Actor: End user
- Success metric: 95% users complete login in < 30 seconds
```

**After (in resources/examples.md):**
```markdown
# OpenSpec Examples

## Phase 1: Questions

### Example 1: Login Form

**Scenario:** User says "I want a login form"

**Agent response:**
- What validation rules apply? (min length, special chars, etc.)
- What's the mobile UX? (full-width input, or compact?)
- Are there accessibility requirements? (WCAG AA minimum?)
- Do users expect "Remember me"?

...
```

**In SKILL.md, reference it:**
```markdown
1. Read requirements
   → See resources/examples.md#phase-1-questions for examples
```

### Step 4: Update Links in Skill File

Replace inline content with links:

```markdown
## If you need more detail
→ resources/examples.md
→ resources/templates.md
```

### Step 5: Verify the Skill File is Now Small

```bash
wc -l [skill-name]/SKILL.md

# Should be < 100 lines
```

### Step 6: Test with an Agent

Ask an agent: "What does [skill-name] do?"

They should:
1. Load just the SKILL.md file
2. See the steps and output
3. Know where to find examples ("→ resources/examples.md")

---

## Splitting a Large Resource File

**When to do this:** A resource file exceeds 500 lines and covers multiple topics

**Time estimate:** 30-45 minutes

### Step 1: Identify the Topics

Read the resource file and list the main sections.

**Example: hotfix-templates.md has these sections:**
- Bug Spec Template
- Hotfix Spec Template
- Follow-up Spec Template
- Incident Response Examples
- Hotfix Governance

**Group related sections:**
- Templates (Bug, Hotfix, Follow-up)
- Examples (Incident Response)
- Governance (Rules and reminders)

### Step 2: Create Split Files

```bash
cd [skill-name]/resources/

# Create:
touch templates.md       # Bug, Hotfix, Follow-up specs
touch incidents.md      # Real incident examples
touch governance.md     # Rules, reminders, principles
```

### Step 3: Move Content to New Files

Copy the relevant sections:

**templates.md:**
```markdown
# Hotfix Templates

## Bug Spec Template

[Full Bug Spec template content]

## Hotfix Spec Template

[Full Hotfix Spec template content]

## Follow-up Spec Template

[Full Follow-up Spec template content]
```

**incidents.md:**
```markdown
# Incident Response Examples

## Example 1: Payment Processor Down

[Full incident walkthrough]

## Example 2: Data Corruption

[Full incident walkthrough]
```

### Step 4: Update the Skill File

Instead of linking to one huge file:

**Before:**
```markdown
## If you need more detail
→ resources/hotfix-templates.md
```

**After:**
```markdown
## If you need more detail
→ resources/templates.md (Bug/Hotfix/Follow-up specs)
→ resources/incidents.md (Real incident examples)
→ resources/governance.md (Hotfix rules and process)
```

### Step 5: Delete the Original File

```bash
rm resources/hotfix-templates.md
```

### Step 6: Verify Links Still Work

```bash
# Find all references to the old file
grep -r "hotfix-templates.md" ../../

# Should find zero results
```

---

## Adding a New Skill

**Time estimate:** 45-60 minutes

### Step 1: Create the Skill Structure

```bash
# Create skill folder
mkdir -p internal/skills/packs/[skill-name]/

# Create SKILL.md file
touch internal/skills/packs/[skill-name]/SKILL.md

# Create resources folder
mkdir -p internal/skills/packs/[skill-name]/resources/
```

### Step 2: Write the Skill File (Use Template)

Edit `internal/skills/packs/[skill-name]/SKILL.md`:

```markdown
---
name: [skill-name]
description: "[One sentence describing what this skill does]"
---

# skill:[skill-name]

## Does exactly this

[One sentence — describe the output.]

## Use this skill when

- [Trigger scenario 1]
- [Trigger scenario 2]
- [Trigger scenario 3]

## Steps — in order, no skipping

1. [Step 1 in 1-2 lines]
   → See resources/[detail].md if needed

2. [Step 2 in 1-2 lines]

3. [Step 3 in 1-2 lines]

## Output

[Describe the exact artifact produced]

## Done when

- [Completion condition 1]
- [Completion condition 2]

## If you need more detail

→ `resources/[topic].md`
```

**Verify size:**
```bash
wc -l internal/skills/packs/[skill-name]/SKILL.md
# Should be 50-80 lines
```

### Step 3: Create Resource Files (if needed)

If your skill file approaches 100 lines, extract detail:

```bash
# Create resources
touch internal/skills/packs/[skill-name]/resources/guide.md
```

Add sections with clear anchors for linking.

### Step 4: Update the Router

Edit the appropriate router file:

**For flat-pack skills, edit `internal/skills/packs/SKILLS.md`:**

```markdown
| `[skill-name]` | `./[skill-name]/SKILL.md` | [Description] |
```

**For SDD skills, edit `internal/skills/packs/sdd/SKILLS.md`:**

```markdown
| `sdd:[skill-name]` | [Role] | `./[skill-name]/SKILL.md` | [Purpose] |
```

**Verify router size:**
```bash
wc -l internal/skills/packs/SKILLS.md internal/skills/packs/sdd/SKILLS.md
# Both should still be < 80 lines
```

### Step 5: Verify the New Skill Works

Test with an agent:

```
User: "[trigger] [task]"

Expected response:
1. Agent loads: internal/skills/packs/SKILLS.md
2. Agent loads: internal/skills/packs/[skill-name]/SKILL.md
3. Agent never loads resources/ unless you ask "show me examples"
```

### Step 6: Document in MEMORY.md

Add entry to your project memory:

```markdown
## New Skill: [skill-name] (Date)

**Purpose:** [One sentence]
**Router:** SKILLS.md
**Location:** `internal/skills/packs/[skill-name]/SKILL.md`
**Resources:** `resources/[topic].md` (if any)
**Size:** [N] lines
**Trigger:** [skill-name]
```

---

## Responding to "This Skill Isn't Working"

**When someone reports a problem with a skill:**

### Step 1: Determine the Issue Type

Ask:
- "What did the skill output?"
- "What were you expecting?"
- "Did it load the right resources?"

**Common issues:**

| Issue | Cause | Fix |
|-------|-------|-----|
| "Skill file is huge and confusing" | Tier 2 bloat | Extract detail to Tier 3 |
| "I can't find what I need" | Resource file is unclear | Add section anchors + better titles |
| "The skill loaded too much" | Tier 2 isn't slim | Audit the skill file for inline content |
| "Link to resources is broken" | Resource file doesn't exist | Create missing file or fix link |
| "Too many steps to follow" | Tier 2 is too complex | Simplify steps, move details to Tier 3 |

### Step 2: Fix the Issue

Apply the appropriate remedy from the table above.

### Step 3: Verify the Fix

Test with an agent again to confirm resolution.

### Step 4: Document Lessons Learned

Add to `docs/TIER-MAINTENANCE.md`:

```markdown
## Issue Resolved: [Issue Description]

**Date:** [date]
**Root cause:** [What caused it]
**Fix:** [What was changed]
**Prevention:** [How to prevent in future]
```

---

## Git Workflow for Tier Maintenance

### Committing Skill Changes

Use semantic commits to track maintenance:

```bash
# Extracting detail from an existing skill
git commit -m "refactor: Extract detail from [skill-name] SKILL.md to resources/"

# Adding a new skill
git commit -m "feat: Add skill:[skill-name]"

# Splitting a large resource file
git commit -m "refactor: Split [skill-name] resources into multiple files"

# Monthly maintenance audit
git commit -m "docs: Update tier maintenance report for [month]"
```

### Reviewing Maintenance PRs

When reviewing a PR that touches skills:

- [ ] Routers didn't grow (new skills = +1 line only)
- [ ] New skill files are < 130 lines
- [ ] Resource files have clear sections and anchors
- [ ] All links exist and work
- [ ] Skill file follows the template exactly
- [ ] No duplicate content between Tier 2 and Tier 3

---

## Quarterly Deep Dive (60 minutes)

Run this comprehensive review each quarter.

### Step 1: Size Audit (15 min)

```bash
# Generate a comprehensive report
echo "# Tier Maintenance Report — $(date +%Y-%m-%d)" > /tmp/audit.md
echo "" >> /tmp/audit.md

echo "## Routers" >> /tmp/audit.md
wc -l internal/skills/packs/SKILLS.md internal/skills/packs/sdd/SKILLS.md >> /tmp/audit.md

echo "" >> /tmp/audit.md
echo "## Skill Files" >> /tmp/audit.md
find internal/skills/packs -name "SKILL.md" -exec wc -l {} + | sort -n >> /tmp/audit.md

echo "" >> /tmp/audit.md
echo "## Resource Files" >> /tmp/audit.md
find internal/skills/packs -path "*/resources/*.md" -exec wc -l {} + | sort -n >> /tmp/audit.md

cat /tmp/audit.md
```

### Step 2: Content Review (20 min)

For each skill file > 100 lines:
- Read the entire file
- Identify any sections that could be extracted
- Create an extraction plan if needed

### Step 3: Link Verification (15 min)

```bash
# Find all resource links
grep -r "→ resources/" internal/skills/packs --include="SKILL.md"

# For each link, verify the file exists
# (Script below)
```

**Link verification script:**

```bash
#!/bin/bash
echo "Checking for broken links..."
grep -r "→ resources/" internal/skills/packs --include="SKILL.md" | while read -r line; do
    file=$(echo "$line" | awk -F: '{print $1}')
    link=$(echo "$line" | sed 's/.*resources\//resources\//' | awk '{print $1}')

    dir=$(dirname "$file")
    resource_path="$dir/$link"

    if [ ! -f "$resource_path" ]; then
        echo "❌ BROKEN: $file → $resource_path"
    fi
done
echo "Done."
```

### Step 4: Duplication Check (10 min)

Look for the same content in multiple resource files:

```bash
# Find common patterns (example: "## Steps" appears multiple times)
grep -r "^## Steps" internal/skills/packs/*/resources/ --files-with-matches

# Are these duplicates or variations?
# If duplicates → consolidate into one resource file
```

### Step 5: Staleness Check (5 min)

For each resource file > 1 year old:

```bash
# Find files with no recent commits
git log --format="%ai" -- internal/skills/packs/*/resources/*.md | head -20
```

If a resource file hasn't been touched in 12+ months:
- Review its content for accuracy
- Add timestamp: `<!-- Last reviewed: 2026-Q1 -->`
- Flag for update if content seems outdated

### Step 6: Report and Plan

Update `docs/TIER-MAINTENANCE.md` with findings:

```markdown
# Quarterly Maintenance Report — Q[N] 2026

**Date:** [YYYY-MM-DD]
**Auditor:** [Your name]
**Status:** ✅ HEALTHY

## Summary
- 2 routers, 19 skills, 16 resources
- Average skill size: 85 lines ✅
- 0 broken links ✅
- 0 files over 150 lines ✅

## Issues Found
- [ ] Issue 1: [Description] → Fix by [date]
- [ ] Issue 2: [Description] → Fix by [date]

## Next Steps
1. [Action 1]
2. [Action 2]

Next audit: [date in 3 months]
```

---

## Troubleshooting Problems

### "The model is breaking down — skills are getting bloated"

**Diagnosis:**
1. Run the monthly maintenance routine (see above)
2. Identify which skill files exceed 130 lines
3. Check if any routers exceed 100 lines

**Recovery:**
1. Extract detail from bloated skill files immediately
2. Split resource files if they exceed 500 lines
3. Review git history to find when drift started
4. Implement weekly size checks in CI (optional)

### "Agents keep loading the wrong resources"

**Root causes:**
- Resource links are ambiguous (should point to `#anchor`)
- Resource file sections aren't clearly labeled
- Multiple resource files have conflicting content

**Fix:**
1. Audit all resource files for section headers
2. Add anchors (e.g., `## Steps → #steps` anchor)
3. Update skill files to link with anchors: `→ resources/file.md#section`
4. Consolidate duplicate content

### "New team members can't find what they need"

**Root cause:** The system isn't documented well enough

**Fix:**
1. Create a "Skills Orientation" doc:
   - How to read a SKILL.md file
   - How to find resources
   - When to use which skill
2. Add to onboarding checklist
3. Create a visual map: `docs/SKILLS-MAP.md`

---

## Template: Maintenance Checklist

Copy this to your project and run monthly:

```markdown
# Monthly Tier Maintenance Checklist

**Month:** [YYYY-MM]
**Auditor:** [Name]
**Date:** [YYYY-MM-DD]

## Router Size Check
- [ ] SKILLS.md < 70 lines
- [ ] sdd/SKILLS.md < 70 lines

## Skill File Audit
- [ ] No SKILL.md > 150 lines
- [ ] Average size reasonable
- [ ] Largest file identified: _________

## Resource File Check
- [ ] All linked resources exist
- [ ] No broken links found
- [ ] Largest resource identified: _________

## Content Quality
- [ ] No duplicate content across resources
- [ ] All resource files have section anchors
- [ ] All skills use consistent template

## Health Status
- [ ] ✅ HEALTHY — All checks passed
- [ ] ⚠️  NEEDS ATTENTION — See issues below
- [ ] 🚨 CRITICAL — Requires immediate action

## Issues Found
- [ ] [Issue 1]
- [ ] [Issue 2]

## Fixes Applied
- [X] [Fix 1]
- [X] [Fix 2]

## Next Steps
1. [Action 1] — Due [date]
2. [Action 2] — Due [date]

## Sign-off
Reviewed by: [Name]
Date: [YYYY-MM-DD]
Next audit: [YYYY-MM-DD]
```

---

## Summary

The 3-tier model stays healthy through **consistent, lightweight maintenance**. You don't need big refactors — just monthly reviews and immediate extraction when files grow too large.

**Key metrics to watch:**
- Routers: < 70 lines
- Skills: < 130 lines (target 60-80)
- Resources: < 500 lines (split if larger)
- Broken links: 0

When in doubt, **extract detail and link from Tier 2 to Tier 3**. It's always the right move.
