package profile

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var MapCommand = cli.Command{
	Name:      "map",
	Usage:     "Map a folder path to a profile for auto-activation",
	ArgsUsage: "<profile-name> <path>",
	Action:    mapProfile,
}

var UnmapCommand = cli.Command{
	Name:      "unmap",
	Usage:     "Remove a folder path mapping",
	ArgsUsage: "<path>",
	Action:    unmapProfile,
}

func mapProfile(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 2 {
		return cli.Exit(colors.Red("usage: gixy profile map <profile-name> <path>"), 1)
	}
	profileName := cmd.Args().Get(0)
	rawPath := cmd.Args().Get(1)

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if _, exists := cfg.Profiles[profileName]; !exists {
		return cli.Exit(colors.Red(fmt.Sprintf("profile %q not found", profileName)), 1)
	}

	absPath, err := filepath.Abs(rawPath)
	if err != nil {
		return cli.Exit(fmt.Errorf("resolve path: %w", err), 1)
	}

	cfg.PathMappings[absPath] = profileName

	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Println(colors.Green(fmt.Sprintf("Mapped %s → profile %q", absPath, profileName)))
	return nil
}

func unmapProfile(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit(colors.Red("usage: gixy profile unmap <path>"), 1)
	}
	rawPath := cmd.Args().Get(0)

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	absPath, err := filepath.Abs(rawPath)
	if err != nil {
		return cli.Exit(fmt.Errorf("resolve path: %w", err), 1)
	}

	if _, exists := cfg.PathMappings[absPath]; !exists {
		return cli.Exit(colors.Red(fmt.Sprintf("no mapping found for %s", absPath)), 1)
	}

	delete(cfg.PathMappings, absPath)

	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Println(colors.Green(fmt.Sprintf("Removed mapping for %s", absPath)))
	return nil
}
