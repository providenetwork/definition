/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
