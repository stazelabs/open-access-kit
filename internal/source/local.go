package source

import (
	"context"

	"github.com/stazelabs/open-access-kit/internal/config"
)

type localSource struct {
	name string
	cfg  config.SourceConfig
}

func newLocal(name string, cfg config.SourceConfig) *localSource {
	return &localSource{name: name, cfg: cfg}
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
	dst := imageDir + "/" + s.cfg.StagePath
	return copyDir(s.cfg.LocalPath, dst)
}
