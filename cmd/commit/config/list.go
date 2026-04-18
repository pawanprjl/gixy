package commitconfig

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
	Usage: "List all configured commit generation providers",
	Action: func(_ context.Context, _ *cli.Command) error {
		cfg, err := config.LoadConfig()
		if err != nil {
			return cli.Exit(fmt.Errorf("load config: %w", err), 1)
		}

		if cfg.CommitGen == nil || len(cfg.CommitGen.Providers) == 0 {
			fmt.Println(colors.Yellow("No commit providers configured."))
			return nil
		}

		keys := make([]string, 0, len(cfg.CommitGen.Providers))
		for k := range cfg.CommitGen.Providers {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			entry := cfg.CommitGen.Providers[k]
			marker := "  "
			name := colors.Cyan(k)
			if k == cfg.CommitGen.Active {
				marker = colors.Green("* ")
			}
			fmt.Printf("%s%-20s %-10s %s\n", marker, name, entry.Provider, entry.Model)
		}
		return nil
	},
}
