package parcello

import (
	"os"
	"path/filepath"
)

// Resource keeps track of all resources
var Resource = NewManager()

// NewManager creates a FileSystemManager based on whether dev mode is enabled
func NewManager() FileSystemManager {
	mode := os.Getenv("PARCELLO_DEV_ENABLED")

	if mode != "" {
		return Dir(getenv("PARCELLO_RESOURCE_DIR", "."))
	}

	manager, err := NewResourceManager()
	if err != nil {
		panic(err)
	}

	return manager
}

// AddResource adds resource to the default resource manager
// Note that the method may panic if the resource not exists
func AddResource(resource Binary) {
	if err := Resource.Add(resource); err != nil {
		panic(err)
	}
}

// Open opens an embedded resource for read
func Open(name string) (ReadOnlyFile, error) {
	return Resource.Open(name)
}

// Root returns a sub-manager for given path, if the path does not exist
// Note that the method may panic if the resource path is not found
func Root(name string) FileSystemManager {
	manager, err := Resource.Root(name)
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

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
