// Command Line Interface of Embedo.
package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/phogolabs/parcel"
	"github.com/urfave/cli"
)

const (
	ErrCodeArg = 101
)

func main() {
	app := &cli.App{
		Name:                 "parcel",
		HelpName:             "embedo",
		Usage:                "Golang Resource Embedding",
		UsageText:            "parcel [global options]",
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
				Name:  "resource-dir, d",
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
	dir, err := directory(ctx)
	if err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	pkg := ctx.String("package")
	if pkg == "" {
		pkg = filepath.Base(dir)
	}

	logger := ioutil.Discard
	if !ctx.Bool("quite") {
		logger = os.Stdout
	}

	generator := &parcel.Emitter{
		Logger:     logger,
		FileSystem: parcel.Dir(dir),
		Composer: &parcel.Generator{
			Config: &parcel.GeneratorConfig{
				InlcudeDocs: ctx.BoolT("include-docs"),
			},
		},
		Compressor: &parcel.TarGZipCompressor{
			Config: &parcel.CompressorConfig{
				Logger:         logger,
				Name:           pkg,
				IgnorePatterns: ctx.StringSlice("ignore"),
				Recurive:       ctx.Bool("recursive"),
			},
		},
	}

	if err := generator.Emit(); err != nil {
		return err
	}

	return nil
}

func directory(ctx *cli.Context) (string, error) {
	var err error
	dir := ctx.String("resource-dir")

	if dir == "" {
		if dir, err = os.Getwd(); err != nil {
			return "", err
		}
	}

	dir, err = filepath.Abs(dir)
	if err != nil {
		return "", err
	}

	return dir, nil
}
