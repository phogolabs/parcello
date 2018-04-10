// Command Line Interface of Embedo.
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/phogolabs/embedo"
	"github.com/urfave/cli"
)

const (
	ErrCodeArg = 101
)

func main() {
	app := &cli.App{
		Name:                 "embedder",
		HelpName:             "embedder",
		Usage:                "Golang Resource Embedder",
		UsageText:            "embedder [global options]",
		Version:              "0.1",
		BashComplete:         cli.DefaultAppComplete,
		EnableBashCompletion: true,
		Writer:               os.Stdout,
		ErrWriter:            os.Stderr,
		Action:               run,
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "quite, q",
				Usage: "Disable logging",
			},
			cli.BoolFlag{
				Name:  "recursive, r",
				Usage: "Embed the resources recursively",
			},
			cli.StringFlag{
				Name:  "dir, d",
				Usage: "Path to directory",
			},
			cli.StringFlag{
				Name:  "package, pkg",
				Usage: "Package name",
			},
			cli.StringSliceFlag{
				Name:  "ignore, i",
				Usage: "Ignore file name",
			},
			cli.BoolTFlag{
				Name:  "include-docs",
				Usage: "Include API documentation in generated source code",
			},
		},
	}

	app.Run(os.Args)
}

func run(ctx *cli.Context) error {
	dir := ctx.String("dir")
	if dir == "" {
		err := fmt.Errorf("Directory is not provided")
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	pkg := ctx.String("pkg")
	if pkg == "" {
		err := fmt.Errorf("Package name is not provided")
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	logger := ioutil.Discard
	if !ctx.Bool("quite") {
		logger = os.Stdout
	}

	dir, err := filepath.Abs(dir)
	if err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	generator := &embedo.Generator{
		Logger:     logger,
		FileSystem: embedo.Dir(dir),
		Config: &embedo.GeneratorConfig{
			IgnorePatterns: ctx.StringSlice("ignore"),
			Recurive:       ctx.Bool("recursive"),
			InlcudeDocs:    ctx.BoolT("include-docs"),
		},
	}

	if err := generator.Generate(pkg); err != nil {
		return err
	}

	return nil
}
