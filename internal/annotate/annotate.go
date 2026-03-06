package annotate

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/stazelabs/open-access-kit/internal/source"
)

// Run generates VERSION.txt, MANIFEST.txt, and README.txt into imageDir.
// It must be called after staging is complete.
func Run(ctx context.Context, imageDir, release string, sources []source.Source, mirrorDir string) error {
	if err := os.MkdirAll(imageDir, 0755); err != nil {
		return fmt.Errorf("creating image dir: %w", err)
	}

	if err := writeVersion(ctx, imageDir, release, sources, mirrorDir); err != nil {
		return fmt.Errorf("writing VERSION.txt: %w", err)
	}
	if err := writeReadme(imageDir, release); err != nil {
		return fmt.Errorf("writing README.txt: %w", err)
	}
	// MANIFEST must be last — it hashes everything including VERSION.txt and README.txt
	if err := writeManifest(imageDir); err != nil {
		return fmt.Errorf("writing MANIFEST.txt: %w", err)
	}
	return nil
}

func writeVersion(ctx context.Context, imageDir, release string, sources []source.Source, mirrorDir string) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "Release: %s\n", release)
	fmt.Fprintf(&sb, "Built:   %s\n", time.Now().UTC().Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(&sb, "\nSources:\n")

	for _, s := range sources {
		ver, err := s.DetectVersion(ctx)
		if err != nil || ver == "" {
			ver = "(unknown)"
		}
		fmt.Fprintf(&sb, "  %-24s %s\n", s.Name(), ver)
	}

	return os.WriteFile(filepath.Join(imageDir, "VERSION.txt"), []byte(sb.String()), 0644)
}

func writeReadme(imageDir, release string) error {
	content := fmt.Sprintf(`Open Access Kit — %s
=====================

This drive contains privacy tools and educational resources for people
facing censorship or surveillance.

START HERE
----------
  Open guides/index.html in any web browser for full documentation,
  getting started instructions, and a directory of onion sites.

CONTENTS
--------
  guides/                 Offline HTML documentation (start here)
  software/tor-browser/   Tor Browser for Windows, macOS, Linux, Android
  software/tails/         Tails OS bootable image (32GB+ drives only)
  guides/resources/       Additional bundled resources and onion site directories
  keys/                   GPG public keys used to verify this software

QUICK START
-----------
  1. Open guides/index.html in any browser
  2. Install Tor Browser from software/tor-browser/
  3. Launch Tor Browser and click Connect

VERIFYING THIS DRIVE
--------------------
  sha256sum -c MANIFEST.txt       (verify file integrity)
  gpg --verify OAK-*.zip.asc     (verify package signature, if present)

More information: https://github.com/stazelabs/open-access-kit
`, release)

	return os.WriteFile(filepath.Join(imageDir, "README.txt"), []byte(content), 0644)
}

func writeManifest(imageDir string) error {
	var sb strings.Builder

	err := filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(imageDir, path)
		if err != nil {
			return err
		}
		// Skip MANIFEST.txt itself — can't hash a file while writing it
		if rel == "MANIFEST.txt" {
			return nil
		}

		h, err := sha256File(path)
		if err != nil {
			return fmt.Errorf("hashing %s: %w", rel, err)
		}
		// sha256sum-compatible format: hash + two spaces + path
		fmt.Fprintf(&sb, "%s  %s\n", h, rel)
		return nil
	})
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(imageDir, "MANIFEST.txt"), []byte(sb.String()), 0644)
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
