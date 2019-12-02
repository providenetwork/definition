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
	"github.com/whiteblock/definition/command/biome"
)

type ResourceDist []PhaseDist

func (rd ResourceDist) ToBiomeCommand(provider biome.CloudProvider, testnetID string, orgID int64) biome.CreateBiome {

	finalDist := rd[len(rd)-1]
	out := biome.CreateBiome{
		TestnetID: testnetID,
		OrgID:     orgID,
		Instances: make([]biome.Instance, len(finalDist)),
	}
	for i, bucket := range finalDist {
		out.Instances[i] = biome.Instance{
			Provider: provider,
			CPUs:     bucket.CPUs,
			Memory:   bucket.Memory,
			Storage:  bucket.Storage,
		}
	}
	return out
}

func (rd *ResourceDist) Add(buckets []Bucket) {
	if rd == nil {
		rd = &ResourceDist{}
	}
	tmp := ResourceDist(append([]PhaseDist(*rd), PhaseDist(buckets)))
	*rd = tmp
}

func (rd ResourceDist) GetPhase(index int) (PhaseDist, error) {
	if rd == nil || len(rd) <= index {
		return nil, fmt.Errorf("index out of bounds")
	}
	return rd[index], nil
}
