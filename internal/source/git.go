package source

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/stazelabs/open-access-kit/internal/config"
)

type gitSource struct {
	name string
	cfg  config.SourceConfig
}

func newGit(name string, cfg config.SourceConfig) *gitSource {
	return &gitSource{name: name, cfg: cfg}
}

func (s *gitSource) Name() string { return s.name }

func (s *gitSource) DetectVersion(ctx context.Context) (string, error) {
	// Git sources track HEAD; no discrete version.
	return "", nil
}

func (s *gitSource) Download(ctx context.Context, mirrorDir string, opts DownloadOptions) error {
	dest := filepath.Join(mirrorDir, s.name)
	gitDir := filepath.Join(dest, ".git")

	if !opts.Force {
		if _, err := os.Stat(gitDir); err == nil {
			// Already cloned — skip
			return nil
		}
	}

	if _, err := os.Stat(gitDir); err == nil {
		// Already cloned — pull
		cmd := exec.CommandContext(ctx, "git", "-C", dest, "pull", "--ff-only", "--quiet")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}

	// Fresh clone
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	args := []string{"clone", "--quiet"}
	if s.cfg.Shallow {
		args = append(args, "--depth=1")
	}
	args = append(args, s.cfg.GitURL, dest)

	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (s *gitSource) Verify(ctx context.Context, mirrorDir string) error {
	// Git transport uses HTTPS; content integrity is enforced by the protocol.
	return nil
}

func (s *gitSource) Size(mirrorDir string) (int64, error) {
	return dirSize(filepath.Join(mirrorDir, s.name))
}

func (s *gitSource) Stage(ctx context.Context, mirrorDir, imageDir string, tier config.TierConfig) error {
	src := filepath.Join(mirrorDir, s.name)
	dst := filepath.Join(imageDir, s.cfg.StagePath)
	// Exclude .git metadata from the staged image
	return copyDirExclude(src, dst, ".git")
}

// copyDirExclude copies src to dst, skipping any entry named excludeDir.
func copyDirExclude(src, dst, excludeDir string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		// Skip the excluded directory and all its contents
		if info.IsDir() && info.Name() == excludeDir {
			return filepath.SkipDir
		}
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}
