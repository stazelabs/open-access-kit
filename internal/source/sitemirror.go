package source

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/stazelabs/open-access-kit/internal/config"
)

type siteMirrorSource struct {
	name string
	cfg  config.SourceConfig
}

func newSiteMirror(name string, cfg config.SourceConfig) *siteMirrorSource {
	return &siteMirrorSource{name: name, cfg: cfg}
}

func (s *siteMirrorSource) Name() string { return s.name }

func (s *siteMirrorSource) DetectVersion(ctx context.Context) (string, error) {
	// Site mirrors are live snapshots; no discrete version.
	return "", nil
}

func (s *siteMirrorSource) Download(ctx context.Context, mirrorDir string, opts DownloadOptions) error {
	dest := filepath.Join(mirrorDir, s.name)

	if !opts.Force {
		if _, err := os.Stat(dest); err == nil {
			// Already mirrored — skip
			return nil
		}
	}

	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	// --mirror            = -r -N -l inf --no-remove-listing
	// --convert-links     = rewrite links for offline use
	// --page-requisites   = fetch CSS, images, etc. needed to display each page
	// --no-parent         = don't crawl above the given URL path
	// --no-host-directories = don't create a hostname subdirectory
	// -P                  = output prefix directory
	cmd := exec.CommandContext(ctx, "wget",
		"--mirror",
		"--convert-links",
		"--page-requisites",
		"--no-parent",
		"--quiet",
		"--no-host-directories",
		"-P", dest,
		s.cfg.MirrorURL,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("wget mirror of %s: %w", s.cfg.MirrorURL, err)
	}
	return nil
}

func (s *siteMirrorSource) Verify(ctx context.Context, mirrorDir string) error {
	// HTTPS transport provides integrity; no additional verification needed.
	return nil
}

func (s *siteMirrorSource) Size(mirrorDir string) (int64, error) {
	return dirSize(filepath.Join(mirrorDir, s.name))
}

func (s *siteMirrorSource) Stage(ctx context.Context, mirrorDir, imageDir string, tier config.TierConfig) error {
	if s.cfg.StagePath == "" {
		return nil
	}
	src := filepath.Join(mirrorDir, s.name)
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("mirror not found at %s — run \"oak download\" first", src)
	}
	return copyDir(src, filepath.Join(imageDir, s.cfg.StagePath))
}
