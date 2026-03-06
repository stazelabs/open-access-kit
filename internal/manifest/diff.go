package manifest

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/stazelabs/open-access-kit/internal/source"
)

// Diff describes the differences between two release manifests.
type Diff struct {
	From    string      `json:"from"`
	To      string      `json:"to"`
	Sources SourcesDiff `json:"sources"`
	Tiers   TiersDiff   `json:"tiers"`
}

// SourcesDiff categorises source-level changes between releases.
type SourcesDiff struct {
	Added     []AddedSource  `json:"added"`
	Removed   []AddedSource  `json:"removed"`
	Updated   []SourceUpdate `json:"updated"`
	Unchanged []string       `json:"unchanged"`
}

// AddedSource is a source that was added or removed between releases.
type AddedSource struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Version     string `json:"version,omitempty"`
	Commit      string `json:"commit,omitempty"`
}

// SourceUpdate describes a version or commit change for a source present in both releases.
type SourceUpdate struct {
	Name        string           `json:"name"`
	Type        string           `json:"type"`
	FromVersion string           `json:"from_version,omitempty"`
	ToVersion   string           `json:"to_version,omitempty"`
	FromCommit  string           `json:"from_commit,omitempty"`
	ToCommit    string           `json:"to_commit,omitempty"`
	ZimAdded    []string         `json:"zim_added,omitempty"`
	ZimRemoved  []string         `json:"zim_removed,omitempty"`
	ZimUpdated  []ZimFileUpdate  `json:"zim_updated,omitempty"`
}

// ZimFileUpdate describes a version change for a ZIM file present in both releases.
type ZimFileUpdate struct {
	Name        string `json:"name"`
	FromVersion string `json:"from_version"`
	ToVersion   string `json:"to_version"`
}

// TiersDiff describes tier composition changes between releases.
type TiersDiff struct {
	Added   map[string][]string `json:"added,omitempty"`
	Removed map[string][]string `json:"removed,omitempty"`
}

// Compare produces a Diff describing what changed from old to new.
func Compare(old, new *Manifest) *Diff {
	d := &Diff{
		From: old.Release,
		To:   new.Release,
	}

	// Source-level diff.
	for name, ne := range new.Sources {
		oe, existed := old.Sources[name]
		if !existed {
			d.Sources.Added = append(d.Sources.Added, AddedSource{
				Name:        name,
				Type:        ne.Type,
				Description: ne.Description,
				Version:     ne.Version,
				Commit:      ne.Commit,
			})
			continue
		}
		zimAdded, zimRemoved, zimUpdated := diffZimFiles(oe.ZimFiles, ne.ZimFiles)
		if oe.Version != ne.Version || oe.Commit != ne.Commit ||
			len(zimAdded) > 0 || len(zimRemoved) > 0 || len(zimUpdated) > 0 {
			d.Sources.Updated = append(d.Sources.Updated, SourceUpdate{
				Name:        name,
				Type:        ne.Type,
				FromVersion: oe.Version,
				ToVersion:   ne.Version,
				FromCommit:  oe.Commit,
				ToCommit:    ne.Commit,
				ZimAdded:    zimAdded,
				ZimRemoved:  zimRemoved,
				ZimUpdated:  zimUpdated,
			})
		} else {
			d.Sources.Unchanged = append(d.Sources.Unchanged, name)
		}
	}
	for name, oe := range old.Sources {
		if _, exists := new.Sources[name]; !exists {
			d.Sources.Removed = append(d.Sources.Removed, AddedSource{
				Name:        name,
				Type:        oe.Type,
				Description: oe.Description,
				Version:     oe.Version,
				Commit:      oe.Commit,
			})
		}
	}

	// Sort for stable output.
	sort.Slice(d.Sources.Added, func(i, j int) bool { return d.Sources.Added[i].Name < d.Sources.Added[j].Name })
	sort.Slice(d.Sources.Removed, func(i, j int) bool { return d.Sources.Removed[i].Name < d.Sources.Removed[j].Name })
	sort.Slice(d.Sources.Updated, func(i, j int) bool { return d.Sources.Updated[i].Name < d.Sources.Updated[j].Name })
	sort.Strings(d.Sources.Unchanged)

	// Tier-level diff.
	d.Tiers.Added = make(map[string][]string)
	d.Tiers.Removed = make(map[string][]string)

	allTiers := make(map[string]bool)
	for k := range old.Tiers {
		allTiers[k] = true
	}
	for k := range new.Tiers {
		allTiers[k] = true
	}

	for tierKey := range allTiers {
		oldSet := toSet(old.Tiers[tierKey].Sources)
		newSet := toSet(new.Tiers[tierKey].Sources)

		var added, removed []string
		for s := range newSet {
			if !oldSet[s] {
				added = append(added, s)
			}
		}
		for s := range oldSet {
			if !newSet[s] {
				removed = append(removed, s)
			}
		}
		sort.Strings(added)
		sort.Strings(removed)
		if len(added) > 0 {
			d.Tiers.Added[tierKey] = added
		}
		if len(removed) > 0 {
			d.Tiers.Removed[tierKey] = removed
		}
	}

	return d
}

// JSON returns the diff as indented JSON.
func (d *Diff) JSON() ([]byte, error) {
	return json.MarshalIndent(d, "", "  ")
}

// Text returns a human-readable markdown summary of the diff.
func (d *Diff) Text() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# Changes: %s → %s\n", d.From, d.To)

	if len(d.Sources.Updated) > 0 {
		sb.WriteString("\n## Updated\n")
		for _, u := range d.Sources.Updated {
			if u.FromVersion != "" || u.ToVersion != "" {
				fmt.Fprintf(&sb, "- **%s**: %s → %s\n", u.Name, displayVersion(u.FromVersion), displayVersion(u.ToVersion))
			} else {
				fmt.Fprintf(&sb, "- **%s**: %s → %s\n", u.Name, shortCommit(u.FromCommit), shortCommit(u.ToCommit))
			}
			for _, z := range u.ZimUpdated {
				fmt.Fprintf(&sb, "  - zim %s: %s → %s\n", z.Name, z.FromVersion, z.ToVersion)
			}
			for _, z := range u.ZimAdded {
				fmt.Fprintf(&sb, "  - zim added: %s\n", z)
			}
			for _, z := range u.ZimRemoved {
				fmt.Fprintf(&sb, "  - zim removed: %s\n", z)
			}
		}
	}

	if len(d.Sources.Added) > 0 {
		sb.WriteString("\n## Added\n")
		for _, a := range d.Sources.Added {
			ver := a.Version
			if ver == "" {
				ver = shortCommit(a.Commit)
			}
			if ver == "" {
				ver = "new"
			}
			fmt.Fprintf(&sb, "- **%s** %s\n", a.Name, ver)
		}
	}

	if len(d.Sources.Removed) > 0 {
		sb.WriteString("\n## Removed\n")
		for _, r := range d.Sources.Removed {
			fmt.Fprintf(&sb, "- **%s**\n", r.Name)
		}
	}

	if len(d.Sources.Unchanged) > 0 {
		sb.WriteString("\n## Unchanged\n")
		sb.WriteString("- " + strings.Join(d.Sources.Unchanged, ", ") + "\n")
	}

	// Tier changes
	hasTierChanges := len(d.Tiers.Added) > 0 || len(d.Tiers.Removed) > 0
	if hasTierChanges {
		sb.WriteString("\n## Tier Changes\n")
		tierKeys := sortedKeys(d.Tiers.Added, d.Tiers.Removed)
		for _, tk := range tierKeys {
			if added := d.Tiers.Added[tk]; len(added) > 0 {
				fmt.Fprintf(&sb, "- **%s**: added %s\n", tk, strings.Join(added, ", "))
			}
			if removed := d.Tiers.Removed[tk]; len(removed) > 0 {
				fmt.Fprintf(&sb, "- **%s**: removed %s\n", tk, strings.Join(removed, ", "))
			}
		}
	}

	return sb.String()
}

// diffZimFiles compares two ZIM file lists and returns added names, removed names,
// and version changes for files present in both.
func diffZimFiles(old, new []source.ZimFileEntry) (added, removed []string, updated []ZimFileUpdate) {
	oldByName := make(map[string]source.ZimFileEntry, len(old))
	for _, z := range old {
		oldByName[z.Name] = z
	}
	newByName := make(map[string]source.ZimFileEntry, len(new))
	for _, z := range new {
		newByName[z.Name] = z
	}
	for _, z := range new {
		if o, exists := oldByName[z.Name]; !exists {
			added = append(added, z.Name)
		} else if o.Version != z.Version {
			updated = append(updated, ZimFileUpdate{Name: z.Name, FromVersion: o.Version, ToVersion: z.Version})
		}
	}
	for _, z := range old {
		if _, exists := newByName[z.Name]; !exists {
			removed = append(removed, z.Name)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	sort.Slice(updated, func(i, j int) bool { return updated[i].Name < updated[j].Name })
	return
}

func toSet(ss []string) map[string]bool {
	m := make(map[string]bool, len(ss))
	for _, s := range ss {
		m[s] = true
	}
	return m
}

func shortCommit(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

func displayVersion(v string) string {
	if v == "" {
		return "(unknown)"
	}
	return v
}

func sortedKeys(maps ...map[string][]string) []string {
	seen := make(map[string]bool)
	for _, m := range maps {
		for k := range m {
			seen[k] = true
		}
	}
	keys := make([]string, 0, len(seen))
	for k := range seen {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
