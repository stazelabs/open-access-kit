package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/config"
	"github.com/stazelabs/open-access-kit/internal/manifest"
	"github.com/stazelabs/open-access-kit/internal/source"
)

var manifestOutput string

var manifestCmd = &cobra.Command{
	Use:   "manifest",
	Short: "Generate a release manifest (releases/{release}.json)",
	Long: `Generate a machine-readable JSON manifest capturing the release name,
build timestamp, tier composition, and detected version or commit for
every source. The manifest is saved to releases/{release}.json by default.

Use "oak diff" to compare two manifests and see what changed.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		// Build Source objects for all defined sources (not just one tier).
		allSources, err := buildAllSources(cfg)
		if err != nil {
			return err
		}

		mdir := effectiveMirrorDir(cfg)

		fmt.Println("Generating release manifest...")
		m, err := manifest.Generate(cmd.Context(), cfg, allSources, mdir)
		if err != nil {
			return fmt.Errorf("generating manifest: %w", err)
		}

		outPath := manifestOutput
		if outPath == "" {
			outPath = filepath.Join("releases", cfg.Release+".json")
		}

		if dryRun {
			fmt.Printf("[dry-run] would write manifest to %s\n", outPath)
			return nil
		}

		if err := manifest.Write(m, outPath); err != nil {
			return err
		}
		fmt.Printf("Manifest written to %s\n", outPath)
		return nil
	},
}

// buildAllSources constructs Source objects for every source in the config.
func buildAllSources(cfg *config.Config) ([]source.Source, error) {
	names := make([]string, 0, len(cfg.Sources))
	for name := range cfg.Sources {
		names = append(names, name)
	}
	return buildSources(cfg, names)
}

func init() {
	manifestCmd.Flags().StringVar(&manifestOutput, "output", "", "output path (default: releases/{release}.json)")
	rootCmd.AddCommand(manifestCmd)
}
