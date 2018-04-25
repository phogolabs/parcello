package cmd

import (
	"path/filepath"

	"github.com/phogolabs/parcello"
	"github.com/urfave/cli"
)

// ResourceEmbedder is a command that generate compilable resources
type ResourceEmbedder struct {
	emitter *parcello.Emitter
}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (r *ResourceEmbedder) CreateCommand() cli.Command {
	return cli.Command{
		Name:         "embed",
		Usage:        "Generates compilable resources",
		Description:  "Generates compilable resources",
		BashComplete: cli.DefaultAppComplete,
		Before:       r.before,
		Action:       r.generate,
		Flags: []cli.Flag{
			cli.BoolTFlag{
				Name:  "include-docs",
				Usage: "Include API documentation in generated source code",
			},
		},
	}
}

func (r *ResourceEmbedder) before(ctx *cli.Context) error {
	resourceDir, err := filepath.Abs(ctx.GlobalString("resource-dir"))
	if err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	bundleDir, err := filepath.Abs(ctx.GlobalString("bundle-dir"))
	if err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	_, packageName := filepath.Split(bundleDir)

	r.emitter = &parcello.Emitter{
		Logger:     logger(ctx),
		FileSystem: parcello.Dir(resourceDir),
		Composer: &parcello.Generator{
			FileSystem: parcello.Dir(bundleDir),
			Config: &parcello.GeneratorConfig{
				Package:     packageName,
				InlcudeDocs: ctx.BoolT("include-docs"),
			},
		},
		Compressor: &parcello.ZipCompressor{
			Config: &parcello.CompressorConfig{
				Logger:         logger(ctx),
				Filename:       "resource",
				IgnorePatterns: ctx.GlobalStringSlice("ignore"),
				Recurive:       ctx.GlobalBool("recursive"),
			},
		},
	}

	return nil
}

func (r *ResourceEmbedder) generate(ctx *cli.Context) error {
	if err := r.emitter.Emit(); err != nil {
		return cli.NewExitError(err.Error(), ErrCodeArg)
	}

	return nil
}
