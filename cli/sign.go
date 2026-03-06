package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/packaging"
)

var signKey string

var signCmd = &cobra.Command{
	Use:   "sign [zipfile]",
	Short: "GPG-sign a packaged ZIP, producing a detached .asc signature",
	Long: `GPG-sign a packaged ZIP file.

If no zipfile argument is given, the path is derived from --config, --tier,
and the release name (e.g. dist/OAK-Q126-M.zip).

Uses the default GPG key unless --sign-key is specified.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var zipPath string

		if len(args) == 1 {
			zipPath = args[0]
		} else {
			cfg, err := loadConfig()
			if err != nil {
				return err
			}
			tierCfg, ok := cfg.Tiers[tier]
			if !ok {
				return fmt.Errorf("unknown tier %q", tier)
			}
			zipName := fmt.Sprintf("OAK-%s-%s.zip", cfg.Release, tierCfg.Label)
			zipPath = filepath.Join(effectiveOutputDir(cfg), zipName)
		}

		if dryRun {
			fmt.Printf("==> [dry-run] would sign %s -> %s.asc\n", zipPath, zipPath)
			return nil
		}

		fmt.Printf("==> Signing %s\n", zipPath)
		if err := packaging.Sign(cmd.Context(), zipPath, signKey); err != nil {
			return err
		}
		fmt.Printf("    %s.asc\n", zipPath)
		return nil
	},
}

func init() {
	signCmd.Flags().StringVar(&signKey, "sign-key", "", "GPG key ID to sign with (uses default key if empty)")
	rootCmd.AddCommand(signCmd)
}
