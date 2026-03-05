package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show mirror state, sizes, and detected upstream versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		mdir := effectiveMirrorDir(cfg)

		fmt.Printf("Release:    %s\n", cfg.Release)
		fmt.Printf("Mirror:     %s\n", mdir)
		fmt.Printf("Image:      %s\n", effectiveImageDir(cfg))
		fmt.Println()

		// Build all sources (not just for one tier)
		for name, scfg := range cfg.Sources {
			s, err := buildSources(cfg, []string{name})
			if err != nil {
				fmt.Printf("  %-22s  error: %v\n", name, err)
				continue
			}
			src := s[0]

			sz, _ := src.Size(mdir)

			var verStr string
			if scfg.VersionDetect.Method != "" {
				ver, err := src.DetectVersion(cmd.Context())
				if err != nil {
					verStr = fmt.Sprintf("(detect error: %v)", err)
				} else {
					verStr = ver
				}
			} else {
				verStr = "N/A"
			}

			cached := "not cached"
			if sz > 0 {
				cached = fmt.Sprintf("%s cached", humanBytes(sz))
			}

			fmt.Printf("  %-22s  [%s]  upstream: %s\n", name, cached, verStr)
			if verbose {
				fmt.Printf("    type: %s  desc: %s\n", scfg.Type, scfg.Description)
			}
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func humanBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
