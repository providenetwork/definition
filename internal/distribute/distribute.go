/*
	Copyright 2019 whiteblock Inc.
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

package distribute

import (
	"fmt"

	"github.com/whiteblock/definition/schema"
)

type ResourceDist [][]Bucket

func (rd *ResourceDist) Add(buckets []Bucket) {
	if rd == nil {
		rd = &ResourceDist{}
	}
	tmp := ResourceDist(append([][]Bucket(*rd), buckets))
	*rd = tmp
}

func (rd ResourceDist) GetPhase(index int) ([]Bucket, error) {
	if rd == nil || len(rd) <= index {
		return nil, fmt.Errorf("index out of bounds")
	}
	return rd[index], nil
}

type Distributor interface {
	Distribute(spec schema.RootSchema) ([]*ResourceDist, error)
}

type distributor struct {
	calculator BiomeCalculator
}

func NewDistributor(calculator BiomeCalculator) Distributor {
	return &distributor{
		calculator: calculator,
	}
}

func (dist *distributor) Distribute(spec schema.RootSchema) ([]*ResourceDist, error) {
	out := []*ResourceDist{}
	for _, test := range spec.Tests {
		sp := dist.calculator.NewStatePack()
		testResources := &ResourceDist{}
		for _, phase := range test.Phases {
			err := dist.calculator.AddNextPhase(sp, phase)
			if err != nil {
				return nil, err
			}
			testResources.Add(dist.calculator.Resources(sp))
		}
		out = append(out, testResources)
	}
	return out, nil
}
