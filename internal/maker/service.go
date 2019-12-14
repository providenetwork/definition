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
	"time"

	"github.com/whiteblock/definition/internal/converter"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/internal/search"
	"github.com/whiteblock/definition/schema"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

type Service interface {
	FromSystemDiff(spec schema.RootSchema, system schema.SystemComponent,
		merged schema.SystemComponent) ([]entity.Service, error)

	FromSystem(spec schema.RootSchema, system schema.SystemComponent) ([]entity.Service, error)
	FromTask(spec schema.RootSchema, task schema.Task, index int) (entity.Service, error)
}

type serviceMaker struct {
	namer    parser.Names
	searcher search.Schema
	convert  converter.Service
	log      logrus.Ext1FieldLogger
}

func NewService(
	namer parser.Names,
	searcher search.Schema,
	convert converter.Service,
	log logrus.Ext1FieldLogger) Service {
	return &serviceMaker{namer: namer, searcher: searcher,
		convert: convert, log: log}
}

func (sp *serviceMaker) FromSystemDiff(spec schema.RootSchema,
	system schema.SystemComponent, merged schema.SystemComponent) ([]entity.Service, error) {

	if merged.Count == system.Count {
		return nil, nil
	}
	if merged.Count < system.Count {
		out := []entity.Service{} //We are removing nodes, so only need name
		for i := merged.Count; i < system.Count; i++ {
			out = append(out, entity.Service{Name: sp.namer.SystemService(merged, int(i))})
		}
		return out, nil
	}

	services, err := sp.FromSystem(spec, merged)
	if err != nil {
		return nil, err
	}
	services = services[int(merged.Count-system.Count):]
	return services, nil
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

	base := entity.Service{
		Name:            "",
		Bucket:          -1,
		SquashedService: squashed,
		Networks:        system.Resources.Networks,
		Sidecars:        sp.searcher.FindSidecarsByService(spec, system.Type),
		Timeout:         0,
		Ports:           portMapping,
	}
	if len(base.Networks) == 0 {
		base.Networks = []schema.Network{
			schema.Network{
				Name: sp.namer.DefaultNetwork(system),
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
		out[i].Name = sp.namer.SystemService(system, i)
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
		task.Networks = []schema.Network{
			schema.Network{Name: "none"},
		}
	}
	to := strings.Replace(task.Timeout, " ", "", -1)
	timeout, err := time.ParseDuration(to)
	return entity.Service{
		Name:            sp.namer.Task(task, index),
		Networks:        task.Networks,
		SquashedService: service,
		Sidecars:        nil,
		IgnoreExitCode:  task.IgnoreExitCode,
		Timeout:         timeout,
		IsTask:          true,
	}, err
}
