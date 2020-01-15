/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package distribute

import (
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
)

//BiomeCalculator is a calculator for the state for the testnet
//  As time goes on.
type BiomeCalculator interface {
	NewStatePack(spec schema.RootSchema, conf config.Bucket) *entity.StatePack
	AddNextPhase(sp *entity.StatePack, phase schema.Phase) error
	Resources(sp *entity.StatePack) []entity.Bucket
}

type biomeCalculator struct {
	parser parser.Resources
	state  SystemState
	logger logrus.Ext1FieldLogger
}

func NewBiomeCalculator(
	parser parser.Resources,
	state SystemState,
	logger logrus.Ext1FieldLogger) BiomeCalculator {

	return &biomeCalculator{parser: parser, state: state, logger: logger}
}

func (bc *biomeCalculator) NewStatePack(spec schema.RootSchema, conf config.Bucket) *entity.StatePack {
	return entity.NewStatePack(spec, conf, bc.logger)
}

func (bc *biomeCalculator) AddNextPhase(sp *entity.StatePack, phase schema.Phase) error {

	changedSystems, systems, hasChanged := bc.state.GetAlreadyExists(sp, phase.System)

	if sp.PrevTasks != nil {
		err := sp.Buckets.Remove(sp.PrevTasks)
		if err != nil {
			return err
		}
	}

	addSysSegs, err := bc.state.Add(sp, sp.Spec, systems)
	if err != nil {
		return err
	}

	toRemove, err := bc.state.Remove(sp, phase.Remove)
	if err != nil {
		return err
	}

	if hasChanged {
		add, remove, err := bc.state.UpdateChanged(sp, sp.Spec, changedSystems)
		if err != nil {
			return err
		}
		addSysSegs = append(addSysSegs, add...)
		toRemove = append(toRemove, remove...)
	}
	bc.logger.WithFields(logrus.Fields{
		"adding":   addSysSegs,
		"removing": toRemove,
		"systems":  systems,
	}).Debug("calculating distribution diff")
	err = sp.Buckets.Remove(toRemove)
	if err != nil {
		return err
	}

	err = sp.Buckets.Add(addSysSegs)
	if err != nil {
		return err
	}

	sp.PrevTasks, err = bc.parser.Tasks(sp.Spec, phase.Tasks)
	if err != nil {
		return err
	}

	return sp.Buckets.Add(sp.PrevTasks)
}

func (bc *biomeCalculator) Resources(sp *entity.StatePack) []entity.Bucket {
	return sp.Buckets.Resources()
}
