package main

import (
	"context"
	"fmt"
	"os"

	initcmd "github.com/pawanprjl/gixy/cmd/init"
	"github.com/pawanprjl/gixy/cmd/profile"
	"github.com/urfave/cli/v3"
)

const version = "1.3.0"

func main() {
	app := &cli.Command{
		Name:    "gixy",
		Usage:   "A CLI companion for git",
		Version: version,
		Commands: []*cli.Command{
			{
				Name:  "profile",
				Usage: "Manage git profiles",
				Commands: []*cli.Command{
					&profile.AddCommand,
					&profile.ListCommand,
					&profile.ShowCommand,
					&profile.EditCommand,
					&profile.DeleteCommand,
					&profile.KeysCommand,
					&profile.GlobalCommand,
					&profile.MapGroupCommand,
					&profile.StatusCommand,
					&profile.ResolveCommand,
					&profile.UseCommand,
				},
			},
			&initcmd.InitCommand,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31merror:\033[0m %s\n", err)
		os.Exit(1)
	}
}
