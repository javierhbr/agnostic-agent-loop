package sdd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// SpecGraph represents the entire spec dependency graph.
type SpecGraph struct {
	Nodes map[string]SpecGraphNode `json:"nodes" yaml:"nodes"`
}

// Load reads a spec graph from disk. Auto-detects JSON or YAML by file extension.
func (g *SpecGraph) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return an empty graph if file doesn't exist
			g.Nodes = make(map[string]SpecGraphNode)
			return nil
		}
		return fmt.Errorf("failed to read spec graph: %w", err)
	}

	// Determine format by file extension
	if strings.HasSuffix(path, ".json") {
		if err := json.Unmarshal(data, &g.Nodes); err != nil {
			return fmt.Errorf("failed to parse JSON spec graph: %w", err)
		}
	} else {
		if err := yaml.Unmarshal(data, &g.Nodes); err != nil {
			return fmt.Errorf("failed to parse YAML spec graph: %w", err)
		}
	}

	if g.Nodes == nil {
		g.Nodes = make(map[string]SpecGraphNode)
	}

	return nil
}

// Save writes the spec graph to disk. Format determined by file extension.
func (g *SpecGraph) Save(path string) error {
	// Ensure parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create spec graph directory: %w", err)
	}

	var data []byte
	var err error

	if strings.HasSuffix(path, ".json") {
		data, err = json.MarshalIndent(g.Nodes, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON spec graph: %w", err)
		}
	} else {
		data, err = yaml.Marshal(g.Nodes)
		if err != nil {
			return fmt.Errorf("failed to marshal YAML spec graph: %w", err)
		}
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write spec graph: %w", err)
	}

	return nil
}

// Upsert adds or updates a node in the spec graph. Sets UpdatedAt to now.
func (g *SpecGraph) Upsert(node SpecGraphNode) {
	node.UpdatedAt = time.Now()
	g.Nodes[node.ID] = node
}

// Get retrieves a node by ID. Returns the node and a bool indicating if it was found.
func (g *SpecGraph) Get(id string) (SpecGraphNode, bool) {
	node, ok := g.Nodes[id]
	return node, ok
}

// ListBlocked returns all nodes with non-empty BlockedBy.
func (g *SpecGraph) ListBlocked() []SpecGraphNode {
	var blocked []SpecGraphNode
	for _, node := range g.Nodes {
		if len(node.BlockedBy) > 0 {
			blocked = append(blocked, node)
		}
	}
	return blocked
}

// SyncToRemote copies the local spec graph to the platform repo location.
// Typically used to sync .agentic/spec-graph.json to graph/index.yaml in the platform repo.
func (g *SpecGraph) SyncToRemote(localPath, remotePath string) error {
	// Read from local path
	if err := g.Load(localPath); err != nil {
		return fmt.Errorf("failed to load local spec graph: %w", err)
	}

	// Write to remote path
	if err := g.Save(remotePath); err != nil {
		return fmt.Errorf("failed to save to remote spec graph: %w", err)
	}

	return nil
}
