package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cleanMirror bool

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove staged images and packaged output (preserves mirror and keys)",
	Long: `Remove build artifacts that can be safely regenerated:

  image/   staged tier images produced by "oak stage"
  dist/    packaged ZIPs, checksums, and signatures produced by "oak package"/"oak sign"

The mirror (mirror/) and signing keys (keys/) are never removed.
Use --mirror to also wipe the mirror when reclaiming disk space; note that
the next build will require a full re-download of all mirrored sources.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		toRemove := []string{
			effectiveImageDir(cfg),
			effectiveOutputDir(cfg),
		}
		if cleanMirror {
			toRemove = append(toRemove, effectiveMirrorDir(cfg))
		}

		for _, dir := range toRemove {
			fmt.Printf("    rm -rf %s\n", dir)
			if !dryRun {
				if err := os.RemoveAll(dir); err != nil {
					return err
				}
			}
		}

		if dryRun {
			fmt.Println("\n[dry-run] no changes made")
		}
		return nil
	},
}

func init() {
	cleanCmd.Flags().BoolVar(&cleanMirror, "mirror", false,
		"also remove the mirror directory (requires full re-download on next build)")
	rootCmd.AddCommand(cleanCmd)
}
