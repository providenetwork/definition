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
	"github.com/whiteblock/definition/internal/config"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/schema"
)

type StatePack struct {
	state     SystemState
	buckets   ResourceBuckets
	prevTasks []entity.Segment
}

//BiomeCalculator is a calculator for the state for the testnet
//as time goes on.
type BiomeCalculator interface {
	NewStatePack() *StatePack
	AddNextPhase(sp *StatePack, phase schema.Phase) error
	Resources(sp *StatePack) []Bucket
}

type biomeCalculator struct {
	conf   config.Bucket
	parser parser.Resources
	namer  parser.Names
}

func NewBiomeCalculator(
	conf config.Bucket,
	parser parser.Resources,
	namer parser.Names) BiomeCalculator {

	return &biomeCalculator{conf: conf, parser: parser, namer: namer}
}

func (bc *biomeCalculator) NewStatePack() *StatePack {
	return &StatePack{
		state:     NewSystemState(bc.parser, bc.namer),
		buckets:   NewResourceBuckets(bc.conf),
		prevTasks: nil,
	}
}

func (bc *biomeCalculator) AddNextPhase(sp *StatePack, phase schema.Phase) error {

	if sp.prevTasks != nil {
		err := sp.buckets.Remove(sp.prevTasks)
		if err != nil {
			return err
		}
	}

	addSysSegs, err := sp.state.Add(phase.System)
	if err != nil {
		return err
	}

	err = sp.buckets.Add(addSysSegs)
	if err != nil {
		return err
	}

	toRemove, err := sp.state.Remove(phase.Remove)
	if err != nil {
		return err
	}

	err = sp.buckets.Remove(toRemove)
	if err != nil {
		return err
	}
	sp.prevTasks, err = bc.parser.Tasks(phase.Tasks)
	if err != nil {
		return err
	}

	return sp.buckets.Add(sp.prevTasks)
}

func (bc *biomeCalculator) Resources(sp *StatePack) []Bucket {
	return sp.buckets.Resources()
}
