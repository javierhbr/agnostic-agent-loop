package platform

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"platform-truth-mcp-server/internal/config"
)

var (
	platformRefPattern = regexp.MustCompile(`\b(?:principles|capabilities|contracts|adrs)\.[A-Za-z0-9._-]+\b`)
	jiraKeyPattern     = regexp.MustCompile(`\b[A-Z][A-Z0-9]+-\d+\b`)
	versionPattern     = regexp.MustCompile(`(?i)\bversion\s*[:=]\s*[` + "`\"" + `]?([0-9]{4}\.[0-9]{2}(?:\.[0-9]+)?)`)
)

type Repository struct {
	cfg *config.Config
}

type RefMatch struct {
	ID   string `json:"id"`
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text"`
	Kind string `json:"kind"`
}

type JiraMatch struct {
	Key  string `json:"key"`
	File string `json:"file"`
	Line int    `json:"line"`
	Text string `json:"text"`
}

func NewRepository(cfg *config.Config) *Repository {
	return &Repository{cfg: cfg}
}

func (r *Repository) LatestVersion(ctx context.Context) (string, string, error) {
	if version, source, err := r.latestVersionFromGit(ctx); err == nil && version != "" {
		return version, source, nil
	}
	matches, err := r.searchText(versionPattern)
	if err != nil {
		return "", "", err
	}
	versions := map[string]string{}
	for _, match := range matches {
		extracted := versionPattern.FindStringSubmatch(match.Text)
		if len(extracted) == 2 {
			versions[extracted[1]] = match.File
		}
	}
	if len(versions) == 0 {
		return "", "", fmt.Errorf("could not determine platform version")
	}
	var ordered []string
	for version := range versions {
		ordered = append(ordered, version)
	}
	sort.Strings(ordered)
	latest := ordered[len(ordered)-1]
	return latest, versions[latest], nil
}

func (r *Repository) latestVersionFromGit(ctx context.Context) (string, string, error) {
	cmd := exec.CommandContext(ctx, "git", "-C", r.cfg.PlatformRepoAbs, "tag", "--sort=-refname")
	output, err := cmd.Output()
	if err != nil {
		return "", "", err
	}
	for _, line := range strings.Split(string(output), "\n") {
		candidate := strings.TrimSpace(line)
		if candidate == "" {
			continue
		}
		return candidate, "git-tag", nil
	}
	return "", "", fmt.Errorf("no git tags found")
}

type textMatch struct {
	File string
	Line int
	Text string
}

func (r *Repository) ListRefs() ([]RefMatch, error) {
	files, err := r.candidateFiles()
	if err != nil {
		return nil, err
	}
	seen := map[string]RefMatch{}
	for _, file := range files {
		matches, err := collectPatternMatches(file, platformRefPattern)
		if err != nil {
			continue
		}
		for _, match := range matches {
			for _, id := range platformRefPattern.FindAllString(match.Text, -1) {
				key := id + "|" + match.File + fmt.Sprintf("|%d", match.Line)
				seen[key] = RefMatch{ID: id, File: match.File, Line: match.Line, Text: match.Text, Kind: kindForRef(id)}
			}
		}
	}
	results := make([]RefMatch, 0, len(seen))
	for _, item := range seen {
		results = append(results, item)
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].ID == results[j].ID {
			if results[i].File == results[j].File {
				return results[i].Line < results[j].Line
			}
			return results[i].File < results[j].File
		}
		return results[i].ID < results[j].ID
	})
	return results, nil
}

func (r *Repository) FindRef(id string) ([]RefMatch, error) {
	all, err := r.ListRefs()
	if err != nil {
		return nil, err
	}
	var filtered []RefMatch
	for _, item := range all {
		if item.ID == id {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (r *Repository) GetJiraMapping() ([]JiraMatch, error) {
	var files []string
	if r.cfg.Inputs.JiraMappingFile != "" {
		if _, err := os.Stat(r.cfg.Inputs.JiraMappingFile); err == nil {
			files = append(files, r.cfg.Inputs.JiraMappingFile)
		}
	}
	if len(files) == 0 {
		candidateFiles, err := r.candidateFiles()
		if err != nil {
			return nil, err
		}
		files = candidateFiles
	}

	seen := map[string]JiraMatch{}
	for _, file := range files {
		matches, err := collectPatternMatches(file, jiraKeyPattern)
		if err != nil {
			continue
		}
		for _, match := range matches {
			for _, key := range jiraKeyPattern.FindAllString(match.Text, -1) {
				signature := key + "|" + match.File + fmt.Sprintf("|%d", match.Line)
				seen[signature] = JiraMatch{Key: key, File: match.File, Line: match.Line, Text: match.Text}
			}
		}
	}
	results := make([]JiraMatch, 0, len(seen))
	for _, item := range seen {
		results = append(results, item)
	}
	sort.Slice(results, func(i, j int) bool {
		if results[i].Key == results[j].Key {
			if results[i].File == results[j].File {
				return results[i].Line < results[j].Line
			}
			return results[i].File < results[j].File
		}
		return results[i].Key < results[j].Key
	})
	return results, nil
}

func (r *Repository) Serialize(value any) string {
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", value)
	}
	return string(encoded)
}

func (r *Repository) searchText(pattern *regexp.Regexp) ([]textMatch, error) {
	files, err := r.candidateFiles()
	if err != nil {
		return nil, err
	}
	var out []textMatch
	for _, file := range files {
		matches, err := collectPatternMatches(file, pattern)
		if err != nil {
			continue
		}
		out = append(out, matches...)
	}
	return out, nil
}

func (r *Repository) candidateFiles() ([]string, error) {
	roots := dedupeNonEmpty(
		r.cfg.Inputs.SpecsDir,
		r.cfg.Inputs.ContractsDir,
		r.cfg.Inputs.ADRsDir,
		r.cfg.Inputs.RefsIndex,
		r.cfg.Inputs.JiraMappingFile,
	)
	if len(roots) == 0 {
		roots = append(roots, r.cfg.PlatformRepoAbs)
	}
	seen := map[string]bool{}
	var files []string
	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil {
			continue
		}
		if info.IsDir() {
			err = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
				if walkErr != nil {
					return nil
				}
				if d.IsDir() {
					name := d.Name()
					if name == ".git" || name == "node_modules" || name == "vendor" || name == ".cache" {
						return filepath.SkipDir
					}
					return nil
				}
				if !isTextFile(path) || seen[path] {
					return nil
				}
				seen[path] = true
				files = append(files, path)
				return nil
			})
			if err != nil {
				return nil, err
			}
			continue
		}
		if !seen[root] {
			seen[root] = true
			files = append(files, root)
		}
	}
	sort.Strings(files)
	return files, nil
}

func collectPatternMatches(path string, pattern *regexp.Regexp) ([]textMatch, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var out []textMatch
	lines := strings.Split(string(data), "\n")
	for idx, line := range lines {
		if pattern.MatchString(line) {
			out = append(out, textMatch{File: path, Line: idx + 1, Text: strings.TrimSpace(line)})
		}
	}
	return out, nil
}

func dedupeNonEmpty(values ...string) []string {
	seen := map[string]bool{}
	var out []string
	for _, value := range values {
		if value == "" || seen[value] {
			continue
		}
		seen[value] = true
		out = append(out, value)
	}
	return out
}

func isTextFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".md", ".mdx", ".txt", ".yaml", ".yml", ".json":
		return true
	default:
		return ext == ""
	}
}

func kindForRef(id string) string {
	parts := strings.SplitN(id, ".", 2)
	if len(parts) == 0 {
		return "unknown"
	}
	return parts[0]
}
