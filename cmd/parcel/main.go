// Command Line Interface of Embedo.
package main

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/phogolabs/parcel"
	"github.com/urfave/cli"
)

const (
	// ErrCodeArg is returned when an invalid argument is passed to CLI
	ErrCodeArg = 101
)

func main() {
	app := &cli.App{
		Name:                 "parcel",
		HelpName:             "parcel",
		Usage:                "Golang Resource Bundler",
		UsageText:            "parcel [global options]",
		Version:              "0.2",
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

	generator := &parcel.Emitter{
		Logger:     logger(ctx),
		FileSystem: parcel.Dir(resourceDir),
		Composer: &parcel.Generator{
			FileSystem: parcel.Dir(bundleDir),
			Config: &parcel.GeneratorConfig{
				Package:     packageName,
				InlcudeDocs: ctx.BoolT("include-docs"),
			},
		},
		Compressor: &parcel.TarGZipCompressor{
			Config: &parcel.CompressorConfig{
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
