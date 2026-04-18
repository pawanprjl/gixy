package profile

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var AddCommand = cli.Command{
	Name:      "add",
	Usage:     "Add a new profile",
	ArgsUsage: "<profile-name>",
	Action:    addProfile,
}

func addProfile(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit("usage: gixy profile add <profile-name>", 1)
	}
	profileName := cmd.Args().Get(0)

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if _, exists := cfg.Profiles[profileName]; exists {
		return cli.Exit(fmt.Sprintf("profile %q already exists", profileName), 1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Git name: ")
	name, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read name: %w", err), 1)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return cli.Exit("name cannot be empty", 1)
	}

	fmt.Print("Git email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read email: %w", err), 1)
	}
	email = strings.TrimSpace(email)
	if email == "" {
		return cli.Exit("email cannot be empty", 1)
	}

	cfg.Profiles[profileName] = config.Profile{Name: name, Email: email}

	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Printf("Profile %q added.\n", profileName)
	return nil
}
