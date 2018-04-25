// Command Line Interface of Embedo.
package main

import (
	"os"
	"sort"

	"github.com/phogolabs/parcello/cmd"
	"github.com/urfave/cli"
)

func main() {
	generator := &cmd.ResourceGenerator{}

	commands := []cli.Command{
		generator.CreateCommand(),
	}

	app := &cli.App{
		Name:                 "parcello",
		HelpName:             "parcello",
		Usage:                "Golang Resource Bundler",
		UsageText:            "parcello [global options]",
		Version:              "0.7",
		BashComplete:         cli.DefaultAppComplete,
		EnableBashCompletion: true,
		Writer:               os.Stdout,
		ErrWriter:            os.Stderr,
		Commands:             commands,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "quiet, q",
				Usage: "Disable logging",
			},
			cli.BoolFlag{
				Name:  "recursive, r",
				Usage: "Embed the resources recursively",
			},
			cli.StringFlag{
				Name:  "resource-dir, d",
				Usage: "Path to directory",
				Value: ".",
			},
			cli.StringFlag{
				Name:  "bundle-dir, b",
				Usage: "Path to bundle directory",
				Value: ".",
			},
			cli.StringSliceFlag{
				Name:  "ignore, i",
				Usage: "Ignore file name",
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	for _, command := range commands {
		sort.Sort(cli.FlagsByName(command.Flags))
		sort.Sort(cli.CommandsByName(command.Subcommands))
	}

	app.Run(os.Args)
}
