// Package site renders Markdown content from content/guides/ into HTML.
// It is used in two contexts:
//   - oak build: renders guides into the staged image
//   - oak site: renders guides into docs/ for GitHub Pages
package site

import (
	"bytes"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

// Options controls rendering behaviour.
type Options struct {
	// TemplatePath is the path to base.html. If empty, a built-in minimal
	// template is used (suitable for offline viewing).
	TemplatePath string
	// BaseURL is prepended to absolute hrefs in the rendered output. Leave
	// empty for offline/local rendering.
	BaseURL string
	// ExcludeDirs lists directory names to skip during the walk (e.g. ".git").
	ExcludeDirs []string
}

// Render walks srcDir for Markdown files, converts them to HTML, and writes
// the results into dstDir (created if it does not exist). Non-Markdown files
// are copied as-is. Existing files in dstDir are overwritten.
func Render(srcDir, dstDir string, opts Options) error {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)

	tmpl, err := loadTemplate(opts.TemplatePath)
	if err != nil {
		return err
	}

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		dst := filepath.Join(dstDir, rel)

		if d.IsDir() {
			for _, ex := range opts.ExcludeDirs {
				if d.Name() == ex {
					return fs.SkipDir
				}
			}
			return os.MkdirAll(dst, 0o755)
		}

		if strings.ToLower(filepath.Ext(path)) == ".md" {
			return renderFile(md, tmpl, path, htmlDst(dst), opts.BaseURL,
				rootFileURL(rel, "index.html"), rootFileURL(rel, "license.html"))
		}
		return copyFile(path, dst)
	})
}

// RenderFile renders a single Markdown file at srcPath into dstPath (HTML).
// The destination directory is created if it does not exist.
func RenderFile(srcPath, dstPath string, opts Options) error {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithRendererOptions(html.WithUnsafe()),
	)
	tmpl, err := loadTemplate(opts.TemplatePath)
	if err != nil {
		return err
	}
	return renderFile(md, tmpl, srcPath, dstPath, opts.BaseURL, "index.html", "license.html")
}

// htmlDst replaces the .md extension with .html.
func htmlDst(mdPath string) string {
	return strings.TrimSuffix(mdPath, filepath.Ext(mdPath)) + ".html"
}

// rootFileURL returns the relative path from a rendered file back to a root-level
// HTML file, based on how many directory levels deep the rendered file is.
func rootFileURL(rel, filename string) string {
	dir := filepath.Dir(rel)
	if dir == "." {
		return filename
	}
	depth := len(strings.Split(dir, string(filepath.Separator)))
	return strings.Repeat("../", depth) + filename
}

func renderFile(md goldmark.Markdown, tmpl *template.Template, src, dst, baseURL, homeURL, licenseURL string) error {
	raw, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Rewrite .md links to .html before parsing so goldmark renders them correctly.
	raw = rewriteMdLinks(raw)

	var body bytes.Buffer
	if err := md.Convert(raw, &body); err != nil {
		return err
	}

	title := titleFromMarkdown(raw)

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	f, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, map[string]any{
		"Title":      title,
		"Body":       template.HTML(body.String()),
		"BaseURL":    baseURL,
		"HomeURL":    homeURL,
		"LicenseURL": licenseURL,
	})
}

// rewriteMdLinks replaces href="foo.md" and href='foo.md' with href="foo.html"
// in the raw Markdown source (affects inline HTML) and also handles Markdown
// link syntax [text](foo.md) by replacing the .md extension.
func rewriteMdLinks(src []byte) []byte {
	// Handle Markdown links: ](something.md)
	s := string(src)
	var out strings.Builder
	i := 0
	for i < len(s) {
		idx := strings.Index(s[i:], "](")
		if idx == -1 {
			out.WriteString(s[i:])
			break
		}
		out.WriteString(s[i : i+idx+2])
		i += idx + 2
		// Find closing )
		end := strings.IndexByte(s[i:], ')')
		if end == -1 {
			out.WriteString(s[i:])
			break
		}
		href := s[i : i+end]
		if strings.HasSuffix(href, ".md") && !strings.Contains(href, "://") {
			href = strings.TrimSuffix(href, ".md") + ".html"
		}
		out.WriteString(href)
		out.WriteByte(')')
		i += end + 1
	}
	return []byte(out.String())
}

// titleFromMarkdown extracts the first H1 heading from the Markdown source,
// falling back to "Open Access Kit" if none is found.
func titleFromMarkdown(src []byte) string {
	for _, line := range strings.Split(string(src), "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return "Open Access Kit"
}

func copyFile(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

const builtinTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>{{.Title}}</title>
<style>
body{font-family:Georgia,serif;max-width:760px;margin:2em auto;padding:0 1em;line-height:1.6;color:#222;background:#fff}
h1,h2,h3{font-family:system-ui,sans-serif;line-height:1.2}
a{color:#005b99}
code,pre{background:#f4f4f4;border-radius:3px;padding:.1em .3em;font-size:.9em}
pre{padding:.8em;overflow-x:auto}
header{margin-bottom:2em;font-family:system-ui,sans-serif;font-size:.9em;border-bottom:1px solid #ddd;padding-bottom:.5em}
header a{text-decoration:none;color:#222;font-weight:600}
header a:hover{color:#005b99}
</style>
</head>
<body>
<header><a href="{{.HomeURL}}">🌳 Open Access Kit</a></header>
<article>
{{.Body}}
</article>
<footer style="margin-top:3em;font-size:.8em;color:#666;border-top:1px solid #ddd;padding-top:1em">
Open Access Kit &mdash; Code: <a href="https://www.gnu.org/licenses/gpl-3.0.html">GPL v3</a> &middot; Content: <a href="https://creativecommons.org/licenses/by-sa/4.0/">CC BY-SA 4.0</a> &middot; <a href="{{.LicenseURL}}">License</a>
</footer>
</body>
</html>`

func loadTemplate(path string) (*template.Template, error) {
	if path == "" {
		return template.New("page").Parse(builtinTemplate)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return template.New("page").Parse(string(data))
}
