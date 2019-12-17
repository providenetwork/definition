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
	"github.com/whiteblock/definition/internal/namer"
	"github.com/whiteblock/definition/schema"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"
)

type Resolve interface {
	CreateNetworks(networkState entity.NetworkState,
		networks []schema.Network, meta map[string]string) ([]command.Command, error)
	CreateSystemNetworks(systems []schema.SystemComponent,
		networkState entity.NetworkState) ([]command.Command, error)

	CreateServices(spec schema.RootSchema, networkState entity.NetworkState, dist entity.PhaseDist,
		services []entity.Service) ([][]command.Command, error)

	RemoveServices(dist entity.PhaseDist, services []entity.Service) ([][]command.Command, error)
	UpdateServices(dist entity.PhaseDist, services []entity.ServiceDiff) ([][]command.Command, error)
}

var (
	ErrBucketNotFound = errors.New("could not find bucket")
)

const (
	FirstInstance = "0"
)

type resolve struct {
	cmdMaker maker.Command
	deps     Dependency
	log      logrus.Ext1FieldLogger
}

func NewResolve(
	cmdMaker maker.Command,
	deps Dependency,
	log logrus.Ext1FieldLogger) Resolve {
	return &resolve{cmdMaker: cmdMaker, deps: deps, log: log}
}

func (resolver resolve) CreateNetworks(networkState entity.NetworkState,
	networks []schema.Network, meta map[string]string) ([]command.Command, error) {
	out := []command.Command{}
	for _, network := range networks {
		subnet, err := networkState.GetNextGlobal()
		if err != nil {
			return nil, err
		}
		order := resolver.cmdMaker.CreateNetwork(network.Name, subnet)
		cmd, err := command.NewCommand(order, FirstInstance)
		if err != nil {
			return nil, err
		}
		err = mergo.Map(&cmd.Meta, meta)
		if err != nil {
			return nil, err
		}
		out = append(out, cmd)
	}
	return out, nil
}

func (resolver resolve) CreateSystemNetworks(systems []schema.SystemComponent,
	networkState entity.NetworkState) ([]command.Command, error) {

	out := []command.Command{}
	for _, system := range systems {
		networks := system.Resources.Networks
		if len(networks) == 0 {
			networks = []schema.Network{
				schema.Network{Name: namer.DefaultNetwork(system)},
			}
		}
		cmds, err := resolver.CreateNetworks(networkState, networks, map[string]string{
			"system": system.Name,
		})
		if err != nil {
			return nil, err
		}
		out = append(out, cmds...)
	}
	return out, nil
}

//instance -> images -> meta
func (resolver resolve) pullImages(images *entity.ImageStore) ([]command.Command, error) {
	out := []command.Command{}
	return out, images.ForEach(func(instance, image string, meta map[string]string) error {
		order := resolver.cmdMaker.PullImage(image)
		cmd, err := command.NewCommand(order, instance)
		if err != nil {
			return err
		}
		mergo.Map(&cmd.Meta, meta)
		out = append(out, cmd)
		return nil
	})
}

func (resolver resolve) CreateServices(spec schema.RootSchema, networkState entity.NetworkState,
	dist entity.PhaseDist, services []entity.Service) ([][]command.Command, error) {

	out := make([][]command.Command, 5)
	images := &entity.ImageStore{}
	for _, service := range services {
		bucket := dist.FindBucket(service.Name)
		if bucket == -1 {
			return nil, ErrBucketNotFound
		}

		createCmd, startCmd, err := resolver.deps.Container(bucket, service)
		if err != nil {
			return nil, err
		}
		images.Insert(createCmd.Target.IP, service.SquashedService.Image, service.Labels)

		if !service.IsTask {
			if len(service.Sidecars) > 0 {
				sidecarCmds, err := resolver.deps.Sidecars(bucket, service, service.Sidecars)
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
					images.Insert(cmd.Target.IP, payload.Image, service.Labels)
				}
				out[3] = append(out[3], sidecarCmds[0]...)
				out[4] = append(out[4], sidecarCmds[1]...)
			}
			sidecarNetworkCmd, err := resolver.deps.SidecarNetwork(bucket, networkState, service)
			if err != nil {
				return nil, err
			}
			out[0] = append(out[0], sidecarNetworkCmd)
		}

		volumeCmds, err := resolver.deps.Volumes(bucket, service)
		if err != nil {
			return nil, err
		}
		attachCmds, err := resolver.deps.AttachNetworks(bucket, service.Name, service.Networks)
		if err != nil {
			return nil, err
		}
		out[3] = append(out[3], attachCmds...)

		if !service.IsTask {
			emulationCmds, err := resolver.deps.Emulation(bucket, service.Name, service.Networks)
			if err != nil {
				return nil, err
			}
			out[4] = append(out[4], emulationCmds...)
		}

		fileCmds, err := resolver.deps.Files(bucket, service)
		if err != nil {
			return nil, err
		}

		out[0] = append(out[0], volumeCmds...)
		out[3] = append(out[3], fileCmds...)
		out[1] = append(out[1], createCmd)
		out[2] = append(out[2], startCmd)

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
		bucket := dist.FindBucket(service.Name)
		if bucket == -1 {
			return nil, ErrBucketNotFound
		}

		order := resolver.cmdMaker.RemoveContainer(service.Name)
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

func (resolver resolve) UpdateServices(dist entity.PhaseDist,
	services []entity.ServiceDiff) ([][]command.Command, error) {

	out := make([][]command.Command, 3)

	for _, service := range services {
		bucket := dist.FindBucket(service.Name)
		if bucket == -1 {
			return nil, ErrBucketNotFound
		}

		if len(service.AddNetworks) > 0 {
			addNetworkCmds, err := resolver.deps.AttachNetworks(bucket, service.Name, service.AddNetworks)
			if err != nil {
				return nil, err
			}
			out[0] = append(out[0], addNetworkCmds...)

			emulationCmds, err := resolver.deps.Emulation(bucket, service.Name, service.AddNetworks)
			if err != nil {
				return nil, err
			}
			out[1] = append(out[1], emulationCmds...)
		}

		if len(service.UpdateNetworks) > 0 {
			resolver.log.WithField("service", service.Name).Info("updating networks")
			emulationCmds, err := resolver.deps.Emulation(bucket, service.Name, service.UpdateNetworks)
			if err != nil {
				return nil, err
			}
			out[0] = append(out[0], emulationCmds...)
		}

		if len(service.DetachNetworks) > 0 {
			addNetworkCmds, err := resolver.deps.DetachNetworks(bucket, service.Name, service.DetachNetworks)
			if err != nil {
				return nil, err
			}
			out[0] = append(out[0], addNetworkCmds...)
		}

		if len(service.AddSidecars) > 0 {
			images := &entity.ImageStore{}
			sidecarCmds, err := resolver.deps.Sidecars(bucket, *service.Parent, service.AddSidecars)
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
				images.Insert(cmd.Target.IP, payload.Image, service.Parent.Labels)
			}

			cmds, err := resolver.pullImages(images)
			if err != nil {
				return nil, err
			}
			out[0] = append(out[0], cmds...)
			out[1] = append(out[1], sidecarCmds[0]...)
			out[2] = append(out[2], sidecarCmds[1]...)
		}
		for _, sidecarToRemove := range service.RemoveSidecars {
			cmd, err := resolver.deps.RemoveContainer(bucket,
				namer.Sidecar(*service.Parent, sidecarToRemove))
			if err != nil {
				return nil, err
			}
			out[1] = append(out[1], cmd)
		}
	}
	return out, nil
}
