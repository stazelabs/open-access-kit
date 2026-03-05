package source

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/stazelabs/open-access-kit/internal/config"
	"github.com/stazelabs/open-access-kit/internal/verify"
	"github.com/stazelabs/open-access-kit/internal/version"
)

type rsyncSource struct {
	name string
	cfg  config.SourceConfig
}

func newRsync(name string, cfg config.SourceConfig) *rsyncSource {
	return &rsyncSource{name: name, cfg: cfg}
}

func (s *rsyncSource) Name() string { return s.name }

func (s *rsyncSource) DetectVersion(ctx context.Context) (string, error) {
	vd := s.cfg.VersionDetect
	if vd.Method == "" {
		return "", nil
	}
	switch vd.Method {
	case "http-scrape":
		return version.HTTPScrape(ctx, vd.URL, vd.Pattern, vd.Select)
	case "rsync-list":
		return version.RsyncList(ctx, vd.URL, vd.Pattern, vd.Select)
	default:
		return "", fmt.Errorf("unknown version_detect method: %s", vd.Method)
	}
}

func (s *rsyncSource) Download(ctx context.Context, mirrorDir string, opts DownloadOptions) error {
	dest := filepath.Join(mirrorDir, s.name)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("creating mirror dir: %w", err)
	}

	if len(s.cfg.Files) > 0 {
		// Versioned file list download
		ver, err := s.DetectVersion(ctx)
		if err != nil {
			return fmt.Errorf("version detection: %w", err)
		}
		base := strings.ReplaceAll(s.cfg.RsyncBase, "{version}", ver)

		for _, f := range s.cfg.Files {
			filename := strings.ReplaceAll(f, "{version}", ver)
			dst := filepath.Join(dest, filename)

			if !opts.Force && fileExists(dst) {
				continue
			}

			src := base + filename
			if err := rsyncFile(ctx, src, dst); err != nil {
				return fmt.Errorf("downloading %s: %w", filename, err)
			}
		}
	} else {
		// Whole-directory sync (e.g. Tails)
		base := s.cfg.RsyncBase
		if !strings.HasSuffix(base, "/") {
			base += "/"
		}
		if !opts.Force {
			// For directory syncs, only skip if mirror already has content
			info, err := os.Stat(dest)
			if err == nil && info.IsDir() {
				entries, _ := os.ReadDir(dest)
				if len(entries) > 0 {
					return nil
				}
			}
		}
		if err := rsyncDir(ctx, base, dest+"/"); err != nil {
			return err
		}
	}
	return nil
}

func (s *rsyncSource) Verify(ctx context.Context, mirrorDir string) error {
	if s.cfg.Verify.Method != "gpg" {
		return nil
	}
	dest := filepath.Join(mirrorDir, s.name)

	if len(s.cfg.Files) > 0 {
		ver, err := s.DetectVersion(ctx)
		if err != nil {
			return err
		}
		return s.verifyFileList(ctx, dest, ver)
	}
	return s.verifyDir(ctx, dest)
}

// verifyFileList verifies GPG signatures for a known list of files.
func (s *rsyncSource) verifyFileList(ctx context.Context, dest, ver string) error {
	for _, f := range s.cfg.Files {
		filename := strings.ReplaceAll(f, "{version}", ver)
		// Skip signature files themselves
		if strings.HasSuffix(filename, ".asc") || strings.HasSuffix(filename, ".sig") {
			continue
		}
		dataFile := filepath.Join(dest, filename)
		if !fileExists(dataFile) {
			return fmt.Errorf("file missing from mirror: %s", dataFile)
		}
		sigFile, err := findSig(dest, filename)
		if err != nil {
			return err
		}
		if err := verify.GPG(ctx, s.cfg.Verify.Keyring, sigFile, dataFile); err != nil {
			return err
		}
	}
	return nil
}

// verifyDir walks the mirror directory verifying all signable files.
func (s *rsyncSource) verifyDir(ctx context.Context, dest string) error {
	return filepath.Walk(dest, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		ext := filepath.Ext(path)
		if ext == ".asc" || ext == ".sig" {
			return nil // skip signatures themselves
		}
		sigFile, err := findSig(filepath.Dir(path), filepath.Base(path))
		if err != nil {
			return nil // no sig found, skip silently
		}
		return verify.GPG(ctx, s.cfg.Verify.Keyring, sigFile, path)
	})
}

// findSig looks for a .asc or .sig file next to filename in dir.
func findSig(dir, filename string) (string, error) {
	for _, ext := range []string{".asc", ".sig"} {
		candidate := filepath.Join(dir, filename+ext)
		if fileExists(candidate) {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("no signature file found for %s", filename)
}

func (s *rsyncSource) Size(mirrorDir string) (int64, error) {
	return dirSize(filepath.Join(mirrorDir, s.name))
}

func (s *rsyncSource) Stage(ctx context.Context, mirrorDir, imageDir string, tier config.TierConfig) error {
	src := filepath.Join(mirrorDir, s.name)
	dst := filepath.Join(imageDir, s.cfg.StagePath)
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	if len(s.cfg.Files) > 0 {
		ver, err := s.DetectVersion(ctx)
		if err != nil {
			return err
		}
		for _, f := range s.cfg.Files {
			filename := strings.ReplaceAll(f, "{version}", ver)
			if err := copyFile(filepath.Join(src, filename), filepath.Join(dst, filename)); err != nil {
				return fmt.Errorf("staging %s: %w", filename, err)
			}
		}
		return nil
	}
	return copyDir(src, dst)
}

func rsyncFile(ctx context.Context, src, dst string) error {
	cmd := exec.CommandContext(ctx, "rsync",
		"--archive",
		"--compress",
		"--progress",
		src, dst,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func rsyncDir(ctx context.Context, src, dst string) error {
	cmd := exec.CommandContext(ctx, "rsync",
		"--archive",
		"--compress",
		"--progress",
		"--delete",
		src, dst,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
