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
	"path/filepath"
	"strings"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config/defaults"
	"github.com/whiteblock/definition/internal/converter"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/namer"
)

type Service interface {
	GetArgs(service entity.Service) []string
	GetCPUs(service entity.Service) int64
	GetEntrypoint(service entity.Service) string
	GetImage(service entity.Service) string

	GetNetwork(service entity.Service) string
	GetIP(state *entity.State, service entity.Service) string
	GetMemory(service entity.Service) int64
	GetVolumes(service entity.Service) []command.Mount
	GetDirectories(service entity.Service) []command.Mount
}

type serviceParser struct {
	defaults defaults.Service
	conv     converter.Resource
}

func NewService(defaults defaults.Service, conv converter.Resource) Service {
	return &serviceParser{defaults: defaults, conv: conv}
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

func (sp *serviceParser) GetCPUs(service entity.Service) int64 {
	if service.SquashedService.Resources.Cpus == 0 {
		return sp.defaults.CPUs
	}
	return service.SquashedService.Resources.Cpus
}

func (sp *serviceParser) GetMemory(service entity.Service) int64 {
	res, err := sp.conv.FromResources(service.SquashedService.Resources)
	if err != nil || res.Memory == 0 {
		return sp.defaults.Memory
	}
	return res.Memory
}

func (sp *serviceParser) GetImage(service entity.Service) string {
	if service.SquashedService.Image == "" {
		return sp.defaults.Image
	}

	return service.SquashedService.Image
}

func (sp *serviceParser) GetNetwork(service entity.Service) string {
	if !service.IsTask {
		return namer.SidecarNetwork(service)
	}
	return "none"
}

func (sp *serviceParser) GetIP(state *entity.State, service entity.Service) string {
	if service.IsTask {
		return ""
	}
	out := state.Subnets[service.Name].Next().String()
	state.IPs[service.Name] = out
	return out
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

	out = append(out, sp.GetDirectories(service)...)

	return out
}

func (sp *serviceParser) GetDirectories(service entity.Service) []command.Mount {
	dirs := GetServiceDirectories(service)
	out := []command.Mount{}
	for _, dir := range dirs {
		out = append(out, command.Mount{
			Name:      namer.InputFileVolume(service, dir),
			Directory: dir,
			ReadOnly:  false,
		})
	}
	return out
}

func GetServiceDirectories(service entity.Service) []string {
	dirs := map[string]bool{}
	for _, inputFiles := range service.SquashedService.InputFiles {
		dst := inputFiles.Destination()
		if strings.HasSuffix(dst, "/") {
			dirs[dst] = true
		} else {
			dirs[filepath.Dir(dst)] = true
		}
	}

	out := []string{}
	for dir := range dirs {
		out = append(out, dir)
	}

	return out
}
