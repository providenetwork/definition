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

package distribute

import (
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
)

type Distributor interface {
	Distribute(spec schema.RootSchema) ([]*entity.ResourceDist, error)
}

type distributor struct {
	calculator BiomeCalculator
	conf       config.Bucket
	log        logrus.Ext1FieldLogger
}

func NewDistributor(conf config.Bucket, calculator BiomeCalculator, log logrus.Ext1FieldLogger) Distributor {
	return &distributor{
		calculator: calculator,
		conf:       conf,
		log:        log,
	}
}

func (dist *distributor) Distribute(spec schema.RootSchema) ([]*entity.ResourceDist, error) {
	out := []*entity.ResourceDist{}
	for _, test := range spec.Tests {
		sp := dist.calculator.NewStatePack(spec, dist.conf)
		testResources := &entity.ResourceDist{}
		err := dist.calculator.AddNextPhase(sp, schema.Phase{
			System: test.System,
		})
		if err != nil {
			return nil, err
		}
		testResources.Add(dist.calculator.Resources(sp))
		for _, phase := range test.Phases {
			err = dist.calculator.AddNextPhase(sp, phase)
			if err != nil {
				return nil, err
			}
			testResources.Add(dist.calculator.Resources(sp))
		}
		out = append(out, testResources)
	}
	return out, nil
}
