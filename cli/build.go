package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/annotate"
	"github.com/stazelabs/open-access-kit/internal/packaging"
	"github.com/stazelabs/open-access-kit/internal/source"
)

var (
	buildSign         bool
	buildSignKey      string
	buildSkipDownload bool
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Run the full pipeline: download -> verify -> stage -> annotate -> package -> sign",
	Long: `Run the full OAK quarterly build pipeline:

  1. Download   — fetch/mirror all sources for the tier (skips cached files)
  2. Verify     — check GPG signatures of mirrored content
  3. Stage      — copy files from mirror into image/OAK-{release}/
  4. Annotate   — generate VERSION.txt, MANIFEST.txt, README.txt
  5. Package    — zip the image into dist/OAK-{release}-{tier}.zip + .sha256
  6. Sign       — GPG-sign the ZIP (only if --sign is set)

Use --skip-download to start from step 2 when the mirror is already fresh.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}

		tierCfg, ok := cfg.Tiers[tier]
		if !ok {
			return fmt.Errorf("unknown tier %q", tier)
		}

		mdir := effectiveMirrorDir(cfg)
		idir := effectiveReleaseImageDir(cfg, tierCfg.Label)
		odir := effectiveOutputDir(cfg)

		sources, err := sourcesForTier(cfg, tier)
		if err != nil {
			return err
		}

		// 1. Download
		if !buildSkipDownload {
			fmt.Println("==> Step 1/6: Download")
			for _, s := range sources {
				fmt.Printf("    %s\n", s.Name())
				if dryRun {
					continue
				}
				if err := s.Download(cmd.Context(), mdir, source.DownloadOptions{}); err != nil {
					return fmt.Errorf("download %s: %w", s.Name(), err)
				}
			}
		} else {
			fmt.Println("==> Step 1/6: Download (skipped)")
		}

		// 2. Verify
		fmt.Println("==> Step 2/6: Verify")
		for _, s := range sources {
			fmt.Printf("    %s\n", s.Name())
			if dryRun {
				continue
			}
			if err := s.Verify(cmd.Context(), mdir); err != nil {
				return fmt.Errorf("verify %s: %w", s.Name(), err)
			}
		}

		// 3. Stage
		fmt.Printf("==> Step 3/6: Stage -> %s\n", idir)
		for _, s := range sources {
			fmt.Printf("    %s\n", s.Name())
			if dryRun {
				continue
			}
			if err := s.Stage(cmd.Context(), mdir, idir, tierCfg); err != nil {
				return fmt.Errorf("stage %s: %w", s.Name(), err)
			}
		}

		// 4. Annotate
		fmt.Println("==> Step 4/6: Annotate")
		if !dryRun {
			if err := annotate.Run(cmd.Context(), idir, cfg.Release, sources, mdir); err != nil {
				return fmt.Errorf("annotate: %w", err)
			}
		}

		// 5. Package
		zipName := fmt.Sprintf("OAK-%s-%s.zip", cfg.Release, tierCfg.Label)
		fmt.Printf("==> Step 5/6: Package -> %s/%s\n", odir, zipName)
		var zipPath string
		if !dryRun {
			zipPath, err = packaging.Run(cmd.Context(), packaging.Options{
				ImageDir:    idir,
				OutputDir:   odir,
				Release:     cfg.Release,
				TierLabel:   tierCfg.Label,
				ZipRootName: zipRootName(cfg),
			})
			if err != nil {
				return fmt.Errorf("package: %w", err)
			}
		}

		// 6. Sign
		if buildSign {
			fmt.Printf("==> Step 6/6: Sign\n")
			if !dryRun {
				if zipPath == "" {
					zipPath = odir + "/" + zipName
				}
				if err := packaging.Sign(cmd.Context(), zipPath, buildSignKey); err != nil {
					return fmt.Errorf("sign: %w", err)
				}
				fmt.Printf("    %s.asc\n", zipPath)
			}
		} else {
			fmt.Println("==> Step 6/6: Sign (skipped — use --sign to enable)")
		}

		if dryRun {
			fmt.Println("\n[dry-run] no changes made")
		} else {
			fmt.Printf("\nBuild complete: %s/%s\n", odir, zipName)
		}
		return nil
	},
}

func init() {
	buildCmd.Flags().BoolVar(&buildSign, "sign", false, "GPG-sign the output ZIP after packaging")
	buildCmd.Flags().StringVar(&buildSignKey, "sign-key", "", "GPG key ID to use for signing (uses default key if empty)")
	buildCmd.Flags().BoolVar(&buildSkipDownload, "skip-download", false, "skip the download step (use when mirror is already current)")
	rootCmd.AddCommand(buildCmd)
}
