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
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config/defaults"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/schema"

	"github.com/docker/docker/api/types/strslice"
	"github.com/jinzhu/copier"
)

type Sidecar interface {
	GetArgs(sidecar schema.Sidecar) []string
	GetEntrypoint(sidecar schema.Sidecar) string
	GetImage(sidecar schema.Sidecar) string
	GetLabels(parent entity.Service, sidecar schema.Sidecar) map[string]string
	GetNetwork(parent entity.Service) strslice.StrSlice
	GetVolumes(sidecar schema.Sidecar) []command.Mount
}

type sidecarParser struct {
	defaults defaults.Service
	namer    Names
}

func NewSidecar(defaults defaults.Service, namer Names) Sidecar {
	return &sidecarParser{defaults: defaults, namer: namer}
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

func (sp sidecarParser) GetImage(sidecar schema.Sidecar) string {
	if sidecar.Image == "" {
		return sp.defaults.Image
	}

	return sidecar.Image
}

func (sp sidecarParser) GetLabels(parent entity.Service, sidecar schema.Sidecar) map[string]string {
	var labels map[string]string
	copier.Copy(&labels, parent.Labels)
	labels["name"] = sidecar.Name
	labels["service"] = parent.Name
	return labels
}

func (sp sidecarParser) GetNetwork(parent entity.Service) strslice.StrSlice {
	return strslice.StrSlice([]string{sp.namer.SidecarNetwork(parent)})
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

	for _, inputVol := range sidecar.InputFiles {
		out = append(out, command.Mount{
			Name:      sp.namer.InputFileVolume(inputVol),
			Directory: inputVol.DestinationPath,
			ReadOnly:  false,
		})
	}
	return out
}
