package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pawanprjl/gixy/cmd/branch"
	"github.com/pawanprjl/gixy/cmd/commit"
	commitconfig "github.com/pawanprjl/gixy/cmd/commit/config"
	"github.com/pawanprjl/gixy/cmd/gixylog"
	"github.com/pawanprjl/gixy/cmd/profile"
	gixyprovider "github.com/pawanprjl/gixy/cmd/provider"
	"github.com/pawanprjl/gixy/cmd/status"
	"github.com/urfave/cli/v3"
)

const version = "0.1.0"

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
					{
						Name:  "config",
						Usage: "Manage commit generation providers",
						Commands: []*cli.Command{
							&commitconfig.AddCommand,
							&commitconfig.UseCommand,
							&commitconfig.ListCommand,
							&commitconfig.RemoveCommand,
						},
					},
					&commit.GenerateCommand,
				},
			},
			{
				Name:  "provider",
				Usage: "Set up an AI provider for commit generation",
				Commands: []*cli.Command{
					&gixyprovider.AddCommand,
					&gixyprovider.ListCommand,
					&gixyprovider.RemoveCommand,
				},
			},
			&branch.BranchCommand,
			&gixylog.LogCommand,
			&status.StatusCommand,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31merror:\033[0m %s\n", err)
		os.Exit(1)
	}
}
