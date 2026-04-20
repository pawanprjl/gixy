package provider

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/pawanprjl/gixy/internal/colors"
	"github.com/pawanprjl/gixy/internal/config"
	"github.com/urfave/cli/v3"
)

var ListCommand = cli.Command{
	Name:  "list",
	Usage: "List configured providers and optionally switch the active one",
	Action: func(_ context.Context, _ *cli.Command) error {
		return runList()
	},
}

func runList() error {
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
	for i, k := range keys {
		entry := cfg.CommitGen.Providers[k]
		marker := "  "
		nameStr := colors.Cyan(k)
		if k == cfg.CommitGen.Active {
			marker = colors.Green("*")
			nameStr = colors.Green(k)
		}
		host := ""
		if entry.Host != "" {
			host = " (" + entry.Host + ")"
		}
		fmt.Printf("  %s %d) %-20s %-10s %s%s\n", marker, i+1, nameStr, entry.Provider, entry.Model, host)
	}
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s (Enter to keep current, or number to switch): ", colors.Cyan("Select provider"))
	line, err := reader.ReadString('\n')
	if err != nil {
		return cli.Exit(fmt.Errorf("read input: %w", err), 1)
	}
	line = strings.TrimSpace(line)

	if line == "" {
		fmt.Println(colors.Cyan(fmt.Sprintf("Keeping %q as active provider.", cfg.CommitGen.Active)))
		return nil
	}

	idx, err := strconv.Atoi(line)
	if err != nil || idx < 1 || idx > len(keys) {
		return cli.Exit(colors.Red(fmt.Sprintf("invalid choice %q; enter a number between 1 and %d", line, len(keys))), 1)
	}

	chosen := keys[idx-1]
	if chosen == cfg.CommitGen.Active {
		fmt.Println(colors.Cyan(fmt.Sprintf("%q is already the active provider.", chosen)))
		return nil
	}

	cfg.CommitGen.Active = chosen
	if err := config.SaveConfig(cfg); err != nil {
		return cli.Exit(fmt.Errorf("save config: %w", err), 1)
	}

	fmt.Println(colors.Green(fmt.Sprintf("Active provider switched to %q.", chosen)))
	return nil
}
