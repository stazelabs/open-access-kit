package cli

import (
	"fmt"

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
		return nil, fmt.Errorf("unknown tier %q (available: 16, 32, 64, max)", tierKey)
	}
	return buildSources(cfg, t.Sources)
}

// sourceByName returns a single Source for the named config entry.
func sourceByName(cfg *config.Config, name string) (source.Source, error) {
	scfg, ok := cfg.Sources[name]
	if !ok {
		return nil, fmt.Errorf("unknown source %q", name)
	}
	return source.New(name, scfg)
}

// buildSources constructs Source objects for the given list of source names.
func buildSources(cfg *config.Config, names []string) ([]source.Source, error) {
	sources := make([]source.Source, 0, len(names))
	for _, name := range names {
		scfg, ok := cfg.Sources[name]
		if !ok {
			return nil, fmt.Errorf("source %q referenced in tier but not defined in sources", name)
		}
		s, err := source.New(name, scfg)
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

// effectiveImageDir returns the image directory, preferring the config
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
