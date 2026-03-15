package component

import (
	"fmt"
	"os"
	"path/filepath"

	"platform-truth-mcp-server/internal/config"
	"platform-truth-mcp-server/internal/yamlish"
)

type PlatformRef struct {
	PlatformID              string   `json:"platformId"`
	PlatformVersion         string   `json:"platformVersion"`
	BaselineRef             string   `json:"baselineRef"`
	ComponentName           string   `json:"componentName"`
	OwnerTeam               string   `json:"ownerTeam"`
	ChangePackageID         string   `json:"changePackageId"`
	AlignmentType           string   `json:"alignmentType"`
	RequiresPlatformChange  bool     `json:"requiresPlatformChange"`
	PlatformChangePackageID string   `json:"platformChangePackageId"`
	RefIDs                  []string `json:"refIds"`
}

type JiraTraceability struct {
	PlatformIssueKey    string   `json:"platformIssueKey"`
	PlatformVersion     string   `json:"platformVersion"`
	ComponentEpicKey    string   `json:"componentEpicKey"`
	LinkedPlatformIssue string   `json:"linkedPlatformIssue"`
	ChangePackageID     string   `json:"changePackageId"`
	StoryKeys           []string `json:"storyKeys"`
}

func LoadPlatformRef(componentRepo string, cfg *config.Config) (*PlatformRef, error) {
	path := filepath.Join(componentRepo, cfg.ComponentAlignment.PlatformRefFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read platform ref: %w", err)
	}
	parsed, err := yamlish.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse platform ref: %w", err)
	}

	ref := &PlatformRef{
		PlatformID:              yamlish.LookupString(parsed, "platform", "id"),
		PlatformVersion:         yamlish.LookupString(parsed, "platform", "version"),
		BaselineRef:             yamlish.LookupString(parsed, "platform", "baseline_ref"),
		ComponentName:           yamlish.LookupString(parsed, "component", "name"),
		OwnerTeam:               yamlish.LookupString(parsed, "component", "owner_team"),
		ChangePackageID:         yamlish.LookupString(parsed, "change", "change_package_id"),
		AlignmentType:           yamlish.LookupString(parsed, "change", "alignment_type"),
		RequiresPlatformChange:  yamlish.LookupBool(parsed, "change", "requires_platform_change"),
		PlatformChangePackageID: yamlish.LookupString(parsed, "change", "platform_change_package_id"),
		RefIDs:                  extractRefIDs(parsed),
	}
	return ref, nil
}

func LoadJiraTraceability(componentRepo string, cfg *config.Config) (*JiraTraceability, error) {
	path := filepath.Join(componentRepo, cfg.ComponentAlignment.JiraTraceabilityFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read jira traceability: %w", err)
	}
	parsed, err := yamlish.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("parse jira traceability: %w", err)
	}

	stories := yamlish.AsSlice(yamlish.Lookup(parsed, "stories"))
	storyKeys := make([]string, 0, len(stories))
	for _, item := range stories {
		key := yamlish.LookupString(item, "key")
		if key != "" {
			storyKeys = append(storyKeys, key)
		}
	}

	trace := &JiraTraceability{
		PlatformIssueKey:    yamlish.LookupString(parsed, "platform_issue", "key"),
		PlatformVersion:     yamlish.LookupString(parsed, "platform_issue", "platform_version"),
		ComponentEpicKey:    yamlish.LookupString(parsed, "component_epic", "key"),
		LinkedPlatformIssue: yamlish.LookupString(parsed, "component_epic", "linked_platform_issue"),
		ChangePackageID:     yamlish.LookupString(parsed, "component_epic", "change_package_id"),
		StoryKeys:           storyKeys,
	}
	return trace, nil
}

func extractRefIDs(parsed any) []string {
	platformRefs := yamlish.AsMap(yamlish.Lookup(parsed, "platform_refs"))
	if platformRefs == nil {
		return nil
	}
	seen := map[string]bool{}
	var ids []string
	for _, value := range platformRefs {
		items := yamlish.AsSlice(value)
		for _, item := range items {
			id := yamlish.LookupString(item, "id")
			if id == "" || seen[id] {
				continue
			}
			seen[id] = true
			ids = append(ids, id)
		}
	}
	return ids
}
