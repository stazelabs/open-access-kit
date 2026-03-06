package source

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stazelabs/open-access-kit/internal/config"
	"github.com/stazelabs/open-access-kit/internal/verify"
)

type githubSource struct {
	name string
	cfg  config.SourceConfig
}

func newGitHub(name string, cfg config.SourceConfig) *githubSource {
	return &githubSource{name: name, cfg: cfg}
}

type ghRelease struct {
	TagName    string    `json:"tag_name"`
	PreRelease bool      `json:"prerelease"`
	Assets     []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func (s *githubSource) Name() string { return s.name }

func (s *githubSource) fetchLatestRelease(ctx context.Context) (*ghRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases", s.cfg.GitHubRepo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API HTTP %d for %s", resp.StatusCode, url)
	}
	var releases []ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&releases); err != nil {
		return nil, fmt.Errorf("decoding GitHub releases: %w", err)
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("no releases found for %s", s.cfg.GitHubRepo)
	}
	for _, r := range releases {
		if !r.PreRelease {
			return &r, nil
		}
	}
	return &releases[0], nil
}

func (s *githubSource) DetectVersion(ctx context.Context) (string, error) {
	rel, err := s.fetchLatestRelease(ctx)
	if err != nil {
		return "", err
	}
	return rel.TagName, nil
}

func (s *githubSource) Download(ctx context.Context, mirrorDir string, opts DownloadOptions) error {
	dest := filepath.Join(mirrorDir, s.name)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	rel, err := s.fetchLatestRelease(ctx)
	if err != nil {
		return fmt.Errorf("fetching release: %w", err)
	}
	patterns, err := compilePatterns(s.cfg.AssetPatterns)
	if err != nil {
		return err
	}
	for _, asset := range rel.Assets {
		if !matchesAny(asset.Name, patterns) {
			continue
		}
		dst := filepath.Join(dest, asset.Name)
		if !opts.Force && fileExists(dst) {
			continue
		}
		if err := downloadHTTP(ctx, asset.BrowserDownloadURL, dst); err != nil {
			return fmt.Errorf("downloading %s: %w", asset.Name, err)
		}
	}
	return nil
}

func (s *githubSource) Verify(ctx context.Context, mirrorDir string) error {
	if s.cfg.Verify.Method != "gpg" {
		return nil
	}
	dest := filepath.Join(mirrorDir, s.name)
	return filepath.WalkDir(dest, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		name := d.Name()
		if strings.HasSuffix(name, ".asc") || strings.HasSuffix(name, ".sig") {
			return nil
		}
		sigFile, err := findSig(dest, name)
		if err != nil {
			return nil // no sig found, skip silently
		}
		return verify.GPG(ctx, s.cfg.Verify.Keyring, sigFile, path)
	})
}

func (s *githubSource) Size(mirrorDir string) (int64, error) {
	return dirSize(filepath.Join(mirrorDir, s.name))
}

func (s *githubSource) Stage(ctx context.Context, mirrorDir, imageDir string, tier config.TierConfig) error {
	src := filepath.Join(mirrorDir, s.name)
	dst := filepath.Join(imageDir, s.cfg.StagePath)
	return copyDir(src, dst)
}

func compilePatterns(patterns []string) ([]*regexp.Regexp, error) {
	compiled := make([]*regexp.Regexp, 0, len(patterns))
	for _, p := range patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, fmt.Errorf("invalid asset_pattern %q: %w", p, err)
		}
		compiled = append(compiled, re)
	}
	return compiled, nil
}

func matchesAny(name string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(name) {
			return true
		}
	}
	return false
}
