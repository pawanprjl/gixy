package profile

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var ShowCommand = cli.Command{
	Name:   "show",
	Usage:  "Show the active git identity for the current repository",
	Action: showProfile,
}

func showProfile(_ context.Context, _ *cli.Command) error {
	nameOut, err := exec.Command("git", "config", "--local", "user.name").Output()
	if err != nil {
		return cli.Exit(colors.Red("no local git identity set; run `gixy profile use <name>` first"), 1)
	}
	emailOut, err := exec.Command("git", "config", "--local", "user.email").Output()
	if err != nil {
		return cli.Exit(colors.Red("no local git identity set; run `gixy profile use <name>` first"), 1)
	}

	name := strings.TrimSpace(string(nameOut))
	email := strings.TrimSpace(string(emailOut))

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	matchedProfile := ""
	for profileName, p := range cfg.Profiles {
		if p.Name == name && p.Email == email {
			matchedProfile = profileName
			break
		}
	}

	if matchedProfile != "" {
		fmt.Printf("%s %s  %s <%s>\n", colors.Green("●"), colors.Green(matchedProfile), name, email)
	} else {
		fmt.Printf("%s %s <%s>\n", colors.Yellow("●"), name, email)
		fmt.Println(colors.Yellow("(no matching gixy profile)"))
	}
	return nil
}
