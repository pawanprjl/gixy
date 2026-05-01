package profile

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/pawanprjl/gixy/internal/sshutil"
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
		return cli.Exit(colors.Red("usage: gixy profile add <profile-name>"), 1)
	}
	profileName := cmd.Args().Get(0)

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if _, exists := cfg.Profiles[profileName]; exists {
		return cli.Exit(colors.Red(fmt.Sprintf("profile %q already exists", profileName)), 1)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(colors.Cyan("Git name: "))
	name, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read name: %w", err), 1)
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return cli.Exit(colors.Red("name cannot be empty"), 1)
	}

	fmt.Print(colors.Cyan("Git email: "))
	email, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read email: %w", err), 1)
	}
	email = strings.TrimSpace(email)
	if email == "" {
		return cli.Exit(colors.Red("email cannot be empty"), 1)
	}
	if !validEmail(email) {
		return cli.Exit(colors.Red("invalid email address"), 1)
	}

	cfg.Profiles[profileName] = config.Profile{Name: name, Email: email}

	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Printf("Generating SSH keypair for profile %q...\n", profileName)
	if err := sshutil.GenerateKeypair(profileName, email); err != nil {
		return cli.Exit(fmt.Errorf("generate SSH keypair: %w", err), 1)
	}

	keyDir, err := sshutil.KeyDir(profileName)
	if err != nil {
		return cli.Exit(fmt.Errorf("resolve key dir: %w", err), 1)
	}
	pubKeyPath := keyDir + "/id_ed25519.pub"
	pubKeyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return cli.Exit(fmt.Errorf("read public key: %w", err), 1)
	}

	fmt.Println(colors.Green(fmt.Sprintf("Profile %q added.", profileName)))
	fmt.Printf("\n%s\n", colors.Cyan("Public key (add this to GitHub / GitLab):"))
	fmt.Print(string(pubKeyBytes))
	return nil
}

func validEmail(email string) bool {
	at := strings.Index(email, "@")
	if at < 1 {
		return false
	}
	dot := strings.LastIndex(email[at:], ".")
	return dot > 1 && at+dot < len(email)-1
}
