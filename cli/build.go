package cli

import (
	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Run the full pipeline: download -> verify -> stage -> annotate -> package -> sign",
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO: implement via internal/pipeline
		return nil
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
