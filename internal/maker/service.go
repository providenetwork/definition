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
	"strconv"
	"strings"

	"github.com/whiteblock/definition/config/defaults"
	"github.com/whiteblock/definition/internal/converter"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/namer"
	"github.com/whiteblock/definition/internal/search"
	"github.com/whiteblock/definition/schema"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

type Service interface {
	FromSystemDiff(spec schema.RootSchema, system schema.SystemComponent,
		merged schema.SystemComponent) (*entity.SystemDiff, error)

	FromSystem(spec schema.RootSchema, system schema.SystemComponent) ([]entity.Service, error)
	FromTask(spec schema.RootSchema, task schema.Task, index int) (entity.Service, error)
}

type serviceMaker struct {
	searcher search.Schema
	convert  converter.Service
	log      logrus.Ext1FieldLogger
	defaults defaults.Defaults
}

func NewService(
	defaults defaults.Defaults,
	searcher search.Schema,
	convert converter.Service,
	log logrus.Ext1FieldLogger) Service {
	return &serviceMaker{searcher: searcher, defaults: defaults,
		convert: convert, log: log}
}

func (sp *serviceMaker) FromSystemDiff(spec schema.RootSchema,
	system schema.SystemComponent, merged schema.SystemComponent) (*entity.SystemDiff, error) {

	out := &entity.SystemDiff{}

	networkStatus := map[string]int{}
	networks := map[string]schema.Network{}
	if len(system.Resources.Networks) == 0 {
		networkStatus[sp.defaults.Network.Name] |= 0x01
		networks[sp.defaults.Network.Name] = schema.Network{Name: sp.defaults.Network.Name}
	}
	for _, network := range system.Resources.Networks {
		networkStatus[network.Name] |= 0x01
		networks[network.Name] = network
	}

	for _, network := range merged.Resources.Networks {
		networkStatus[network.Name] |= 0x02
		networks[network.Name] = network
	}
	sp.log.WithField("networkStatus", networkStatus).Trace("got the network status")

	for networkName, status := range networkStatus {
		switch status {
		case 0x03:
			continue
		case 0x02:
			out.AddedNetworks = append(out.AddedNetworks, networks[networkName])
		case 0x01:
			out.RemovedNetworks = append(out.RemovedNetworks, networks[networkName])
		}
	}

	services, err := sp.FromSystem(spec, merged)
	if err != nil {
		return nil, err
	}

	oldServices, err := sp.FromSystem(spec, system)
	if err != nil {
		return nil, err
	}
	oldLength := len(oldServices)
	newLength := len(services)

	max := newLength
	if oldLength > max {
		max = oldLength
	}

	for i := 0; i < max; i++ {
		if i >= oldLength {
			out.Added = append(out.Added, services[i])
		} else if i >= newLength {
			out.Removed = append(out.Removed, oldServices[i])
		} else if !services[i].Equal(oldServices[i]) {
			out.Modified = append(out.Modified, services[i].CalculateDiff(oldServices[i]))
		}
	}

	return out, nil
}

func (sp *serviceMaker) FromSystem(spec schema.RootSchema,
	system schema.SystemComponent) ([]entity.Service, error) {
	sp.log.WithField("system", system).Debug("creating service entities from system")
	squashed, err := sp.searcher.FindServiceByType(spec, system.Type)
	if err != nil {
		return nil, err
	}

	err = mergo.Map(&squashed.Environment, system.Environment, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	if system.Args != nil {
		squashed.Args = system.Args
	}

	if system.Resources.Cpus != 0 {
		squashed.Resources.Cpus = system.Resources.Cpus
	}

	if system.Resources.Memory != "" {
		squashed.Resources.Memory = system.Resources.Memory
	}

	if system.Resources.Storage != "" {
		squashed.Resources.Storage = system.Resources.Storage
	}

	portMapping := map[int]int{}
	for _, pm := range system.PortMappings {
		ports := strings.Split(pm, ":")
		if len(ports) != 2 {
			return nil, fmt.Errorf(`invalid port mapping "%s"`, pm)
		}

		hostPort, err := strconv.Atoi(ports[0])
		if err != nil {
			return nil, err
		}

		cntrPort, err := strconv.Atoi(ports[1])
		if err != nil {
			return nil, err
		}
		sp.log.WithFields(logrus.Fields{
			"host":      hostPort,
			"container": cntrPort,
		}).Debug("processed port mapping")
		portMapping[hostPort] = cntrPort
	}
	base := entity.GetDefaultService(sp.defaults)

	base.Networks = system.Resources.Networks
	base.Sidecars = sp.searcher.FindSidecarsByService(spec, system.Type)
	base.Ports = portMapping
	err = mergo.Map(&base.SquashedService, squashed, mergo.WithOverride)
	if err != nil {
		return nil, err
	}

	if len(base.Networks) == 0 {
		base.Networks = []schema.Network{
			schema.Network{
				Name: namer.DefaultNetwork(system),
			},
		}
	}

	for _, sidecar := range system.Sidecars {
		realSidecar, err := sp.searcher.FindSidecarByType(spec, sidecar.Type)
		if err != nil {
			return nil, err
		}

		err = mergo.Map(&realSidecar.Environment, sidecar.Environment, mergo.WithOverride)
		if err != nil {
			return nil, err
		}

		err = mergo.Map(&realSidecar.Resources, sidecar.Resources, mergo.WithOverride)
		if err != nil {
			return nil, err
		}

		if sidecar.Args != nil {
			realSidecar.Args = sidecar.Args
		}
		base.Sidecars = append(base.Sidecars, realSidecar)
	}

	out := make([]entity.Service, system.Count)

	for i := range out {
		copier.Copy(&out[i], base)
		out[i].Name = namer.SystemService(system, i)
	}
	return out, nil
}

func (sp *serviceMaker) FromTask(spec schema.RootSchema,
	task schema.Task, index int) (entity.Service, error) {

	taskRunner, err := sp.searcher.FindTaskRunnerByType(spec, task.Type)
	if err != nil {
		return entity.Service{}, err
	}

	service := sp.convert.FromTaskRunner(taskRunner)
	if task.Args != nil {
		copier.Copy(&service.Args, task.Args)
	}

	if task.Environment != nil {
		err = mergo.Map(&service.Environment, task.Environment, mergo.WithOverride)
		if err != nil {
			return entity.Service{}, err
		}
	}
	if len(task.Networks) == 0 {
		task.Networks = []schema.Network{}
	}
	out := entity.GetDefaultService(sp.defaults)
	return out, mergo.Map(&out, entity.Service{
		Name:            namer.Task(task, index),
		Networks:        task.Networks,
		SquashedService: service,
		Sidecars:        nil,
		IgnoreExitCode:  task.IgnoreExitCode,
		Timeout:         task.Timeout,
		IsTask:          true,
	}, mergo.WithOverride)
}
