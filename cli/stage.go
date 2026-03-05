package cli

import (
	"github.com/spf13/cobra"
)

var stageCmd = &cobra.Command{
	Use:   "stage",
	Short: "Build image directory from mirror, applying tier size budget",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implement via internal/stage
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stageCmd)
}
