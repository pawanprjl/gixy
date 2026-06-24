package profile

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var StatusCommand = cli.Command{
	Name:   "status",
	Usage:  "Show which profile applies to the current directory and why",
	Action: showStatus,
}

func showStatus(_ context.Context, _ *cli.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return cli.Exit(fmt.Errorf("get working directory: %w", err), 1)
	}
	cwd = filepath.Clean(cwd)

	matchedPath, name := resolveProfileName(cwd, cfg)

	fmt.Printf("%s %s\n", colors.Cyan("Directory:"), cwd)

	switch {
	case matchedPath != "":
		fmt.Printf("%s   %s → %s\n", colors.Cyan("Mapping:"), matchedPath, name)
	case name != "":
		fmt.Printf("%s   %s\n", colors.Cyan("Mapping:"), colors.Dim("(default profile)"))
	default:
		fmt.Printf("%s   %s\n", colors.Cyan("Mapping:"), colors.Dim("(none — git runs unchanged)"))
	}

	if name != "" {
		if p, ok := cfg.Profiles[name]; ok {
			fmt.Printf("%s   %s — %s <%s>\n", colors.Cyan("Applies:"), colors.Green(name), p.Name, p.Email)
		} else {
			fmt.Printf("%s   %s\n", colors.Cyan("Applies:"), colors.Yellow(fmt.Sprintf("%q (profile no longer exists)", name)))
		}
	}

	if bName, bEmail := globalGitIdentity(); bName != "" || bEmail != "" {
		fmt.Printf("%s    %s <%s> %s\n", colors.Cyan("Global:"), bName, bEmail, colors.Dim("(baseline for non-shell git)"))
	}

	return nil
}

// globalGitIdentity returns the global git user.name/email baseline (empty if unset).
func globalGitIdentity() (name, email string) {
	nameOut, err := exec.Command("git", "config", "--global", "user.name").Output()
	if err != nil {
		return "", ""
	}
	emailOut, err := exec.Command("git", "config", "--global", "user.email").Output()
	if err != nil {
		return "", ""
	}
	return strings.TrimSpace(string(nameOut)), strings.TrimSpace(string(emailOut))
}
