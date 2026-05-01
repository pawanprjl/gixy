package profile

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/pawanprjl/gixy/internal/sshutil"
	"github.com/urfave/cli/v3"
)

var UseCommand = cli.Command{
	Name:      "use",
	Usage:     "Activate a profile globally (git identity + SSH keys)",
	ArgsUsage: "<profile-name>",
	Action:    useProfile,
}

func useProfile(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit(colors.Red("usage: gixy profile use <profile-name>"), 1)
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

	if err := exec.Command("git", "config", "--global", "user.name", p.Name).Run(); err != nil {
		return cli.Exit(fmt.Errorf("set user.name: %w", err), 1)
	}

	if err := exec.Command("git", "config", "--global", "user.email", p.Email).Run(); err != nil {
		return cli.Exit(fmt.Errorf("set user.email: %w", err), 1)
	}

	if err := sshutil.ActivateKeys(profileName); err != nil {
		return cli.Exit(fmt.Errorf("activate SSH keys: %w", err), 1)
	}

	fmt.Println(colors.Green(fmt.Sprintf("Switched to profile %q — %s <%s>", profileName, p.Name, p.Email)))
	return nil
}
