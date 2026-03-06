package source

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/stazelabs/open-access-kit/internal/config"
	"github.com/stazelabs/open-access-kit/internal/verify"
)

// fetchBody performs a GET request and returns the response body.
func fetchBody(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}
	return io.ReadAll(resp.Body)
}

const kiwixOPDS = "https://library.kiwix.org/catalog/v2/entries?name="

type kiwixSource struct {
	name string
	cfg  config.SourceConfig
}

func newKiwix(name string, cfg config.SourceConfig) *kiwixSource {
	return &kiwixSource{name: name, cfg: cfg}
}

func (s *kiwixSource) Name() string { return s.name }

func (s *kiwixSource) DetectVersion(_ context.Context) (string, error) {
	return "", nil // version is per-file, detected during Download
}

// zimInfo holds the resolved version and the full download URL base for one ZIM file.
type zimInfo struct {
	version     string // e.g. "2024-08"
	downloadURL string // e.g. "https://download.kiwix.org/zim/zimgit/zimgit-medicine_en_2024-08.zim"
}

// opdsNameCandidates returns OPDS lookup names to try, from most specific to
// least specific, by progressively dropping trailing _xxx segments.
// e.g. "librepathology_en_all_maxi" → ["librepathology_en_all_maxi",
// "librepathology_en_all", "librepathology_en", "librepathology"]
func opdsNameCandidates(name string) []string {
	candidates := []string{name}
	for {
		idx := strings.LastIndex(name, "_")
		if idx < 0 {
			break
		}
		name = name[:idx]
		candidates = append(candidates, name)
	}
	return candidates
}

// resolveRemoteInfo queries the Kiwix OPDS catalog for the latest version and
// full download URL of a ZIM file. If the exact name isn't found, progressively
// shorter name prefixes are tried — the download URL regex still requires the
// full filename so the wrong variant is never selected.
// If download_url is set in the config, OPDS is skipped entirely.
func resolveRemoteInfo(ctx context.Context, zf config.ZimFileConfig) (zimInfo, error) {
	if zf.DownloadURL != "" {
		re := regexp.MustCompile(regexp.QuoteMeta(zf.Name) + `_(\d{4}-\d{2})\.zim`)
		m := re.FindStringSubmatch(zf.DownloadURL)
		if m == nil {
			return zimInfo{}, fmt.Errorf("could not parse YYYY-MM date from download_url %s", zf.DownloadURL)
		}
		return zimInfo{version: m[1], downloadURL: zf.DownloadURL}, nil
	}

	// Regex that matches the exact full filename in any OPDS response body.
	fileRe := regexp.MustCompile(
		`(https://download\.kiwix\.org/zim/\S+/` + regexp.QuoteMeta(zf.Name) + `_(\d{4}-\d{2})\.zim)`,
	)

	for _, candidate := range opdsNameCandidates(zf.Name) {
		body, err := fetchBody(ctx, kiwixOPDS+candidate)
		if err != nil {
			continue
		}
		matches := fileRe.FindAllSubmatch(body, -1)
		if len(matches) == 0 {
			continue
		}
		var bestURL, bestDate string
		for _, m := range matches {
			if date := string(m[2]); date > bestDate {
				bestDate = date
				bestURL = string(m[1])
			}
		}
		return zimInfo{version: bestDate, downloadURL: bestURL}, nil
	}

	return zimInfo{}, fmt.Errorf("OPDS: no download URL found for %s (tried %d name candidates)", zf.Name, len(opdsNameCandidates(zf.Name)))
}

// resolveLocalVersion scans the local mirror directory to find the latest
// YYYY-MM version already downloaded. Used during Verify and Stage.
func resolveLocalVersion(dir, name string) (string, error) {
	re := regexp.MustCompile(`^` + regexp.QuoteMeta(name) + `_(\d{4}-\d{2})\.zim$`)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("reading mirror dir %s: %w", dir, err)
	}
	var best string
	for _, e := range entries {
		m := re.FindStringSubmatch(e.Name())
		if m == nil {
			continue
		}
		if m[1] > best {
			best = m[1]
		}
	}
	if best == "" {
		return "", fmt.Errorf("no local version found for %s in %s", name, dir)
	}
	return best, nil
}

func (s *kiwixSource) Download(ctx context.Context, mirrorDir string, opts DownloadOptions) error {
	dest := filepath.Join(mirrorDir, s.name)
	if err := os.MkdirAll(dest, 0755); err != nil {
		return err
	}
	for _, zf := range s.cfg.ZimFiles {
		info, err := resolveRemoteInfo(ctx, zf)
		if err != nil {
			return err
		}
		filename := zf.Name + "_" + info.version + ".zim"
		dst := filepath.Join(dest, filename)
		sidecar := dst + ".sha256"
		if !opts.Force && fileExists(dst) && fileExists(sidecar) {
			continue
		}
		if err := downloadHTTP(ctx, info.downloadURL, dst); err != nil {
			return fmt.Errorf("downloading %s: %w", filename, err)
		}
		if err := downloadHTTP(ctx, info.downloadURL+".sha256", sidecar); err != nil {
			return fmt.Errorf("downloading %s.sha256: %w", filename, err)
		}
	}
	return nil
}

func (s *kiwixSource) Verify(_ context.Context, mirrorDir string) error {
	dest := filepath.Join(mirrorDir, s.name)
	for _, zf := range s.cfg.ZimFiles {
		ver, err := resolveLocalVersion(dest, zf.Name)
		if err != nil {
			return err
		}
		filename := zf.Name + "_" + ver + ".zim"
		dataFile := filepath.Join(dest, filename)
		sidecar := dataFile + ".sha256"
		if err := verify.SidecarFile(sidecar, dataFile); err != nil {
			return err
		}
	}
	return nil
}

func (s *kiwixSource) Size(mirrorDir string) (int64, error) {
	return dirSize(filepath.Join(mirrorDir, s.name))
}

func (s *kiwixSource) Stage(_ context.Context, mirrorDir, imageDir string, _ config.TierConfig) error {
	dest := filepath.Join(mirrorDir, s.name)
	for _, zf := range s.cfg.ZimFiles {
		ver, err := resolveLocalVersion(dest, zf.Name)
		if err != nil {
			return err
		}
		filename := zf.Name + "_" + ver + ".zim"
		src := filepath.Join(dest, filename)
		dstDir := filepath.Join(imageDir, s.cfg.StagePath, zf.StageSubdir)
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			return err
		}
		if err := copyFile(src, filepath.Join(dstDir, filename)); err != nil {
			return err
		}
	}
	return nil
}
