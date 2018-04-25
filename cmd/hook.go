package cmd

import (
	"io"
	"io/ioutil"
	"os"

	"github.com/urfave/cli"
)

const (
	// ErrCodeArg is returned when an invalid argument is passed to CLI
	ErrCodeArg = 101
)

func logger(ctx *cli.Context) io.Writer {
	if ctx.GlobalBool("quiet") {
		return ioutil.Discard
	}

	return os.Stdout
}
