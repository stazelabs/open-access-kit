package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/packaging"
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Zip the staged image into dist/OAK-{release}-{tier}.zip with a .sha256",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		tierCfg, ok := cfg.Tiers[tier]
		if !ok {
			return fmt.Errorf("unknown tier %q", tier)
		}

		idir := effectiveReleaseImageDir(cfg, tierCfg.Label)
		odir := effectiveOutputDir(cfg)

		opts := packaging.Options{
			ImageDir:    idir,
			OutputDir:   odir,
			Release:     cfg.Release,
			TierLabel:   tierCfg.Label,
			ZipRootName: zipRootName(cfg),
		}

		if dryRun {
			zipName := fmt.Sprintf("OAK-%s-%s.zip", cfg.Release, tierCfg.Label)
			fmt.Printf("==> [dry-run] would package %s -> %s/%s\n", idir, odir, zipName)
			return nil
		}

		fmt.Printf("==> Packaging %s\n", idir)
		zipPath, err := packaging.Run(cmd.Context(), opts)
		if err != nil {
			return err
		}
		fmt.Printf("    %s\n", zipPath)
		fmt.Printf("    %s.sha256\n", zipPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(packageCmd)
}
