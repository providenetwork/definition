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
	"github.com/whiteblock/definition/schema"
)

type Distributor interface {
	Distribute(spec schema.RootSchema) ([][]Bucket, error)
}

type distributor struct {
	calculator BiomeCalculator
}

func NewDistributor(calculator BiomeCalculator) Distributor {
	return &distributor{
		calculator: calculator,
	}
}

func (dist *distributor) Distribute(spec schema.RootSchema) ([][]Bucket, error) {
	out := [][]Bucket{}
	for _, test := range spec.Tests {
		sp := dist.calculator.NewStatePack()
		for _, phase := range test.Phases {
			err := dist.calculator.AddNextPhase(sp, phase)
			if err != nil {
				return nil, err
			}
		}
		out = append(out, dist.calculator.Resources(sp))
	}
	return out, nil
}
