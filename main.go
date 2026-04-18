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
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
