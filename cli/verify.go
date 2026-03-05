package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/source"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [source]",
	Short: "Check GPG signatures and checksums of mirrored content",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := loadConfig()
		if err != nil {
			return err
		}
		mdir := effectiveMirrorDir(cfg)

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

		allOK := true
		for _, s := range sources {
			fmt.Printf("==> Verifying %s\n", s.Name())
			if err := s.Verify(cmd.Context(), mdir); err != nil {
				fmt.Printf("    FAIL: %v\n", err)
				allOK = false
			} else {
				fmt.Printf("    OK\n")
			}
		}
		if !allOK {
			return fmt.Errorf("one or more verification failures")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
