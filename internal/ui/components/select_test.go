package components

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func makeOptions(n int) []SelectOption {
	opts := make([]SelectOption, n)
	for i := range opts {
		opts[i] = NewSelectOption(
			"Option "+string(rune('A'+i)),
			"Description "+string(rune('A'+i)),
			"val-"+string(rune('a'+i)),
		)
	}
	return opts
}

// --- SelectOption ---

func TestSelectOption_Fields(t *testing.T) {
	opt := NewSelectOption("My Title", "My Desc", "my-val")

	if opt.Title() != "My Title" {
		t.Errorf("Title() = %q, want %q", opt.Title(), "My Title")
	}
	if opt.Description() != "My Desc" {
		t.Errorf("Description() = %q, want %q", opt.Description(), "My Desc")
	}
	if opt.Value() != "my-val" {
		t.Errorf("Value() = %q, want %q", opt.Value(), "my-val")
	}
	if opt.FilterValue() != "My Title" {
		t.Errorf("FilterValue() = %q, want %q", opt.FilterValue(), "My Title")
	}
}

// --- SimpleSelect ---

func TestSimpleSelect_Navigation(t *testing.T) {
	opts := makeOptions(3)
	ss := NewSimpleSelect("Pick one", opts)

	if ss.SelectedIdx != 0 {
		t.Fatalf("initial SelectedIdx = %d, want 0", ss.SelectedIdx)
	}

	// Move down
	ss = ss.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if ss.SelectedIdx != 1 {
		t.Errorf("after j: SelectedIdx = %d, want 1", ss.SelectedIdx)
	}

	// Move down again
	ss = ss.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if ss.SelectedIdx != 2 {
		t.Errorf("after j j: SelectedIdx = %d, want 2", ss.SelectedIdx)
	}

	// Move down at boundary — should stay at 2
	ss = ss.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if ss.SelectedIdx != 2 {
		t.Errorf("at bottom after j: SelectedIdx = %d, want 2", ss.SelectedIdx)
	}

	// Move up
	ss = ss.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}})
	if ss.SelectedIdx != 1 {
		t.Errorf("after k: SelectedIdx = %d, want 1", ss.SelectedIdx)
	}
}

func TestSimpleSelect_SelectedValue(t *testing.T) {
	opts := makeOptions(3)
	ss := NewSimpleSelect("Pick", opts)

	if ss.SelectedValue() != "val-a" {
		t.Errorf("initial SelectedValue = %q, want %q", ss.SelectedValue(), "val-a")
	}

	ss = ss.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}})
	if ss.SelectedValue() != "val-b" {
		t.Errorf("after j: SelectedValue = %q, want %q", ss.SelectedValue(), "val-b")
	}
}

func TestSimpleSelect_SelectedOption(t *testing.T) {
	opts := makeOptions(2)
	ss := NewSimpleSelect("Pick", opts)

	opt := ss.SelectedOption()
	if opt == nil {
		t.Fatal("SelectedOption() returned nil")
	}
	if opt.Value() != "val-a" {
		t.Errorf("SelectedOption().Value() = %q, want %q", opt.Value(), "val-a")
	}
}

func TestSimpleSelect_SelectedValue_Empty(t *testing.T) {
	ss := NewSimpleSelect("Pick", nil)
	if ss.SelectedValue() != "" {
		t.Errorf("SelectedValue() on empty = %q, want empty", ss.SelectedValue())
	}
	if ss.SelectedOption() != nil {
		t.Error("SelectedOption() on empty should be nil")
	}
}

func TestSimpleSelect_View(t *testing.T) {
	opts := makeOptions(2)
	ss := NewSimpleSelect("Test Label", opts)
	v := ss.View()

	if !strings.Contains(v, "Test Label") {
		t.Error("View() should contain label")
	}
	if !strings.Contains(v, "Option A") {
		t.Error("View() should contain first option title")
	}
	if !strings.Contains(v, "Option B") {
		t.Error("View() should contain second option title")
	}
}

// --- MultiSelect ---

func TestMultiSelect_Toggle(t *testing.T) {
	opts := makeOptions(3)
	ms := NewMultiSelect("Multi", opts)

	if ms.HasSelection() {
		t.Error("initial HasSelection() should be false")
	}

	// Toggle first item
	ms = ms.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	if !ms.Selected[0] {
		t.Error("after space: item 0 should be selected")
	}
	if !ms.HasSelection() {
		t.Error("after toggle: HasSelection() should be true")
	}

	vals := ms.SelectedValues()
	if len(vals) != 1 || vals[0] != "val-a" {
		t.Errorf("SelectedValues() = %v, want [val-a]", vals)
	}

	// Toggle off
	ms = ms.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}})
	if ms.Selected[0] {
		t.Error("after second space: item 0 should be deselected")
	}
}

func TestMultiSelect_Navigation(t *testing.T) {
	opts := makeOptions(4)
	ms := NewMultiSelect("Multi", opts)

	if ms.CursorIdx != 0 {
		t.Fatalf("initial CursorIdx = %d, want 0", ms.CursorIdx)
	}

	// down arrow
	ms = ms.Update(tea.KeyMsg{Type: tea.KeyDown})
	if ms.CursorIdx != 1 {
		t.Errorf("after down: CursorIdx = %d, want 1", ms.CursorIdx)
	}

	// up arrow
	ms = ms.Update(tea.KeyMsg{Type: tea.KeyUp})
	if ms.CursorIdx != 0 {
		t.Errorf("after up: CursorIdx = %d, want 0", ms.CursorIdx)
	}

	// up at top — stays at 0
	ms = ms.Update(tea.KeyMsg{Type: tea.KeyUp})
	if ms.CursorIdx != 0 {
		t.Errorf("at top after up: CursorIdx = %d, want 0", ms.CursorIdx)
	}
}

func TestMultiSelect_SelectedValues_Multiple(t *testing.T) {
	opts := makeOptions(4)
	ms := NewMultiSelect("Multi", opts)

	// Select items 0 and 2
	ms = ms.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}) // select 0
	ms = ms.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) // move to 1
	ms = ms.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}) // move to 2
	ms = ms.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{' '}}) // select 2

	vals := ms.SelectedValues()
	if len(vals) != 2 {
		t.Fatalf("SelectedValues() has %d items, want 2", len(vals))
	}
	if vals[0] != "val-a" || vals[1] != "val-c" {
		t.Errorf("SelectedValues() = %v, want [val-a val-c]", vals)
	}
}

// --- MultiSelect Viewport Scrolling ---

func TestMultiSelect_SetMaxVisible(t *testing.T) {
	opts := makeOptions(10)
	ms := NewMultiSelect("Multi", opts)

	// Default: no limit
	if ms.MaxVisible != 0 {
		t.Errorf("initial MaxVisible = %d, want 0", ms.MaxVisible)
	}

	// Set with 20 available lines: (20-2)/2 = 9
	ms.SetMaxVisible(20)
	if ms.MaxVisible != 9 {
		t.Errorf("SetMaxVisible(20): MaxVisible = %d, want 9", ms.MaxVisible)
	}

	// Set with small height: (4-2)/2 = 1
	ms.SetMaxVisible(4)
	if ms.MaxVisible != 1 {
		t.Errorf("SetMaxVisible(4): MaxVisible = %d, want 1", ms.MaxVisible)
	}

	// Minimum is 1
	ms.SetMaxVisible(1)
	if ms.MaxVisible < 1 {
		t.Errorf("SetMaxVisible(1): MaxVisible = %d, want >= 1", ms.MaxVisible)
	}
}

func TestMultiSelect_EnsureCursorVisible_ScrollDown(t *testing.T) {
	opts := makeOptions(10)
	ms := NewMultiSelect("Multi", opts)
	ms.MaxVisible = 3

	// Cursor at 0, offset at 0 — visible
	ms.ensureCursorVisible()
	if ms.offset != 0 {
		t.Errorf("offset = %d, want 0", ms.offset)
	}

	// Move cursor beyond visible window
	ms.CursorIdx = 4
	ms.ensureCursorVisible()
	if ms.offset != 2 { // 4 - 3 + 1 = 2
		t.Errorf("after cursor=4: offset = %d, want 2", ms.offset)
	}
}

func TestMultiSelect_EnsureCursorVisible_ScrollUp(t *testing.T) {
	opts := makeOptions(10)
	ms := NewMultiSelect("Multi", opts)
	ms.MaxVisible = 3
	ms.offset = 5
	ms.CursorIdx = 3

	ms.ensureCursorVisible()
	if ms.offset != 3 {
		t.Errorf("after scroll up: offset = %d, want 3", ms.offset)
	}
}

func TestMultiSelect_EnsureCursorVisible_NoopWhenAllFit(t *testing.T) {
	opts := makeOptions(3)
	ms := NewMultiSelect("Multi", opts)
	ms.MaxVisible = 5 // more than len(Options)

	ms.CursorIdx = 2
	ms.ensureCursorVisible()
	if ms.offset != 0 {
		t.Errorf("all fit: offset = %d, want 0", ms.offset)
	}
}

func TestMultiSelect_EnsureCursorVisible_NoopWhenUnlimited(t *testing.T) {
	opts := makeOptions(10)
	ms := NewMultiSelect("Multi", opts)
	// MaxVisible = 0 (default, unlimited)

	ms.CursorIdx = 8
	ms.ensureCursorVisible()
	if ms.offset != 0 {
		t.Errorf("unlimited: offset = %d, want 0", ms.offset)
	}
}

func TestMultiSelect_View_ShowsAllWhenNoLimit(t *testing.T) {
	opts := makeOptions(5)
	ms := NewMultiSelect("Test", opts)

	v := ms.View()
	for _, opt := range opts {
		if !strings.Contains(v, opt.Title()) {
			t.Errorf("View() missing option %q", opt.Title())
		}
	}

	if strings.Contains(v, "↑ more") || strings.Contains(v, "↓ more") {
		t.Error("View() should not show scroll indicators when all items fit")
	}
}

func TestMultiSelect_View_ScrollIndicators(t *testing.T) {
	opts := makeOptions(8)
	ms := NewMultiSelect("Test", opts)
	ms.MaxVisible = 3

	// At top: no ↑, but ↓
	v := ms.View()
	if strings.Contains(v, "↑ more") {
		t.Error("at top: should not show ↑ more")
	}
	if !strings.Contains(v, "↓ more") {
		t.Error("at top: should show ↓ more")
	}

	// Move to middle
	ms.offset = 2
	ms.CursorIdx = 3
	v = ms.View()
	if !strings.Contains(v, "↑ more") {
		t.Error("in middle: should show ↑ more")
	}
	if !strings.Contains(v, "↓ more") {
		t.Error("in middle: should show ↓ more")
	}

	// Move to bottom
	ms.offset = 5 // shows items 5,6,7
	ms.CursorIdx = 7
	v = ms.View()
	if !strings.Contains(v, "↑ more") {
		t.Error("at bottom: should show ↑ more")
	}
	if strings.Contains(v, "↓ more") {
		t.Error("at bottom: should not show ↓ more")
	}
}

func TestMultiSelect_View_WindowContainsCorrectItems(t *testing.T) {
	opts := makeOptions(8)
	ms := NewMultiSelect("Test", opts)
	ms.MaxVisible = 3
	ms.offset = 2
	ms.CursorIdx = 3

	v := ms.View()

	// Items 2, 3, 4 should be visible (offset=2, window=3)
	if !strings.Contains(v, "Option C") {
		t.Error("View() should contain Option C (index 2)")
	}
	if !strings.Contains(v, "Option D") {
		t.Error("View() should contain Option D (index 3)")
	}
	if !strings.Contains(v, "Option E") {
		t.Error("View() should contain Option E (index 4)")
	}

	// Items outside window should not be visible
	if strings.Contains(v, "Option A") {
		t.Error("View() should not contain Option A (index 0, before window)")
	}
	if strings.Contains(v, "Option B") {
		t.Error("View() should not contain Option B (index 1, before window)")
	}
	if strings.Contains(v, "Option F") {
		t.Error("View() should not contain Option F (index 5, after window)")
	}
}

func TestMultiSelect_NavigationWithViewport(t *testing.T) {
	opts := makeOptions(8)
	ms := NewMultiSelect("Test", opts)
	ms.MaxVisible = 3

	// Navigate down past visible window
	for range 5 {
		ms = ms.Update(tea.KeyMsg{Type: tea.KeyDown})
	}

	if ms.CursorIdx != 5 {
		t.Errorf("CursorIdx = %d, want 5", ms.CursorIdx)
	}
	// offset should have scrolled so cursor is visible
	if ms.CursorIdx < ms.offset || ms.CursorIdx >= ms.offset+ms.MaxVisible {
		t.Errorf("cursor %d not in visible window [%d, %d)", ms.CursorIdx, ms.offset, ms.offset+ms.MaxVisible)
	}

	// Navigate back up
	for range 5 {
		ms = ms.Update(tea.KeyMsg{Type: tea.KeyUp})
	}

	if ms.CursorIdx != 0 {
		t.Errorf("CursorIdx = %d, want 0", ms.CursorIdx)
	}
	if ms.offset != 0 {
		t.Errorf("offset = %d, want 0", ms.offset)
	}
}
