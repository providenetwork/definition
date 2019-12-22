/*
	Copyright 2019 Whiteblock Inc.
	This file is a part of the Definitio

	Definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later versio

	Definition is distributed in the hope that it will be useful,
	but dock ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package namer

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/schema"
)

func capString(s string, size int) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))[:size]
}

func InputFileVolume(input schema.InputFile) string {
	return SanitizeVolumeName(input.DestinationPath)
}

func SanitizeVolumeName(volName string) string {
	if strings.HasPrefix(volName, "/") {
		volName = "wb_" + volName[1:]
	}
	return strings.Replace(volName, "/", "-", -1)
}

func DefaultNetwork(sys schema.SystemComponent) string {
	return Network(SystemComponent(sys))
}

func Network(name string) string {
	return "net-" + capString(name, 11)
}

func Sidecar(parent entity.Service, sidecar schema.Sidecar) string {
	return fmt.Sprintf("%s-%s", parent.Name, sidecar.Name)
}

func SidecarNetwork(parent entity.Service) string {
	return "snet-" + capString(parent.Name, 10)
}

func SystemComponent(sys schema.SystemComponent) string {
	if sys.Name != "" {
		return sys.Name
	}
	return sys.Type
}

func SystemService(sys schema.SystemComponent, index int) string {
	return fmt.Sprintf("%s-service%d", SystemComponent(sys), index)
}

func Task(task schema.Task, index int) string {
	return fmt.Sprintf("%s-task%d", task.Type, index)
}
