/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package parser

import (
	"fmt"

	"github.com/whiteblock/definition/config/defaults"
	"github.com/whiteblock/definition/internal/converter"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/pkg/namer"
	"github.com/whiteblock/definition/schema"

	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
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

func (sp sidecarParser) GetIP(state *entity.State, parent entity.Service,
	sidecar schema.Sidecar) string {

	out := state.Subnets[parent.Name].Next().String()
	state.IPs[namer.IPEnvSidecar(parent, sidecar)] = out
	return out
}

func GetSidecarEnv(sidecar schema.Sidecar, parent entity.Service, state *entity.State) map[string]string {
	out := sidecar.Environment
	if out == nil {
		out = map[string]string{}
	}
	out["SERVICE"] = state.IPs[namer.IPEnvService(parent.Name)]
	out["NAME"] = namer.Sidecar(parent, sidecar)
	return out
}
