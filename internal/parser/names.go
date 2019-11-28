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
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/schema"
	"strings"
)

type Names interface {
	InputFileVolume(input schema.InputFile) string
	Sidecar(parent entity.Service, sidecar schema.Sidecar) string
	SidecarNetwork(parent entity.Service) string
	SystemComponent(systemComponent schema.SystemComponent) string
	SystemService(systemComponent schema.SystemComponent, index int) string
	Task(task schema.Task, index int) string
}

type namer struct {
}

func NewNames() Names {
	return &namer{}
}

func (n *namer) InputFileVolume(input schema.InputFile) string {
	return strings.Replace(input.DestinationPath, "/", "-", 0)
}

func (n *namer) Sidecar(parent entity.Service, sidecar schema.Sidecar) string {
	//TODO
	return ""
}

func (n *namer) SidecarNetwork(parent entity.Service) string {
	//TODO
	return ""
}

func (n *namer) SystemComponent(systemComponent schema.SystemComponent) string {
	//TODO
	return ""
}

func (n *namer) SystemService(systemComponent schema.SystemComponent, index int) string {
	//TODO
	return ""
}

func (n *namer) Task(task schema.Task, index int) string {
	//TODO
	return ""
}
