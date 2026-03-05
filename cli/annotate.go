package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/annotate"
)

var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Generate VERSION.txt, MANIFEST.txt, and README.txt into the staged image",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		mdir := effectiveMirrorDir(cfg)
		idir := effectiveReleaseImageDir(cfg)

		sources, err := sourcesForTier(cfg, tier)
		if err != nil {
			return err
		}

		if dryRun {
			fmt.Printf("==> [dry-run] would annotate %s\n", idir)
			fmt.Printf("    VERSION.txt  — release %s, build date, source versions\n", cfg.Release)
			fmt.Printf("    README.txt   — plain-text intro\n")
			fmt.Printf("    MANIFEST.txt — SHA256 of all files in %s\n", idir)
			return nil
		}

		fmt.Printf("==> Annotating %s\n", idir)
		return annotate.Run(cmd.Context(), idir, cfg.Release, sources, mdir)
	},
}

func init() {
	rootCmd.AddCommand(annotateCmd)
}
