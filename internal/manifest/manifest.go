package manifest

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/stazelabs/open-access-kit/internal/config"
	"github.com/stazelabs/open-access-kit/internal/source"
)

// Manifest is a machine-readable snapshot of a release, capturing versions
// and composition so that two releases can be compared programmatically.
type Manifest struct {
	Release    string                  `json:"release"`
	Built      time.Time               `json:"built"`
	OAKVersion string                  `json:"oak_version,omitempty"`
	Tiers      map[string]TierManifest `json:"tiers"`
	Sources    map[string]SourceEntry  `json:"sources"`
}

// TierManifest records which sources belong to a tier.
type TierManifest struct {
	Label    string   `json:"label"`
	BudgetGB float64  `json:"budget_gb"`
	Sources  []string `json:"sources"`
}

// SourceEntry captures the identity and version of a single source at build time.
type SourceEntry struct {
	Type        string                  `json:"type"`
	Description string                  `json:"description"`
	Version     string                  `json:"version"`
	Commit      string                  `json:"commit,omitempty"`
	ZimFiles    []source.ZimFileEntry   `json:"zim_files,omitempty"`
}

// zimInfoProvider is implemented by kiwix sources to expose per-file ZIM metadata.
type zimInfoProvider interface {
	ZimInfo(mirrorDir string) ([]source.ZimFileEntry, error)
}

// Generate builds a Manifest from the current config, detecting versions
// for each source and reading git HEAD commits from the mirror directory.
func Generate(ctx context.Context, cfg *config.Config, allSources []source.Source, mirrorDir string) (*Manifest, error) {
	m := &Manifest{
		Release: cfg.Release,
		Built:   time.Now().UTC(),
		Tiers:   make(map[string]TierManifest, len(cfg.Tiers)),
		Sources: make(map[string]SourceEntry, len(cfg.Sources)),
	}

	// Record tier composition.
	for key, tc := range cfg.Tiers {
		m.Tiers[key] = TierManifest{
			Label:    tc.Label,
			BudgetGB: tc.BudgetGB,
			Sources:  tc.Sources,
		}
	}

	// Build a lookup of source objects by name for version detection.
	srcByName := make(map[string]source.Source, len(allSources))
	for _, s := range allSources {
		srcByName[s.Name()] = s
	}

	// Record each source's type, description, version, and commit.
	for name, scfg := range cfg.Sources {
		entry := SourceEntry{
			Type:        scfg.Type,
			Description: scfg.Description,
		}

		// Detect upstream version if the source object is available.
		if s, ok := srcByName[name]; ok {
			ver, err := s.DetectVersion(ctx)
			if err == nil {
				entry.Version = ver
			}
		}

		// For git sources, read the HEAD commit from the mirror.
		if scfg.Type == "git" {
			commit, err := gitHeadCommit(filepath.Join(mirrorDir, name))
			if err == nil {
				entry.Commit = commit
			}
		}

		// For kiwix-zim sources, collect per-file resolved metadata from the mirror.
		if zp, ok := srcByName[name].(zimInfoProvider); ok {
			if files, err := zp.ZimInfo(mirrorDir); err == nil {
				entry.ZimFiles = files
			}
		}

		m.Sources[name] = entry
	}

	return m, nil
}

// Write serialises the manifest as indented JSON to the given path.
func Write(m *Manifest, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory for manifest: %w", err)
	}
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling manifest: %w", err)
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}

// Read deserialises a manifest from a JSON file.
func Read(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest %s: %w", path, err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest %s: %w", path, err)
	}
	return &m, nil
}

// gitHeadCommit returns the full SHA of HEAD in the given git repository.
func gitHeadCommit(repoDir string) (string, error) {
	cmd := exec.Command("git", "-C", repoDir, "rev-parse", "HEAD")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
