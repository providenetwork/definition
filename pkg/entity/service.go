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

import (
	"fmt"
	"reflect"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config/defaults"
	"github.com/whiteblock/definition/schema"

	"github.com/imdario/mergo"
)

type Service struct {
	Name            string
	Bucket          int64
	SquashedService schema.Service
	Networks        []schema.Network
	Sidecars        []schema.Sidecar
	Labels          map[string]string
	Ports           map[int]int

	Timeout        command.Timeout
	IsTask         bool
	IgnoreExitCode bool
}

func GetDefaultService(def defaults.Defaults) Service {
	return Service{
		Name:   "",
		Bucket: -1,
		SquashedService: schema.Service{
			Name:        "",
			Description: "",
			Volumes:     nil,
			Resources: schema.Resources{
				Cpus:    def.Resources.CPUs,
				Memory:  fmt.Sprintf("%dMB", def.Resources.Memory),
				Storage: fmt.Sprintf("%dMB", def.Resources.Storage),
			},
			Args:        nil,
			Environment: nil,
			Image:       def.Service.Image,
			InputFiles:  nil,
		},
		Networks: nil,
		Sidecars: nil,
	}
}

func (serv Service) Equal(serv2 Service) bool {
	return reflect.DeepEqual(serv, serv2)
}

func (serv Service) CalculateDiff(serv2 Service) ServiceDiff {
	out := ServiceDiff{Name: serv.Name, Parent: &serv2}

	sidecars := map[string]int{}
	allSidecars := map[string]schema.Sidecar{}

	for _, sidecar := range serv.Sidecars {
		sidecars[sidecar.Name] |= 0x01
		allSidecars[sidecar.Name] = sidecar
	}

	for _, sidecar := range serv2.Sidecars {
		sidecars[sidecar.Name] |= 0x02
		allSidecars[sidecar.Name] = sidecar
	}

	for sidecar, status := range sidecars {
		switch status {
		case 0x03: // Sidecar exists in both
			continue
		case 0x02: // Sidecar exists only in the new service
			out.AddSidecars = append(out.AddSidecars, allSidecars[sidecar])
		case 0x01: // Sidecar exists only in the old service
			out.RemoveSidecars = append(out.RemoveSidecars, allSidecars[sidecar])
		}
	}

	networks := map[string]schema.Network{}
	networkStatus := map[string]int{}
	for _, network := range serv.Networks {
		networkStatus[network.Name] |= 0x01
		networks[network.Name] = network
	}

	for _, network := range serv2.Networks {
		if _, exists := networks[network.Name]; exists { //Update
			if !reflect.DeepEqual(networks[network.Name], network) {
				tmp := networks[network.Name]
				mergo.Map(&tmp, network, mergo.WithOverride)
				networks[network.Name] = tmp
				networkStatus[network.Name] = 0x04
				continue
			}
		} else {
			networks[network.Name] = network
		}
		networkStatus[network.Name] |= 0x02
	}

	for networkName, status := range networkStatus {
		switch status {
		case 0x04: // Needs update
			out.UpdateNetworks = append(out.UpdateNetworks, networks[networkName])
		case 0x03: // Both have
			continue
		case 0x02: // Only new has
			out.AddNetworks = append(out.AddNetworks, networks[networkName])
		case 0x01:
			out.DetachNetworks = append(out.DetachNetworks, networks[networkName])
		}
	}
	return out
}
