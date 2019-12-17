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

package entity

import (
	"github.com/whiteblock/definition/schema"
)

type SystemDiff struct {
	Modified        []ServiceDiff
	AddedNetworks   []schema.Network
	RemovedNetworks []schema.Network
	Added           []Service
	Removed         []Service
}

func (diff *SystemDiff) Append(sys *SystemDiff) {
	if sys == nil {
		return
	}
	diff.Modified = append(diff.Modified, sys.Modified...)
	diff.Added = append(diff.Added, sys.Added...)
	diff.Removed = append(diff.Removed, sys.Removed...)
	diff.AddedNetworks = append(diff.AddedNetworks, sys.AddedNetworks...)
	diff.RemovedNetworks = append(diff.RemovedNetworks, sys.RemovedNetworks...)
}
