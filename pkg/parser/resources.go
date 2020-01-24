/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/whiteblock/definition/pkg/converter"
	"github.com/whiteblock/definition/pkg/search"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/pkg/namer"
	"github.com/whiteblock/definition/schema"
)

// Resources presents methods for extracting named resources from parts of the schema
type Resources interface {
	FromSystemDiff(spec schema.RootSchema,
		system schema.SystemComponent, merged schema.SystemComponent) ([]entity.Segment, error)
	SystemComponentNamesOnly(sys schema.SystemComponent) []entity.Segment
	SystemComponent(spec schema.RootSchema, sys schema.SystemComponent) ([]entity.Segment, error)
	Tasks(spec schema.RootSchema, tasks []schema.Task) ([]entity.Segment, error)
}

type resources struct {
	searcher search.Schema
	conv     converter.Resource
}

// NewResources creates a new Resources
func NewResources(
	searcher search.Schema,
	conv converter.Resource) Resources {

	return &resources{
		searcher: searcher,
		conv:     conv,
	}
}

func (res *resources) FromSystemDiff(spec schema.RootSchema,
	system schema.SystemComponent, merged schema.SystemComponent) ([]entity.Segment, error) {

	if merged.GetCount() == system.GetCount() {
		return nil, nil
	}
	if merged.GetCount() < system.GetCount() {
		out := []entity.Segment{} //We are removing nodes, so only need name
		for i := merged.GetCount(); i < system.GetCount(); i++ {
			out = append(out, entity.Segment{Name: namer.SystemService(merged, int(i))})
		}
		return out, nil
	}

	services, err := res.SystemComponent(spec, merged)
	if err != nil {
		return nil, err
	}
	services = services[int(merged.GetCount()-system.GetCount()):]
	return services, nil
}

func (res *resources) SystemComponent(spec schema.RootSchema,
	sys schema.SystemComponent) ([]entity.Segment, error) {

	service, err := res.searcher.FindServiceByType(spec, sys.Type)
	if err != nil {
		return nil, err
	}
	out := make([]entity.Segment, sys.GetCount())
	for i := range out {
		out[i].Name = namer.SystemService(sys, i)
		resource, err := res.conv.FromResources(service.Resources)
		if err != nil {
			return nil, err
		}

		for _, pm := range sys.PortMappings {
			ports := strings.Split(pm, ":")
			if len(ports) != 2 {
				return nil, fmt.Errorf(`invalid port mapping "%s"`, pm)
			}

			hostPort, err := strconv.Atoi(ports[0])
			if err != nil {
				return nil, err
			}
			(&resource).InsertPorts(hostPort)
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

	out := entity.Segment{Name: namer.Task(task, index)}
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
	out := make([]entity.Segment, sys.GetCount())
	for i := range out {
		out[i] = entity.Segment{Name: namer.SystemService(sys, i)}
	}
	return out
}
