package profile

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var PinCommand = cli.Command{
	Name:      "pin",
	Usage:     "Pin the current repo's local git config to a profile (applies to all git tools)",
	ArgsUsage: "[profile-name]",
	Action:    runPin,
}

var UnpinCommand = cli.Command{
	Name:   "unpin",
	Usage:  "Remove gixy-managed identity/SSH settings from the current repo's local config",
	Action: runUnpin,
}

func runPin(_ context.Context, cmd *cli.Command) error {
	if !insideGitRepo() {
		return cli.Exit(colors.Red("not inside a git repository"), 1)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	profileName := cmd.Args().First()
	if profileName == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return cli.Exit(fmt.Errorf("get working directory: %w", err), 1)
		}
		_, profileName = resolveProfileName(filepath.Clean(cwd), cfg)
		if profileName == "" {
			return cli.Exit(colors.Red("no profile maps to this directory — pass a profile name or add a mapping"), 1)
		}
	}

	p, ok := cfg.Profiles[profileName]
	if !ok {
		return cli.Exit(colors.Red(fmt.Sprintf("profile %q not found", profileName)), 1)
	}

	if err := writeRepoProfile(profileName, p); err != nil {
		return cli.Exit(fmt.Errorf("pin profile: %w", err), 1)
	}

	fmt.Println(colors.Green(fmt.Sprintf("Pinned this repo to %q — %s <%s>", profileName, p.Name, p.Email)))
	return nil
}

func runUnpin(_ context.Context, _ *cli.Command) error {
	if !insideGitRepo() {
		return cli.Exit(colors.Red("not inside a git repository"), 1)
	}

	for _, key := range []string{"user.name", "user.email", "core.sshCommand", gixyProfileKey} {
		unsetLocalConfig(key)
	}

	fmt.Println(colors.Green("Unpinned this repo — gixy auto-activation / global baseline applies again"))
	return nil
}
