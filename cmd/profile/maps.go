package profile

import (
	"context"
	"fmt"
	"sort"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var MapsCommand = cli.Command{
	Name:   "maps",
	Usage:  "List all folder path mappings",
	Action: listMaps,
}

func listMaps(_ context.Context, _ *cli.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if cfg.DefaultProfile != "" {
		fmt.Printf("default: %s\n\n", colors.Cyan(cfg.DefaultProfile))
	}

	if len(cfg.PathMappings) == 0 {
		fmt.Println(colors.Yellow("No path mappings configured."))
		fmt.Println(colors.Yellow("Use: gixy profile map <profile-name> <path>"))
		return nil
	}

	paths := make([]string, 0, len(cfg.PathMappings))
	for p := range cfg.PathMappings {
		paths = append(paths, p)
	}
	sort.Strings(paths)

	for _, p := range paths {
		profileName := cfg.PathMappings[p]
		fmt.Printf("  %-50s → %s\n", p, colors.Cyan(profileName))
	}
	return nil
}
