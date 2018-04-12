package parcel

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

var _ Compressor = &TarGZipCompressor{}

// CompressorConfig controls how the code generation happens
type CompressorConfig struct {
	// Logger prints each step of compression
	Logger io.Writer
	// Name of the compressed bundle
	Name string
	// IgnorePatterns provides a list of all files that has to be ignored
	IgnorePatterns []string
	// Recurive enables embedding the resources recursively
	Recurive bool
}

// TarGZipCompressor compresses content as GZip tarball
type TarGZipCompressor struct {
	// Config controls how the compression is made
	Config *CompressorConfig
}

// Compress compresses given source in tar.gz
func (c *TarGZipCompressor) Compress(fileSystem FileSystem) (*Bundle, error) {
	archive, err := ioutil.TempFile("", "parcel")
	if err != nil {
		return nil, err
	}

	processed, err := c.writeTo(fileSystem, archive)
	if err != nil {
		return nil, err
	}

	if _, err = archive.Seek(0, os.SEEK_SET); err != nil {
		archive.Close()
		return nil, err
	}

	bundle := &Bundle{
		Name:   c.Config.Name,
		Length: processed,
		Body:   archive,
	}

	return bundle, nil
}

func (e *TarGZipCompressor) writeTo(fileSystem FileSystem, bundle io.Writer) (int, error) {
	compressor := gzip.NewWriter(bundle)
	bundler := tar.NewWriter(compressor)
	processed := 0

	err := fileSystem.Walk("/", func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return nil
		}

		if info.IsDir() {
			if !e.Config.Recurive && path != "." {
				return filepath.SkipDir
			}
			return nil
		}

		ignore := append(e.Config.IgnorePatterns, "*.go")
		for _, pattern := range ignore {
			matched, err := filepath.Match(pattern, info.Name())
			if err != nil || matched {
				return err
			}
		}

		fmt.Fprintln(e.Config.Logger, fmt.Sprintf("Compressing '%s'", path))

		header, err := tar.FileInfoHeader(info, path)
		if err != nil {
			return nil
		}

		header.Name = path
		if err := bundler.WriteHeader(header); err != nil {
			return nil
		}

		resource, err := fileSystem.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			return err
		}

		defer func() {
			if ioErr := resource.Close(); err == nil {
				err = ioErr
			}
		}()

		if _, err = io.Copy(bundler, resource); err != nil {
			return err
		}

		if err = bundler.Flush(); err != nil {
			return err
		}

		if err = compressor.Flush(); err != nil {
			return err
		}

		processed = processed + 1
		return err
	})

	if ioErr := bundler.Close(); err == nil {
		err = ioErr
	}

	return processed, err
}
