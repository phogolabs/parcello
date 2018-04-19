package parcello

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var _ FileSystem = &Manager{}

// Manager represents a virtual in memory file system
type Manager struct {
	rw   sync.RWMutex
	root *Node
}

// Add adds resource to the manager
func (m *Manager) Add(binary Binary) error {
	m.rw.Lock()
	defer m.rw.Unlock()

	if m.root == nil {
		m.root = &Node{Name: "/", IsDir: true}
	}

	gzipper, err := gzip.NewReader(bytes.NewBuffer(binary))
	if err != nil {
		return err
	}

	reader := tar.NewReader(gzipper)
	return m.uncompress(reader)
}

func (m *Manager) uncompress(reader *tar.Reader) error {
	for {
		header, err := reader.Next()
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil
		}

		if err != nil {
			return err
		}

		path := split(header.Name)
		node := add(path, m.root)

		if node == m.root || node == nil {
			return fmt.Errorf("Invalid path: '%s'", header.Name)
		}

		content, _ := ioutil.ReadAll(reader)
		node.IsDir = false
		node.Content = &content
	}
}

// Root returns a sub-manager for given path
func (m *Manager) Root(name string) (*Manager, error) {
	if node := find(split(name), m.root); node != nil {
		if node.IsDir {
			return &Manager{root: node}, nil
		}
	}

	return nil, fmt.Errorf("Resource hierarchy not found")
}

// Open opens an embedded resource for read
func (m *Manager) Open(name string) (ReadOnlyFile, error) {
	return m.OpenFile(name, 0, 0)
}

// OpenFile is the generalized open call; most users will use Open
func (m *Manager) OpenFile(name string, flag int, perm os.FileMode) (File, error) {
	if node := find(split(name), m.root); node != nil {
		return NewResourceFile(node), nil
	}

	return nil, fmt.Errorf("File '%s' not found", name)
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root.
func (m *Manager) Walk(dir string, fn filepath.WalkFunc) error {
	if node := find(split(dir), m.root); node != nil {
		return walk(dir, node, fn)
	}

	return fmt.Errorf("Directory '%s' not found", dir)
}

func add(path []string, node *Node) *Node {
	if !node.IsDir || node.Content != nil {
		return nil
	}

	if len(path) == 0 {
		return node
	}

	name := path[0]

	for _, child := range node.Children {
		if child.Name == name {
			return add(path[1:], child)
		}
	}

	child := &Node{
		Mutex: &sync.RWMutex{},
		Name:  name,
		IsDir: true,
	}

	node.Children = append(node.Children, child)
	return add(path[1:], child)
}

func split(path string) []string {
	parts := []string{}

	for _, part := range strings.Split(path, string(os.PathSeparator)) {
		if part != "" && part != "/" {
			parts = append(parts, part)
		}
	}

	return parts
}

func find(path []string, node *Node) *Node {
	if len(path) == 0 {
		return node
	}

	for _, child := range node.Children {
		if path[0] == child.Name {
			if len(path) == 1 {
				return child
			}
			return find(path[1:], child)
		}
	}

	return nil
}

func walk(path string, node *Node, fn filepath.WalkFunc) error {
	if err := fn(path, &ResourceFileInfo{Node: node}, nil); err != nil {
		return err
	}

	for _, child := range node.Children {
		if err := walk(filepath.Join(path, child.Name), child, fn); err != nil {
			return err
		}
	}

	return nil
}
