package verify

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
)

// GPG verifies dataFile against its detached sigFile using the given keyring.
// The keyring must be a GPG public keyring file (.gpg) in binary (dearmored) format.
// Uses gpgv, which operates directly on keyring files without the keyboxd
// daemon that modern gpg uses and which ignores --keyring on some systems.
// Verification is strict: any failure is returned as a hard error.
func GPG(ctx context.Context, keyring, sigFile, dataFile string) error {
	absKeyring, err := filepath.Abs(keyring)
	if err != nil {
		return fmt.Errorf("resolving keyring path: %w", err)
	}
	absSig, err := filepath.Abs(sigFile)
	if err != nil {
		return fmt.Errorf("resolving sig path: %w", err)
	}
	absData, err := filepath.Abs(dataFile)
	if err != nil {
		return fmt.Errorf("resolving data path: %w", err)
	}

	cmd := exec.CommandContext(ctx, "gpgv",
		"--keyring", absKeyring,
		absSig, absData,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("gpg verification failed for %s:\n%s", dataFile, out)
	}
	return nil
}

// Checksum verifies that the SHA256 of dataFile matches the expected hash.
func Checksum(expected, dataFile string) error {
	actual, err := sha256File(dataFile)
	if err != nil {
		return fmt.Errorf("computing checksum of %s: %w", dataFile, err)
	}
	if actual != expected {
		return fmt.Errorf("checksum mismatch for %s: got %s, want %s", dataFile, actual, expected)
	}
	return nil
}
