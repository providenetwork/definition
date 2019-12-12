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

package entity

type Resource struct {
	CPUs    int64
	Memory  int64
	Storage int64
	Ports   map[int]bool
}

func (res Resource) GetResources() Resource {
	return res
}

func (res *Resource) UpdateResources(newRes Resource) {
	res.CPUs = newRes.CPUs
	res.Memory = newRes.Memory
	res.Storage = newRes.Storage
	res.Ports = newRes.Ports
}

func (res Resource) NoPortConflicts(ports ...int) bool {
	for _, port := range ports {
		if _, exists := res.Ports[port]; exists {
			return false
		}
	}
	return true
}

func (res Resource) GetPorts() []int {
	out := []int{}
	for port := range res.Ports {
		out = append(out, port)
	}
	return out
}

func (res *Resource) InsertPorts(ports ...int) {
	for _, port := range ports {
		res.Ports[port] = true
	}
}

func (res Resource) HasPort(port int) bool {
	_, exists := res.Ports[port]
	return exists
}

func (res *Resource) RemovePorts(ports ...int) {
	for _, port := range ports {
		res.RemovePort(port)
	}
}

func (res *Resource) RemovePort(port int) {
	if !res.HasPort(port) {
		return
	}
	delete(res.Ports, port)
}
