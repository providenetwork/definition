/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

import (
	"fmt"

	"github.com/whiteblock/definition/command/biome"
)

type ResourceDist []PhaseDist

func (rd ResourceDist) ToBiomeCommand(provider biome.CloudProvider,
	defID string, orgID string, testID string, domain string) biome.CreateBiome {

	finalDist := rd[len(rd)-1]
	out := biome.CreateBiome{
		DefinitionID: defID,
		TestID:       testID,
		OrgID:        orgID,
		Instances:    make([]biome.Instance, len(finalDist)),
	}
	for i, bucket := range finalDist {
		dns := ""
		if domain != "" {
			dns = fmt.Sprintf("%s-%d", domain, i)
		}
		out.Instances[i] = biome.Instance{
			Provider: provider,
			CPUs:     bucket.CPUs,
			Memory:   bucket.Memory,
			Storage:  bucket.Storage,
			Domain:   dns,
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
		return nil, fmt.Errorf("index %d out of bounds", index)
	}
	return rd[index], nil
}

func (rd ResourceDist) Size() int {
	if len(rd) == 0 {
		return 0
	}
	return len(rd[len(rd)-1])
}
