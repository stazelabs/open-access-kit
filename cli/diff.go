package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/manifest"
)

var diffFormat string

var diffCmd = &cobra.Command{
	Use:   "diff <release1> <release2>",
	Short: "Compare two release manifests and show what changed",
	Long: `Compare two release manifests (from releases/) and output a structured
summary of what changed: version bumps, added/removed sources, and tier
composition changes.

Examples:
  oak diff Q125 Q226
  oak diff Q125 Q226 --format text`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldPath := filepath.Join("releases", args[0]+".json")
		newPath := filepath.Join("releases", args[1]+".json")

		old, err := manifest.Read(oldPath)
		if err != nil {
			return fmt.Errorf("reading %s manifest: %w", args[0], err)
		}
		new, err := manifest.Read(newPath)
		if err != nil {
			return fmt.Errorf("reading %s manifest: %w", args[1], err)
		}

		d := manifest.Compare(old, new)

		switch diffFormat {
		case "text":
			fmt.Print(d.Text())
		default:
			data, err := d.JSON()
			if err != nil {
				return fmt.Errorf("formatting diff: %w", err)
			}
			fmt.Println(string(data))
		}
		return nil
	},
}

func init() {
	diffCmd.Flags().StringVar(&diffFormat, "format", "json", "output format: json or text")
	rootCmd.AddCommand(diffCmd)
}
