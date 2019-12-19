// Package biome messages used to perform commands, transiting over RabbitMQ.
// the testexecution service imports this package.
package biome

// CloudProvider is the remote cloud provide to provision the instance on
type CloudProvider string

const (
	// GCPProvider represents the Google Cloud Platform Provider
	GCPProvider = CloudProvider("gcp")
)

// Instance defines the size of each instance to create
type Instance struct {
	Provider CloudProvider `json:"provider,omitonempty"`
	CPUs     int64         `json:"cpus"`
	Memory   int64         `json:"memory"` //MB
	Storage  int64         `json:"storage"`
}

// CreateBiome represents the create biome command
type CreateBiome struct {
	DefinitionID string     `json:"definitionID"`
	TestID       string     `json:"testID"`
	OrgID        string     `json:"orgID"`
	Instances    []Instance `json:"nodes"`
}

// DestroyBiome represents the destroy biome command
type DestroyBiome struct {
	DefinitionID string `json:"definitionID"`
	TestID       string `json:"testID"`
}
