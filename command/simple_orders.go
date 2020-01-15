/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

// SimpleName is a simple order payload with just the container name
type SimpleName struct {
	// Name of the container.
	Name string `json:"name"`
}

// ContainerNetwork is a container and network order payload.
type ContainerNetwork struct {
	// Name of the container.
	Container string `json:"container"`
	// Name of the network.
	Network string `json:"network"`
	// IP is the IP address of the container in this network
	IP string `json:"ip,omitempty"`
}

// FileAndContainer is a file and container order payload.
type FileAndContainer struct {
	// Name of the container.
	ContainerName string `json:"container"`
	// Name of the file.
	File File `json:"file"`
}

// SetupSwarm is the payload to setup a docker swarm
type SetupSwarm struct {
	//Hosts is an array of the hosts to be setup with docker swarm
	Hosts []string `json:"hosts"`
}

// PullImage contains the information necessary to pull a docker image from a registry
type PullImage struct {
	Image string `json:"image"`

	// RegistryAuth is the base64 encoded credentials for the registry (optional)
	RegistryAuth string `json:"registryAuth,omitempty"`
}
