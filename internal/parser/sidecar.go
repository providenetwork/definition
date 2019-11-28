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
	"github.com/whiteblock/definition/schema"

	"github.com/docker/docker/api/types/strslice"
)

type Sidecar interface {
	GetEntrypoint(sidecar schema.Sidecar) string
	GetImage(sidecar schema.Sidecar) string
	GetLabels(parent entity.Service, sidecar schema.Sidecar) map[string]string
	GetNetwork(parent entity.Service) strslice.StrSlice
	GetVolumes(sidecar schema.Sidecar) []command.Mount
}

type sidecarParser struct {
}

func NewSidecar() Sidecar {
	return &sidecarParser{}
}

func (sp sidecarParser) GetEntrypoint(sidecar schema.Sidecar) string {
	//TODO
	return ""
}

func (sp sidecarParser) GetImage(sidecar schema.Sidecar) string {
	//TODO
	return ""
}

func (sp sidecarParser) GetLabels(parent entity.Service, sidecar schema.Sidecar) map[string]string {
	//TODO
	return nil
}

func (sp sidecarParser) GetNetwork(parent entity.Service) strslice.StrSlice {
	//TODO
	return nil
}

func (sp sidecarParser) GetVolumes(sidecar schema.Sidecar) []command.Mount {
	//TODO
	return nil
}
