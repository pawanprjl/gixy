package profile

import (
	"path/filepath"
	"strings"

	"github.com/pawanprjl/gixy/internal/config"
)

// resolveProfileName returns the profile for cwd (longest mapping prefix, else default); matchedPath is empty when defaulted.
func resolveProfileName(cwd string, cfg *config.Config) (matchedPath, name string) {
	path, profile := longestPrefixMatch(cwd, cfg.PathMappings)
	if profile == "" {
		return "", cfg.DefaultProfile
	}
	return path, profile
}

// longestPrefixMatch returns the most-specific mapping prefix of cwd as (path, profile), or ("", "").
func longestPrefixMatch(cwd string, mappings map[string]string) (string, string) {
	bestLen := -1
	bestPath := ""
	bestProfile := ""

	for mappedPath, profileName := range mappings {
		clean := filepath.Clean(mappedPath)
		// cwd == mapping or inside it; trailing slash stops /work matching /workspace.
		if cwd == clean || strings.HasPrefix(cwd+string(filepath.Separator), clean+string(filepath.Separator)) {
			if len(clean) > bestLen {
				bestLen = len(clean)
				bestPath = clean
				bestProfile = profileName
			}
		}
	}

	return bestPath, bestProfile
}
