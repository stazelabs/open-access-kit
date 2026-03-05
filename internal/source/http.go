package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/stazelabs/open-access-kit/internal/config"
	"github.com/stazelabs/open-access-kit/internal/verify"
	"github.com/stazelabs/open-access-kit/internal/version"
)

type httpSource struct {
	name string
	cfg  config.SourceConfig
}

func newHTTP(name string, cfg config.SourceConfig) *httpSource {
	return &httpSource{name: name, cfg: cfg}
}

func (s *httpSource) Name() string { return s.name }

func (s *httpSource) DetectVersion(ctx context.Context) (string, error) {
	vd := s.cfg.VersionDetect
	if vd.Method == "" {
		return "", nil
	}
	switch vd.Method {
	case "http-scrape":
		return version.HTTPScrape(ctx, vd.URL, vd.Pattern, vd.Select)
	default:
		return "", fmt.Errorf("unknown version_detect method: %s", vd.Method)
	}
}

func (s *httpSource) Download(ctx context.Context, mirrorDir string, opts DownloadOptions) error {
	dest := filepath.Join(mirrorDir, s.name)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}

	ver, err := s.DetectVersion(ctx)
	if err != nil {
		return fmt.Errorf("version detection: %w", err)
	}

	for _, f := range s.cfg.Files {
		filename := strings.ReplaceAll(f, "{version}", ver)
		dst := filepath.Join(dest, filename)
		if !opts.Force && fileExists(dst) {
			continue
		}
		url := strings.ReplaceAll(s.cfg.RsyncBase, "{version}", ver) + filename
		if err := downloadHTTP(ctx, url, dst); err != nil {
			return fmt.Errorf("downloading %s: %w", filename, err)
		}
	}
	return nil
}

func (s *httpSource) Verify(ctx context.Context, mirrorDir string) error {
	if s.cfg.Verify.Method != "gpg" {
		return nil
	}
	dest := filepath.Join(mirrorDir, s.name)
	ver, err := s.DetectVersion(ctx)
	if err != nil {
		return err
	}
	for _, f := range s.cfg.Files {
		filename := strings.ReplaceAll(f, "{version}", ver)
		if strings.HasSuffix(filename, ".asc") || strings.HasSuffix(filename, ".sig") {
			continue
		}
		dataFile := filepath.Join(dest, filename)
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

func (s *httpSource) Size(mirrorDir string) (int64, error) {
	return dirSize(filepath.Join(mirrorDir, s.name))
}

func (s *httpSource) Stage(ctx context.Context, mirrorDir, imageDir string, tier config.TierConfig) error {
	src := filepath.Join(mirrorDir, s.name)
	dst := filepath.Join(imageDir, s.cfg.StagePath)
	return copyDir(src, dst)
}

func downloadHTTP(ctx context.Context, url, dst string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}
	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}
