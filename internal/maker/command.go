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

package maker

import (
	"fmt"
	"time"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/schema"

	"github.com/docker/docker/api/types/strslice"
	"github.com/google/uuid"
)

//Command handles the simple schema -> order conversions
type Command interface {
	New(order command.Order, endpoint string, timeout time.Duration) (command.Command, error)

	CreateNetwork(network schema.Network) command.Order
	CreateVolume(volume schema.SharedVolume) command.Order
	CreateContainer(service entity.Service) command.Order
	StartContainer(service entity.Service) command.Order

	CreateSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order
	StartSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order

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

func (cmd commandMaker) New(order command.Order, endpoint string,
	timeout time.Duration) (command.Command, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return command.Command{}, err
	}
	return command.Command{
		ID:        id.String(),
		Timestamp: time.Now().Unix(),
		Timeout:   timeout,
		Target: command.Target{
			IP: endpoint,
		},
		Order: order,
	}, nil
}

func (cmd commandMaker) CreateNetwork(network schema.Network) command.Order {
	return command.Order{
		Type: command.Createnetwork,
		Payload: command.Network{
			Name:    network.Name,
			Subnet:  "",                  //TODO
			Gateway: "",                  //TODO
			Global:  true,                //TODO
			Labels:  map[string]string{}, //TODO
		},
	}
}

func (cmd commandMaker) CreateVolume(volume schema.SharedVolume) command.Order {
	return command.Order{
		Type: command.Createvolume,
		Payload: command.Volume{
			Name:   volume.Name,
			Labels: map[string]string{}, //TODO
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
			Ports:       nil, //NYI
			Volumes:     cmd.service.GetVolumes(service),
			Cpus:        fmt.Sprint(service.SquashedService.Resources.Cpus),
			Memory:      fmt.Sprint(service.SquashedService.Resources.Memory),
			Image:       cmd.service.GetImage(service),
			Args:        service.SquashedService.Args,
		},
	}
}

func (cmd commandMaker) StartContainer(service entity.Service) command.Order {
	return cmd.startContainer(service.Name)
}

func (cmd commandMaker) CreateSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order {
	return command.Order{
		Type: command.Createcontainer,
		Payload: command.Container{
			EntryPoint:  cmd.sidecar.GetEntrypoint(sidecar),
			Environment: sidecar.Environment,
			Labels:      cmd.sidecar.GetLabels(parent, sidecar),
			Name:        sidecar.Name,
			Network:     strslice.StrSlice(cmd.sidecar.GetNetwork(parent)),
			Volumes:     cmd.sidecar.GetVolumes(sidecar),
			Cpus:        fmt.Sprint(sidecar.Resources.Cpus),
			Memory:      fmt.Sprint(sidecar.Resources.Memory),
			Image:       cmd.sidecar.GetImage(sidecar),
			Args:        sidecar.Args,
		},
	}
}

func (cmd commandMaker) StartSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order {
	return cmd.startContainer(cmd.namer.Sidecar(parent, sidecar))
}

func (cmd commandMaker) AttachNetwork(service entity.Service, network schema.Network) command.Order {
	return command.Order{
		Type: command.Attachnetwork,
		Payload: command.ContainerNetwork{
			ContainerName: service.Name,
			Network:       network.Name,
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
			Network:     network.Name,
			Limit:       0, //NYI
			Loss:        loss,
			Delay:       delay,
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

func (cmd commandMaker) startContainer(name string) command.Order {
	return command.Order{
		Type: command.Startcontainer,
		Payload: command.StartContainer{
			Name:   name,
			Attach: false,
		},
	}
}
