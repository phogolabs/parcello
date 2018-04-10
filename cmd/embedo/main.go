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
				Name:  "package-dir, d",
				Usage: "Path to package directory",
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
	dir := ctx.String("package-dir")
	logger := ioutil.Discard

	if dir == "" {
		err := fmt.Errorf("Package directory is not provided")
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

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
			Recurive:    ctx.Bool("recursive"),
			InlcudeDocs: ctx.BoolT("include-docs"),
		},
	}

	if err := generator.Generate(filepath.Base(dir)); err != nil {
		return err
	}

	return nil
}
