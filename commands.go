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
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/command/biome"
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal"
	"github.com/whiteblock/definition/internal/distribute"
	"github.com/whiteblock/definition/internal/process"

	"github.com/spf13/viper"
)

var globalCommands Commands

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
func NewCommands(conf config.Config) (Commands, error) {
	proc, dist, err := internal.GetFunctionality(conf)
	return &commands{proc: proc, dist: dist}, err
}

//GetCommands gets all of the commands, for both provisioner and genesis.
//The genesis commands will be in dependency groups, so that
//res[n+1] is the set of commands which require the execution of the commands
//in res[n]. We get both at once, since we have to compute the commands for provisioning to produce
//the commands for Genesis.
func (cmdParser commands) GetCommands(def Definition) (biome.CreateBiome, [][]command.Command, error) {
	//TODO
	return biome.CreateBiome{}, nil, nil
}

//ConfigureGlobal allows you to provide the global config for this library
func ConfigureGlobal(conf config.Config) (err error) {
	globalCommands, err = NewCommands(conf)
	return
}

//ConfigureGlobalFromViper allows you to tie in configuration for this library from viper.
func ConfigureGlobalFromViper(v *viper.Viper) error {
	err := config.SetupViper(v)
	if err != nil {
		return err
	}
	conf, err := config.New(v)
	if err != nil {
		return err
	}
	return ConfigureGlobal(conf)
}

//GetCommands gets all of the commands, for both provisioner and genesis.
//The genesis commands will be in dependency groups, so that
//res[n+1] is the set of commands which require the execution of the commands
//in res[n].
func GetCommands(def Definition) (biome.CreateBiome, [][]command.Command, error) {
	return globalCommands.GetCommands(def)
}

func init() {
	//This may fail if the default configuration is bad, perhaps we might want to just
	//error out if ConfigureGlobal is not called.
	err := ConfigureGlobalFromViper(viper.New())
	if err != nil {
		panic(err)
	}
}
