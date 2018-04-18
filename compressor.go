package parcello

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var _ Compressor = &TarGZipCompressor{}

// ErrSkipResource skips a particular file from processing
var ErrSkipResource = fmt.Errorf("Skip Resource Error")

// CompressorConfig controls how the code generation happens
type CompressorConfig struct {
	// Logger prints each step of compression
	Logger io.Writer
	// Filename is the name of the compressed bundle
	Filename string
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
func (e *TarGZipCompressor) Compress(fileSystem FileSystem) (*Bundle, error) {
	buffer := &bytes.Buffer{}
	count, err := e.write(fileSystem, buffer)

	if err != nil {
		return nil, err
	}

	if count == 0 {
		return nil, nil
	}

	bundle := &Bundle{
		Name:   e.Config.Filename,
		Length: count,
		Body:   buffer.Bytes(),
	}

	return bundle, nil
}

func (e *TarGZipCompressor) write(fileSystem FileSystem, bundle io.Writer) (int, error) {
	compressor := gzip.NewWriter(bundle)
	bundler := tar.NewWriter(compressor)
	count := 0

	err := fileSystem.Walk("/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		err = e.filter(path, info)

		switch err {
		case ErrSkipResource:
			return nil
		default:
			if err != nil {
				return err
			}
		}

		if err = e.walk(bundler, fileSystem, path, info); err != nil {
			return err
		}

		count = count + 1
		return nil
	})

	if err != nil {
		return count, err
	}

	bundler.Flush()
	compressor.Flush()

	if ioErr := bundler.Close(); err == nil {
		err = ioErr
	}

	return count, err
}

func (e *TarGZipCompressor) walk(bundler *tar.Writer, fileSystem FileSystem, path string, info os.FileInfo) error {
	fmt.Fprintln(e.Config.Logger, fmt.Sprintf("Compressing '%s'", path))

	header, err := tar.FileInfoHeader(info, path)
	if err != nil {
		return err
	}

	header.Name = path
	if err = bundler.WriteHeader(header); err != nil {
		return err
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

	_, err = io.Copy(bundler, resource)
	return err
}

func (e *TarGZipCompressor) filter(path string, info os.FileInfo) error {
	if info == nil {
		return ErrSkipResource
	}

	if err := e.ignore(path, info); err != nil {
		return err
	}

	if !info.IsDir() {
		return nil
	}

	if !e.Config.Recurive && path != "." {
		return filepath.SkipDir
	}

	return ErrSkipResource
}

func (e *TarGZipCompressor) ignore(path string, info os.FileInfo) error {
	ignore := append(e.Config.IgnorePatterns, "*.go")

	for _, pattern := range ignore {
		matched, err := filepath.Match(pattern, path)
		if err != nil {
			return err
		}

		if !matched {
			matched, err = filepath.Match(pattern, info.Name())
			if err != nil {
				return err
			}

			if !matched {
				continue
			}
		}

		if info.IsDir() {
			return filepath.SkipDir
		}

		return ErrSkipResource
	}

	return nil
}
