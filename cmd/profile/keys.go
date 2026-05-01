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
	"github.com/pawanprjl/gixy/internal/sshutil"
	"github.com/urfave/cli/v3"
)

var KeysCommand = cli.Command{
	Name:      "keys",
	Usage:     "Show SSH keys for a profile",
	ArgsUsage: "<profile-name>",
	Action:    showKeys,
}

func showKeys(_ context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() != 1 {
		return cli.Exit(colors.Red("usage: gixy profile keys <profile-name>"), 1)
	}
	profileName := cmd.Args().Get(0)

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if _, exists := cfg.Profiles[profileName]; !exists {
		return cli.Exit(colors.Red(fmt.Sprintf("profile %q not found", profileName)), 1)
	}

	keyDir, err := sshutil.KeyDir(profileName)
	if err != nil {
		return cli.Exit(fmt.Errorf("resolve key dir: %w", err), 1)
	}

	privPath := filepath.Join(keyDir, "id_ed25519")
	pubPath := filepath.Join(keyDir, "id_ed25519.pub")

	if _, err := os.Stat(privPath); os.IsNotExist(err) {
		return cli.Exit(colors.Yellow(fmt.Sprintf(
			"No SSH keys found for profile %q. Re-run `gixy profile add %s` to generate them.",
			profileName, profileName,
		)), 0)
	}

	pubBytes, err := os.ReadFile(pubPath)
	if err != nil {
		return cli.Exit(fmt.Errorf("read public key: %w", err), 1)
	}

	// Fingerprint via ssh-keygen
	fingerprintOut, err := exec.Command("ssh-keygen", "-lf", pubPath).Output()
	fingerprint := ""
	if err == nil {
		fingerprint = strings.TrimSpace(string(fingerprintOut))
	}

	fmt.Printf("%s  %s\n", colors.Cyan("Private key:"), privPath)
	fmt.Printf("%s %s\n", colors.Cyan("Public key: "), pubPath)
	if fingerprint != "" {
		fmt.Printf("%s %s\n", colors.Cyan("Fingerprint:"), fingerprint)
	}
	fmt.Printf("\n%s\n", colors.Cyan("Public key contents:"))
	fmt.Print(string(pubBytes))
	return nil
}
