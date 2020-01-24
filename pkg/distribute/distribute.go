/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package distribute

import (
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/pkg/entity"
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
