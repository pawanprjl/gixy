package main

import (
	"context"
	"fmt"
	"os"

	"github.com/pawanprjl/gixy/cmd/profile"
	"github.com/urfave/cli/v3"
)

func main() {
	app := &cli.Command{
		Name:  "gixy",
		Usage: "A CLI companion for git",
		Commands: []*cli.Command{
			{
				Name:  "profile",
				Usage: "Manage git profiles",
				Commands: []*cli.Command{
					&profile.AddCommand,
					&profile.ListCommand,
					&profile.DeleteCommand,
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31merror:\033[0m %s\n", err)
		os.Exit(1)
	}
}
