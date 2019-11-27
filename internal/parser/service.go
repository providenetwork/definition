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
	//"github.com/imdario/mergo"
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
	//TODO: nate left off here
	if len(squashed.SharedVolumes) == 0 { //just make it compile
		return nil, nil
	}

	return nil, nil

}

func (sp *serviceParser) FromTask(spec schema.RootSchema,
	task schema.Task, index int) ([]entity.Service, error) {

	return nil, nil
}

func (sp *serviceParser) FromSidecar(parent entity.Service,
	sidecar schema.Sidecar) (entity.Service, error) {
	return entity.Service{}, nil
}
