package parcel

import (
	"bufio"
	"fmt"
	"io"
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
		fmt.Fprintln(w, "// File contains an embedded resources")
		fmt.Fprintln(w, "// Auto-generated at", time.Now().Format(time.UnixDate))
	}

	fmt.Fprintf(w, "package %s\n\n", bundle.Name)
	fmt.Fprintf(w, "import \"github.com/phogolabs/parcel\"\n\n")
	fmt.Fprintf(w, "func init() {\n")
	fmt.Fprintf(w, "  parcel.AddResource([]byte{")

	buffer := bufio.NewReader(bundle.Body)
	written := 0
	pattern := ""

	for {
		bit, err := buffer.ReadByte()
		if err == io.EOF {
			_, err = fmt.Fprintf(w, "})\n}\n")
			return nil
		}

		if written > 0 {
			pattern = ", %d"
		} else {
			pattern = "%d"
		}

		if _, err = fmt.Fprintf(w, pattern, bit); err != nil {
			return err
		}

		written = written + 1
	}
}
