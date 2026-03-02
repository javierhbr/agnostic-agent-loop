#!/bin/bash
set -e

echo "=== Tier Audit ==="
echo ""

FAIL=0
WARN=0
FAIL_COUNT=0
WARN_COUNT=0

# 1. Router sizes
echo "1. Router sizes..."
router1=$(wc -l < internal/skills/packs/SKILLS.md)
router2=$(wc -l < internal/skills/packs/sdd/SKILLS.md)

if [ "$router1" -gt 100 ] || [ "$router2" -gt 100 ]; then
  echo "   ❌ FAIL: Router exceeds 100 lines"
  FAIL=1
  ((FAIL_COUNT++))
elif [ "$router1" -gt 70 ] || [ "$router2" -gt 70 ]; then
  echo "   ⚠ WARN: packs/SKILLS.md: $router1 lines | packs/sdd/SKILLS.md: $router2 lines (warn threshold: 70)"
  ((WARN_COUNT++))
else
  echo "   ✅ packs/SKILLS.md: $router1 lines | packs/sdd/SKILLS.md: $router2 lines"
fi
echo ""

# 2. Skill file sizes
echo "2. Skill files..."
skill_fail=0
skill_warn=0
skill_count=0

while IFS= read -r skillfile; do
  ((skill_count++))
  lines=$(wc -l < "$skillfile")
  if [ "$lines" -gt 200 ]; then
    echo "   ❌ $(basename $(dirname $skillfile)): $lines lines (fail threshold: 200)"
    skill_fail=1
    ((FAIL_COUNT++))
  elif [ "$lines" -gt 130 ]; then
    echo "   ⚠ $(basename $(dirname $skillfile)): $lines lines (warn threshold: 130)"
    ((skill_warn++))
    ((WARN_COUNT++))
  fi
done < <(find internal/skills/packs -name "SKILL.md")

if [ "$skill_fail" -eq 0 ] && [ "$skill_warn" -eq 0 ]; then
  echo "   ✅ $skill_count files checked"
fi
echo ""

# 3. Resource file sizes
echo "3. Resource files..."
resource_warn=0
resource_count=0

while IFS= read -r resfile; do
  ((resource_count++))
  lines=$(wc -l < "$resfile")
  if [ "$lines" -gt 500 ]; then
    echo "   ⚠ $(echo $resfile | sed 's|internal/skills/packs/||'): $lines lines (warn threshold: 500)"
    ((resource_warn++))
    ((WARN_COUNT++))
  fi
done < <(find internal/skills/packs -path "*/resources/*.md")

if [ "$resource_warn" -eq 0 ]; then
  echo "   ✅ $resource_count files checked"
fi
echo ""

# 4. Resource link check
echo "4. Resource links..."
link_fail=0

while IFS= read -r skillfile; do
  dir=$(dirname "$skillfile")
  # Look for actual resource links: "→ `resources/..." (arrow followed by backtick)
  grep "→ \`resources/" "$skillfile" 2>/dev/null | \
    sed 's/.*→ `//' | \
    sed 's/`.*//' | \
    sed 's/#.*//' | \
    while IFS= read -r link; do
      [ -z "$link" ] && continue
      target="$dir/$link"
      if [ ! -f "$target" ]; then
        echo "   ❌ $(basename $(dirname $skillfile)): broken link → $link"
        link_fail=1
        FAIL=1
      fi
    done
done < <(find internal/skills/packs -name "SKILL.md")

if [ "$link_fail" -ne 1 ]; then
  echo "   ✅ All links resolve"
fi
echo ""

# 5. Summary
echo "=== Summary ==="
if [ "$FAIL_COUNT" -gt 0 ]; then
  echo "❌ FAIL: $FAIL_COUNT violations found"
  exit 1
elif [ "$WARN_COUNT" -gt 0 ]; then
  echo "⚠ WARN: $WARN_COUNT warnings | PASS"
  exit 0
else
  echo "✅ PASS"
  exit 0
fi
