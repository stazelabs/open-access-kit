package packaging

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Options configures a packaging run.
type Options struct {
	ImageDir  string // e.g. image/OAK-Q126  (the directory that becomes the ZIP root)
	OutputDir string // e.g. dist/
	Release   string // e.g. Q126
	TierLabel string // e.g. 64GB
}

// Run zips ImageDir into OutputDir/OAK-{release}-{tierLabel}.zip and writes
// a companion .sha256 file. Returns the path to the ZIP.
func Run(_ context.Context, opts Options) (string, error) {
	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return "", fmt.Errorf("creating output dir: %w", err)
	}

	zipName := fmt.Sprintf("OAK-%s-%s.zip", opts.Release, opts.TierLabel)
	zipPath := filepath.Join(opts.OutputDir, zipName)

	if err := createZip(opts.ImageDir, zipPath); err != nil {
		return "", fmt.Errorf("creating zip: %w", err)
	}

	if err := writeSHA256(zipPath); err != nil {
		return "", fmt.Errorf("writing sha256: %w", err)
	}

	return zipPath, nil
}

// Sign creates a detached ASCII-armored GPG signature for zipPath.
// If keyID is non-empty it selects that key; otherwise the default key is used.
func Sign(_ context.Context, zipPath, keyID string) error {
	sigPath := zipPath + ".asc"
	args := []string{"--detach-sign", "--armor", "--output", sigPath}
	if keyID != "" {
		args = append(args, "--local-user", keyID)
	}
	args = append(args, zipPath)

	cmd := exec.Command("gpg", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gpg signing failed:\n%s", out)
	}
	return nil
}

// createZip walks imageDir and writes all files into zipPath.
// Each entry is stored under a top-level directory named after the imageDir basename,
// so the ZIP extracts to OAK-Q126/<files>.
func createZip(imageDir, zipPath string) error {
	f, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := zip.NewWriter(f)
	defer w.Close()

	// The ZIP root dir name matches the image directory name (e.g. OAK-Q126)
	rootName := filepath.Base(imageDir)

	return filepath.Walk(imageDir, func(path string, info os.FileInfo, err error) error {
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
		// Use forward slashes inside the ZIP regardless of host OS
		entry := rootName + "/" + filepath.ToSlash(rel)

		fw, err := w.Create(entry)
		if err != nil {
			return fmt.Errorf("creating zip entry %s: %w", entry, err)
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		_, err = io.Copy(fw, src)
		return err
	})
}

// writeSHA256 computes the SHA256 of zipPath and writes it to zipPath+".sha256".
func writeSHA256(zipPath string) error {
	f, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	hash := hex.EncodeToString(h.Sum(nil))

	shaPath := zipPath + ".sha256"
	return os.WriteFile(shaPath, []byte(hash+"  "+filepath.Base(zipPath)+"\n"), 0644)
}
