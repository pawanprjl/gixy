package commitconfig

import (
	"context"
	"fmt"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var supportedProviders = []string{"gemini"}

var AddCommand = cli.Command{
	Name:      "add",
	Usage:     "Add a new commit generation provider",
	ArgsUsage: "<name>",
	Flags: []cli.Flag{
		&cli.StringFlag{Name: "provider", Usage: "AI provider (gemini)", Required: true},
		&cli.StringFlag{Name: "model", Usage: "Model name (e.g. gemini-2.0-flash)", Required: true},
		&cli.StringFlag{Name: "api-key", Usage: "API key for the provider", Required: true},
	},
	Action: func(_ context.Context, cmd *cli.Command) error {
		if cmd.Args().Len() != 1 {
			return cli.Exit(colors.Red("usage: gixy commit config add <name> --provider <p> --model <m> --api-key <k>"), 1)
		}
		name := cmd.Args().Get(0)
		provider := cmd.String("provider")
		model := cmd.String("model")
		apiKey := cmd.String("api-key")

		if provider != "gemini" {
			return cli.Exit(colors.Red(fmt.Sprintf("unsupported provider %q; supported: gemini", provider)), 1)
		}

		cfg, err := config.LoadConfig()
		if err != nil {
			return cli.Exit(fmt.Errorf("load config: %w", err), 1)
		}

		if cfg.CommitGen == nil {
			cfg.CommitGen = &config.CommitGenConfig{
				Providers: make(map[string]config.CommitGenEntry),
			}
		}

		if _, exists := cfg.CommitGen.Providers[name]; exists {
			return cli.Exit(colors.Red(fmt.Sprintf("provider %q already exists", name)), 1)
		}

		cfg.CommitGen.Providers[name] = config.CommitGenEntry{
			Provider: provider,
			Model:    model,
			APIKey:   apiKey,
		}

		// Auto-activate if this is the first provider.
		if cfg.CommitGen.Active == "" {
			cfg.CommitGen.Active = name
		}

		if err := config.SaveConfig(cfg); err != nil {
			return cli.Exit(fmt.Errorf("save config: %w", err), 1)
		}

		fmt.Println(colors.Green(fmt.Sprintf("Provider %q added.", name)))
		if cfg.CommitGen.Active == name {
			fmt.Println(colors.Cyan(fmt.Sprintf("Active provider set to %q.", name)))
		}
		return nil
	},
}
