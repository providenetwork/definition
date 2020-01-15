/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
	if res.Ports == nil {
		res.Ports = map[int]bool{}
	}
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
