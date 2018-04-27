package parcello

import (
	"os"
	"path/filepath"
)

// ResourceManager keeps track of all resources
var ResourceManager FileSystemManager

func init() {
	var err error

	if ResourceManager, err = manager(); err != nil {
		panic(err)
	}
}

func manager() (FileSystemManager, error) {
	if mode := os.Getenv("PARCELLO_DEV"); mode != "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}

		return Dir(dir), nil
	}

	manager, err := NewResourceManager()
	if err != nil {
		return nil, err
	}

	return manager, nil
}

// AddResource adds resource to the default resource manager
// Note that the method may panic if the resource not exists
func AddResource(resource Binary) {
	if err := ResourceManager.Add(resource); err != nil {
		panic(err)
	}
}

// Open opens an embedded resource for read
func Open(name string) (ReadOnlyFile, error) {
	return ResourceManager.Open(name)
}

// Root returns a sub-manager for given path, if the path does not exist
// Note that the method may panic if the resource path is not found
func Root(name string) FileSystemManager {
	manager, err := ResourceManager.Root(name)
	if err != nil {
		panic(err)
	}

	return manager
}

func match(pattern, path, name string) (bool, error) {
	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return false, err
	}

	try, _ := filepath.Match(pattern, name)
	return matched || try, nil
}
