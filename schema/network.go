/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package schema

type Network struct {
	Name       string `yaml:"name,omitempty" json:"name,omitempty"`
	Bandwidth  string `yaml:"bandwidth,omitempty" json:"bandwidth,omitempty"`
	Latency    string `yaml:"latency,omitempty" json:"latency,omitempty"`
	PacketLoss string `yaml:"packet-loss,omitempty" json:"packet-loss,omitempty"`
}

// HasEmulation checks if the network struct has any emulation included
func (net Network) HasEmulation() bool {
	return net.Bandwidth != "" || net.Latency != "" || net.PacketLoss != ""
}
