package profile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var ActivateCommand = cli.Command{
	Name:  "activate",
	Usage: "Auto-activate the profile for the current directory (used by shell hooks)",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "silent",
			Usage: "Suppress all output",
		},
	},
	Action: activateProfile,
}

func activateProfile(_ context.Context, cmd *cli.Command) error {
	silent := cmd.Bool("silent")

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return cli.Exit(fmt.Errorf("get working directory: %w", err), 1)
	}
	cwd = filepath.Clean(cwd)

	matched := longestPrefixMatch(cwd, cfg.PathMappings)

	if matched != "" {
		return applyProfileMaybeSilent(matched, cfg, silent)
	}

	if cfg.DefaultProfile != "" {
		return applyProfileMaybeSilent(cfg.DefaultProfile, cfg, silent)
	}

	// No mapping and no default — do nothing.
	return nil
}

// longestPrefixMatch returns the profile name for the most-specific path mapping
// that is a prefix of cwd, or "" if none match.
func longestPrefixMatch(cwd string, mappings map[string]string) string {
	bestLen := -1
	bestProfile := ""

	for mappedPath, profileName := range mappings {
		clean := filepath.Clean(mappedPath)
		// Require cwd == mappedPath or cwd is inside mappedPath (trailing slash prevents
		// /work matching /workspace).
		if cwd == clean || strings.HasPrefix(cwd+string(filepath.Separator), clean+string(filepath.Separator)) {
			if len(clean) > bestLen {
				bestLen = len(clean)
				bestProfile = profileName
			}
		}
	}

	return bestProfile
}

func applyProfileMaybeSilent(profileName string, cfg *config.Config, silent bool) error {
	if silent {
		// Redirect stdout by temporarily swapping os.Stdout; applyProfile writes there.
		// Simpler: wrap in a version that skips the Println.
		return applyProfileSilent(profileName, cfg)
	}
	return applyProfile(profileName, cfg)
}
