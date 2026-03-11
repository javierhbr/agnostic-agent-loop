package service

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"platform-truth-mcp-server/internal/component"
	"platform-truth-mcp-server/internal/config"
	"platform-truth-mcp-server/internal/platform"
)

type Tool struct {
	Name        string
	Description string
	InputSchema map[string]any
	Handler     func(context.Context, map[string]any) (ToolResult, error)
}

type ToolResult struct {
	StructuredContent map[string]any `json:"structuredContent,omitempty"`
	Text              string         `json:"-"`
	IsError           bool           `json:"isError,omitempty"`
}

type Service struct {
	cfg     *config.Config
	repo    *platform.Repository
	version string
	tools   map[string]Tool
}

func New(cfg *config.Config, version string) (*Service, error) {
	svc := &Service{
		cfg:     cfg,
		repo:    platform.NewRepository(cfg),
		version: version,
		tools:   map[string]Tool{},
	}
	for _, tool := range svc.buildTools() {
		svc.tools[tool.Name] = tool
	}
	return svc, nil
}

func (s *Service) ToolDefinitions() []map[string]any {
	names := make([]string, 0, len(s.tools))
	for name := range s.tools {
		names = append(names, name)
	}
	sort.Strings(names)
	defs := make([]map[string]any, 0, len(names))
	for _, name := range names {
		tool := s.tools[name]
		defs = append(defs, map[string]any{
			"name":        tool.Name,
			"description": tool.Description,
			"inputSchema": tool.InputSchema,
		})
	}
	return defs
}

func (s *Service) Call(ctx context.Context, name string, args map[string]any) (map[string]any, error) {
	tool, ok := s.tools[name]
	if !ok {
		return nil, fmt.Errorf("unknown tool %q", name)
	}
	result, err := tool.Handler(ctx, args)
	if err != nil {
		result = ToolResult{
			StructuredContent: map[string]any{"tool": name, "error": err.Error()},
			Text:              err.Error(),
			IsError:           true,
		}
	}
	response := map[string]any{
		"content": []map[string]any{{"type": "text", "text": result.Text}},
	}
	if result.StructuredContent != nil {
		response["structuredContent"] = result.StructuredContent
	}
	if result.IsError {
		response["isError"] = true
	}
	return response, nil
}

func (s *Service) buildTools() []Tool {
	return []Tool{
		{
			Name:        "get_platform_version",
			Description: "Return the platform version. Defaults to the component's pinned version when componentRepoPath is provided.",
			InputSchema: schema(
				prop("componentRepoPath", "string", "Absolute or relative path to the component repository."),
				prop("mode", "string", "Version mode: pinned or latest."),
			),
			Handler: s.getPlatformVersion,
		},
		{
			Name:        "list_platform_refs",
			Description: "List discoverable platform refs from the platform repository.",
			InputSchema: schema(prop("kind", "string", "Optional ref kind filter: principles, capabilities, contracts, adrs.")),
			Handler:     s.listPlatformRefs,
		},
		{
			Name:        "get_platform_ref",
			Description: "Find a specific platform ref and return its locations.",
			InputSchema: schemaRequired([]string{"id"}, prop("id", "string", "Exact platform ref id to find.")),
			Handler:     s.getPlatformRef,
		},
		{
			Name:        "get_jira_mapping",
			Description: "Return JIRA-linked issue keys that are discoverable from platform artifacts.",
			InputSchema: schema(),
			Handler:     s.getJiraMapping,
		},
		{
			Name:        "validate_component_alignment",
			Description: "Validate a component's platform refs against the platform repository and the pinned version.",
			InputSchema: schemaRequired([]string{"componentRepoPath"}, prop("componentRepoPath", "string", "Path to the component repository.")),
			Handler:     s.validateComponentAlignment,
		},
		{
			Name:        "validate_component_jira_chain",
			Description: "Validate a component's JIRA traceability chain and cross-links.",
			InputSchema: schemaRequired([]string{"componentRepoPath"}, prop("componentRepoPath", "string", "Path to the component repository.")),
			Handler:     s.validateComponentJiraChain,
		},
		{
			Name:        "detect_platform_drift_from_pinned_version",
			Description: "Compare the component's pinned platform version with the latest known version.",
			InputSchema: schemaRequired([]string{"componentRepoPath"}, prop("componentRepoPath", "string", "Path to the component repository.")),
			Handler:     s.detectPlatformDrift,
		},
	}
}

func (s *Service) getPlatformVersion(ctx context.Context, args map[string]any) (ToolResult, error) {
	mode := stringArg(args, "mode")
	if mode == "" {
		mode = s.cfg.PlatformRepo.DefaultRefMode
	}
	if mode == "latest" && !s.cfg.PlatformRepo.AllowLatestQuery {
		return ToolResult{}, fmt.Errorf("latest platform queries are disabled by config")
	}

	result := map[string]any{"mode": mode}
	if componentRepo := stringArg(args, "componentRepoPath"); componentRepo != "" && mode == "pinned" {
		ref, err := component.LoadPlatformRef(resolveArgPath(componentRepo), s.cfg)
		if err != nil {
			return ToolResult{}, err
		}
		result["version"] = ref.PlatformVersion
		result["source"] = "component/platform-ref.yaml"
		result["component"] = ref.ComponentName
		text := s.repo.Serialize(result)
		return ToolResult{StructuredContent: result, Text: text}, nil
	}
	if mode == "pinned" {
		return ToolResult{}, fmt.Errorf("componentRepoPath is required when mode is pinned")
	}

	latest, source, err := s.repo.LatestVersion(ctx)
	if err != nil {
		return ToolResult{}, err
	}
	result["version"] = latest
	result["source"] = source
	text := s.repo.Serialize(result)
	return ToolResult{StructuredContent: result, Text: text}, nil
}

func (s *Service) listPlatformRefs(_ context.Context, args map[string]any) (ToolResult, error) {
	kind := strings.TrimSpace(stringArg(args, "kind"))
	refs, err := s.repo.ListRefs()
	if err != nil {
		return ToolResult{}, err
	}
	var filtered []platform.RefMatch
	for _, ref := range refs {
		if kind == "" || ref.Kind == kind {
			filtered = append(filtered, ref)
		}
	}
	result := map[string]any{"count": len(filtered), "refs": filtered}
	return ToolResult{StructuredContent: result, Text: s.repo.Serialize(result)}, nil
}

func (s *Service) getPlatformRef(_ context.Context, args map[string]any) (ToolResult, error) {
	id := stringArg(args, "id")
	if id == "" {
		return ToolResult{}, fmt.Errorf("id is required")
	}
	matches, err := s.repo.FindRef(id)
	if err != nil {
		return ToolResult{}, err
	}
	result := map[string]any{"id": id, "count": len(matches), "matches": matches}
	return ToolResult{StructuredContent: result, Text: s.repo.Serialize(result)}, nil
}

func (s *Service) getJiraMapping(_ context.Context, _ map[string]any) (ToolResult, error) {
	matches, err := s.repo.GetJiraMapping()
	if err != nil {
		return ToolResult{}, err
	}
	result := map[string]any{"count": len(matches), "issues": matches}
	return ToolResult{StructuredContent: result, Text: s.repo.Serialize(result)}, nil
}

func (s *Service) validateComponentAlignment(ctx context.Context, args map[string]any) (ToolResult, error) {
	componentRepoPath := resolveArgPath(requiredPath(args, "componentRepoPath"))
	if componentRepoPath == "" {
		return ToolResult{}, fmt.Errorf("componentRepoPath is required")
	}
	platformRef, err := component.LoadPlatformRef(componentRepoPath, s.cfg)
	if err != nil {
		return ToolResult{}, err
	}
	latestVersion, latestSource, latestErr := s.repo.LatestVersion(ctx)

	missing := []string{}
	found := map[string][]platform.RefMatch{}
	for _, id := range platformRef.RefIDs {
		matches, err := s.repo.FindRef(id)
		if err != nil || len(matches) == 0 {
			missing = append(missing, id)
			continue
		}
		found[id] = matches
	}

	result := map[string]any{
		"component":               platformRef.ComponentName,
		"platformVersion":         platformRef.PlatformVersion,
		"latestVersion":           latestVersion,
		"latestVersionSource":     latestSource,
		"alignmentType":           platformRef.AlignmentType,
		"requiresPlatformChange":  platformRef.RequiresPlatformChange,
		"platformChangePackageId": platformRef.PlatformChangePackageID,
		"missingRefs":             missing,
		"resolvedRefs":            found,
		"valid":                   len(missing) == 0,
	}
	if latestErr != nil {
		result["latestVersionError"] = latestErr.Error()
	}
	text := s.repo.Serialize(result)
	return ToolResult{StructuredContent: result, Text: text, IsError: len(missing) > 0}, nil
}

func (s *Service) validateComponentJiraChain(_ context.Context, args map[string]any) (ToolResult, error) {
	componentRepoPath := resolveArgPath(requiredPath(args, "componentRepoPath"))
	if componentRepoPath == "" {
		return ToolResult{}, fmt.Errorf("componentRepoPath is required")
	}
	trace, err := component.LoadJiraTraceability(componentRepoPath, s.cfg)
	if err != nil {
		return ToolResult{}, err
	}
	platformRef, _ := component.LoadPlatformRef(componentRepoPath, s.cfg)
	platformIssues, _ := s.repo.GetJiraMapping()

	issuesPresent := false
	for _, issue := range platformIssues {
		if issue.Key == trace.PlatformIssueKey {
			issuesPresent = true
			break
		}
	}

	warnings := []string{}
	if trace.LinkedPlatformIssue != trace.PlatformIssueKey {
		warnings = append(warnings, "component epic linked_platform_issue does not match platform_issue.key")
	}
	if platformRef != nil && platformRef.ChangePackageID != trace.ChangePackageID {
		warnings = append(warnings, "change package id differs between platform-ref.yaml and jira-traceability.yaml")
	}
	if platformRef != nil && platformRef.PlatformVersion != trace.PlatformVersion {
		warnings = append(warnings, "platform version differs between platform-ref.yaml and jira-traceability.yaml")
	}
	if !issuesPresent {
		warnings = append(warnings, "platform issue key was not found in the platform JIRA mapping inputs")
	}

	result := map[string]any{
		"platformIssueKey": trace.PlatformIssueKey,
		"componentEpicKey": trace.ComponentEpicKey,
		"storyKeys":        trace.StoryKeys,
		"storyCount":       len(trace.StoryKeys),
		"warnings":         warnings,
		"valid":            len(warnings) == 0,
	}
	return ToolResult{StructuredContent: result, Text: s.repo.Serialize(result), IsError: len(warnings) > 0}, nil
}

func (s *Service) detectPlatformDrift(ctx context.Context, args map[string]any) (ToolResult, error) {
	componentRepoPath := resolveArgPath(requiredPath(args, "componentRepoPath"))
	if componentRepoPath == "" {
		return ToolResult{}, fmt.Errorf("componentRepoPath is required")
	}
	platformRef, err := component.LoadPlatformRef(componentRepoPath, s.cfg)
	if err != nil {
		return ToolResult{}, err
	}
	latestVersion, source, err := s.repo.LatestVersion(ctx)
	if err != nil {
		return ToolResult{}, err
	}
	result := map[string]any{
		"component":       platformRef.ComponentName,
		"pinnedVersion":   platformRef.PlatformVersion,
		"latestVersion":   latestVersion,
		"latestSource":    source,
		"isDriftDetected": platformRef.PlatformVersion != latestVersion,
	}
	return ToolResult{StructuredContent: result, Text: s.repo.Serialize(result), IsError: platformRef.PlatformVersion != latestVersion}, nil
}

func schema(properties ...map[string]any) map[string]any {
	return schemaRequired(nil, properties...)
}

func schemaRequired(required []string, properties ...map[string]any) map[string]any {
	propMap := map[string]any{}
	for _, property := range properties {
		propMap[property["name"].(string)] = map[string]any{
			"type":        property["type"],
			"description": property["description"],
		}
	}
	result := map[string]any{
		"type":       "object",
		"properties": propMap,
	}
	if len(required) > 0 {
		result["required"] = required
	}
	return result
}

func prop(name, typ, description string) map[string]any {
	return map[string]any{"name": name, "type": typ, "description": description}
}

func stringArg(args map[string]any, key string) string {
	if args == nil {
		return ""
	}
	value, _ := args[key].(string)
	return strings.TrimSpace(value)
}

func requiredPath(args map[string]any, key string) string {
	return stringArg(args, key)
}

func resolveArgPath(path string) string {
	if path == "" {
		return ""
	}
	if filepath.IsAbs(path) {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	return abs
}
