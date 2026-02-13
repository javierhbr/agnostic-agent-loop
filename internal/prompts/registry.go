package prompts

// Category represents a prompt category.
type Category string

const (
	CategoryAgentPrompt    Category = "agent-prompt"
	CategoryCLIExample     Category = "cli-example"
	CategoryWorkflowRecipe Category = "workflow-recipe"
)

// CategoryInfo maps category slugs to display names.
var CategoryInfo = map[Category]string{
	CategoryAgentPrompt:    "Agent Prompts",
	CategoryCLIExample:     "CLI Examples",
	CategoryWorkflowRecipe: "Workflow Recipes",
}

// AllCategories returns categories in display order.
func AllCategories() []Category {
	return []Category{CategoryAgentPrompt, CategoryCLIExample, CategoryWorkflowRecipe}
}

// Prompt represents a single example prompt.
type Prompt struct {
	Slug        string
	Title       string
	Category    Category
	Description string
	Content     string
	Tags        []string
}

// Registry holds all built-in prompts.
type Registry struct {
	prompts []Prompt
}

// NewRegistry creates a registry with all built-in prompts.
func NewRegistry() *Registry {
	return &Registry{prompts: builtinPrompts}
}

// All returns all prompts.
func (r *Registry) All() []Prompt {
	return r.prompts
}

// ByCategory returns prompts filtered by category.
func (r *Registry) ByCategory(cat Category) []Prompt {
	var result []Prompt
	for _, p := range r.prompts {
		if p.Category == cat {
			result = append(result, p)
		}
	}
	return result
}

// FindBySlug returns a prompt by slug, or nil if not found.
func (r *Registry) FindBySlug(slug string) *Prompt {
	for _, p := range r.prompts {
		if p.Slug == slug {
			return &p
		}
	}
	return nil
}
