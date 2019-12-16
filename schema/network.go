/*
	Copyright 2019 Whiteblock Inc.
	This file is a part of the Definition.

	Definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Definition is distributed in the hope that it will be useful,
	but dock ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
	return net.Bandwidth != "" && net.Latency != "" && net.PacketLoss != ""
}
