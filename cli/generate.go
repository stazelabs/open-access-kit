package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/generate"
)

var generateGuidesDir string

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate derived guide pages from mirrored upstream sources",
	Long: `Generate Markdown pages in content/guides/ from mirrored upstream data.

Generated files are .gitignored and regenerated each release. Run this after
"oak download" and before "oak stage" (or use "oak build" which runs it automatically).

Currently generates:
  resources/onion-sites/real-world-onion-sites.md  — from real-world-onion-sites master.csv and securedrop-api.csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		mdir := effectiveMirrorDir(cfg)
		out := filepath.Join(generateGuidesDir, "resources", "onion-sites", "real-world-onion-sites.md")

		fmt.Printf("==> Generating onion-sites -> %s\n", out)
		if dryRun {
			fmt.Printf("    [dry-run] would write %s\n", out)
			return nil
		}

		if err := generate.OnionSites(generate.OnionSitesOptions{
			MirrorDir: mdir,
			OutPath:   out,
		}); err != nil {
			return fmt.Errorf("generate onion-sites: %w", err)
		}

		fmt.Printf("    done\n")
		return nil
	},
}

func init() {
	generateCmd.Flags().StringVar(&generateGuidesDir, "guides-dir", "./content/guides", "destination directory for generated guide pages")
	rootCmd.AddCommand(generateCmd)
}
