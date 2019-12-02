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

package process

import (
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/internal/distribute"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/schema"
)

type TestCalculator interface {
	Commands(spec schema.RootSchema, dist *distribute.ResourceDist, index int) (entity.TestCommands, error)
}

type testCalculator struct {
	sys      System
	resolver Resolve
}

func NewTestCalculator(sys System, resolver Resolve) TestCalculator {
	return &testCalculator{sys: sys, resolver: resolver}
}

func (calc *testCalculator) handlePhase(state *entity.State, spec schema.RootSchema,
	phase schema.Phase, dist *distribute.ResourceDist, index int) ([][]command.Command, error) {

	servicesToAdd, err := calc.sys.Add(state, spec, phase.System)
	if err != nil {
		return nil, err
	}

	servicesToRemove, err := calc.sys.Remove(state, spec, phase.Remove)
	if err != nil {
		return nil, err
	}

	servicesForTasks, err := calc.sys.Tasks(state, spec, phase.Tasks)
	if err != nil {
		return nil, err
	}
	servicesToAdd = append(servicesToAdd, servicesForTasks...)
	//Break it down into commands now

	networkCommands, err := calc.resolver.CreateNetworks(phase.System)
	if err != nil {
		return nil, err
	}
	out := [][]command.Command{networkCommands}

	phaseDist, err := dist.GetPhase(index)
	if err != nil {
		return nil, err
	}

	removalCommands, err := calc.resolver.RemoveServices(phaseDist, servicesToRemove)
	if err != nil {
		return nil, err
	}
	out = append(out, removalCommands...)

	addCommands, err := calc.resolver.CreateServices(spec, phaseDist, servicesToRemove)
	if err != nil {
		return nil, err
	}
	out = append(out, addCommands...)
	return out, nil
}

func (calc *testCalculator) Commands(spec schema.RootSchema,
	dist *distribute.ResourceDist, index int) (entity.TestCommands, error) {

	state := entity.NewState()
	phase := schema.Phase{System: spec.Tests[index].System}
	out := entity.TestCommands{}
	cmds, err := calc.handlePhase(state, spec, phase, dist, 0)
	if err != nil {
		return nil, err
	}
	out = out.Append(cmds)
	for i, phase := range spec.Tests[index].Phases {
		cmds, err = calc.handlePhase(state, spec, phase, dist, i+1)
		if err != nil {
			return nil, err
		}
		out = out.Append(cmds)
	}
	return out, nil
}
