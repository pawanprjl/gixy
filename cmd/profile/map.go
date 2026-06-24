package profile

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var MapGroupCommand = cli.Command{
	Name:  "map",
	Usage: "Manage folder path mappings for auto-activation",
	Commands: []*cli.Command{
		{
			Name:      "add",
			Usage:     "Map a folder path to a profile",
			ArgsUsage: "<profile-name> <path>",
			Action:    mapAdd,
		},
		{
			Name:      "remove",
			Usage:     "Remove a folder path mapping",
			ArgsUsage: "<path>",
			Action:    mapRemove,
		},
		{
			Name:   "list",
			Usage:  "List all folder path mappings",
			Action: mapList,
		},
		{
			Name:      "default",
			Usage:     "Set the fallback profile used when no path mapping matches",
			ArgsUsage: "<profile-name>",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "clear",
					Usage: "Unset the default profile",
				},
			},
			Action: mapDefault,
		},
	},
}

func mapDefault(_ context.Context, cmd *cli.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if cmd.Bool("clear") {
		cfg.DefaultProfile = ""
		if err := config.SaveConfig(cfg); err != nil {
			return cli.Exit(fmt.Errorf("save config: %w", err), 1)
		}
		fmt.Println(colors.Green("Default profile cleared."))
		return nil
	}

	if cmd.Args().Len() != 1 {
		return cli.Exit(colors.Red("usage: gixy profile map default <profile-name>"), 1)
	}
	profileName := cmd.Args().Get(0)

	if _, exists := cfg.Profiles[profileName]; !exists {
		return cli.Exit(colors.Red(fmt.Sprintf("profile %q not found", profileName)), 1)
	}

	cfg.DefaultProfile = profileName

	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Println(colors.Green(fmt.Sprintf("Default profile set to %q", profileName)))
	return nil
}

func mapAdd(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 2 {
		return cli.Exit(colors.Red("usage: gixy profile map add <profile-name> <path>"), 1)
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

func mapRemove(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit(colors.Red("usage: gixy profile map remove <path>"), 1)
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

func mapList(_ context.Context, _ *cli.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if cfg.DefaultProfile != "" {
		fmt.Printf("default: %s\n\n", colors.Cyan(cfg.DefaultProfile))
	}

	if len(cfg.PathMappings) == 0 {
		fmt.Println(colors.Yellow("No path mappings configured."))
		fmt.Println(colors.Yellow("Use: gixy profile map add <profile-name> <path>"))
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
