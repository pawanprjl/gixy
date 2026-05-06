package profile

import (
	"context"
	"fmt"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var DefaultCommand = cli.Command{
	Name:      "default",
	Usage:     "Set the default profile used when no path mapping matches",
	ArgsUsage: "<profile-name>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "clear",
			Usage: "Unset the default profile",
		},
	},
	Action: setDefault,
}

func setDefault(_ context.Context, cmd *cli.Command) error {
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
		return cli.Exit(colors.Red("usage: gixy profile default <profile-name>"), 1)
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
