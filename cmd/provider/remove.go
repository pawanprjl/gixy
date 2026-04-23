package provider

import (
"context"
"fmt"

"github.com/pawanprjl/gixy/internal/colors"
"github.com/pawanprjl/gixy/internal/config"
"github.com/urfave/cli/v3"
)

var RemoveCommand = cli.Command{
Name:      "remove",
Aliases:   []string{"delete"},
Usage:     "Remove a configured AI provider",
ArgsUsage: "<name>",
Action: func(_ context.Context, cmd *cli.Command) error {
if cmd.Args().Len() != 1 {
return cli.Exit(colors.Red("usage: gixy provider remove <name>"), 1)
}
name := cmd.Args().Get(0)

cfg, err := config.LoadConfig()
if err != nil {
return cli.Exit(fmt.Errorf("load config: %w", err), 1)
}

if cfg.CommitGen == nil {
return cli.Exit(colors.Red(fmt.Sprintf("provider %q not found", name)), 1)
}
if _, exists := cfg.CommitGen.Providers[name]; !exists {
return cli.Exit(colors.Red(fmt.Sprintf("provider %q not found", name)), 1)
}

delete(cfg.CommitGen.Providers, name)
if cfg.CommitGen.Active == name {
cfg.CommitGen.Active = ""
}

if err := config.SaveConfig(cfg); err != nil {
return cli.Exit(fmt.Errorf("save config: %w", err), 1)
}

fmt.Println(colors.Green(fmt.Sprintf("Provider %q removed.", name)))
if cfg.CommitGen.Active == "" && len(cfg.CommitGen.Providers) > 0 {
fmt.Println(colors.Yellow("No active provider set. Run `gixy provider use <name>` to set one."))
}
return nil
},
}
