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

var DeleteCommand = cli.Command{
	Name:      "delete",
	Aliases:   []string{"remove"},
	Usage:     "Delete a profile",
	ArgsUsage: "<profile-name>",
	Action:    deleteProfile,
}

func deleteProfile(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit(colors.Red("usage: gixy profile delete <profile-name>"), 1)
	}
	profileName := cmd.Args().Get(0)

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if _, exists := cfg.Profiles[profileName]; !exists {
		return cli.Exit(colors.Red(fmt.Sprintf("profile %q not found", profileName)), 1)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Delete profile %s? [y/N]: ", colors.Cyan(profileName))
	answer, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read input: %w", err), 1)
	}
	answer = strings.TrimSpace(strings.ToLower(answer))
	if answer != "y" && answer != "yes" {
		fmt.Println("Aborted.")
		return nil
	}

	delete(cfg.Profiles, profileName)

	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Println(colors.Green(fmt.Sprintf("Profile %q deleted.", profileName)))
	return nil
}
