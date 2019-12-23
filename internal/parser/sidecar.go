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
	"fmt"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config/defaults"
	"github.com/whiteblock/definition/internal/converter"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/namer"
	"github.com/whiteblock/definition/schema"

	"github.com/jinzhu/copier"
)

type Sidecar interface {
	GetArgs(sidecar schema.Sidecar) []string
	GetCPUs(sidecar schema.Sidecar) string
	GetEntrypoint(sidecar schema.Sidecar) string
	GetImage(sidecar schema.Sidecar) string
	GetLabels(parent entity.Service, sidecar schema.Sidecar) map[string]string
	GetMemory(sidecar schema.Sidecar) string
	GetNetwork(parent entity.Service) string
	GetIP(state *entity.State, parent entity.Service, sidecar schema.Sidecar) string
	GetVolumes(sidecar schema.Sidecar) []command.Mount
}

type sidecarParser struct {
	defaults defaults.Service
	conv     converter.Resource
}

func NewSidecar(defaults defaults.Service, conv converter.Resource) Sidecar {
	return &sidecarParser{defaults: defaults, conv: conv}
}

func (sp sidecarParser) GetArgs(sidecar schema.Sidecar) []string {
	if sidecar.Script.Inline != "" {
		return []string{"-c", sidecar.Script.Inline}
	}
	return sidecar.Args
}

func (sp sidecarParser) GetEntrypoint(sidecar schema.Sidecar) string {
	if sidecar.Script.SourcePath != "" {
		return sidecar.Script.SourcePath
	}
	if sidecar.Script.Inline != "" {
		return "/bin/sh"
	}
	return ""
}

func (sp sidecarParser) GetCPUs(sidecar schema.Sidecar) string {
	if sidecar.Resources.Cpus == 0 {
		return fmt.Sprint(sp.defaults.CPUs)
	}
	return fmt.Sprint(sidecar.Resources.Cpus)
}

func (sp sidecarParser) GetMemory(sidecar schema.Sidecar) string {
	res, err := sp.conv.FromResources(sidecar.Resources)
	if err != nil || res.Memory == 0 {
		return fmt.Sprint(sp.defaults.Memory)
	}
	return fmt.Sprint(res.Memory)

}

func (sp sidecarParser) GetImage(sidecar schema.Sidecar) string {
	if sidecar.Image == "" {
		return sp.defaults.Image
	}

	return sidecar.Image
}

func (sp sidecarParser) GetLabels(parent entity.Service, sidecar schema.Sidecar) map[string]string {
	var labels map[string]string
	copier.Copy(&labels, parent.Labels)
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["name"] = sidecar.Name
	labels["service"] = parent.Name
	return labels
}

func (sp sidecarParser) GetNetwork(parent entity.Service) string {
	return namer.SidecarNetwork(parent)
}

func (sp sidecarParser) GetVolumes(sidecar schema.Sidecar) []command.Mount {
	out := []command.Mount{}

	for _, mntVol := range sidecar.MountedVolumes {
		readOnly := false
		if mntVol.Permissions == "r" || mntVol.Permissions == "read" {
			readOnly = true
		}
		out = append(out, command.Mount{
			Name:      mntVol.VolumeName,
			Directory: mntVol.DestinationPath,
			ReadOnly:  readOnly,
		})
	}

	return out
}

func (sp sidecarParser) GetIP(state *entity.State, parent entity.Service,
	sidecar schema.Sidecar) string {

	out := state.Subnets[parent.Name].Next().String()
	state.IPs[namer.Sidecar(parent, sidecar)+"_"+parent.Name] = out
	return out
}
