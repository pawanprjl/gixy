package commit

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

var SetupCommand = cli.Command{
	Name:  "setup",
	Usage: "Configure the AI provider for commit message generation",
	Action: func(_ context.Context, _ *cli.Command) error {
		return runSetup()
	},
}

func runSetup() error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println(colors.Cyan("Available providers: gemini"))
	fmt.Print(colors.Cyan("Provider: "))
	provider, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read provider: %w", err), 1)
	}
	provider = strings.TrimSpace(provider)
	if provider == "" {
		return cli.Exit(colors.Red("provider cannot be empty"), 1)
	}
	if provider != "gemini" {
		return cli.Exit(colors.Red(fmt.Sprintf("unsupported provider %q; supported: gemini", provider)), 1)
	}

	fmt.Print(colors.Cyan("Model (e.g. gemini-2.0-flash): "))
	model, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read model: %w", err), 1)
	}
	model = strings.TrimSpace(model)
	if model == "" {
		return cli.Exit(colors.Red("model cannot be empty"), 1)
	}

	fmt.Print(colors.Cyan("API key: "))
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read api key: %w", err), 1)
	}
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return cli.Exit(colors.Red("api key cannot be empty"), 1)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	cfg.CommitGen = &config.CommitGenConfig{
		Provider: provider,
		Model:    model,
		APIKey:   apiKey,
	}

	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Println(colors.Green("Commit generation configured successfully."))
	return nil
}
