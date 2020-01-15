/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

//Volume represents a docker volume which may be shared among multiple containers
type Volume struct {
	// Name is the name of the docker volume
	Name string `json:"name"`
	// Labels to be attached to the volume
	Labels map[string]string `json:"labels"`

	Global bool `json:"global"`

	Hosts []string `json:"hosts,omitempty"`
}

// Mount represents the information needed for the mounting of a volume
type Mount struct {
	Name      string `json:"name"`
	Directory string `json:"directory"`
	ReadOnly  bool   `json:"readOnly"`
}

// VolumeShare prepares the environment for global volumes
type VolumeShare struct {
	Hosts []string `json:"hosts"`
}
