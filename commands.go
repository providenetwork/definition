/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
	"github.com/whiteblock/definition/pkg/entity"

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
	GetTests(def Definition, meta Meta) ([]command.Test, error)

	// GetEnvs gets the environment variables which will be supplied to each of the containers
	GetEnvs(def Definition) ([]map[string]string, error)

	// GetDist gets the resource distribution
	GetDist(def Definition) ([]*entity.ResourceDist, error)
}

type commands struct {
	proc process.Commands
	dist distribute.Distributor
	conf config.Config
}

// Meta is addition data for functionalities not covered in the
// scope of the test definition schema. Such as the data on the referenced
// files themselves
type Meta struct {
	// Files is the metadata of the files provided for the user. This field is optional.
	Files []common.Metadata

	// Domains are domain names for the created instances. This field is optional.
	Domains []string
}

// NewCommands creates a new command extractor from the given viper config
func NewCommands(conf config.Config) (Commands, error) {
	proc, dist, err := internal.GetFunctionality(conf)
	return &commands{conf: conf, proc: proc, dist: dist}, err
}

// GetDist gets the resource distribution
func (cmdParser commands) GetDist(def Definition) ([]*entity.ResourceDist, error) {
	return cmdParser.dist.Distribute(def.Spec)
}

// GetTests gets all of the commands, for both provisioner and genesis.
// The genesis commands will be in dependency groups, so that
// res[n+1] is the set of commands which require the execution of the commands
// In res[n]. We get both at once, since we have to compute the commands for provisioning to produce
// the commands for Genesis.
func (cmdParser commands) GetTests(def Definition, meta Meta) ([]command.Test, error) {
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
		domain := ""
		if len(meta.Domains) > i {
			domain = meta.Domains[i]
		}
		id := utils.GetUUIDString()
		testCmds[i].MetaInject(
			command.OrgIDKey, def.OrgID,
			command.DefinitionIDKey, def.ID,
			command.TestIDKey, id)

		phases, global := parse.Timeouts(def.Spec.Tests[i])

		out[i] = command.Test{
			ProvisionCommand: resDist[i].ToBiomeCommand(biome.GCPProvider, def.ID, def.OrgID, id, domain),
			Instructions: command.Instructions{
				ID:            id,
				OrgID:         def.OrgID,
				DefinitionID:  def.ID,
				Timestamp:     time.Now(),
				Commands:      [][]command.Command(testCmds[i]),
				PhaseTimeouts: phases,
				GlobalTimeout: global,
			},
		}
		out[i].PlaceInProperIDs(logger, meta.Files)
	}
	return out, nil
}

// GetEnvs gets the environment variables which will be supplied to each of the containers
func (cmdParser commands) GetEnvs(def Definition) ([]map[string]string, error) {
	resDist, err := cmdParser.dist.Distribute(def.Spec)
	if err != nil {
		return nil, errors.Wrap(err, "distribute")
	}

	return cmdParser.proc.Env(def.Spec, resDist)
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
func GetTests(def Definition, meta Meta) ([]command.Test, error) {
	return globalCommands.GetTests(def, meta)
}

// GetEnvs gets the environment variables which will be supplied to each of the containers
func GetEnvs(def Definition) ([]map[string]string, error) {
	return globalCommands.GetEnvs(def)
}

// GetDist gets the resource distribution
func GetDist(def Definition) ([]*entity.ResourceDist, error) {
	return globalCommands.GetDist(def)
}

func init() {
	// This may fail if the default configuration is bad, perhaps we might want to just
	// Error out if ConfigureGlobal is not called.
	err := ConfigureGlobalFromViper(viper.New())
	if err != nil {
		panic(err)
	}
}
