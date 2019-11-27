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

package process

import (
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/internal/distribute"
	"github.com/whiteblock/definition/internal/entity"
	"github.com/whiteblock/definition/schema"
)

type Resolve interface {
	CreateNetworks(systems []schema.SystemComponent) ([]command.Command, error)

	CreateServices(spec schema.RootSchema, dist *distribute.ResourceDist, index int,
		services []entity.Service) ([][]command.Command, error)

	RemoveServices(dist *distribute.ResourceDist, services []entity.Service) ([][]command.Command, error)
}

type resolve struct {
}
