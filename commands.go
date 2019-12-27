/*
	Copyright 2019 Whiteblock Inc.
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
	"time"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/command/biome"
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal"
	"github.com/whiteblock/definition/internal/distribute"
	parse "github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/internal/process"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/whiteblock/utility/common"
	"github.com/whiteblock/utility/utils"
)

var globalCommands Commands

// Commands is the interface of a parser that extracts commands from a definition
type Commands interface {
	// GetTests gets all of the commands, for both provisioner and genesis.
	// The genesis commands will be in dependency groups, so that
	// res[n+1] is the set of commands which require the execution of the commands
	// In res[n].
	GetTests(def Definition, files ...common.Metadata) ([]command.Test, error)
}

type commands struct {
	proc process.Commands
	dist distribute.Distributor
	conf config.Config
}

// NewCommands creates a new command extractor from the given viper config
func NewCommands(conf config.Config) (Commands, error) {
	proc, dist, err := internal.GetFunctionality(conf)
	return &commands{conf: conf, proc: proc, dist: dist}, err
}

// GetTests gets all of the commands, for both provisioner and genesis.
// The genesis commands will be in dependency groups, so that
// res[n+1] is the set of commands which require the execution of the commands
// In res[n]. We get both at once, since we have to compute the commands for provisioning to produce
// the commands for Genesis.
func (cmdParser commands) GetTests(def Definition, files ...common.Metadata) ([]command.Test, error) {
	resDist, err := cmdParser.dist.Distribute(def.Spec)
	if err != nil {
		return nil, errors.Wrap(err, "distribute")
	}

	testCmds, err := cmdParser.proc.Interpret(def.Spec, resDist)
	if err != nil {
		return nil, errors.Wrap(err, "interpret")
	}
	logger, err := cmdParser.conf.Logger.GetLogger()
	if err != nil {
		return nil, err
	}
	out := make([]command.Test, len(testCmds))
	for i := range testCmds {

		id := utils.GetUUIDString()
		testCmds[i].MetaInject(
			command.OrgIDKey, def.Metadata.OrgID,
			command.DefinitionIDKey, def.ID,
			command.TestIDKey, id)

		phases, global := parse.Timeouts(def.Spec.Tests[i])

		out[i] = command.Test{
			ProvisionCommand: resDist[i].ToBiomeCommand(biome.GCPProvider, def.ID, def.Metadata.OrgID, id),
			Instructions: command.Instructions{
				ID:            id,
				OrgID:         def.Metadata.OrgID,
				DefinitionID:  def.ID,
				Timestamp:     time.Now(),
				Commands:      [][]command.Command(testCmds[i]),
				PhaseTimeouts: phases,
				GlobalTimeout: global,
			},
		}
		out[i].PlaceInProperIDs(logger, files)
	}
	return out, nil
}

// ConfigureGlobal allows you to provide the global config for this library
func ConfigureGlobal(conf config.Config) (err error) {
	globalCommands, err = NewCommands(conf)
	return
}

// ConfigureGlobalFromViper allows you to tie in configuration for this library from viper.
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

// GetTests gets all of the commands, for both provisioner and genesis.
// The genesis commands will be in dependency groups, so that
// res[n+1] is the set of commands which require the execution of the commands
// In res[n].
func GetTests(def Definition, files ...common.Metadata) ([]command.Test, error) {
	return globalCommands.GetTests(def, files...)
}

func init() {
	// This may fail if the default configuration is bad, perhaps we might want to just
	// Error out if ConfigureGlobal is not called.
	err := ConfigureGlobalFromViper(viper.New())
	if err != nil {
		panic(err)
	}
}
