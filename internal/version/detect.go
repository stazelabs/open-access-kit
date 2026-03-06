package version

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

// RsyncList runs `rsync <url>` to list a remote directory, finds all matches
// of pattern (first capture group), and returns the selected version.
func RsyncList(ctx context.Context, url, pattern, selectMode string) (string, error) {
	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, "rsync", url)
	cmd.Stdout = &out
	// rsync exits non-zero if the listing is empty; stderr noise is fine to ignore
	_ = cmd.Run()

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("compiling pattern %q: %w", pattern, err)
	}

	matches := re.FindAllSubmatch(out.Bytes(), -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("no versions found in rsync listing of %s with pattern %q", url, pattern)
	}

	var versions []string
	seen := map[string]bool{}
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		v := string(m[1])
		if !seen[v] {
			seen[v] = true
			versions = append(versions, v)
		}
	}

	switch selectMode {
	case "highest-semver", "":
		return highestSemver(versions)
	case "latest-date":
		return latestDate(versions)
	case "first":
		return versions[0], nil
	default:
		return "", fmt.Errorf("unknown select mode: %s", selectMode)
	}
}

// HTTPScrape fetches url, finds all matches of pattern (first capture group),
// and returns the selected version. selectMode is "highest-semver" or "first".
func HTTPScrape(ctx context.Context, url, pattern, selectMode string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", fmt.Errorf("building request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetching %s: %w", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("reading response: %w", err)
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("compiling pattern %q: %w", pattern, err)
	}

	matches := re.FindAllSubmatch(body, -1)
	if len(matches) == 0 {
		return "", fmt.Errorf("no versions found at %s with pattern %q", url, pattern)
	}

	var versions []string
	seen := map[string]bool{}
	for _, m := range matches {
		if len(m) < 2 {
			continue
		}
		v := string(m[1])
		if !seen[v] {
			seen[v] = true
			versions = append(versions, v)
		}
	}

	switch selectMode {
	case "highest-semver", "":
		return highestSemver(versions)
	case "latest-date":
		return latestDate(versions)
	case "first":
		return versions[0], nil
	default:
		return "", fmt.Errorf("unknown select mode: %s", selectMode)
	}
}

// latestDate returns the lexicographically greatest value from a list of
// "YYYY-MM" date strings. Since the format is zero-padded, lex order == date order.
func latestDate(dates []string) (string, error) {
	if len(dates) == 0 {
		return "", fmt.Errorf("no versions to compare")
	}
	best := dates[0]
	for _, d := range dates[1:] {
		if d > best {
			best = d
		}
	}
	return best, nil
}

// highestSemver returns the highest version from a list of "X.Y.Z" strings.
func highestSemver(versions []string) (string, error) {
	if len(versions) == 0 {
		return "", fmt.Errorf("no versions to compare")
	}
	best := versions[0]
	for _, v := range versions[1:] {
		if cmpSemver(v, best) > 0 {
			best = v
		}
	}
	return best, nil
}

// cmpSemver compares two dotted-version strings. Returns >0 if a > b.
func cmpSemver(a, b string) int {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	maxLen := len(aParts)
	if len(bParts) > maxLen {
		maxLen = len(bParts)
	}
	for i := 0; i < maxLen; i++ {
		var av, bv int
		if i < len(aParts) {
			av, _ = strconv.Atoi(aParts[i])
		}
		if i < len(bParts) {
			bv, _ = strconv.Atoi(bParts[i])
		}
		if av != bv {
			return av - bv
		}
	}
	return 0
}
