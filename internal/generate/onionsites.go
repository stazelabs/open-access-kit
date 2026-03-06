// Package generate produces derived Markdown content in content/guides/ from
// mirrored upstream sources. Generated files are .gitignored and regenerated
// each release build.
package generate

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

// OnionSitesOptions controls generation of the onion-sites guide page.
type OnionSitesOptions struct {
	// MirrorDir is the root mirror directory (e.g. ./mirror).
	MirrorDir string
	// OutPath is the destination Markdown file (e.g. ./content/guides/onion-sites.md).
	OutPath string
}

// entry holds a single onion site row normalised from either CSV source.
type entry struct {
	Category string
	SiteName string
	OnionURL string
	ProofURL string
	Comment  string
}

// OnionSites reads master.csv and securedrop-api.csv from the mirrored
// real-world-onion-sites repository and writes a curated, attributed
// Markdown page to opts.OutPath.
func OnionSites(opts OnionSitesOptions) error {
	base := filepath.Join(opts.MirrorDir, "onion-sites")

	// Read preamble and footnotes verbatim from the upstream source.
	preamble, err := os.ReadFile(filepath.Join(base, "01-preamble.md"))
	if err != nil {
		return fmt.Errorf("reading preamble: %w", err)
	}
	footnotes, err := os.ReadFile(filepath.Join(base, "02-footnotes.md"))
	if err != nil {
		return fmt.Errorf("reading footnotes: %w", err)
	}

	mainEntries, err := readMasterCSV(filepath.Join(base, "master.csv"))
	if err != nil {
		return fmt.Errorf("reading master.csv: %w", err)
	}
	sdEntries, err := readSecureDropCSV(filepath.Join(base, "securedrop-api.csv"))
	if err != nil {
		return fmt.Errorf("reading securedrop-api.csv: %w", err)
	}

	// Group main entries by category, preserving a stable order.
	categoryOrder, byCategory := groupByCategory(mainEntries)

	if err := os.MkdirAll(filepath.Dir(opts.OutPath), 0o755); err != nil {
		return err
	}
	f, err := os.Create(opts.OutPath)
	if err != nil {
		return err
	}
	defer f.Close()

	return renderOnionSites(f, renderData{
		Generated:     time.Now().UTC().Format("2006-01-02"),
		Preamble:      string(preamble),
		Footnotes:     string(footnotes),
		CategoryOrder: categoryOrder,
		ByCategory:    byCategory,
		SecureDrop:    sdEntries,
	})
}

// readMasterCSV parses master.csv (columns: category,flaky,site_name,onion_url,onion_name,proof_url,comment).
func readMasterCSV(path string) ([]entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.LazyQuotes = true
	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	idx := colIndex(header)

	var entries []entry
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		onionURL := col(row, idx, "onion_url")
		if onionURL == "" {
			continue
		}
		entries = append(entries, entry{
			Category: col(row, idx, "category"),
			SiteName: col(row, idx, "site_name"),
			OnionURL: onionURL,
			ProofURL: col(row, idx, "proof_url"),
			Comment:  col(row, idx, "comment"),
		})
	}
	return entries, nil
}

// readSecureDropCSV parses securedrop-api.csv (columns: flaky,category,site_name,onion_name,onion_url,proof_url,comment).
func readSecureDropCSV(path string) ([]entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	r.LazyQuotes = true
	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	idx := colIndex(header)

	var entries []entry
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		onionURL := col(row, idx, "onion_url")
		if onionURL == "" {
			continue
		}
		entries = append(entries, entry{
			SiteName: col(row, idx, "site_name"),
			OnionURL: onionURL,
			ProofURL: col(row, idx, "proof_url"),
		})
	}
	return entries, nil
}

func groupByCategory(entries []entry) ([]string, map[string][]entry) {
	seen := map[string]bool{}
	var order []string
	byCategory := map[string][]entry{}

	for _, e := range entries {
		cat := e.Category
		if cat == "" {
			cat = "Other"
		}
		if !seen[cat] {
			seen[cat] = true
			order = append(order, cat)
		}
		byCategory[cat] = append(byCategory[cat], e)
	}

	// Sort categories alphabetically, but keep "Other" last.
	sort.Slice(order, func(i, j int) bool {
		if order[i] == "Other" {
			return false
		}
		if order[j] == "Other" {
			return true
		}
		return order[i] < order[j]
	})
	// Sort entries within each category by site name.
	for cat := range byCategory {
		sort.Slice(byCategory[cat], func(i, j int) bool {
			return byCategory[cat][i].SiteName < byCategory[cat][j].SiteName
		})
	}
	return order, byCategory
}

// colIndex builds a column-name → index map from the CSV header row.
func colIndex(header []string) map[string]int {
	m := map[string]int{}
	for i, h := range header {
		m[strings.TrimSpace(h)] = i
	}
	return m
}

func col(row []string, idx map[string]int, name string) string {
	i, ok := idx[name]
	if !ok || i >= len(row) {
		return ""
	}
	return strings.TrimSpace(row[i])
}

type renderData struct {
	Generated     string
	Preamble      string
	Footnotes     string
	CategoryOrder []string
	ByCategory    map[string][]entry
	SecureDrop    []entry
}

var pageTmpl = template.Must(template.New("onion-sites").Funcs(template.FuncMap{
	"mdlink": func(text, url string) string {
		if url == "" {
			return text
		}
		return fmt.Sprintf("[%s](%s)", text, url)
	},
}).Parse(`# Onion Sites Directory

<!-- AUTO-GENERATED FILE — DO NOT EDIT DIRECTLY -->
<!--
  Generated by oak from the real-world-onion-sites project.
  Source: https://github.com/alecmuffett/real-world-onion-sites
  Re-run "oak generate" to refresh.
-->

> **Source:** This page is generated from the
> [real-world-onion-sites](https://github.com/alecmuffett/real-world-onion-sites)
> project by [Alec Muffett](https://github.com/alecmuffett),
> licensed [CC-BY-SA 4.0](https://creativecommons.org/licenses/by-sa/4.0/).
> **Generated: {{.Generated}}. Do not edit this file — it is regenerated each release.**

> **Disclaimer:** While every effort is made to include only legitimate, established
> websites in this directory, Open Access Kit cannot be responsible for the content
> of third-party onion sites. Use your judgement and stay safe.

---

{{.Preamble}}

---
{{range .CategoryOrder}}
## {{.}}

| Site | Onion URL | Notes |
|------|-----------|-------|
{{- $cat := .}}{{range (index $.ByCategory $cat)}}
| {{.SiteName}} | {{mdlink .OnionURL .OnionURL}} | {{.Comment}} |
{{- end}}
{{end}}

## SecureDrop Submissions

*Entries sourced automatically from [securedrop.org/api/v1/directory/](https://securedrop.org/api/v1/directory/).
To add or amend a SecureDrop entry, contact [securedrop.org](https://securedrop.org) directly.*

| Organization | Onion URL | Proof |
|---|---|---|
{{- range .SecureDrop}}
| {{.SiteName}} | {{mdlink .OnionURL .OnionURL}} | {{mdlink "link" .ProofURL}} |
{{- end}}

---

{{.Footnotes}}

---

→ [Onion Sites](index.md) · [Home](../../index.md) · [Resources](../../resources.md)
`))

func renderOnionSites(w io.Writer, data renderData) error {
	return pageTmpl.Execute(w, data)
}
