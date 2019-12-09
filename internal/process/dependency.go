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
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/maker"
	"github.com/whiteblock/definition/schema"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
)

type Dependency interface {
	Container(spec schema.RootSchema, dist entity.PhaseDist,
		service entity.Service) (create command.Command, start command.Command, err error)

	Emulation(spec schema.RootSchema, dist entity.PhaseDist,
		service entity.Service) ([]command.Command, error)

	Sidecars(spec schema.RootSchema, dist entity.PhaseDist,
		service entity.Service) ([][]command.Command, error)

	SidecarNetwork(spec schema.RootSchema, networkState entity.NetworkState,
		dist entity.PhaseDist, service entity.Service) (command.Command, error)

	Volumes(spec schema.RootSchema, dist entity.PhaseDist,
		service entity.Service) ([]command.Command, error)
}

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

func (dep dependency) Emulation(spec schema.RootSchema, dist entity.PhaseDist,
	service entity.Service) ([]command.Command, error) {
	if service.IsTask {
		return []command.Command{}, nil
	}
	bucket := dist.FindBucket(service.Name)
	if bucket == -1 {
		return nil, fmt.Errorf("could not find bucket")
	}
	out := []command.Command{}
	for _, network := range service.Networks {
		order, err := dep.cmdMaker.Emulation(service, network)
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

func (dep dependency) Container(spec schema.RootSchema, dist entity.PhaseDist,
	service entity.Service) (create command.Command, start command.Command, err error) {

	bucket := dist.FindBucket(service.Name)
	if bucket == -1 {
		err = fmt.Errorf("could not find bucket")
		return
	}

	order := dep.cmdMaker.CreateContainer(service)

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

func (dep dependency) Sidecars(spec schema.RootSchema, dist entity.PhaseDist,
	service entity.Service) ([][]command.Command, error) {

	bucket := dist.FindBucket(service.Name)
	if bucket == -1 {
		return nil, fmt.Errorf("could not find bucket")
	}
	out := make([][]command.Command, 2)
	for _, sidecar := range service.Sidecars {

		order := dep.cmdMaker.CreateSidecar(service, sidecar)
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

func (dep dependency) SidecarNetwork(spec schema.RootSchema, networkState entity.NetworkState,
	dist entity.PhaseDist, service entity.Service) (command.Command, error) {

	bucket := dist.FindBucket(service.Name)
	if bucket == -1 {
		return command.Command{}, fmt.Errorf("could not find bucket")
	}
	subnet, err := networkState.GetNextLocal(bucket)
	if err != nil {
		return command.Command{}, err
	}
	order := dep.cmdMaker.CreateSidecarNetwork(service, subnet)
	return command.NewCommand(order, fmt.Sprint(bucket))
}

func (dep dependency) Volumes(spec schema.RootSchema, dist entity.PhaseDist,
	service entity.Service) ([]command.Command, error) {

	bucket := dist.FindBucket(service.Name)
	if bucket == -1 {
		return nil, fmt.Errorf("could not find bucket")
	}
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
