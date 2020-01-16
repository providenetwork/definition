/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package process

import (
	"fmt"
	"strings"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
)

const FirstPhaseName = "init"

type TestCalculator interface {
	Commands(spec schema.RootSchema, dist *entity.ResourceDist, index int) (entity.TestCommands, error)
	Env(spec schema.RootSchema, dist *entity.ResourceDist, index int) (map[string]string, error)
}

type testCalculator struct {
	conf     config.Config
	sys      System
	resolver Resolve
	log      logrus.Ext1FieldLogger
}

func NewTestCalculator(
	conf config.Config,
	sys System,
	resolver Resolve,
	log logrus.Ext1FieldLogger) TestCalculator {
	return &testCalculator{
		conf:     conf,
		sys:      sys,
		resolver: resolver,
		log:      log}
}

func (calc testCalculator) handlePhase(state *entity.State,
	spec schema.RootSchema,
	phase schema.Phase,
	dist *entity.ResourceDist,
	index int) ([][]command.Command, error) {

	out := [][]command.Command{}
	changedSystems, systems, _ := calc.sys.GetAlreadyExists(state, phase.System)
	calc.log.WithFields(logrus.Fields{
		"systems": systems,
		"changed": changedSystems,
		"phase":   phase.Name,
		"index":   index}).Info("adding these systems")

	networkCommands, err := calc.resolver.CreateSystemNetworks(state, phase.System)
	if err != nil {
		return nil, err
	}

	servicesToAdd, err := calc.sys.Add(state, spec, systems)
	if err != nil {
		return nil, err
	}

	servicesToRemove, err := calc.sys.Remove(state, spec, phase.Remove)
	if err != nil {
		return nil, err
	}

	diff, err := calc.sys.UpdateChanged(state, spec, changedSystems)
	if err != nil {
		return nil, err
	}

	additionalNetworkCmds, err := calc.resolver.CreateNetworks(state, diff.AddedNetworks, nil)
	if err != nil {
		return nil, err
	}

	networkCommands = append(networkCommands, additionalNetworkCmds...)
	if len(networkCommands) > 0 {
		out = append(out, networkCommands)
		calc.log.WithFields(logrus.Fields{"count": len(networkCommands)}).Trace(
			"got the network commands")
	}

	servicesToAdd = append(servicesToAdd, diff.Added...)
	servicesToRemove = append(servicesToRemove, diff.Removed...)

	servicesForTasks, err := calc.sys.Tasks(state, spec, phase.Tasks)
	if err != nil {
		return nil, err
	}
	servicesToAdd = append(servicesToAdd, servicesForTasks...)
	//Break it down into commands now

	calc.log.WithFields(logrus.Fields{
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
		calc.log.WithFields(logrus.Fields{"count": len(removalCommands)}).Trace(
			"got the removal commands")
		out = append(out, removalCommands...)
	}

	updateCommands, err := calc.resolver.UpdateServices(state, phaseDist, diff.Modified)
	if err != nil {
		return nil, err
	}

	if len(updateCommands) > 0 {
		for i, set := range updateCommands {
			calc.log.WithFields(logrus.Fields{
				"count": len(set),
				"set":   i,
			}).Trace("got the update commands set")
			if len(set) > 0 {
				out = append(out, set)
			}
		}
	}

	addCommands, err := calc.resolver.CreateServices(state, spec, phaseDist, servicesToAdd)
	if err != nil {
		return nil, err
	}

	if len(addCommands) > 0 {
		for i, set := range addCommands {
			calc.log.WithFields(logrus.Fields{
				"count": len(set),
				"set":   i,
			}).Trace("got the add commands set")
			if len(set) > 0 {
				out = append(out, set)
			}
		}
	}

	tmp := entity.TestCommands(out)
	if !phase.Duration.Empty() {
		order := command.Order{
			Type:    command.Pauseexecution,
			Payload: phase.Duration,
		}
		pauseCmd, err := command.NewCommand(order, FirstInstance)
		if err != nil {
			return nil, err
		}
		tmp = tmp.Append([][]command.Command{{pauseCmd}})
	}

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
	cmd, err := command.NewCommand(order, FirstInstance)
	return [][]command.Command{{cmd}}, err
}

func volumeCommands(spec schema.RootSchema, dist *entity.ResourceDist) ([][]command.Command, error) {
	volumes := parser.ExtractGlobalVolumes(spec)
	if len(volumes) == 0 {
		return nil, nil
	}

	hosts := make([]string, dist.Size()) //dont initialize glusterfs if there is only one instance
	for i := range hosts {
		hosts[i] = fmt.Sprint(i)
	}

	out := []command.Command{}
	for _, volume := range volumes {
		order := command.Order{
			Type: command.Createvolume,
			Payload: command.Volume{
				Name:   volume,
				Global: len(hosts) > 1,
				Hosts:  hosts,
			},
		}
		cmd, err := command.NewCommand(order, FirstInstance)
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}

	if len(hosts) < 2 {
		return [][]command.Command{out}, nil
	}
	order := command.Order{
		Type: command.Volumeshare,
		Payload: command.VolumeShare{
			Hosts: hosts,
		},
	}

	initCmd, err := command.NewCommand(order, FirstInstance)
	if err != nil {
		return nil, err
	}

	return [][]command.Command{{initCmd}, out}, nil
}

func (calc testCalculator) breakUpCommands(in entity.TestCommands) entity.TestCommands {
	if calc.conf.Output.AsIsCommands {
		return in
	}
	out := [][]command.Command{}
	for _, segment := range in {
		var current []command.Command
		for i, cmd := range segment {
			if i%calc.conf.Output.MaxParallelCommands == 0 {
				if current != nil {
					out = append(out, current)
				}
				current = []command.Command{}
			}
			current = append(current, cmd)
		}
		if len(current) > 0 {
			out = append(out, current)
		}
	}
	return entity.TestCommands(out)
}

func (calc testCalculator) processTest(spec schema.RootSchema,
	dist *entity.ResourceDist, index int) (entity.TestCommands, map[string]string, error) {

	network, err := entity.NewNetworkState(calc.conf.Network.GlobalNetwork,
		calc.conf.Network.SidecarNetwork, calc.conf.Network.MaxNodesPerNetwork)
	if err != nil {
		return nil, nil, err
	}
	state := entity.NewState(network)

	network.GetNextGlobal() // don't use the first entry

	out := entity.TestCommands{}

	sCmds, err := calc.swarmInit(dist)
	if err != nil {
		return nil, nil, err
	}
	out = out.Append(calc.breakUpCommands(sCmds))

	vCmds, err := volumeCommands(spec, dist)
	if err != nil {
		return nil, nil, err
	}
	out = out.Append(calc.breakUpCommands(vCmds))
	out.MetaInject(command.PhaseKey, FirstPhaseName)

	cmds, err := calc.handlePhase(state, spec, schema.Phase{
		System:   spec.Tests[index].System,
		Name:     FirstPhaseName,
		Duration: spec.Tests[index].Duration,
	}, dist, 0)

	if err != nil {
		return nil, nil, err
	}
	out = out.Append(calc.breakUpCommands(cmds))
	for i, phase := range spec.Tests[index].Phases {
		cmds, err = calc.handlePhase(state, spec, phase, dist, i+1)
		if err != nil {
			return nil, nil, err
		}
		out = out.Append(calc.breakUpCommands(cmds))
	}
	out.MetaInject(command.TestNameKey, spec.Tests[index].Name)

	envVars := map[string]string{}
	for name, ip := range state.IPs {
		name = strings.Replace(name, "-", "_", -1)
		name = strings.ToUpper(name)
		envVars[name] = ip
	}
	return out, envVars, err
}

func (calc testCalculator) Env(spec schema.RootSchema,
	dist *entity.ResourceDist, index int) (map[string]string, error) {

	_, envVars, err := calc.processTest(spec, dist, index)
	return envVars, err
}

func (calc testCalculator) Commands(spec schema.RootSchema,
	dist *entity.ResourceDist, index int) (entity.TestCommands, error) {

	out, envVars, err := calc.processTest(spec, dist, index)
	if err != nil {
		return nil, err
	}
	for i := range out {
		for j := range out[i] {
			if out[i][j].Order.Type != command.Createcontainer {
				continue
			}
			var order command.Container

			err = out[i][j].ParseOrderPayloadInto(&order)
			if err != nil {
				return nil, err
			}
			if order.Environment == nil {
				order.Environment = map[string]string{}
			}
			for key, val := range envVars {
				order.Environment[key] = val
			}

			out[i][j].Order.Payload = order
		}
	}
	return out, nil
}
