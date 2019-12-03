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
	"github.com/sirupsen/logrus"
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/schema"
)

type TestCalculator interface {
	Commands(spec schema.RootSchema, dist *entity.ResourceDist, index int) (entity.TestCommands, error)
}

type testCalculator struct {
	netConf  config.Network
	sys      System
	resolver Resolve
	logger   logrus.Ext1FieldLogger
}

func NewTestCalculator(
	netConf config.Network,
	sys System,
	resolver Resolve,
	logger logrus.Ext1FieldLogger) TestCalculator {
	return &testCalculator{
		netConf:  netConf,
		sys:      sys,
		resolver: resolver,
		logger:   logger}
}

func (calc *testCalculator) handlePhase(state *entity.State, networkState entity.NetworkState,
	spec schema.RootSchema,
	phase schema.Phase, dist *entity.ResourceDist, index int) ([][]command.Command, error) {

	changedSystems, systems, hasChanged := calc.sys.GetAlreadyExists(state, phase.System)

	servicesToAdd, err := calc.sys.Add(state, spec, systems)
	if err != nil {
		return nil, err
	}

	servicesToRemove, err := calc.sys.Remove(state, spec, phase.Remove)
	if err != nil {
		return nil, err
	}

	if hasChanged {
		add, remove, err := calc.sys.UpdateChanged(state, spec, changedSystems)
		if err != nil {
			return nil, err
		}
		servicesToAdd = append(servicesToAdd, add...)
		servicesToRemove = append(servicesToRemove, remove...)
	}

	servicesForTasks, err := calc.sys.Tasks(state, spec, phase.Tasks)
	if err != nil {
		return nil, err
	}
	servicesToAdd = append(servicesToAdd, servicesForTasks...)
	//Break it down into commands now

	calc.logger.WithFields(logrus.Fields{
		"adding":   servicesToAdd,
		"removing": servicesToRemove,
		"systems":  systems,
		"changed":  changedSystems,
	}).Info("calculating command diff")

	networkCommands, err := calc.resolver.CreateNetworks(phase.System, networkState)
	if err != nil {
		return nil, err
	}
	calc.logger.WithFields(logrus.Fields{
		"count": len(networkCommands),
	}).Info("got the network commands")
	out := [][]command.Command{}

	if len(networkCommands) > 0 {
		out = append(out, networkCommands)
	}

	phaseDist, err := dist.GetPhase(index)
	if err != nil {
		return nil, err
	}

	removalCommands, err := calc.resolver.RemoveServices(phaseDist, servicesToRemove)
	if err != nil {
		return nil, err
	}

	calc.logger.WithFields(logrus.Fields{"count": len(removalCommands)}).Trace("got the removal commands")
	if len(removalCommands) > 0 {
		out = append(out, removalCommands...)
	}

	addCommands, err := calc.resolver.CreateServices(spec, networkState, phaseDist, servicesToAdd)
	if err != nil {
		return nil, err
	}

	if len(addCommands) > 0 {
		for i, set := range addCommands {
			calc.logger.WithFields(logrus.Fields{
				"count": len(set),
				"set":   i,
			}).Trace("got the add commands set")
			if len(set) > 0 {
				out = append(out, set)
			}
		}
	}

	return out, nil
}

func (calc *testCalculator) Commands(spec schema.RootSchema,
	dist *entity.ResourceDist, index int) (entity.TestCommands, error) {

	state := entity.NewState()
	network, err := entity.NewNetworkState(calc.netConf.GlobalNetwork,
		calc.netConf.SidecarNetwork, calc.netConf.MaxNodesPerNetwork)
	if err != nil {
		return nil, err
	}

	phase := schema.Phase{System: spec.Tests[index].System}
	out := entity.TestCommands{}
	cmds, err := calc.handlePhase(state, network, spec, phase, dist, 0)
	if err != nil {
		return nil, err
	}
	out = out.Append(cmds)
	for i, phase := range spec.Tests[index].Phases {
		cmds, err = calc.handlePhase(state, network, spec, phase, dist, i+1)
		if err != nil {
			return nil, err
		}
		out = out.Append(cmds)
	}
	return out, nil
}
