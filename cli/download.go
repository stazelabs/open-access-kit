package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/source"
)

var downloadForce bool

var downloadCmd = &cobra.Command{
	Use:   "download [source]",
	Short: "Fetch and mirror all sources, or a specific source by name",
	Long: `Downloads content into the mirror directory.

By default, files that already exist in the mirror are skipped. Use --force
to re-download everything regardless.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		mdir := effectiveMirrorDir(cfg)
		opts := source.DownloadOptions{Force: downloadForce}

		var sources []source.Source
		if len(args) == 1 {
			s, err := sourceByName(cfg, args[0])
			if err != nil {
				return err
			}
			sources = []source.Source{s}
		} else {
			sources, err = sourcesForTier(cfg, tier)
			if err != nil {
				return err
			}
		}

		for _, s := range sources {
			fmt.Printf("==> Downloading %s\n", s.Name())
			if dryRun {
				fmt.Printf("    [dry-run] would download to %s/%s\n", mdir, s.Name())
				continue
			}
			if err := s.Download(cmd.Context(), mdir, opts); err != nil {
				return fmt.Errorf("download %s: %w", s.Name(), err)
			}
		}
		return nil
	},
}

func init() {
	downloadCmd.Flags().BoolVar(&downloadForce, "force", false, "re-download files even if they already exist in the mirror")
	rootCmd.AddCommand(downloadCmd)
}
