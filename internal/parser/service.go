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
	"github.com/whiteblock/definition/internal/config/defaults"
	"github.com/whiteblock/definition/internal/entity"
	//shell "github.com/kballard/go-shellquote"
)

type Service interface {
	GetArgs(service entity.Service) []string
	GetEntrypoint(service entity.Service) string
	GetImage(service entity.Service) string
	GetSidecarNetwork(service entity.Service) string
	GetNetworks(service entity.Service) []string
	GetVolumes(service entity.Service) []command.Mount
}

type serviceParser struct {
	defaults defaults.Service
	namer    Names
}

func NewService(defaults defaults.Service, namer Names) Service {
	return &serviceParser{defaults: defaults, namer: namer}
}

func (sp *serviceParser) GetArgs(service entity.Service) []string {
	if service.SquashedService.Script.Inline != "" {
		return []string{"-c", service.SquashedService.Script.Inline}
	}
	return service.SquashedService.Args
}

func (sp *serviceParser) GetEntrypoint(service entity.Service) string {
	if service.SquashedService.Script.SourcePath != "" {
		return service.SquashedService.Script.SourcePath
	}
	if service.SquashedService.Script.Inline != "" {
		return "/bin/sh"
	}
	return ""
}

func (sp *serviceParser) GetImage(service entity.Service) string {
	if service.SquashedService.Image == "" {
		return sp.defaults.Image
	}

	return service.SquashedService.Image
}

func (sp *serviceParser) GetNetworks(service entity.Service) []string {
	out := make([]string, len(service.Networks)+1)
	out[0] = GetSidecarNetwork(service)
	for i := range service.Networks {
		out[i+1] = service.Networks[i].Name
	}
	return out
}

func (sp *serviceParser) GetSidecarNetwork(service entity.Service) string {
	return sp.namer.SidecarNetwork(service)
}

func (sp *serviceParser) GetVolumes(service entity.Service) []command.Mount {

	out := []command.Mount{}

	for _, sharedVol := range service.SquashedService.SharedVolumes {
		out = append(out, command.Mount{
			Name:      sharedVol.Name,
			Directory: sharedVol.SourcePath,
			ReadOnly:  false,
		})
	}

	for _, inputVol := range service.SquashedService.InputFiles {
		out = append(out, command.Mount{
			Name:      sp.namer.InputFileVolume(inputVol),
			Directory: inputVol.DestinationPath,
			ReadOnly:  false,
		})
	}
	return out
}
