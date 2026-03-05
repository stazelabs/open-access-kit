package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the oak CLI version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("oak %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
