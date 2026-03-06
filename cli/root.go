package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "oak",
	Short: "Open Access Kit — build offline privacy tool images",
	Long: `OAK curates, downloads, verifies, stages, and packages quarterly
releases of privacy tools and educational content onto removable media.`,
}

// Global flags
var (
	cfgFile   string
	tier      string
	mirrorDir string
	imageDir  string
	release   string
	dryRun    bool
	verbose   bool
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "./oak.yaml", "path to configuration file")
	rootCmd.PersistentFlags().StringVarP(&tier, "tier", "t", "64", "target tier: 16, 32, 64, max")
	rootCmd.PersistentFlags().StringVar(&mirrorDir, "mirror-dir", "./mirror", "path to mirror directory")
	rootCmd.PersistentFlags().StringVar(&imageDir, "image-dir", "./image", "path to image output directory")
	rootCmd.PersistentFlags().StringVar(&release, "release", "", "release name override (e.g. Q126)")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false, "show what would happen without doing it")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
