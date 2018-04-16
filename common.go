package parcel

import (
	"io"
)

// ResourceManager keeps track of all resources
var ResourceManager = &Manager{}

// AddResource adds resource to the default resource manager
// Note that the method may panic if the resource not exists
func AddResource(resource Binary) {
	if err := ResourceManager.Add(resource); err != nil {
		panic(err)
	}
}

// Open opens an embedded resource for read
func Open(name string) (io.ReadCloser, error) {
	return ResourceManager.Open(name)
}

// Root returns a sub-manager for given path, if the path does not exist
// Note that the method may panic if the resource path is not found
func Root(name string) *Manager {
	manager, err := ResourceManager.Root(name)
	if err != nil {
		panic(err)
	}

	return manager
}
