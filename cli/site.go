package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/stazelabs/open-access-kit/internal/site"
)

// rootDocs are repo-root Markdown files copied into docs/ as a publishing step.
// They are kept at the repo root so GitHub can find them, but also rendered
// into the site so they appear on the companion website.
var rootDocs = []string{
	"CONTRIBUTING.md",
	"ARCHITECTURE.md",
	"LICENSE-CODE",
	"LICENSE-CONTENT",
}

var (
	siteBaseURL      string
	siteTemplatePath string
	siteSrcDir       string
	siteDstDir       string
)

var siteCmd = &cobra.Command{
	Use:   "site",
	Short: "Render content/guides/ into docs/ for GitHub Pages",
	Long: `Render Markdown guides from content/guides/ (or --src) into HTML at
docs/ (or --dst), ready for GitHub Pages publishing.

The same content is also rendered into each tier's image during oak build.
This command lets you preview or publish the web version independently.

Example:
  oak site                         # render content/guides/ -> docs/
  oak site --base-url /open-access-kit  # for a GitHub Pages project sub-path`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if dryRun {
			fmt.Printf("[dry-run] would render %s -> %s\n", siteSrcDir, siteDstDir)
			return nil
		}

		fmt.Printf("==> Rendering %s -> %s\n", siteSrcDir, siteDstDir)

		opts := site.Options{
			TemplatePath: siteTemplatePath,
			BaseURL:      siteBaseURL,
		}
		if err := site.Render(siteSrcDir, siteDstDir, opts); err != nil {
			return fmt.Errorf("site render: %w", err)
		}

		// Render repo-root docs (CONTRIBUTING.md, ARCHITECTURE.md) into docs/
		// so they appear on the companion website without moving the canonical files.
		for _, f := range rootDocs {
			dst := filepath.Join(siteDstDir, strings.ToLower(strings.TrimSuffix(f, ".md"))+".html")
			fmt.Printf("    publishing %s -> %s\n", f, dst)
			if err := site.RenderFile(f, dst, opts); err != nil {
				return fmt.Errorf("rendering %s: %w", f, err)
			}
		}

		fmt.Printf("    done\n")
		return nil
	},
}

func init() {
	siteCmd.Flags().StringVar(&siteBaseURL, "base-url", "", "base URL prefix for web links (e.g. /open-access-kit)")
	siteCmd.Flags().StringVar(&siteTemplatePath, "template", "./content/templates/site/base.html", "HTML template file (uses built-in if file is absent)")
	siteCmd.Flags().StringVar(&siteSrcDir, "src", "./content/guides", "source directory of Markdown guides")
	siteCmd.Flags().StringVar(&siteDstDir, "dst", "./docs", "destination directory for rendered HTML")
	rootCmd.AddCommand(siteCmd)
}
