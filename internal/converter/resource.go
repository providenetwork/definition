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

package converter

import (
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/internal/util"
	"github.com/whiteblock/definition/schema"
)

type Resource interface {
	FromResources(sRes schema.Resources) (entity.Resource, error)
}

type resourceCovnverter struct {
}

func NewResource() Resource {
	return &resourceCovnverter{}
}

func (rc resourceCovnverter) FromResources(sRes schema.Resources) (entity.Resource, error) {
	out := entity.Resource{CPUs: int64(sRes.Cpus)}
	mem, err := util.Memconv(sRes.Memory)
	if err != nil {
		return entity.Resource{}, err
	}
	out.Memory = mem

	storage, err := util.Memconv(sRes.Storage)
	if err != nil {
		return entity.Resource{}, err
	}
	out.Storage = storage
	return out, nil
}
