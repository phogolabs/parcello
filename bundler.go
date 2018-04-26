package parcello

import (
	"fmt"
	"io"
	"os"
)

// BundlerContext the context of this bundler
type BundlerContext struct {
	// Name of the binary
	Name string
	// FileSystem represents the underlying file system
	FileSystem FileSystem
}

// Bundler bundles the resources to the provided binary
type Bundler struct {
	// Logger prints each step of compression
	Logger io.Writer
	// Compressor compresses the resources
	Compressor Compressor
	// FileSystem represents the underlying file system
	FileSystem FileSystem
}

// Bundle bundles the resources to the provided binary
func (e *Bundler) Bundle(ctx *BundlerContext) error {
	fmt.Fprintf(e.Logger, "Bundling resource(s) at '%s'", ctx.Name)

	file, err := ctx.FileSystem.OpenFile(ctx.Name, os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer file.Close()

	cctx := &CompressorContext{
		FileSystem: e.FileSystem,
		Writer:     file,
		Offset:     0,
	}

	finfo, ferr := file.Stat()
	if ferr != nil {
		return ferr
	}

	if finfo.IsDir() {
		return fmt.Errorf("'%s' is not a regular file", ctx.Name)
	}

	cctx.Offset = finfo.Size()
	_, err = e.Compressor.Compress(cctx)
	return err
}
