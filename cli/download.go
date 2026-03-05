package cli

import (
	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download [source]",
	Short: "Fetch and mirror all sources, or a specific source by name",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implement via internal/source
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
