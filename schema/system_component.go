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
