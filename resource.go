package embedo

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Resource represents a single resource
type Resource struct {
	items map[string][]byte
}

// Add data to the resource
func (r *Resource) Add(key string, data []byte) {
	if r.items == nil {
		r.items = make(map[string][]byte)
	}
	r.items[key] = data
}

// ResourceManager represents a virtual in memory file system
type ResourceManager struct {
	root *Node
}

// Open creates an instance of ResourceManager
func Open(resouce *Resource) *ResourceManager {
	root := &Node{name: "/", dir: true}

	for key, data := range resouce.items {
		path := split(key)
		add(path, data, root)
	}

	return &ResourceManager{root: root}
}

// Open opens an embeded resource for read
func (fs *ResourceManager) Open(name string) (io.Reader, error) {
	return fs.OpenFile(name, 0, 0)
}

// OpenFile is the generalized open call; most users will use Open
func (fs *ResourceManager) OpenFile(name string, flag int, perm os.FileMode) (io.ReadWriteCloser, error) {
	if node := find(split(name), fs.root); node != nil {
		if node.dir {
			return nil, fmt.Errorf("Cannot open directory '%s'", name)
		}
		return NewBufferCloser(node.content), nil
	}

	return nil, fmt.Errorf("File '%s' not found", name)
}

// Walk walks the file tree rooted at root, calling walkFn for each file or
// directory in the tree, including root.
func (fs *ResourceManager) Walk(dir string, fn filepath.WalkFunc) error {
	if node := find(split(dir), fs.root); node != nil {
		return walk(dir, node, fn)
	}

	return fmt.Errorf("Directory '%s' not found", dir)
}

func add(path []string, content []byte, node *Node) {
	if len(path) == 0 {
		return
	}

	name := path[0]

	for _, child := range node.children {
		if child.name == name {
			add(path[1:], content, child)
			return
		}
	}

	child := &Node{
		name: name,
		dir:  true,
	}

	node.children = append(node.children, child)

	if len(path) == 1 {
		child.dir = false
		child.content = content
	}

	add(path[1:], content, child)
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

	for _, child := range node.children {
		if path[0] == child.name {
			if len(path) == 1 {
				return child
			}
			return find(path[1:], child)
		}
	}

	return nil
}

func walk(path string, node *Node, fn filepath.WalkFunc) error {
	if err := fn(path, node, nil); err != nil {
		return err
	}

	for _, child := range node.children {
		if err := walk(filepath.Join(path, child.name), child, fn); err != nil {
			return err
		}
	}

	return nil
}
