package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config is the top-level structure parsed from oak.yaml.
type Config struct {
	Release      string                  `yaml:"release"`
	DownloadRoot string                  `yaml:"download_root"`
	Paths        PathConfig              `yaml:"paths"`
	Tiers        map[string]TierConfig   `yaml:"tiers"`
	Sources      map[string]SourceConfig `yaml:"sources"`
	Signing      SigningConfig           `yaml:"signing"`
}

type PathConfig struct {
	Mirror string `yaml:"mirror"`
	Image  string `yaml:"image"`
	Output string `yaml:"output"`
}

type TierConfig struct {
	Label    string   `yaml:"label"`
	BudgetGB float64  `yaml:"budget_gb"`
	Sources  []string `yaml:"sources"`
}

type SourceConfig struct {
	Type        string `yaml:"type"`
	Description string `yaml:"description"`

	// rsync sources
	RsyncBase string   `yaml:"rsync_base"`
	Files     []string `yaml:"files"`

	// git sources
	GitURL         string `yaml:"git_url"`
	GitBranch      string `yaml:"git_branch"`
	Shallow        bool   `yaml:"shallow"`
	RenderMarkdown bool   `yaml:"render_markdown"`

	// site-mirror sources
	MirrorURL string `yaml:"mirror_url"`

	// local sources
	LocalPath string `yaml:"local_path"`

	// github-release sources
	GitHubRepo    string   `yaml:"github_repo"`
	AssetPatterns []string `yaml:"asset_patterns"`

	// kiwix-zim sources
	ZimFiles []ZimFileConfig `yaml:"zim_files"`

	// common
	ExcludeDirs   []string            `yaml:"exclude_dirs"`
	StagePath     string              `yaml:"stage_path"`
	Verify        VerifyConfig        `yaml:"verify"`
	VersionDetect VersionDetectConfig `yaml:"version_detect"`
}

type VerifyConfig struct {
	Method  string `yaml:"method"`
	Keyring string `yaml:"keyring"`
}

type VersionDetectConfig struct {
	Method  string `yaml:"method"`
	URL     string `yaml:"url"`
	Pattern string `yaml:"pattern"`
	Select  string `yaml:"select"`
}

type ZimFileConfig struct {
	Name        string `yaml:"name"`         // base name without date, e.g. "zimgit-medicine_en"
	Category    string `yaml:"category"`     // Kiwix download category (unused when download_url is set)
	StageSubdir string `yaml:"stage_subdir"` // subdirectory under stage_path, e.g. "medical"
	DownloadURL string `yaml:"download_url"` // optional: full URL, bypasses OPDS lookup
}

type SigningConfig struct {
	Enabled   bool   `yaml:"enabled"`
	KeyID     string `yaml:"key_id"`
	PublicKey string `yaml:"public_key"`
}

// Load reads and parses the config file at path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config %s: %w", path, err)
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config %s: %w", path, err)
	}
	if cfg.Paths.Mirror == "" {
		cfg.Paths.Mirror = "./mirror"
	}
	if cfg.Paths.Image == "" {
		cfg.Paths.Image = "./image"
	}
	if cfg.Paths.Output == "" {
		cfg.Paths.Output = "./dist"
	}
	if cfg.Release == "" {
		cfg.Release = autoRelease()
	}
	return &cfg, nil
}

// autoRelease generates a release name like "Q126" from the current date.
func autoRelease() string {
	t := time.Now()
	q := (int(t.Month())-1)/3 + 1
	year := t.Year() % 100
	return fmt.Sprintf("Q%d%02d", q, year)
}
