/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package schema

import "fmt"

type SystemComponentSidecar struct {
	Type        string            `yaml:"type,omitempty" json:"type,omitempty"`
	Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
	Resources   Resources         `yaml:"resources,omitempty" json:"resources,omitempty"`
	Args        []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
}

type SystemComponentResources struct {
	Cpus     int64     `yaml:"cpus,omitempty" json:"cpus,omitempty"`
	Memory   string    `yaml:"memory,omitempty" json:"memory,omitempty"`
	Storage  string    `yaml:"storage,omitempty" json:"storage,omitempty"`
	Networks []Network `yaml:"networks,omitempty" json:"networks,omitempty"`
}

type SystemComponent struct {
	Type         string                   `yaml:"type,omitempty" json:"type,omitempty"`
	Name         string                   `yaml:"name,omitempty" json:"name,omitempty"`
	Count        int64                    `yaml:"count,omitempty" json:"count,omitempty"`
	Resources    SystemComponentResources `yaml:"resources,omitempty" json:"resources,omitempty"`
	PortMappings []string                 `yaml:"port-mappings,omitempty" json:"portMappings,omitempty"`
	Sidecars     []SystemComponentSidecar `yaml:"sidecars,omitempty" json:"sidecars,omitempty"`
	Args         []string                 `yaml:"args,omitempty" json:"args,omitempty"`
	Environment  map[string]string        `yaml:"environment,omitempty" json:"environment,omitempty"`
}

func (sc SystemComponent) String() string {
	type tmp SystemComponent
	return fmt.Sprintf("%+v", tmp(sc))
}

func (sc SystemComponent) GetCount() int64 {
	if sc.Count == 0 {
		return 1
	}
	return sc.Count
}
