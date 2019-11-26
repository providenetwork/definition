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

package distribute

import (
	"fmt"

	"github.com/whiteblock/definition/schema"
)

type SystemState interface {
	Add(systems []schema.SystemComponent) ([]Segment, error)
	Remove(systems []string) ([]Segment, error)
}

type systemState struct {
	totalSystemState map[string]schema.SystemComponent
	parser           SchemaParser
}

func NewSystemState(parser SchemaParser) SystemState {
	return &systemState{
		totalSystemState: map[string]schema.SystemComponent{},
		parser:           parser}
}

func (state *systemState) Add(systems []schema.SystemComponent) ([]Segment, error) {
	out := []Segment{}
	for _, system := range systems {
		name := state.parser.NameSystemComponent(system)
		_, exists := state.totalSystemState[name]
		if exists {
			return nil, fmt.Errorf("already have a system with the name \"%s\"", name)
		}
		segments, err := state.parser.ParseSystemComponent(system)
		if err != nil {
			return nil, err
		}
		out = append(out, segments...)
	}

	for _, system := range systems {
		name := state.parser.NameSystemComponent(system)
		state.totalSystemState[name] = system
	}

	return out, nil
}

func (state *systemState) Remove(systems []string) ([]Segment, error) {
	out := []Segment{}
	for _, toRemove := range systems {
		system, exists := state.totalSystemState[toRemove]
		if !exists {
			return nil, fmt.Errorf("system not found")
		}
		segments, err := state.parser.ParseSystemComponent(system)
		if err != nil {
			return nil, err
		}
		out = append(out, segments...)
	}
	for _, toRemove := range systems {
		delete(state.totalSystemState, toRemove)
	}
	return out, nil
}
