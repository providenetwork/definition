/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
