package models

import (
	"strings"
	"testing"
)

func TestSimpleTaskSelectModel_Done_InitiallyFalse(t *testing.T) {
	m := SimpleTaskSelectModel{}
	if m.Done() {
		t.Error("Done() should be false for zero-value model")
	}
}

func TestSimpleTaskSelectModel_Done_TrueWhenSet(t *testing.T) {
	m := SimpleTaskSelectModel{done: true}
	if !m.Done() {
		t.Error("Done() should be true when done=true")
	}
}

func TestSimpleTaskSelectModel_ResultMessage_Success(t *testing.T) {
	m := SimpleTaskSelectModel{
		success: true,
		message: "Task claimed!",
	}
	result := m.ResultMessage()
	if !strings.Contains(result, "Task claimed!") {
		t.Errorf("ResultMessage() = %q, should contain %q", result, "Task claimed!")
	}
}

func TestSimpleTaskSelectModel_ResultMessage_Error(t *testing.T) {
	m := SimpleTaskSelectModel{
		success: false,
		error:   "something went wrong",
	}
	result := m.ResultMessage()
	if !strings.Contains(result, "something went wrong") {
		t.Errorf("ResultMessage() = %q, should contain %q", result, "something went wrong")
	}
}

func TestSimpleTaskSelectModel_ResultMessage_SuccessEmpty(t *testing.T) {
	m := SimpleTaskSelectModel{success: true, message: ""}
	result := m.ResultMessage()
	// Should not panic and should return something (styled empty string)
	if result == "" {
		// RenderSuccess("") may still produce styled output; just verify no panic
	}
	_ = result
}

func TestSimpleTaskSelectModel_ActionConstants(t *testing.T) {
	if ActionClaim != "claim" {
		t.Errorf("ActionClaim = %q, want %q", ActionClaim, "claim")
	}
	if ActionComplete != "complete" {
		t.Errorf("ActionComplete = %q, want %q", ActionComplete, "complete")
	}
	if ActionShow != "show" {
		t.Errorf("ActionShow = %q, want %q", ActionShow, "show")
	}
}
