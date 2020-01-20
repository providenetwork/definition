/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package namer

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"
)

func capString(s string, size int) string {
	h := sha256.New()
	h.Write([]byte(s))
	return fmt.Sprintf("%x", h.Sum(nil))[:size]
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
	return SidecarS(parent.Name, sidecar.Name)
}

func SidecarS(parent, sidecar string) string {
	return fmt.Sprintf("%s-%s", parent, sidecar)
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

func InputFileVolume(containerName string, dir string) string {
	return "input-" + capString(containerName+dir, 14)
}

func LocalVolume(containerName string, volumeName string) string {
	return "local-" + capString(containerName+volumeName, 14)
}

func toEnv(name string) {
	name = strings.Replace(name, "-", "_", -1)
	name = strings.ToUpper(name)
	return name
}

func IPEnvSidecar(parent entity.Service, sidecar schema.Sidecar) string {
	return IPEnvSidecarS(parent.Name, sidecar.Name)
}

func IPEnvSidecarS(parent, sidecar string) string {
	return toEnv(SidecarS(parent, sidecar) + "_" + parent)
}

func IPEnvService(name string) string {
	return toEnv(name)
}

func IPEnvServiceNet(container, network string) string {
	return toEnv(container + "_" + network.Name)
}
