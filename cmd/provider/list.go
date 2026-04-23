package provider

import (
	"context"
	"fmt"
	"sort"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var ListCommand = cli.Command{
	Name:  "list",
	Usage: "List all configured AI providers",
	Action: func(_ context.Context, _ *cli.Command) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return cli.Exit(fmt.Errorf("load config: %w", err), 1)
		}

		if cfg.CommitGen == nil || len(cfg.CommitGen.Providers) == 0 {
			fmt.Println(colors.Yellow("No providers configured. Run `gixy provider add` to add one."))
			return nil
		}

		keys := make([]string, 0, len(cfg.CommitGen.Providers))
		for k := range cfg.CommitGen.Providers {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		fmt.Println()
		fmt.Printf("    %-15s %-15s %-25s %s\n", "NAME", "PROVIDER", "MODEL", "HOST")
		fmt.Println()
		for _, k := range keys {
			entry := cfg.CommitGen.Providers[k]
			marker := " "
			paddedName := fmt.Sprintf("%-15s", k)
			nameStr := colors.Cyan(paddedName)
			if k == cfg.CommitGen.Active {
				marker = colors.Green("*")
				nameStr = colors.Green(paddedName)
			}
			fmt.Printf("  %s %s %-15s %-25s %s\n", marker, nameStr, entry.Provider, entry.Model, entry.Host)
		}
		fmt.Println()
		return nil
	},
}
