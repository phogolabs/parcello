package parcel

import (
	"fmt"
	"io"
)

// Emitter emits the resources to the provided package
type Emitter struct {
	// Logger prints each step of compression
	Logger io.Writer
	// Composer composes the resources
	Composer Composer
	// Compressor compresses the resources
	Compressor Compressor
	// FileSystem represents the underlying file system
	FileSystem FileSystem
}

// Emit emits the resources to the provided package
func (e *Emitter) Emit() error {
	bundle, err := e.Compressor.Compress(e.FileSystem)
	if err != nil {
		return err
	}

	if bundle == nil {
		return nil
	}

	fmt.Fprintf(e.Logger, "Bundling %d resource(s) at 'resource.go'\n", bundle.Length)
	err = e.Composer.Compose(bundle)
	return err
}
