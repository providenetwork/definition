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

package parser

import (
	"github.com/whiteblock/definition/internal/converter"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/search"
	"github.com/whiteblock/definition/schema"
)

// Resources presents methods for extracting named resources from parts of the schema
type Resources interface {
	SystemComponentNamesOnly(sys schema.SystemComponent) []entity.Segment
	SystemComponent(spec schema.RootSchema, sys schema.SystemComponent) ([]entity.Segment, error)
	Tasks(spec schema.RootSchema, tasks []schema.Task) ([]entity.Segment, error)
}

type resources struct {
	namer    Names
	searcher search.Schema
	conv     converter.Resource
}

// NewResources creates a new Resources
func NewResources(
	namer Names,
	searcher search.Schema,
	conv converter.Resource) Resources {

	return &resources{
		namer:    namer,
		searcher: searcher,
		conv:     conv,
	}
}

func (res *resources) SystemComponent(spec schema.RootSchema,
	sys schema.SystemComponent) ([]entity.Segment, error) {

	service, err := res.searcher.FindServiceByType(spec, sys.Type)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Segment, sys.Count)
	for i := range out {
		out[i].Name = res.namer.SystemService(sys, i)
		resource, err := res.conv.FromResources(service.Resources)
		if err != nil {
			return nil, err
		}
		out[i].UpdateResources(resource)
	}

	return out, nil
}

func (res *resources) task(spec schema.RootSchema, task schema.Task, index int) (entity.Segment, error) {
	taskRunner, err := res.searcher.FindTaskRunnerByType(spec, task.Type)
	if err != nil {
		return entity.Segment{}, err
	}

	out := entity.Segment{Name: res.namer.Task(task, index)}
	resource, err := res.conv.FromResources(taskRunner.Resources)
	if err != nil {
		return entity.Segment{}, err
	}
	out.UpdateResources(resource)
	return out, nil
}

func (res *resources) Tasks(spec schema.RootSchema, tasks []schema.Task) ([]entity.Segment, error) {
	out := make([]entity.Segment, len(tasks))
	for i := range out {
		segment, err := res.task(spec, tasks[i], i)
		if err != nil {
			return nil, err
		}
		out[i] = segment
	}
	return out, nil
}

func (res *resources) SystemComponentNamesOnly(sys schema.SystemComponent) []entity.Segment {
	out := make([]entity.Segment, sys.Count)
	for i := range out {
		out[i] = entity.Segment{Name: res.namer.SystemService(sys, i)}
	}
	return out
}
