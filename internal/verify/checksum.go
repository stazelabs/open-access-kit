package verify

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"
)

// SidecarFile reads a .sha256 sidecar file (sha256sum format: "{hash}  {filename}")
// and verifies that the SHA-256 of dataFile matches the expected hash.
func SidecarFile(sidecarPath, dataFile string) error {
	content, err := os.ReadFile(sidecarPath)
	if err != nil {
		return fmt.Errorf("reading sha256 sidecar %s: %w", sidecarPath, err)
	}
	fields := strings.Fields(string(content))
	if len(fields) < 1 {
		return fmt.Errorf("invalid sha256 sidecar format in %s", sidecarPath)
	}
	return Checksum(fields[0], dataFile)
}

func sha256File(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("opening file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hashing file: %w", err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
