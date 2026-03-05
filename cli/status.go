package cli

import (
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show mirror state, sizes, and detected upstream versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implement via internal/source + internal/version
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
