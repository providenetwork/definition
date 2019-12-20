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

package command

import (
	"encoding/json"

	"github.com/whiteblock/definition/command/biome"
)

// Test contains the instructions necessary for the execution of a single test
type Test struct {
	Instructions
	ProvisionCommand biome.CreateBiome
}

// Instructions contains all of the execution based information, for use in anything that executes the
// Commands
type Instructions struct {
	ID           string              `json:"id"`
	OrgID        string              `json:"orgID"`
	DefinitionID string              `json:"definitionID"`
	Commands     [][]command.Command `json:"commands"`
}

// NextRound pops the first element off of Commmands. If this results in Commands being
// empty, it returns false
func (instruct *Instructions) NextRound() bool {
	if len(instruct.Commands) < 2 {
		instruct.Commands = [][]command.Command{}
		return false
	}
	instruct.Commands = instruct.Commands[1:]
	return true
}
func (instruct Instructions) GetNextCommands() []command.Commands {
	if len(instruct.Commands) == 0 {
		return nil
	}
	return instruct.Commands[0]
}
