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
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/internal/entity"
)

type Service interface {
	GetEntrypoint(service entity.Service) string
	GetImage(service entity.Service) string
	GetNetworks(service entity.Service) []string
	GetVolumes(service entity.Service) []command.Mount
}

type serviceParser struct {
}

func NewService() Service {
	return &serviceParser{}
}

func (sp *serviceParser) GetEntrypoint(service entity.Service) string {
	//TODO
	return ""
}

func (sp *serviceParser) GetImage(service entity.Service) string {
	//TODO
	return ""
}

func (sp *serviceParser) GetNetworks(service entity.Service) []string {
	//TODO
	return nil
}

func (sp *serviceParser) GetVolumes(service entity.Service) []command.Mount {
	//TODO
	return nil
}

/*
type Script struct {
	SourcePath string `yaml:"source-path,omitempty" json:"source-path,omitempty"`
	Inline     string `yaml:"inline,omitempty" json:"inline,omitempty"`
}

type Service struct {
	Name          string            `yaml:"name,omitempty" json:"name,omitempty"`
	Description   string            `yaml:"description,omitempty" json:"description,omitempty"`
	SharedVolumes []SharedVolume    `yaml:"shared-volumes,omitempty" json:"shared-volumes,omitempty"`
	Resources     Resources         `yaml:"resources,omitempty" json:"resources,omitempty"`
	Args          []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Environment   map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Image         string            `yaml:"image,omitempty" json:"image,omitempty"`
	Script        Script            `yaml:"script,omitempty" json:"script,omitempty"`
	InputFiles    []InputFile       `yaml:"input-files,omitempty" json:"input-files,omitempty"`
}
*/
