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
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/schema"
)

//BiomeCalculator is a calculator for the state for the testnet
//as time goes on.
type BiomeCalculator interface {
	AddNextPhase(phase schema.Phase) error
	Resources() []entity.Resource
}

type biomeCalculator struct {
	state   SystemState
	buckets ResourceBuckets
	parser  parser.Schema
}

func NewBiomeCalculator(
	state SystemState,
	buckets ResourceBuckets,
	parser parser.Schema) BiomeCalculator {
	return &biomeCalculator{state: state, buckets: buckets, parser: parser}
}

func (bc *biomeCalculator) AddNextPhase(phase schema.Phase) error {

	addSysSegs, err := bc.state.Add(phase.System)
	if err != nil {
		return err
	}

	err = bc.buckets.Add(addSysSegs)
	if err != nil {
		return err
	}

	toRemove, err := bc.state.Remove(phase.Remove)
	if err != nil {
		return err
	}

	err = bc.buckets.Remove(addSysSegs)
	if err != nil {
		return err
	}
	taskSegments := bc.parser.ParseTasks(phase.Tasks)

	err = bc.buckets.Add(taskSegments)
	if err != nil {
		return err
	}

	return bc.buckets.Remove(taskSegments)
}

func (bc *biomeCalculator) Resources() []entity.Resource {
	return bc.buckets.Resources()
}
