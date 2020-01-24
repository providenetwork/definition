/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package converter

import (
	"strings"

	"github.com/whiteblock/definition/config/defaults"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"

	util "github.com/whiteblock/utility/utils"
)

type Resource interface {
	FromResources(sRes schema.Resources) (entity.Resource, error)
}

type resourceConverter struct {
	def defaults.Resources
}

func NewResource(def defaults.Resources) Resource {
	return &resourceConverter{def: def}
}

func (rc resourceConverter) FromResources(sRes schema.Resources) (out entity.Resource, err error) {
	out.CPUs = int64(sRes.Cpus)
	if out.CPUs == 0 {
		out.CPUs = rc.def.CPUs
	}

	memory := strings.Trim(sRes.Memory, " ")
	if memory != "" {
		out.Memory, err = util.Memconv(sRes.Memory, util.Mibi)
		if err != nil {
			return
		}
		out.Memory /= util.Mibi
	} else {
		out.Memory = rc.def.Memory
	}
	storage := strings.Trim(sRes.Storage, " ")
	if storage != "" {
		out.Storage, err = util.Memconv(sRes.Storage, util.Gibi)
		if err != nil {
			return
		}
		out.Storage /= util.Mibi
	} else {
		out.Storage = rc.def.Storage
	}

	return
}
