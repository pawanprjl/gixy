package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pawanprjl/gixy/cmd/commit"
	"github.com/pawanprjl/gixy/cmd/profile"
	gixyprovider "github.com/pawanprjl/gixy/cmd/provider"
	"github.com/urfave/cli/v3"
)

const version = "1.1.0"

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
					&profile.DeleteCommand,
					&profile.UseCommand,
					&profile.ShowCommand,
					&profile.EditCommand,
				},
			},
			{
				Name:  "commit",
				Usage: "AI-powered commit message generation",
				Commands: []*cli.Command{
					&commit.GenerateCommand,
				},
			},
			{
				Name:  "provider",
				Usage: "Manage AI providers",
				Commands: []*cli.Command{
					&gixyprovider.AddCommand,
					&gixyprovider.ListCommand,
					&gixyprovider.UseCommand,
					&gixyprovider.RemoveCommand,
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31merror:\033[0m %s\n", err)
		os.Exit(1)
	}
}
