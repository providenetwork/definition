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
package process

import (
	"fmt"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/internal/distribute"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/maker"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
)

type Dependency interface {
	Container(spec schema.RootSchema, dist distribute.PhaseDist,
		service entity.Service) (create command.Command, start command.Command, err error)

	Emulation(spec schema.RootSchema, dist distribute.PhaseDist,
		service entity.Service) ([]command.Command, error)

	Sidecars(spec schema.RootSchema, dist distribute.PhaseDist,
		service entity.Service) ([][]command.Command, error)

	Volumes(spec schema.RootSchema, dist distribute.PhaseDist,
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

func (dep dependency) Emulation(spec schema.RootSchema, dist distribute.PhaseDist,
	service entity.Service) ([]command.Command, error) {

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
		cmd, err := dep.cmdMaker.New(order, fmt.Sprint(bucket), 0)
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}
	return out, nil
}

func (dep dependency) Container(spec schema.RootSchema, dist distribute.PhaseDist,
	service entity.Service) (create command.Command, start command.Command, err error) {

	bucket := dist.FindBucket(service.Name)
	if bucket == -1 {
		err = fmt.Errorf("could not find bucket")
		return
	}

	order := dep.cmdMaker.CreateContainer(service)

	create, err = dep.cmdMaker.New(order, fmt.Sprint(bucket), 0)
	if err != nil {
		return
	}

	order = dep.cmdMaker.StartContainer(service)

	start, err = dep.cmdMaker.New(order, fmt.Sprint(bucket), service.Timeout)
	return
}

func (dep dependency) Sidecars(spec schema.RootSchema, dist distribute.PhaseDist,
	service entity.Service) ([][]command.Command, error) {

	bucket := dist.FindBucket(service.Name)
	if bucket == -1 {
		return nil, fmt.Errorf("could not find bucket")
	}
	out := make([][]command.Command, 2)
	for _, sidecar := range service.Sidecars {

		order := dep.cmdMaker.CreateSidecar(service, sidecar)
		create, err := dep.cmdMaker.New(order, fmt.Sprint(bucket), 0)
		if err != nil {
			return nil, err
		}
		out[0] = append(out[0], create)
		order = dep.cmdMaker.StartSidecar(service, sidecar)

		start, err := dep.cmdMaker.New(order, fmt.Sprint(bucket), 0)
		if err != nil {
			return nil, err
		}
		out[1] = append(out[1], start)
	}
	return out, nil
}

func (dep dependency) Volumes(spec schema.RootSchema, dist distribute.PhaseDist,
	service entity.Service) ([]command.Command, error) {

	bucket := dist.FindBucket(service.Name)
	if bucket == -1 {
		return nil, fmt.Errorf("could not find bucket")
	}
	out := []command.Command{}
	for _, volume := range service.SquashedService.SharedVolumes {
		order := dep.cmdMaker.CreateVolume(volume)
		cmd, err := dep.cmdMaker.New(order, fmt.Sprint(bucket), 0)
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}
	return out, nil
}
