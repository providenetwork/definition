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
package definition

import(
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/command/biome"
	"github.com/whiteblock/definition/internal"
)

//Definition is the representation of the test definition format. 
type Definition interface {
	//Gets the UUID which uniquely identifies this Definition
	GetID() string

	//GetCommands gets the commands in dependency groups, so that
	//res[n+1] is the set of commands which require the execution of the commands
	//in res[n].
	GetCommands() ([][]command.Command,error)

	//GetProvisioningRequest calculates the biome resources necessary to support
	//this test
	GetProvisioningRequest() (biome.CreateBiome, error)

	//Validate returns nil if it is a valid test definition, other the returned 
	//error will contain an explanation for the issue
	Validate() error
}

// ParseYAML takes a raw set of bytes and
// de-serializes them to a Definition structure
func ParseYAML(raw []byte) (Definition, error) {
	return internal.ParseYAML(raw)
}