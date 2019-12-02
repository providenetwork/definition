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
	"github.com/whiteblock/definition/internal/distribute"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/maker"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
)

type Resolve interface {
	CreateNetworks(systems []schema.SystemComponent) ([]command.Command, error)

	CreateServices(spec schema.RootSchema, dist distribute.PhaseDist,
		services []entity.Service) ([][]command.Command, error)

	RemoveServices(dist distribute.PhaseDist, services []entity.Service) ([][]command.Command, error)
}

type resolve struct {
	cmdMaker maker.Command
	deps     Dependency
	log      logrus.Ext1FieldLogger
}

func NewResolve(cmdMaker maker.Command, deps Dependency, log logrus.Ext1FieldLogger) Resolve {
	return &resolve{cmdMaker: cmdMaker, deps: deps, log: log}
}

func (resolver resolve) CreateNetworks(systems []schema.SystemComponent) ([]command.Command, error) {
	out := []command.Command{}
	for _, system := range systems {
		for _, network := range system.Resources.Networks {
			order := resolver.cmdMaker.CreateNetwork(network, true)
			cmd, err := resolver.cmdMaker.New(order, "0", 0)
			if err != nil {
				return nil, err
			}
			out = append(out, cmd)
		}
	}
	return out, nil
}

func (resolver resolve) CreateServices(spec schema.RootSchema,
	dist distribute.PhaseDist, services []entity.Service) ([][]command.Command, error) {

	out := make([][]command.Command, 5)
	for _, service := range services {

		createCmd, startCmd, err := resolver.deps.Container(spec, dist, service)
		if err != nil {
			return nil, err
		}

		sidecarCmds, err := resolver.deps.Sidecars(spec, dist, service)
		if err != nil {
			return nil, err
		}

		volumeCmds, err := resolver.deps.Volumes(spec, dist, service)
		if err != nil {
			return nil, err
		}

		sidecarNetworkCmd, err := resolver.deps.SidecarNetwork(spec, dist, service)
		if err != nil {
			return nil, err
		}

		emulationCmds, err := resolver.deps.Emulation(spec, dist, service)
		if err != nil {
			return nil, err
		}

		out[0] = append(out[0], volumeCmds...)
		out[0] = append(out[0], sidecarNetworkCmd)
		out[1] = append(out[1], createCmd)
		out[2] = append(out[2], startCmd)
		out[3] = append(out[3], sidecarCmds[0]...)
		out[3] = append(out[3], emulationCmds...)
		out[4] = append(out[4], sidecarCmds[1]...)
	}

	return out, nil
}

func (resolver resolve) RemoveServices(dist distribute.PhaseDist,
	services []entity.Service) ([][]command.Command, error) {

	out := []command.Command{}
	for _, service := range services {
		order := resolver.cmdMaker.RemoveContainer(service)
		bucket := dist.FindBucket(service.Name)
		if bucket == -1 {
			return nil, fmt.Errorf("could not find bucket")
		}

		cmd, err := resolver.cmdMaker.New(order, fmt.Sprint(bucket), 0)
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}
	//If needed, we can also add commands for removing volumes and networks
	return [][]command.Command{out}, nil
}
