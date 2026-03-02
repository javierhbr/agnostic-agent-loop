package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestTaskStatus_Constants(t *testing.T) {
	assert.Equal(t, TaskStatus("pending"), StatusPending)
	assert.Equal(t, TaskStatus("in-progress"), StatusInProgress)
	assert.Equal(t, TaskStatus("done"), StatusDone)
}

func TestTask_YAMLMarshal(t *testing.T) {
	task := Task{
		ID:          "TASK-001",
		Title:       "Test Task",
		Description: "Test description",
		Status:      StatusPending,
		AssignedTo:  "user1",
		Scope:       []string{"src/module"},
		SpecRefs:    []string{"spec1.md", "spec2.md"},
		Inputs:      []string{"input.txt"},
		Outputs:     []string{"output.txt"},
		Acceptance:  []string{"Works", "Has tests"},
		SubTasks: []SubTask{
			{ID: "TASK-001.1", Title: "Subtask", Status: StatusPending},
		},
	}

	data, err := yaml.Marshal(&task)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Verify YAML contains expected fields
	yamlStr := string(data)
	assert.Contains(t, yamlStr, "id: TASK-001")
	assert.Contains(t, yamlStr, "title: Test Task")
	assert.Contains(t, yamlStr, "status: pending")
	assert.Contains(t, yamlStr, "assigned_to: user1")
	assert.Contains(t, yamlStr, "spec1.md")
	assert.Contains(t, yamlStr, "spec2.md")
	assert.Contains(t, yamlStr, "input.txt")
	assert.Contains(t, yamlStr, "output.txt")
	assert.Contains(t, yamlStr, "Works")
	assert.Contains(t, yamlStr, "TASK-001.1")
}

func TestTask_YAMLUnmarshal(t *testing.T) {
	yamlData := `
id: TASK-002
title: Unmarshal Test
description: Testing unmarshal
status: in-progress
assigned_to: user2
scope:
  - src/auth
  - src/core
spec_refs:
  - spec/architecture.md
inputs:
  - context.md
outputs:
  - module.go
  - module_test.go
acceptance:
  - Feature works
  - Tests pass
subtasks:
  - id: TASK-002.1
    title: Sub 1
    status: done
    assigned_to: user3
`

	var task Task
	err := yaml.Unmarshal([]byte(yamlData), &task)
	assert.NoError(t, err)

	assert.Equal(t, "TASK-002", task.ID)
	assert.Equal(t, "Unmarshal Test", task.Title)
	assert.Equal(t, "Testing unmarshal", task.Description)
	assert.Equal(t, StatusInProgress, task.Status)
	assert.Equal(t, "user2", task.AssignedTo)
	assert.Equal(t, []string{"src/auth", "src/core"}, task.Scope)
	assert.Equal(t, []string{"spec/architecture.md"}, task.SpecRefs)
	assert.Equal(t, []string{"context.md"}, task.Inputs)
	assert.Equal(t, []string{"module.go", "module_test.go"}, task.Outputs)
	assert.Equal(t, []string{"Feature works", "Tests pass"}, task.Acceptance)
	assert.Len(t, task.SubTasks, 1)
	assert.Equal(t, "TASK-002.1", task.SubTasks[0].ID)
	assert.Equal(t, "Sub 1", task.SubTasks[0].Title)
	assert.Equal(t, StatusDone, task.SubTasks[0].Status)
	assert.Equal(t, "user3", task.SubTasks[0].AssignedTo)
}

func TestTask_YAMLMarshalEmpty(t *testing.T) {
	task := Task{
		ID:     "TASK-MIN",
		Title:  "Minimal",
		Status: StatusPending,
	}

	data, err := yaml.Marshal(&task)
	assert.NoError(t, err)
	assert.NotEmpty(t, data)

	// Verify optional fields are omitted
	yamlStr := string(data)
	assert.NotContains(t, yamlStr, "assigned_to")
	assert.NotContains(t, yamlStr, "scope")
	assert.NotContains(t, yamlStr, "spec_refs")
}

func TestTask_YAMLRoundTrip(t *testing.T) {
	original := Task{
		ID:          "TASK-RT",
		Title:       "Round Trip",
		Description: "Test round trip",
		Status:      StatusInProgress,
		AssignedTo:  "tester",
		Scope:       []string{"src/test"},
		SpecRefs:    []string{"spec.md"},
		Inputs:      []string{"in.txt"},
		Outputs:     []string{"out.txt"},
		Acceptance:  []string{"It works"},
	}

	// Marshal
	data, err := yaml.Marshal(&original)
	assert.NoError(t, err)

	// Unmarshal
	var decoded Task
	err = yaml.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Compare
	assert.Equal(t, original.ID, decoded.ID)
	assert.Equal(t, original.Title, decoded.Title)
	assert.Equal(t, original.Description, decoded.Description)
	assert.Equal(t, original.Status, decoded.Status)
	assert.Equal(t, original.AssignedTo, decoded.AssignedTo)
	assert.Equal(t, original.Scope, decoded.Scope)
	assert.Equal(t, original.SpecRefs, decoded.SpecRefs)
	assert.Equal(t, original.Inputs, decoded.Inputs)
	assert.Equal(t, original.Outputs, decoded.Outputs)
	assert.Equal(t, original.Acceptance, decoded.Acceptance)
}

func TestSubTask_YAMLMarshal(t *testing.T) {
	subtask := SubTask{
		ID:         "TASK-001.1",
		Title:      "Subtask Title",
		Status:     StatusDone,
		AssignedTo: "worker",
	}

	data, err := yaml.Marshal(&subtask)
	assert.NoError(t, err)

	yamlStr := string(data)
	assert.Contains(t, yamlStr, "TASK-001.1")
	assert.Contains(t, yamlStr, "Subtask Title")
	assert.Contains(t, yamlStr, "done")
	assert.Contains(t, yamlStr, "worker")
}

func TestSubTask_YAMLUnmarshal(t *testing.T) {
	yamlData := `
id: TASK-003.2
title: Sub Unmarshal
status: pending
`

	var subtask SubTask
	err := yaml.Unmarshal([]byte(yamlData), &subtask)
	assert.NoError(t, err)

	assert.Equal(t, "TASK-003.2", subtask.ID)
	assert.Equal(t, "Sub Unmarshal", subtask.Title)
	assert.Equal(t, StatusPending, subtask.Status)
	assert.Empty(t, subtask.AssignedTo)
}

func TestTask_WithMultipleSubtasks(t *testing.T) {
	task := Task{
		ID:     "TASK-MULTI",
		Title:  "Multi Subtask",
		Status: StatusInProgress,
		SubTasks: []SubTask{
			{ID: "TASK-MULTI.1", Title: "Sub 1", Status: StatusDone},
			{ID: "TASK-MULTI.2", Title: "Sub 2", Status: StatusInProgress},
			{ID: "TASK-MULTI.3", Title: "Sub 3", Status: StatusPending},
		},
	}

	data, err := yaml.Marshal(&task)
	assert.NoError(t, err)

	var decoded Task
	err = yaml.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	assert.Len(t, decoded.SubTasks, 3)
	assert.Equal(t, "Sub 1", decoded.SubTasks[0].Title)
	assert.Equal(t, StatusDone, decoded.SubTasks[0].Status)
	assert.Equal(t, "Sub 2", decoded.SubTasks[1].Title)
	assert.Equal(t, StatusInProgress, decoded.SubTasks[1].Status)
	assert.Equal(t, "Sub 3", decoded.SubTasks[2].Title)
	assert.Equal(t, StatusPending, decoded.SubTasks[2].Status)
}

func TestTask_EmptyArrays(t *testing.T) {
	task := Task{
		ID:         "TASK-EMPTY",
		Title:      "Empty Arrays",
		Status:     StatusPending,
		Scope:      []string{},
		SpecRefs:   []string{},
		Inputs:     []string{},
		Outputs:    []string{},
		Acceptance: []string{},
		SubTasks:   []SubTask{},
	}

	data, err := yaml.Marshal(&task)
	assert.NoError(t, err)

	var decoded Task
	err = yaml.Unmarshal(data, &decoded)
	assert.NoError(t, err)

	// Empty slices should decode as nil in Go
	assert.Empty(t, decoded.Scope)
	assert.Empty(t, decoded.SpecRefs)
	assert.Empty(t, decoded.Inputs)
	assert.Empty(t, decoded.Outputs)
	assert.Empty(t, decoded.Acceptance)
	assert.Empty(t, decoded.SubTasks)
}
