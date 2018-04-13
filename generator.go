package parcel

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"time"
)

var _ Composer = &Generator{}

// GeneratorConfig controls how the code generation happens
type GeneratorConfig struct {
	// InlcudeDocs determines whether to include documentation
	InlcudeDocs bool
}

// Generator generates an embedable resource
type Generator struct {
	// Config controls how the code generation happens
	Config *GeneratorConfig
}

// Generate generates an embedable resource for given directory
func (g *Generator) WriteTo(w io.Writer, bundle *Bundle) error {
	if g.Config.InlcudeDocs {
		fmt.Fprintln(w, "// Package", bundle.Name, "contains embedded resources")
		fmt.Fprintln(w, "// Auto-generated at", time.Now().Format(time.UnixDate))
	}

	fmt.Fprintln(w, "package", bundle.Name)
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

	_, err := io.Copy(w, buffer)
	return err
}
