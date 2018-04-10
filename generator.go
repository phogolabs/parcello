package embedo

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/tools/imports"
)

// FileSystem provides with primitives to work with the underlying file system
type FileSystem interface {
	// Walk walks the file tree rooted at root, calling walkFn for each file or
	// directory in the tree, including root.
	Walk(dir string, fn filepath.WalkFunc) error
	// OpenFile is the generalized open call; most users will use Open
	OpenFile(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error)
}

// GeneratorConfig controls how the code generation happens
type GeneratorConfig struct {
	// Recurive enables embedding the resources recursively
	Recurive bool
	// InlcudeDocs determines whether to include documentation
	InlcudeDocs bool
}

// Generator generates an embedable resource
type Generator struct {
	// Logger prints each step of generation
	Logger io.Writer
	// FileSystem provides with primitives to work with the underlying file system
	FileSystem FileSystem
	// Config controls how the code generation happens
	Config *GeneratorConfig
}

// Generate generates an embedable resource for given directory
func (g *Generator) Generate(pkg string) error {
	if g.Logger == nil {
		g.Logger = ioutil.Discard
	}

	buffer := &bytes.Buffer{}

	if g.Config.InlcudeDocs {
		fmt.Fprintln(buffer, "// File contains an embedded resources")
		fmt.Fprintln(buffer, "// Auto-generated at", time.Now().Format(time.UnixDate))
	}

	fmt.Fprintf(buffer, "package %s", pkg)
	fmt.Fprintln(buffer)

	if g.Config.InlcudeDocs {
		fmt.Fprintln(buffer, "// ResourceManager contains the embeded resources for this package")
	}

	fmt.Fprintln(buffer, "var ResourceManager *embedo.ResourceManager")
	fmt.Fprintln(buffer)
	fmt.Fprintln(buffer, "func init() {")
	fmt.Fprintln(buffer, "  resource := &embedo.Resource{}")
	fmt.Fprintln(buffer)

	err := g.FileSystem.Walk("/", func(path string, info os.FileInfo, err error) error {
		if info == nil {
			return fmt.Errorf("Directory does not exist")
		}

		if info.IsDir() {
			if !g.Config.Recurive && path != "." {
				return filepath.SkipDir
			}
			return nil
		}

		matched, err := filepath.Match("*.go", info.Name())
		if err != nil || matched {
			return err
		}

		fmt.Fprintln(g.Logger, fmt.Sprintf("Embedding '%s'", path))

		file, err := g.FileSystem.OpenFile(path, os.O_RDONLY, 0)
		if err != nil {
			return err
		}

		defer func() {
			if ioErr := file.Close(); err == nil {
				err = ioErr
			}
		}()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return err
		}

		fmt.Fprintf(buffer, "  resource.Add(\"%s\", %s)", path, g.encode(data))
		fmt.Fprintln(buffer)
		return nil
	})

	fmt.Fprintln(buffer)
	fmt.Fprintln(buffer, "  ResourceManager = embedo.Open(resource)")
	fmt.Fprintln(buffer, "}")

	if err != nil {
		return err
	}

	if err := g.format(buffer); err != nil {
		return err
	}

	return g.write("resource.go", buffer.Bytes(), 0600)
}

func (g *Generator) encode(data []byte) string {
	array := []string{}

	for _, bit := range data {
		array = append(array, fmt.Sprintf("%d", bit))
	}

	return fmt.Sprintf("[]byte{%s}", strings.Join(array, ","))
}

func (g *Generator) format(buffer *bytes.Buffer) error {
	data, err := imports.Process("resource", buffer.Bytes(), nil)
	if err != nil {
		return err
	}

	data, err = format.Source(data)
	if err != nil {
		return err
	}

	buffer.Reset()

	_, err = buffer.Write(data)
	return err
}

func (g *Generator) write(filename string, data []byte, perm os.FileMode) error {
	f, err := g.FileSystem.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return err
	}
	n, err := f.Write(data)
	if err == nil && n < len(data) {
		err = io.ErrShortWrite
	}
	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
