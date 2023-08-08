// Package app is the main package for the application.
package app

import (
	"fmt"
	"os"

	"git.sr.ht/~jamesponddotco/shrinkimages/cmd/shrinkctl/internal/meta"
	"github.com/urfave/cli/v2"
)

// Run is the entry point for the application.
func Run() int {
	app := cli.NewApp()
	app.Name = meta.Name
	app.Version = meta.Version
	app.Usage = meta.Description
	app.HideHelpCommand = true

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "config",
			Aliases: []string{"c"},
			Usage:   "path to configuration file",
			Value:   "config.json",
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:   "start",
			Usage:  "start the server for the Shrink Images service",
			Action: StartAction,
		},
		{
			Name:   "stop",
			Usage:  "stop the server for the Shrink Images service",
			Action: StopAction,
		},
		{
			Name:   "generate-key",
			Usage:  "generate a new API key for the Shrink Images service",
			Action: GenerateKeyAction,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)

		return 1
	}

	return 0
}
