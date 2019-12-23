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
	"errors"
	"fmt"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/maker"
	"github.com/whiteblock/definition/schema"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
)

type Dependency interface {
	AttachNetworks(bucket int, state *entity.State, container string,
		networks []schema.Network) ([]command.Command, error)

	//Container returns create, start, error
	Container(bucket int, state *entity.State,
		service entity.Service) (command.Command, command.Command, error)

	DetachNetworks(bucket int, container string,
		networks []schema.Network) ([]command.Command, error)

	Emulation(bucket int, container string, networks []schema.Network) ([]command.Command, error)

	Files(bucket int, service entity.Service) ([]command.Command, error)

	RemoveContainer(bucket int, name string) (command.Command, error)

	Sidecars(bucket int, state *entity.State, service entity.Service,
		sidecars []schema.Sidecar) ([][]command.Command, error)

	SidecarNetwork(bucket int, state *entity.State,
		service entity.Service) (command.Command, error)

	Volumes(bucket int, service entity.Service) ([]command.Command, error)
}

var (
	ErrNoFreeIP = errors.New("out of ip address to allocate")
)

type dependency struct {
	parser   maker.Service
	cmdMaker maker.Command
	log      logrus.Ext1FieldLogger
}

func NewDependency(
	cmdMaker maker.Command,
	parser maker.Service,
	log logrus.Ext1FieldLogger) Dependency {
	return &dependency{cmdMaker: cmdMaker, parser: parser, log: log}
}

func (dep dependency) Emulation(bucket int, container string,
	networks []schema.Network) ([]command.Command, error) {

	out := []command.Command{}
	for _, network := range networks {
		if !network.HasEmulation() {
			dep.log.WithField("network", network).Debug("skipping network which doesn't have emulation")
			continue
		}
		order, err := dep.cmdMaker.Emulation(container, network)
		if err != nil {
			return nil, err
		}
		cmd, err := command.NewCommand(order, fmt.Sprint(bucket))
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}
	return out, nil
}

func (dep dependency) DetachNetworks(bucket int, container string,
	networks []schema.Network) ([]command.Command, error) {
	out := []command.Command{}
	for _, network := range networks {
		order := dep.cmdMaker.DetachNetwork(container, network.Name)
		cmd, err := command.NewCommand(order, fmt.Sprint(bucket))
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}
	return out, nil
}

func (dep dependency) AttachNetworks(bucket int, state *entity.State, container string,
	networks []schema.Network) ([]command.Command, error) {
	out := []command.Command{}
	for _, network := range networks {
		ip := state.Subnets[network.Name].Next()
		if ip == nil {
			return nil, ErrNoFreeIP
		}
		state.IPs[container+"_"+network.Name] = ip.String()
		order := dep.cmdMaker.AttachNetwork(container, network.Name, ip.String())
		cmd, err := command.NewCommand(order, fmt.Sprint(bucket))
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}
	return out, nil
}

func (dep dependency) Container(bucket int, state *entity.State, service entity.Service) (
	create command.Command, start command.Command, err error) {

	order := dep.cmdMaker.CreateContainer(state, service)

	create, err = command.NewCommand(order, fmt.Sprint(bucket))
	if err != nil {
		return
	}

	err = mergo.Map(&create.Meta, service.Labels)
	if err != nil {
		return
	}

	order = dep.cmdMaker.StartContainer(service, service.IsTask, service.Timeout)

	start, err = command.NewCommand(order, fmt.Sprint(bucket))
	if err != nil {
		return
	}

	err = mergo.Map(&start.Meta, service.Labels)
	return
}

func (dep dependency) files(bucket int, service entity.Service, name string,
	inputs []schema.InputFile) ([]command.Command, error) {

	out := []command.Command{}
	for _, input := range inputs {
		order := dep.cmdMaker.File(name, input)
		cmd, err := command.NewCommand(order, fmt.Sprint(bucket))
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}
	return out, nil
}

func (dep dependency) Files(bucket int, service entity.Service) ([]command.Command, error) {
	cmds, err := dep.files(bucket, service, service.Name, service.SquashedService.InputFiles)
	if err != nil {
		return nil, err
	}
	for _, sidecar := range service.Sidecars {
		sCmds, err := dep.files(bucket, service, sidecar.Name, sidecar.InputFiles)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, sCmds...)
	}
	return cmds, nil

}

func (dep dependency) Sidecars(bucket int, state *entity.State, service entity.Service,
	sidecars []schema.Sidecar) ([][]command.Command, error) {

	out := make([][]command.Command, 2)
	for _, sidecar := range sidecars {

		order := dep.cmdMaker.CreateSidecar(state, service, sidecar)
		create, err := command.NewCommand(order, fmt.Sprint(bucket))
		if err != nil {
			return nil, err
		}

		err = mergo.Map(&create.Meta, service.Labels)
		if err != nil {
			return nil, err
		}
		create.Meta["service"] = service.Name

		out[0] = append(out[0], create)
		order = dep.cmdMaker.StartSidecar(service, sidecar)

		start, err := command.NewCommand(order, fmt.Sprint(bucket))
		if err != nil {
			return nil, err
		}
		err = mergo.Map(&start.Meta, service.Labels)
		if err != nil {
			return nil, err
		}

		out[1] = append(out[1], start)
	}
	return out, nil
}

func (dep dependency) SidecarNetwork(bucket int, state *entity.State,
	service entity.Service) (command.Command, error) {

	return command.NewCommand(

		dep.cmdMaker.CreateSidecarNetwork(service,
			state.Subnets[service.Name]),
		fmt.Sprint(bucket))
}

func (dep dependency) RemoveContainer(bucket int, name string) (command.Command, error) {
	order := dep.cmdMaker.RemoveContainer(name)
	return command.NewCommand(order, fmt.Sprint(bucket))
}

func (dep dependency) Volumes(bucket int, service entity.Service) ([]command.Command, error) {
	out := []command.Command{}
	for _, volume := range service.SquashedService.SharedVolumes {
		order := dep.cmdMaker.CreateVolume(volume)
		cmd, err := command.NewCommand(order, fmt.Sprint(bucket))
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}
	return out, nil
}
