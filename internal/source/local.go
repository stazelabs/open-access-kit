package source

import (
	"context"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/stazelabs/open-access-kit/internal/config"
	"github.com/stazelabs/open-access-kit/internal/site"
)

type localSource struct {
	name string
	cfg  config.SourceConfig
	vars map[string]any
}

func newLocal(name string, cfg config.SourceConfig, vars map[string]any) *localSource {
	return &localSource{name: name, cfg: cfg, vars: vars}
}

func (s *localSource) Name() string { return s.name }

func (s *localSource) DetectVersion(ctx context.Context) (string, error) {
	return "", nil
}

func (s *localSource) Download(ctx context.Context, mirrorDir string, opts DownloadOptions) error {
	// Local sources don't need downloading.
	return nil
}

func (s *localSource) Verify(ctx context.Context, mirrorDir string) error {
	// Local (bundled) content is trusted by definition.
	return nil
}

func (s *localSource) Size(mirrorDir string) (int64, error) {
	return dirSize(s.cfg.LocalPath)
}

func (s *localSource) Stage(ctx context.Context, mirrorDir, imageDir string, tier config.TierConfig) error {
	dst := filepath.Join(imageDir, s.cfg.StagePath)
	if hasMarkdown(s.cfg.LocalPath) {
		tmplPath := "./content/templates/site/base.html"
		return site.Render(s.cfg.LocalPath, dst, site.Options{TemplatePath: tmplPath, Vars: s.vars})
	}
	return copyDir(s.cfg.LocalPath, dst)
}

// hasMarkdown reports whether dir contains any .md files (recursively).
func hasMarkdown(dir string) bool {
	found := false
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && strings.ToLower(filepath.Ext(path)) == ".md" {
			found = true
			return fs.SkipAll
		}
		return nil
	})
	return found
}
