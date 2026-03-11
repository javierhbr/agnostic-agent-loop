package adapters

import (
	"context"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// MockOpenSpecAdapter is a mock implementation of OpenSpecAdapter
type MockOpenSpecAdapter struct{}

func (m *MockOpenSpecAdapter) Specify(ctx context.Context, input string) (string, error) {
	return "mock/proposal.md", nil
}

func (m *MockOpenSpecAdapter) Bridge(ctx context.Context, source string, dest string) error {
	return nil
}

// MockSpecKitAdapter is a mock implementation of SpecKitAdapter
type MockSpecKitAdapter struct{}

func (m *MockSpecKitAdapter) HydrateContext(ctx context.Context) error {
	return nil
}

func (m *MockSpecKitAdapter) CheckClarity(ctx context.Context, target string) (int, error) {
	return 100, nil // Perfect clarity mock
}

// MockBMADAdapter is a mock implementation of BMADAdapter
type MockBMADAdapter struct{}

func (m *MockBMADAdapter) Route(ctx context.Context, request string) (string, error) {
	return "Standard Track", nil
}

func (m *MockBMADAdapter) Size(ctx context.Context, request string) (string, error) {
	return "Medium", nil
}

// MockPlannerAdapter is a mock implementation of PlannerAdapter
type MockPlannerAdapter struct{}

func (m *MockPlannerAdapter) Sync(ctx context.Context, tasks []models.Task) error {
	return nil
}

func (m *MockPlannerAdapter) FetchState(ctx context.Context, taskID string) (string, error) {
	return "In Progress", nil
}
