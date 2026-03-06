package source

import (
	"context"

	"github.com/stazelabs/open-access-kit/internal/config"
)

// DownloadOptions controls download behavior.
type DownloadOptions struct {
	// Force re-downloads files even if they already exist in the mirror.
	Force bool
}

// Source represents a single content source that can be mirrored, verified, and staged.
type Source interface {
	// Name returns the source's config key (e.g. "tor-browser").
	Name() string

	// DetectVersion queries upstream for the latest version string.
	// Returns "" if this source does not have version detection.
	DetectVersion(ctx context.Context) (string, error)

	// Download fetches content into mirrorDir. By default, skips files that
	// already exist; set opts.Force to re-download unconditionally.
	Download(ctx context.Context, mirrorDir string, opts DownloadOptions) error

	// Verify checks GPG signatures and/or checksums of mirrored content.
	Verify(ctx context.Context, mirrorDir string) error

	// Size returns the total size in bytes of this source's mirrored content.
	Size(mirrorDir string) (int64, error)

	// Stage copies files from mirrorDir into the image layout under imageDir.
	Stage(ctx context.Context, mirrorDir, imageDir string, tier config.TierConfig) error
}

// New constructs the appropriate Source implementation from a config entry.
func New(name string, cfg config.SourceConfig) (Source, error) {
	switch cfg.Type {
	case "rsync":
		return newRsync(name, cfg), nil
	case "git":
		return newGit(name, cfg), nil
	case "local":
		return newLocal(name, cfg), nil
	case "http":
		return newHTTP(name, cfg), nil
	case "github-release":
		return newGitHub(name, cfg), nil
	case "kiwix-zim":
		return newKiwix(name, cfg), nil
	case "site-mirror":
		return newSiteMirror(name, cfg), nil
	default:
		return nil, &UnknownTypeError{Name: name, Type: cfg.Type}
	}
}

// UnknownTypeError is returned when a source config has an unrecognized type.
type UnknownTypeError struct {
	Name string
	Type string
}

func (e *UnknownTypeError) Error() string {
	return "unknown source type " + e.Type + " for source " + e.Name
}
