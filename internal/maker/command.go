/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package maker

import (
	"fmt"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/internal/namer"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"
)

// Command handles the simple schema -> order conversions
type Command interface {
	CreateNetwork(name string, network entity.Network) command.Order
	CreateVolume(volume schema.Volume, global bool) command.Order
	CreateContainer(state *entity.State, service entity.Service) command.Order
	CreateSidecarNetwork(service entity.Service, network entity.Network) command.Order
	StartContainer(service entity.Service, isTask bool, timeout command.Timeout) command.Order

	CreateSidecar(state *entity.State, parent entity.Service, sidecar schema.Sidecar) command.Order
	StartSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order
	PullImage(image string) command.Order
	File(name string, input schema.InputFile) command.Order
	AttachNetwork(service string, network string, ip string) command.Order
	DetachNetwork(service string, network string) command.Order
	Emulation(serviceName string, network schema.Network) (command.Order, error)

	RemoveContainer(name string) command.Order

	Mkdir(name string) command.Order
}

type commandMaker struct {
	service parser.Service
	sidecar parser.Sidecar
	network parser.Network
}

func NewCommand(
	service parser.Service,
	sidecar parser.Sidecar,
	network parser.Network) Command {

	return &commandMaker{
		service: service,
		sidecar: sidecar,
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
	return cmd.createNetwork(namer.Network(name), network, true)
}

func (cmd commandMaker) PullImage(image string) command.Order {
	return command.Order{
		Type: command.Pullimage,
		Payload: command.PullImage{
			Image: image,
		},
	}
}

func (cmd commandMaker) CreateSidecarNetwork(service entity.Service,
	network entity.Network) command.Order {
	return cmd.createNetwork(namer.SidecarNetwork(service), network, false)
}

func (cmd commandMaker) CreateVolume(volume schema.Volume, global bool) command.Order {
	return CreateVolumeOrder(volume.Name, global)
}

func (cmd commandMaker) Mkdir(name string) command.Order {
	return MkdirOrder(name)
}

func (cmd commandMaker) CreateContainer(state *entity.State, service entity.Service) command.Order {
	return command.Order{
		Type: command.Createcontainer,
		Payload: command.Container{
			BoundCPUs:   nil, //NYI
			EntryPoint:  cmd.service.GetEntrypoint(service),
			Environment: parser.GetServiceEnv(service),
			Labels:      service.Labels,
			Name:        service.Name,
			Network:     cmd.service.GetNetwork(service),
			Ports:       service.Ports,
			Volumes:     parser.GetVolumes(service, service.SquashedService.Volumes),
			Cpus:        fmt.Sprint(cmd.service.GetCPUs(service)),
			Memory:      fmt.Sprint(cmd.service.GetMemory(service)),
			Image:       cmd.service.GetImage(service),

			Args: cmd.service.GetArgs(service),
			IP:   cmd.service.GetIP(state, service),
		},
	}
}

func (cmd commandMaker) StartContainer(service entity.Service, isTask bool,
	timeout command.Timeout) command.Order {
	return cmd.startContainer(service.Name, isTask, timeout)
}

func (cmd commandMaker) CreateSidecar(state *entity.State, parent entity.Service,
	sidecar schema.Sidecar) command.Order {

	return command.Order{
		Type: command.Createcontainer,
		Payload: command.Container{
			EntryPoint:  cmd.sidecar.GetEntrypoint(sidecar),
			Environment: parser.GetSidecarEnv(sidecar, parent, state),
			Labels:      cmd.sidecar.GetLabels(parent, sidecar),
			Name:        namer.Sidecar(parent, sidecar),
			Network:     cmd.sidecar.GetNetwork(parent),
			Volumes:     parser.GetVolumes(parent, sidecar.Volumes),
			Cpus:        cmd.sidecar.GetCPUs(sidecar),
			Memory:      cmd.sidecar.GetMemory(sidecar),
			Image:       cmd.sidecar.GetImage(sidecar),
			Args:        cmd.sidecar.GetArgs(sidecar),
			Ports:       parent.Ports,
			IP:          cmd.sidecar.GetIP(state, parent, sidecar),
		},
	}
}
func (cmd commandMaker) File(name string, input schema.InputFile) command.Order {
	return FileOrder(name, input)
}

func (cmd commandMaker) StartSidecar(parent entity.Service, sidecar schema.Sidecar) command.Order {
	return cmd.startContainer(namer.Sidecar(parent, sidecar), false, command.Timeout{})
}

func (cmd commandMaker) AttachNetwork(service string, network string, ip string) command.Order {
	return command.Order{
		Type: command.Attachnetwork,
		Payload: command.ContainerNetwork{
			Container: service,
			Network:   namer.Network(network),
			IP:        ip,
		},
	}
}

func (cmd commandMaker) DetachNetwork(service string, network string) command.Order {
	return command.Order{
		Type: command.Detachnetwork,
		Payload: command.ContainerNetwork{
			Container: service,
			Network:   namer.Network(network),
		},
	}
}

func (cmd commandMaker) Emulation(serviceName string, network schema.Network) (command.Order, error) {

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
			Container:   serviceName,
			Network:     network.Name,
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

func (cmd commandMaker) RemoveContainer(name string) command.Order {
	return command.Order{
		Type: command.Removecontainer,
		Payload: command.SimpleName{
			Name: name,
		},
	}
}

func (cmd commandMaker) startContainer(name string, isTask bool,
	timeout command.Timeout) command.Order {

	return command.Order{
		Type: command.Startcontainer,
		Payload: command.StartContainer{
			Name:    name,
			Attach:  isTask,
			Timeout: timeout,
		},
	}
}

func FileOrder(name string, input schema.InputFile) command.Order {
	return command.Order{
		Type: command.Putfileincontainer,
		Payload: command.FileAndContainer{
			ContainerName: name,
			File: command.File{
				Mode:        input.Mode(),
				Destination: input.Destination(),
				ID:          input.Source(),
			},
		},
	}
}

func CreateVolumeOrder(name string, global bool) command.Order {
	return command.Order{
		Type: command.Createvolume,
		Payload: command.Volume{
			Name:   name,
			Global: global,
		},
	}
}

func MkdirOrder(name string) command.Order {
	return CreateVolumeOrder(name, false)
}
