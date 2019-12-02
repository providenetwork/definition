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

package search

import (
	"fmt"

	"github.com/whiteblock/definition/schema"
)

type Schema interface {
	FindServiceByType(spec schema.RootSchema, serviceType string) (schema.Service, error)
	FindSidecarByType(spec schema.RootSchema, sidecarType string) (schema.Sidecar, error)
	FindTaskRunnerByType(spec schema.RootSchema, taskRunnerType string) (schema.TaskRunner, error)
	FindSidecarsByService(spec schema.RootSchema, name string) []schema.Sidecar
}

type schemaSearcher struct {
}

func NewSchema() Schema {
	return &schemaSearcher{}
}

func (searcher schemaSearcher) FindServiceByType(spec schema.RootSchema,
	serviceType string) (schema.Service, error) {

	for _, service := range spec.Services {
		if service.Name == serviceType {
			return service, nil
		}
	}

	return schema.Service{}, fmt.Errorf(`could not find service "%s"`, serviceType)
}

func (searcher schemaSearcher) FindSidecarByType(spec schema.RootSchema,
	sidecarType string) (schema.Sidecar, error) {

	for _, sidecar := range spec.Sidecars {
		if sidecar.Name == sidecarType {
			return sidecar, nil
		}
	}
	return schema.Sidecar{}, fmt.Errorf(`could not find sidecar "%s"`, sidecarType)
}

func (searcher schemaSearcher) FindTaskRunnerByType(spec schema.RootSchema,
	taskRunnerType string) (schema.TaskRunner, error) {

	for _, taskRunner := range spec.TaskRunners {
		if taskRunner.Name == taskRunnerType {
			return taskRunner, nil
		}
	}

	return schema.TaskRunner{}, fmt.Errorf(`could not find task runner "%s"`, taskRunnerType)
}

func (searcher schemaSearcher) FindSidecarsByService(spec schema.RootSchema,
	name string) []schema.Sidecar {

	out := []schema.Sidecar{}
	for _, sidecar := range spec.Sidecars {
		for _, serviceName := range sidecar.SidecarTo {
			if serviceName == name {
				out = append(out, sidecar)
				break
			}
		}
	}
	return out
}
