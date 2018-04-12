package parcel

import (
	"os"
)

// Emitter emits the resources to the provided package
type Emitter struct {
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

	defer func() {
		if ioErr := bundle.Body.Close(); err == nil {
			err = ioErr
		}
	}()

	if bundle.Length == 0 {
		return nil
	}

	resource, err := e.FileSystem.OpenFile("resource.go", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer func() {
		if ioErr := resource.Close(); err == nil {
			err = ioErr
		}
	}()

	err = e.Composer.WriteTo(resource, bundle)
	return err
}
