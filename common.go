package parcel

// ResourceManager keeps track of all resources
var ResourceManager *Manager = &Manager{}

// AddResource adds resource to the default resource manager
func AddResource(resource Binary) {
	if err := ResourceManager.Add(resource); err != nil {
		panic(err)
	}
}
