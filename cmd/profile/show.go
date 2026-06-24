package profile

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/pawanprjl/gixy/internal/sshutil"
	"github.com/urfave/cli/v3"
)

var ShowCommand = cli.Command{
	Name:      "show",
	Usage:     "Show a profile's identity, SSH key, and mappings",
	ArgsUsage: "<profile-name>",
	Action:    showProfile,
}

func showProfile(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit(colors.Red("usage: gixy profile show <profile-name>"), 1)
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

	fmt.Printf("%s     %s\n", colors.Cyan("Profile:"), colors.Green(profileName))
	fmt.Printf("%s    %s <%s>\n", colors.Cyan("Identity:"), p.Name, p.Email)
	if cfg.DefaultProfile == profileName {
		fmt.Printf("%s     %s\n", colors.Cyan("Default:"), "yes (used when no mapping matches)")
	}

	// SSH key
	if keyDir, err := sshutil.KeyDir(profileName); err == nil {
		privPath := filepath.Join(keyDir, "id_ed25519")
		pubPath := privPath + ".pub"
		if _, err := os.Stat(privPath); err == nil {
			fmt.Printf("%s     %s\n", colors.Cyan("SSH key:"), privPath)
			if out, err := exec.Command("ssh-keygen", "-lf", pubPath).Output(); err == nil {
				fmt.Printf("%s %s\n", colors.Cyan("Fingerprint:"), strings.TrimSpace(string(out)))
			}
		} else {
			fmt.Printf("%s     %s\n", colors.Cyan("SSH key:"), colors.Yellow("none (run `gixy profile add` to generate)"))
		}
	}

	// Mappings pointing at this profile
	var paths []string
	for path, name := range cfg.PathMappings {
		if name == profileName {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	if len(paths) > 0 {
		fmt.Printf("%s\n", colors.Cyan("Mappings:"))
		for _, path := range paths {
			fmt.Printf("  %s\n", path)
		}
	}

	return nil
}
