// Command Line Interface of Embedo.
package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"

	"github.com/phogolabs/parcello"
	"github.com/urfave/cli"
)

const (
	// ErrCodeArg is returned when an invalid argument is passed to CLI
	ErrCodeArg = 101
)

func main() {
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
		Action:               run,
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
			cli.BoolTFlag{
				Name:  "include-docs",
				Usage: "Include API documentation in generated source code",
			},
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	app.Run(os.Args)
}

func run(ctx *cli.Context) error {
	resourceDir, err := filepath.Abs(ctx.String("resource-dir"))
	if err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	bundleDir, err := filepath.Abs(ctx.String("bundle-dir"))
	if err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	_, packageName := filepath.Split(bundleDir)

	generator := &parcello.Emitter{
		Logger:     logger(ctx),
		FileSystem: parcello.Dir(resourceDir),
		Composer: &parcello.Generator{
			FileSystem: parcello.Dir(bundleDir),
			Config: &parcello.GeneratorConfig{
				Package:     packageName,
				InlcudeDocs: ctx.BoolT("include-docs"),
			},
		},
		Compressor: &parcello.TarGZipCompressor{
			Config: &parcello.CompressorConfig{
				Logger:         logger(ctx),
				Filename:       "resource",
				IgnorePatterns: ctx.StringSlice("ignore"),
				Recurive:       ctx.Bool("recursive"),
			},
		},
	}

	if err := generator.Emit(); err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	return nil
}

func logger(ctx *cli.Context) io.Writer {
	if ctx.Bool("quiet") {
		return ioutil.Discard
	}

	return os.Stdout
}
