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

package maker

import (
	"fmt"
	"time"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/schema"

	"github.com/docker/docker/api/types/strslice"
)

//  Command handles the simple schema -> order conversions
type Command interface {
	CreateNetwork(name string, network entity.Network) command.Order
	CreateVolume(volume schema.SharedVolume) command.Order
	CreateContainer(service entity.Service) command.Order
	CreateSidecarNetwork(service entity.Service, network entity.Network) command.Order
	StartContainer(service entity.Service, isTask bool, timeout time.Duration) command.Order

	CreateSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order
	StartSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order
	PullImage(image string) command.Order
	File(name string, input schema.InputFile) command.Order
	AttachNetwork(service entity.Service, network schema.Network) command.Order
	Emulation(service entity.Service, network schema.Network) (command.Order, error)

	RemoveContainer(service entity.Service) command.Order
}

type commandMaker struct {
	service parser.Service
	sidecar parser.Sidecar
	network parser.Network
	namer   parser.Names
}

func NewCommand(
	service parser.Service,
	sidecar parser.Sidecar,
	network parser.Network,
	namer parser.Names) Command {

	return &commandMaker{
		service: service,
		sidecar: sidecar,
		namer:   namer,
		network: network,
	}
}

func (cmd commandMaker) createNetwork(name string, network entity.Network, global bool) command.Order {
	return command.Order{
		Type: command.Createnetwork,
		Payload: command.Network{
			Name:    name,
			Global:  global,
			Gateway: network.Gateway(),
			Subnet:  network.Subnet(),
		},
	}
}

func (cmd commandMaker) CreateNetwork(name string, network entity.Network) command.Order {
	return cmd.createNetwork(name, network, true)
}

func (cmd commandMaker) PullImage(image string) command.Order {
	return command.Order{
		Type: command.Pullimage,
		Payload: command.PullImage{
			Image: image,
		},
	}
}

func (cmd commandMaker) CreateSidecarNetwork(service entity.Service, network entity.Network) command.Order {
	return cmd.createNetwork(cmd.namer.SidecarNetwork(service), network, false)
}

func (cmd commandMaker) CreateVolume(volume schema.SharedVolume) command.Order {
	return command.Order{
		Type: command.Createvolume,
		Payload: command.Volume{
			Name:   volume.Name,
			Labels: map[string]string{},
		},
	}
}

func (cmd commandMaker) CreateContainer(service entity.Service) command.Order {
	return command.Order{
		Type: command.Createcontainer,
		Payload: command.Container{
			BoundCPUs:   nil, //NYI
			EntryPoint:  cmd.service.GetEntrypoint(service),
			Environment: service.SquashedService.Environment,
			Labels:      service.Labels,
			Name:        service.Name,
			Network:     strslice.StrSlice(cmd.service.GetNetworks(service)),
			Ports:       service.Ports,
			Volumes:     cmd.service.GetVolumes(service),
			Cpus:        fmt.Sprint(cmd.service.GetCPUs(service)),
			Memory:      fmt.Sprint(cmd.service.GetMemory(service)),
			Image:       cmd.service.GetImage(service),
			Args:        cmd.service.GetArgs(service),
		},
	}
}

func (cmd commandMaker) StartContainer(service entity.Service, isTask bool, timeout time.Duration) command.Order {
	return cmd.startContainer(service.Name, isTask, timeout)
}

func (cmd commandMaker) CreateSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order {
	return command.Order{
		Type: command.Createcontainer,
		Payload: command.Container{
			EntryPoint:  cmd.sidecar.GetEntrypoint(sidecar),
			Environment: sidecar.Environment,
			Labels:      cmd.sidecar.GetLabels(parent, sidecar),
			Name:        cmd.namer.Sidecar(parent, sidecar),
			Network:     strslice.StrSlice(cmd.sidecar.GetNetwork(parent)),
			Volumes:     cmd.sidecar.GetVolumes(sidecar),
			Cpus:        cmd.sidecar.GetCPUs(sidecar),
			Memory:      cmd.sidecar.GetMemory(sidecar),
			Image:       cmd.sidecar.GetImage(sidecar),
			Args:        cmd.sidecar.GetArgs(sidecar),
			Ports:       parent.Ports,
		},
	}
}

func (cmd commandMaker) File(name string, input schema.InputFile) command.Order {
	return command.Order{

		Type: command.Putfileincontainer,
		Payload: command.FileAndContainer{
			ContainerName: name,
			File: command.File{
				Mode:        0644,
				Destination: input.DestinationPath,
				ID:          input.GetSource(),
			},
		},
	}
}

func (cmd commandMaker) StartSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order {
	return cmd.startContainer(cmd.namer.Sidecar(parent, sidecar), false, 0)
}

func (cmd commandMaker) AttachNetwork(service entity.Service, network schema.Network) command.Order {
	return command.Order{
		Type: command.Attachnetwork,
		Payload: command.ContainerNetwork{
			ContainerName: service.Name,
			Network:       cmd.namer.Network(network),
		},
	}
}

func (cmd commandMaker) Emulation(service entity.Service, network schema.Network) (command.Order, error) {

	loss, err := cmd.network.GetPacketLoss(network)
	if err != nil {
		return command.Order{}, err
	}
	delay, err := cmd.network.GetLatency(network)
	if err != nil {
		return command.Order{}, err
	}
	return command.Order{
		Type: command.Emulation,
		Payload: command.Netconf{
			Container:   service.Name,
			Network:     cmd.namer.Network(network),
			Limit:       0, //NYI
			Loss:        loss,
			Delay:       int(delay),
			Rate:        cmd.network.GetBandwidth(network),
			Duplication: 0.0, //NYI
			Corrupt:     0.0, //NYI
			Reorder:     0.0, //NYI
		},
	}, nil
}

func (cmd commandMaker) RemoveContainer(service entity.Service) command.Order {
	return command.Order{
		Type: command.Removecontainer,
		Payload: command.SimpleName{
			Name: service.Name,
		},
	}
}

func (cmd commandMaker) startContainer(name string, isTask bool, timeout time.Duration) command.Order {
	return command.Order{
		Type: command.Startcontainer,
		Payload: command.StartContainer{
			Name:    name,
			Attach:  isTask,
			Timeout: timeout,
		},
	}
}
