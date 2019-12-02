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
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/schema"
)

type StatePack struct {
	State
	Buckets   ResourceBuckets
	PrevTasks []Segment
	Spec      schema.RootSchema
}

func NewStatePack(spec schema.RootSchema, conf config.Bucket) *StatePack {
	out := &StatePack{
		Buckets:   NewResourceBuckets(conf),
		PrevTasks: nil,
		Spec:      spec,
	}
	out.SystemState = map[string]schema.SystemComponent{}
	return out
}
