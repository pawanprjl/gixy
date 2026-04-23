package profile

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var EditCommand = cli.Command{
	Name:      "edit",
	Usage:     "Edit an existing profile's name and email",
	ArgsUsage: "<profile-name>",
	Action:    editProfile,
}

func editProfile(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit(colors.Red("usage: gixy profile edit <profile-name>"), 1)
	}
	profileName := cmd.Args().Get(0)

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	p, exists := cfg.Profiles[profileName]
	if !exists {
		return cli.Exit(colors.Red(fmt.Sprintf("profile %q not found", profileName)), 1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [%s]: ", colors.Cyan("Git name"), p.Name)
	name, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read name: %w", err), 1)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		name = p.Name
	}

	fmt.Printf("%s [%s]: ", colors.Cyan("Git email"), p.Email)
	email, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read email: %w", err), 1)
	}
	email = strings.TrimSpace(email)
	if email == "" {
		email = p.Email
	} else if !validEmail(email) {
		return cli.Exit(colors.Red("invalid email address"), 1)
	}

	cfg.Profiles[profileName] = config.Profile{Name: name, Email: email}

	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Println(colors.Green(fmt.Sprintf("Profile %q updated: %s <%s>", profileName, name, email)))
	return nil
}
