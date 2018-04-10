package embedo

import (
	"bytes"
	"os"
	"time"
)

// Binary represents a resource blob content
type Binary = []byte

// Node represents a hierarchy node in the resource manager
type Node struct {
	name     string
	dir      bool
	content  Binary
	children []*Node
}

// Name returns the base name of the file
func (n *Node) Name() string {
	return n.name
}

// Size is the length in bytes for regular files
func (n *Node) Size() int64 {
	return int64(len(n.content))
}

// FileMode are the file mode bits
func (n *Node) Mode() os.FileMode {
	return 0
}

// ModTime returns the modification time
func (n *Node) ModTime() time.Time {
	return time.Now()
}

// IsDir returns true if the node is directory
func (n *Node) IsDir() bool {
	return n.dir
}

// Underlying data source
func (n *Node) Sys() interface{} {
	return nil
}

// BufferCloser represents a *bytes.Buffer that can be closed
type BufferCloser struct {
	buffer *bytes.Buffer
}

// NewBufferCloser creates a new BufferCloser
func NewBufferCloser(data []byte) *BufferCloser {
	return &BufferCloser{
		buffer: bytes.NewBuffer(data),
	}
}

// Read reads the next len(p) bytes from the buffer or until the buffer is drainged
func (b *BufferCloser) Read(p []byte) (n int, err error) {
	return b.buffer.Read(p)
}

// Write appends the contents of p to the buffer, growing the buffer as needed.
func (b *BufferCloser) Write(p []byte) (n int, err error) {
	return b.buffer.Write(p)
}

// Close closes the buffer (noop).
func (buff *BufferCloser) Close() error {
	return nil
}

type item struct {
	key  string
	data Binary
}
