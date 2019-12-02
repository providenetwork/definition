/*
	Copyright 2019 whiteblock Inc.
	This file is a part of the definition.

	definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	definition is distributed in the hope that it will be useful,
	but dock ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package definition

import (
	"encoding/json"
	"fmt"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/command/biome"
	"github.com/whiteblock/definition/internal"
	"github.com/whiteblock/definition/internal/distribute"
	"github.com/whiteblock/definition/internal/process"
	"github.com/whiteblock/definition/schema"
	"github.com/whiteblock/definition/validator"

	"github.com/spf13/viper"
)

//Commands is the interface of a parser that extracts commands from a definition
type Commands interface {
	//GetCommands gets all of the commands, for both provisioner and genesis.
	//The genesis commands will be in dependency groups, so that
	//res[n+1] is the set of commands which require the execution of the commands
	//in res[n].
	GetCommands(def Definition) (biome.CreateBiome, [][]command.Command, error)
}

type commands struct {
	proc process.Commands
	dist distribute.Distributor
}

//NewCommands creates a new command extractor from the given viper config
func NewCommands(v *viper.Viper) (Commands, error) {
	GetFunctionality(v*viper.Viper)(process.Commands, distribute.Distributor, error)
}

//GetCommands gets all of the commands, for both provisioner and genesis.
//The genesis commands will be in dependency groups, so that
//res[n+1] is the set of commands which require the execution of the commands
//in res[n].
func (cmdParser commands) GetCommands(def Definition) (biome.CreateBiome, [][]command.Command, error) {
	//TODO
	return biome.CreateBiome{}, nil, nil
}

//ConfigureGlobal allows you to tie in configuration for this library from viper.
func ConfigureGlobal(v *viper.Viper) error {
	//TODO
	return nil
}

//GetCommands gets all of the commands, for both provisioner and genesis.
//The genesis commands will be in dependency groups, so that
//res[n+1] is the set of commands which require the execution of the commands
//in res[n].
func GetCommands(def Definition) (biome.CreateBiome, [][]command.Command, error) {
	//TODO
	return biome.CreateBiome{}, nil, nil
}
