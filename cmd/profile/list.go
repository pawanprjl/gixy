package profile

import (
	"context"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var ListCommand = cli.Command{
	Name:   "list",
	Usage:  "List all profiles",
	Action: listProfiles,
}

func activeGitIdentity() (name, email string) {
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

func listProfiles(_ context.Context, _ *cli.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if len(cfg.Profiles) == 0 {
		fmt.Println(colors.Yellow("No profiles found."))
		return nil
	}

	activeName, activeEmail := activeGitIdentity()

	keys := make([]string, 0, len(cfg.Profiles))
	for k := range cfg.Profiles {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		p := cfg.Profiles[k]
		marker := "  "
		if p.Name == activeName && p.Email == activeEmail {
			marker = colors.Green("* ")
		}
		fmt.Printf("%s%-20s %s <%s>\n", marker, colors.Cyan(k), p.Name, p.Email)
	}
	return nil
}
