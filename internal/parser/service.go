/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package parser

import (
	"path/filepath"
	"strings"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config/defaults"
	"github.com/whiteblock/definition/internal/converter"
	"github.com/whiteblock/definition/internal/namer"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"
)

type Service interface {
	GetArgs(service entity.Service) []string
	GetCPUs(service entity.Service) int64
	GetEntrypoint(service entity.Service) string
	GetImage(service entity.Service) string

	GetNetwork(service entity.Service) string
	GetIP(state *entity.State, service entity.Service) string
	GetMemory(service entity.Service) int64
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

func getVolumes(service entity.Service, volumes []schema.Volume) []command.Mount {
	out := []command.Mount{}

	for _, vol := range volumes {
		readOnly := false
		if vol.Permissions == "r" || vol.Permissions == "read" || vol.Permissions == "ro" {
			readOnly = true
		}
		name := vol.Name
		if vol.Local() {
			name = namer.LocalVolume(service.Name, vol.Name)
		}
		out = append(out, command.Mount{
			Name:      name,
			Directory: vol.Path,
			ReadOnly:  readOnly,
		})
	}

	return out
}

func GetVolumes(service entity.Service, volumes []schema.Volume) []command.Mount {
	out := []command.Mount{}
	out = append(out, getVolumes(service, volumes)...)
	out = append(out, GetServiceDirMounts(service)...)
	return out
}

func GetServiceDirMounts(service entity.Service) []command.Mount {
	dirs := GetDirectories(service.SquashedService.InputFiles)
	out := []command.Mount{}
	for _, dir := range dirs {
		out = append(out, command.Mount{
			Name:      namer.InputFileVolume(service.Name, dir),
			Directory: dir,
			ReadOnly:  false,
		})
	}
	return out
}

func GetDirectories(files []schema.InputFile) []string {
	dirs := map[string]bool{}
	for _, file := range files {
		dst := file.Destination()
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
