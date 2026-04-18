package profile

import (
	"context"
	"fmt"
	"sort"

	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var ListCommand = cli.Command{
	Name:   "list",
	Usage:  "List all profiles",
	Action: listProfiles,
}

func listProfiles(_ context.Context, _ *cli.Command) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return cli.Exit(fmt.Errorf("load config: %w", err), 1)
	}

	if len(cfg.Profiles) == 0 {
		fmt.Println("\033[33mNo profiles found.\033[0m")
		return nil
	}

	keys := make([]string, 0, len(cfg.Profiles))
	for k := range cfg.Profiles {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		p := cfg.Profiles[k]
		fmt.Printf("%-20s %s <%s>\n", k, p.Name, p.Email)
	}
	return nil
}
