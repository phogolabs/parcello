package parcello

import (
	"bytes"
	"fmt"
	"io"
)

// Embedder embeds the resources to the provided package
type Embedder struct {
	// Logger prints each step of compression
	Logger io.Writer
	// Composer composes the resources
	Composer Composer
	// Compressor compresses the resources
	Compressor Compressor
	// FileSystem represents the underlying file system
	FileSystem FileSystem
}

// Embed embeds the resources to the provided package
func (e *Embedder) Embed() error {
	buffer := &bytes.Buffer{}

	ctx := &CompressorContext{
		Writer:     buffer,
		FileSystem: e.FileSystem,
	}

	info, err := e.Compressor.Compress(ctx)
	if err != nil {
		return err
	}

	if info == nil {
		return nil
	}

	fmt.Fprintf(e.Logger, "Embedding %d resource(s) at 'resource.go'\n", info.Count)
	err = e.Composer.Compose(&Bundle{Info: info, Body: buffer.Bytes()})
	return err
}
