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
	"time"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/schema"
)

//Command handles the simple schema -> order conversions
type Command interface {
	New(order command.Order, endpoint string, timeout time.Duration) (command.Command, error)
	CreateNetwork(network schema.Network) command.Order
	CreateVolume(volume schema.SharedVolume) command.Order
	CreateContainer(service schema.Service) command.Order
	StartContainer(service schema.Service) command.Order
	AttachNetwork(service schema.Service, network schema.Network) command.Order

	RemoveContainer(service schema.Service) command.Order
}
