package cli

import (
	"github.com/spf13/cobra"
)

var verifyCmd = &cobra.Command{
	Use:   "verify [source]",
	Short: "Check GPG signatures and checksums of mirrored content",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implement via internal/verify
		return nil
	},
}

func init() {
	rootCmd.AddCommand(verifyCmd)
}
