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
	"fmt"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
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

func (calc testCalculator) handlePhase(state *entity.State,
	networkState entity.NetworkState,
	spec schema.RootSchema,
	phase schema.Phase,
	dist *entity.ResourceDist,
	index int) ([][]command.Command, error) {

	out := [][]command.Command{}
	changedSystems, systems, _ := calc.sys.GetAlreadyExists(state, phase.System)
	calc.logger.WithField("systems", systems).Info("adding these systems")
	servicesToAdd, err := calc.sys.Add(state, spec, systems)
	if err != nil {
		return nil, err
	}

	servicesToRemove, err := calc.sys.Remove(state, spec, phase.Remove)
	if err != nil {
		return nil, err
	}

	networkCommands, err := calc.resolver.CreateNetworks(phase.System, networkState)
	if err != nil {
		return nil, err
	}
	if len(networkCommands) > 0 {
		out = append(out, networkCommands)
		calc.logger.WithFields(logrus.Fields{"count": len(networkCommands)}).Trace(
			"got the network commands")
	}

	diff, err := calc.sys.UpdateChanged(state, spec, changedSystems)
	if err != nil {
		return nil, err
	}
	servicesToAdd = append(servicesToAdd, diff.Added...)
	servicesToRemove = append(servicesToRemove, diff.Removed...)

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
	}).Debug("calculating command diff")

	phaseDist, err := dist.GetPhase(index)
	if err != nil {
		return nil, err
	}

	removalCommands, err := calc.resolver.RemoveServices(phaseDist, servicesToRemove)
	if err != nil {
		return nil, err
	}

	if len(removalCommands) > 0 {
		calc.logger.WithFields(logrus.Fields{"count": len(removalCommands)}).Trace(
			"got the removal commands")
		out = append(out, removalCommands...)
	}

	updateCommands, err := calc.resolver.UpdateServices(phaseDist, diff.Modified)
	if err != nil {
		return nil, err
	}

	if len(updateCommands) > 0 {
		for i, set := range updateCommands {
			calc.logger.WithFields(logrus.Fields{
				"count": len(set),
				"set":   i,
			}).Trace("got the update commands set")
			if len(set) > 0 {
				out = append(out, set)
			}
		}
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

	tmp := entity.TestCommands(out)
	tmp.MetaInject(
		"phase", phase.Name,
		"phaseNum", fmt.Sprint(index))
	return [][]command.Command(tmp), nil
}

func (calc testCalculator) swarmInit(dist *entity.ResourceDist) ([][]command.Command, error) {
	hosts := make([]string, dist.Size())
	for i := range hosts {
		hosts[i] = fmt.Sprint(i)
	}

	order := command.Order{
		Type: command.SwarmInit,
		Payload: command.SetupSwarm{
			Hosts: hosts,
		},
	}
	cmd, err := command.NewCommand(order, "0")
	return [][]command.Command{[]command.Command{cmd}}, err
}

func (calc testCalculator) breakUpCommands(in entity.TestCommands) entity.TestCommands {
	out := [][]command.Command{}
	for _, segment := range in {
		for _, cmd := range segment {
			out = append(out, []command.Command{cmd})
		}
	}
	return entity.TestCommands(out)
}

func (calc testCalculator) Commands(spec schema.RootSchema,
	dist *entity.ResourceDist, index int) (entity.TestCommands, error) {

	state := entity.NewState()
	network, err := entity.NewNetworkState(calc.netConf.GlobalNetwork,
		calc.netConf.SidecarNetwork, calc.netConf.MaxNodesPerNetwork)
	if err != nil {
		return nil, err
	}
	network.GetNextGlobal() //don't use the first entry

	phase := schema.Phase{System: spec.Tests[index].System}
	out := entity.TestCommands{}
	sCmds, err := calc.swarmInit(dist)
	if err != nil {
		return nil, err
	}
	out = out.Append(calc.breakUpCommands(sCmds))
	cmds, err := calc.handlePhase(state, network, spec, phase, dist, 0)
	if err != nil {
		return nil, err
	}
	out = out.Append(calc.breakUpCommands(cmds))
	for i, phase := range spec.Tests[index].Phases {
		cmds, err = calc.handlePhase(state, network, spec, phase, dist, i+1)
		if err != nil {
			return nil, err
		}
		out = out.Append(calc.breakUpCommands(cmds))
	}
	out.MetaInject("test", spec.Tests[index].Name)
	return out, nil
}
