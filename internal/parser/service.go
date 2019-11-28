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

package parser

import (
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/search"
	"github.com/whiteblock/definition/schema"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
)

type Service interface {
	FromSystem(spec schema.RootSchema, system schema.SystemComponent) ([]entity.Service, error)
	FromTask(spec schema.RootSchema, task schema.Task, index int) ([]entity.Service, error)
	FromSidecar(parent entity.Service, sidecar schema.Sidecar) (entity.Service, error)
}

type serviceParser struct {
	namer    Names
	searcher search.Schema
}

func NewService(namer Names, searcher search.Schema) Service {
	return &serviceParser{namer: namer, searcher: searcher}
}

func (sp *serviceParser) FromSystem(spec schema.RootSchema,
	system schema.SystemComponent) ([]entity.Service, error) {

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
	systemName := ps.namer.SystemComponent(system)
	externalSidecars, err := sp.searcher.FindSidecarsBySystem(spec, systemName)
	if err != nil {
		return nil, err
	}

	base := entity.Service{
		Name:            "",
		Bucket:          -1,
		SquashedService: squashed,
		Networks:        system.Resources.Networks,
		Sidecars:        externalSidecars,
		Timeout:         0,
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
		out[i].Name = ps.namer.SystemService(system, i)
	}
	return out, nil

}

func (sp *serviceParser) FromTask(spec schema.RootSchema,
	task schema.Task, index int) ([]entity.Service, error) {

	return nil, nil
}

func (sp *serviceParser) FromSidecar(parent entity.Service,
	sidecar schema.Sidecar) (entity.Service, error) {
	return entity.Service{}, nil
}
