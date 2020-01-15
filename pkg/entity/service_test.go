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
	"sort"
	"strings"
	"testing"

	"github.com/whiteblock/definition/schema"

	"github.com/stretchr/testify/assert"
)

func TestService_CalculateDiff(t *testing.T) {
	networks := []schema.Network{
		schema.Network{
			Name: "test",
		},
		schema.Network{
			Name:    "test",
			Latency: "10ms",
		},
		schema.Network{
			Name:      "test",
			Latency:   "10ms",
			Bandwidth: "100 MBps",
		},
		schema.Network{
			Name:    "test2",
			Latency: "10ms",
		},
		schema.Network{
			Name:      "test3",
			Latency:   "10ms",
			Bandwidth: "100 MBps",
		},
	}
	sidecars := []schema.Sidecar{
		schema.Sidecar{
			Name: "test",
		},
		schema.Sidecar{
			Name: "test2",
		},
		schema.Sidecar{
			Name: "test3",
		},
	}

	services := []Service{
		Service{
			Name: "foo",
			Networks: []schema.Network{
				networks[0],
				networks[4],
			},
			Sidecars: []schema.Sidecar{},
		},
		Service{
			Name: "foo",
			Networks: []schema.Network{
				networks[0],
				networks[3],
			},
			Sidecars: []schema.Sidecar{
				sidecars[0],
				sidecars[2],
			},
		},
		Service{
			Name: "foo",
			Networks: []schema.Network{
				networks[1],
			},
			Sidecars: []schema.Sidecar{
				sidecars[0],
				sidecars[1],
			},
		},
	}

	diff1 := services[0].CalculateDiff(services[1])
	expectedDiff1 := ServiceDiff{
		Name: services[0].Name,
		AddNetworks: []schema.Network{
			networks[3],
		},
		DetachNetworks: []schema.Network{
			networks[4],
		},
		UpdateNetworks: nil,
		AddSidecars: []schema.Sidecar{
			sidecars[0],
			sidecars[2],
		},
		RemoveSidecars: nil,
		Parent:         &services[1],
	}

	diff2 := services[0].CalculateDiff(services[2])
	expectedDiff2 := ServiceDiff{
		Name:        services[0].Name,
		AddNetworks: nil,
		DetachNetworks: []schema.Network{
			networks[4],
		},
		UpdateNetworks: []schema.Network{
			networks[1],
		},
		AddSidecars: []schema.Sidecar{
			sidecars[0],
			sidecars[1],
		},
		RemoveSidecars: nil,
		Parent:         &services[2],
	}

	diff3 := services[2].CalculateDiff(services[1])
	expectedDiff3 := ServiceDiff{
		Name: services[2].Name,
		AddNetworks: []schema.Network{
			networks[3],
		},
		DetachNetworks: nil,
		UpdateNetworks: []schema.Network{
			networks[1],
		},
		AddSidecars: []schema.Sidecar{
			sidecars[2],
		},
		RemoveSidecars: []schema.Sidecar{
			sidecars[1],
		},
		Parent: &services[1],
	}

	diff4 := services[2].CalculateDiff(services[0])
	expectedDiff4 := ServiceDiff{
		Name: services[2].Name,
		AddNetworks: []schema.Network{
			networks[4],
		},
		DetachNetworks: nil,
		UpdateNetworks: []schema.Network{
			networks[1],
		},
		AddSidecars: nil,
		RemoveSidecars: []schema.Sidecar{
			sidecars[0],
			sidecars[1],
		},
		Parent: &services[0],
	}

	sort.Slice(diff1.AddSidecars, func(i, j int) bool {
		return strings.Compare(diff1.AddSidecars[i].Name, diff1.AddSidecars[j].Name) < 0
	})

	sort.Slice(diff2.AddSidecars, func(i, j int) bool {
		return strings.Compare(diff2.AddSidecars[i].Name, diff2.AddSidecars[j].Name) < 0
	})

	sort.Slice(diff4.RemoveSidecars, func(i, j int) bool {
		return strings.Compare(diff4.RemoveSidecars[i].Name, diff4.RemoveSidecars[j].Name) < 0
	})

	assert.Equal(t, expectedDiff1, diff1)
	assert.Equal(t, expectedDiff2, diff2)
	assert.Equal(t, expectedDiff3, diff3)
	assert.Equal(t, expectedDiff4, diff4)
}
