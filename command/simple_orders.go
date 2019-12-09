/*
	Copyright 2019 Whiteblock Inc.
	This file is a part of the Definition.

	Definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Definition is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
	ContainerName string `json:"container"`
	// Name of the network.
	Network string `json:"network"`
}

// FileAndContainer is a file and container order payload.
type FileAndContainer struct {
	// Name of the container.
	ContainerName string `json:"container"`
	// Name of the file.
	File File `json:"file"`
}

// FileAndVolume is a file and volume order payload.
type FileAndVolume struct {
	// Name of the volume.
	VolumeName string `json:"volume"`
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
	RegistryAuth string `json:"registryAuth,omitonempty"`
}
