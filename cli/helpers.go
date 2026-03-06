package cli

import (
	"fmt"
	"path/filepath"

	"github.com/stazelabs/open-access-kit/internal/config"
	"github.com/stazelabs/open-access-kit/internal/source"
)

// loadConfig reads and returns the config from the global cfgFile path.
func loadConfig() (*config.Config, error) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, fmt.Errorf("loading config: %w", err)
	}
	return cfg, nil
}

// sourcesForTier returns Source objects for all sources in the given tier,
// ordered as listed in the tier config.
func sourcesForTier(cfg *config.Config, tierKey string) ([]source.Source, error) {
	t, ok := cfg.Tiers[tierKey]
	if !ok {
		return nil, fmt.Errorf("unknown tier %q (available: S, M, L)", tierKey)
	}
	return buildSources(cfg, t.Sources)
}

// globalVars returns the template vars derived from the top-level config,
// used when rendering Markdown guides to HTML during staging.
func globalVars(cfg *config.Config) map[string]any {
	return map[string]any{
		"Release":      cfg.Release,
		"DownloadRoot": cfg.DownloadRoot,
	}
}

// sourceByName returns a single Source for the named config entry.
func sourceByName(cfg *config.Config, name string) (source.Source, error) {
	scfg, ok := cfg.Sources[name]
	if !ok {
		return nil, fmt.Errorf("unknown source %q", name)
	}
	return source.New(name, scfg, globalVars(cfg))
}

// buildSources constructs Source objects for the given list of source names.
func buildSources(cfg *config.Config, names []string) ([]source.Source, error) {
	vars := globalVars(cfg)
	sources := make([]source.Source, 0, len(names))
	for _, name := range names {
		scfg, ok := cfg.Sources[name]
		if !ok {
			return nil, fmt.Errorf("source %q referenced in tier but not defined in sources", name)
		}
		s, err := source.New(name, scfg, vars)
		if err != nil {
			return nil, err
		}
		sources = append(sources, s)
	}
	return sources, nil
}

// effectiveMirrorDir returns the mirror directory, preferring the config
// value unless it was overridden by the --mirror-dir flag.
func effectiveMirrorDir(cfg *config.Config) string {
	if mirrorDir != "./mirror" {
		return mirrorDir
	}
	if cfg.Paths.Mirror != "" {
		return cfg.Paths.Mirror
	}
	return "./mirror"
}

// effectiveImageDir returns the image base directory, preferring the config
// value unless it was overridden by the --image-dir flag.
func effectiveImageDir(cfg *config.Config) string {
	if imageDir != "./image" {
		return imageDir
	}
	if cfg.Paths.Image != "" {
		return cfg.Paths.Image
	}
	return "./image"
}

// effectiveReleaseImageDir returns the tier-specific staging directory,
// e.g. image/OAK-Q126-M. Each tier gets its own directory so multiple
// tiers can coexist without overwriting each other.
// tierLabel is the human label from the tier config (e.g. "S", "M", "L").
func effectiveReleaseImageDir(cfg *config.Config, tierLabel string) string {
	return filepath.Join(effectiveImageDir(cfg), "OAK-"+cfg.Release+"-"+tierLabel)
}

// zipRootName returns the directory name inside the ZIP (what appears on the removable media).
// This is tier-agnostic — every tier extracts to OAK-{release}/.
func zipRootName(cfg *config.Config) string {
	return "OAK-" + cfg.Release
}

// effectiveOutputDir returns the dist directory from config.
func effectiveOutputDir(cfg *config.Config) string {
	if cfg.Paths.Output != "" {
		return cfg.Paths.Output
	}
	return "./dist"
}
