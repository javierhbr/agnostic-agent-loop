package adapters

import (
	"context"

	"github.com/javierbenavides/agentic-agent/pkg/models"
)

// OpenSpecAdapter bridges to the OpenSpec CLI/logic for component proposals and delta-specs.
type OpenSpecAdapter interface {
	Specify(ctx context.Context, input string) (proposalPath string, err error)
	Bridge(ctx context.Context, source string, dest string) error
}

// SpecKitAdapter bridges to the Spec Kit CLI/logic for context and durable rules.
type SpecKitAdapter interface {
	HydrateContext(ctx context.Context) error
	CheckClarity(ctx context.Context, target string) (clarityScore int, err error)
}

// BMADAdapter bridges to the BMAD CLI/logic for track selection and scoping.
type BMADAdapter interface {
	Route(ctx context.Context, request string) (track string, err error)
	Size(ctx context.Context, request string) (size string, err error)
}

// PlannerAdapter interfaces with delivery planners like Jira or Linear.
type PlannerAdapter interface {
	Sync(ctx context.Context, tasks []models.Task) error
	FetchState(ctx context.Context, taskID string) (state string, err error)
}
