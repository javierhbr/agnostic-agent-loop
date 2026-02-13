package prompts

import "testing"

func TestNewRegistry(t *testing.T) {
	r := NewRegistry()
	if len(r.All()) == 0 {
		t.Fatal("registry should have built-in prompts")
	}
}

func TestByCategory(t *testing.T) {
	r := NewRegistry()
	for _, cat := range AllCategories() {
		items := r.ByCategory(cat)
		if len(items) == 0 {
			t.Errorf("category %s should have at least one prompt", cat)
		}
		for _, p := range items {
			if p.Category != cat {
				t.Errorf("prompt %s has category %s, expected %s", p.Slug, p.Category, cat)
			}
		}
	}
}

func TestFindBySlug(t *testing.T) {
	r := NewRegistry()
	p := r.FindBySlug("claim-and-implement")
	if p == nil {
		t.Fatal("should find claim-and-implement prompt")
	}
	if p.Title == "" || p.Content == "" {
		t.Error("prompt should have title and content")
	}
}

func TestFindBySlugNotFound(t *testing.T) {
	r := NewRegistry()
	p := r.FindBySlug("nonexistent-slug")
	if p != nil {
		t.Error("should return nil for unknown slug")
	}
}

func TestSlugsAreUnique(t *testing.T) {
	r := NewRegistry()
	seen := map[string]bool{}
	for _, p := range r.All() {
		if seen[p.Slug] {
			t.Errorf("duplicate slug: %s", p.Slug)
		}
		seen[p.Slug] = true
	}
}

func TestAllPromptsHaveRequiredFields(t *testing.T) {
	r := NewRegistry()
	for _, p := range r.All() {
		if p.Slug == "" {
			t.Error("prompt missing slug")
		}
		if p.Title == "" {
			t.Errorf("prompt %s missing title", p.Slug)
		}
		if p.Category == "" {
			t.Errorf("prompt %s missing category", p.Slug)
		}
		if p.Description == "" {
			t.Errorf("prompt %s missing description", p.Slug)
		}
		if p.Content == "" {
			t.Errorf("prompt %s missing content", p.Slug)
		}
	}
}
