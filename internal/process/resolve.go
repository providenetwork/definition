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

type Resolve interface {
	CreateNetworks(systems []schema.SystemComponent,
		networkState entity.NetworkState) ([]command.Command, error)

	CreateServices(spec schema.RootSchema, networkState entity.NetworkState, dist entity.PhaseDist,
		services []entity.Service) ([][]command.Command, error)

	RemoveServices(dist entity.PhaseDist, services []entity.Service) ([][]command.Command, error)
}

type resolve struct {
	cmdMaker maker.Command
	deps     Dependency
	log      logrus.Ext1FieldLogger
}

func NewResolve(cmdMaker maker.Command, deps Dependency, log logrus.Ext1FieldLogger) Resolve {
	return &resolve{cmdMaker: cmdMaker, deps: deps, log: log}
}

func (resolver resolve) CreateNetworks(systems []schema.SystemComponent,
	networkState entity.NetworkState) ([]command.Command, error) {
	out := []command.Command{}
	for _, system := range systems {
		for _, network := range system.Resources.Networks {
			subnet, err := networkState.GetNextGlobal()
			if err != nil {
				return nil, err
			}
			order := resolver.cmdMaker.CreateNetwork(network.Name, subnet)
			cmd, err := command.NewCommand(order, "0")
			if err != nil {
				return nil, err
			}
			cmd.Meta["system"] = system.Name
			cmd.Meta["network"] = network.Name
			out = append(out, cmd)
		}
	}
	return out, nil
}

//instance -> images -> meta
func (resolver resolve) pullImages(allImages map[string]map[string]map[string]string) ([]command.Command, error) {

	out := []command.Command{}
	for instance, images := range allImages {
		for image, meta := range images {
			order := resolver.cmdMaker.PullImage(image)
			cmd, err := command.NewCommand(order, instance)
			if err != nil {
				return nil, err
			}
			mergo.Map(&cmd.Meta, meta)
			out = append(out, cmd)
		}

	}
	return out, nil
}

func (resolver resolve) CreateServices(spec schema.RootSchema, networkState entity.NetworkState,
	dist entity.PhaseDist, services []entity.Service) ([][]command.Command, error) {

	out := make([][]command.Command, 5)
	images := map[string]map[string]map[string]string{}
	for _, service := range services {

		createCmd, startCmd, err := resolver.deps.Container(spec, dist, service)
		if err != nil {
			return nil, err
		}

		if images[createCmd.Target.IP] == nil {
			images[createCmd.Target.IP] = map[string]map[string]string{}
		}
		images[createCmd.Target.IP][service.SquashedService.Image] = service.Labels

		if !service.IsTask && len(service.Sidecars) > 0 {
			sidecarCmds, err := resolver.deps.Sidecars(spec, dist, service)
			if err != nil {
				return nil, err
			}
			for _, cmd := range sidecarCmds[0] {
				if cmd.Order.Type != command.Createcontainer {
					continue
				}
				payload, ok := cmd.Order.Payload.(command.Container)
				if !ok {
					continue
				}
				if images[cmd.Target.IP] == nil {
					images[cmd.Target.IP] = map[string]map[string]string{}
				}
				images[cmd.Target.IP][payload.Image] = service.Labels
			}
			out[3] = append(out[3], sidecarCmds[0]...)
			out[4] = append(out[4], sidecarCmds[1]...)
			sidecarNetworkCmd, err := resolver.deps.SidecarNetwork(spec, networkState, dist, service)
			if err != nil {
				return nil, err
			}
			out[0] = append(out[0], sidecarNetworkCmd)
		}

		volumeCmds, err := resolver.deps.Volumes(spec, dist, service)
		if err != nil {
			return nil, err
		}

		emulationCmds, err := resolver.deps.Emulation(spec, dist, service)
		if err != nil {
			return nil, err
		}

		fileCmds, err := resolver.deps.Files(dist, service)
		if err != nil {
			return nil, err
		}

		out[0] = append(out[0], volumeCmds...)
		out[0] = append(out[0], fileCmds...)
		out[1] = append(out[1], createCmd)
		out[2] = append(out[2], startCmd)
		out[1] = append(out[1], emulationCmds...)
	}
	cmds, err := resolver.pullImages(images)
	if err != nil {
		return nil, err
	}
	out[0] = append(out[0], cmds...)

	return out, nil
}

func (resolver resolve) RemoveServices(dist entity.PhaseDist,
	services []entity.Service) ([][]command.Command, error) {

	if len(services) == 0 {
		return nil, nil
	}
	out := []command.Command{}
	for _, service := range services {
		order := resolver.cmdMaker.RemoveContainer(service)
		bucket := dist.FindBucket(service.Name)
		if bucket == -1 {
			return nil, fmt.Errorf("could not find bucket")
		}

		cmd, err := command.NewCommand(order, fmt.Sprint(bucket))
		if err != nil {
			return nil, err
		}
		mergo.Map(&cmd.Meta, service.Labels)
		out = append(out, cmd)
	}
	//  If needed, we can also add commands for removing volumes and networks
	return [][]command.Command{out}, nil
}
