package parcel

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

var _ Composer = &Generator{}

// GeneratorConfig controls how the code generation happens
type GeneratorConfig struct {
	// Package determines the name of the package
	Package string
	// InlcudeDocs determines whether to include documentation
	InlcudeDocs bool
}

// Generator generates an embedable resource
type Generator struct {
	// FileSystem represents the underlying file system
	FileSystem FileSystem
	// Config controls how the code generation happens
	Config *GeneratorConfig
}

// Generate generates an embedable resource for given directory
func (g *Generator) Compose(bundle *Bundle) error {
	filename := fmt.Sprintf("%s.go", bundle.Name)

	w, err := g.FileSystem.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	defer func() {
		if ioErr := w.Close(); err == nil {
			err = ioErr
		}
	}()

	if g.Config.InlcudeDocs {
		fmt.Fprintln(w, "// Package", g.Config.Package, "contains embedded resources")
		fmt.Fprintln(w, "// Auto-generated at", time.Now().Format(time.UnixDate))
	}

	fmt.Fprintln(w, "package", g.Config.Package)
	fmt.Fprintln(w)
	fmt.Fprintf(w, `import "github.com/phogolabs/parcel"`)
	fmt.Fprintln(w)
	fmt.Fprintln(w)
	fmt.Fprintln(w, "func init() {")
	fmt.Fprintln(w, "\tparcel.AddResource([]byte{")

	reader := bufio.NewReader(bundle.Body)
	buffer := &bytes.Buffer{}

	for {
		bit, err := reader.ReadByte()
		if err == io.EOF {
			break
		}

		if buffer.Len() == 0 {
			fmt.Fprint(w, "\t\t")
		}

		fmt.Fprintf(buffer, "%d, ", int(bit))

		if buffer.Len() >= 60 {
			line := strings.TrimSpace(buffer.String())
			if _, err := fmt.Fprintln(w, line); err != nil {
				return err
			}

			buffer.Reset()
			continue
		}
	}

	if ln := buffer.Len(); ln > 0 && ln < 70 {
		fmt.Fprintln(buffer)
	}

	line := strings.TrimSpace(buffer.String())

	buffer.Reset()

	fmt.Fprintln(buffer, line)
	fmt.Fprintln(buffer, "\t})")
	fmt.Fprintln(buffer, "}")

	_, err = io.Copy(w, buffer)
	return err
}
