package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stageCmd = &cobra.Command{
	Use:   "stage",
	Short: "Build image directory from mirror, applying tier size budget",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		mdir := effectiveMirrorDir(cfg)
		idir := effectiveImageDir(cfg)

		sources, err := sourcesForTier(cfg, tier)
		if err != nil {
			return err
		}

		tierCfg, ok := cfg.Tiers[tier]
		if !ok {
			return fmt.Errorf("unknown tier %q", tier)
		}

		for _, s := range sources {
			fmt.Printf("==> Staging %s\n", s.Name())
			if dryRun {
				sz, _ := s.Size(mdir)
				fmt.Printf("    [dry-run] would stage %s (%s) -> %s\n", s.Name(), humanBytes(sz), idir)
				continue
			}
			if err := s.Stage(cmd.Context(), mdir, idir, tierCfg); err != nil {
				return fmt.Errorf("staging %s: %w", s.Name(), err)
			}
		}

		// Check size budget if set
		if !dryRun && tierCfg.BudgetGB > 0 {
			var totalBytes int64
			for _, s := range sources {
				sz, err := s.Size(mdir)
				if err == nil {
					totalBytes += sz
				}
			}
			budgetBytes := int64(tierCfg.BudgetGB * 1024 * 1024 * 1024)
			if totalBytes > budgetBytes {
				return fmt.Errorf("tier %s budget exceeded: %s > %.1f GB",
					tier, humanBytes(totalBytes), tierCfg.BudgetGB)
			}
			fmt.Printf("Size check: %s / %.1f GB budget\n", humanBytes(totalBytes), tierCfg.BudgetGB)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(stageCmd)
}
