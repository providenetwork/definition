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
	"errors"
	"time"

	"github.com/whiteblock/definition/command/biome"

	"github.com/imdario/mergo"
)

/**
 * Note: changing any of these constants is a breaking change, so use caution.
 */
const (
	// PhaseKey is the meta key for the phase name
	PhaseKey = "phase"

	// OrgIDKey is the location of the org id in meta
	OrgIDKey = "org"

	// TestIDKey is the location of the test id in the meta
	TestIDKey = "testRun"

	// DefinitionIDKey is the location of the definition in the meta
	DefinitionIDKey = "definition"
)

var (

	// NoExpiration is the placeholder for the expiration to be ignored
	NoExpiration = time.Unix(0, 0)

	NoTimeout time.Duration = 1<<63 - 1

	// ErrNoCommands is when commands are needed but there are none
	ErrNoCommands = errors.New("there are no commands")

	// ErrNoPhase means that the commands meta doesn't include phase data
	ErrNoPhase = errors.New("phase not found")

	// ErrDone is given when isntructions is finished after this round
	ErrDone = errors.New("all done")
)

// Test contains the instructions necessary for the execution of a single test
type Test struct {
	Instructions
	ProvisionCommand biome.CreateBiome `json:"provisionCommand"`
}

// Instructions contains all of the execution based information, for use in anything that executes the
// Commands
type Instructions struct {
	ID           string      `json:"id,omitempty"`
	OrgID        string      `json:"orgID,omitempty"`
	DefinitionID string      `json:"definitionID,omitempty"`
	Timestamp    time.Time   `json:"timestamp,omitempty"`
	Commands     [][]Command `json:"commands,omitempty"`

	Round            int                  `json:"round,omitempty"`
	GlobalTimeout    Timeout              `json:"globalTimeout,omitempty"`
	GlobalExpiration time.Time            `json:"globalExpiration,omitempty"`
	PhaseTimeouts    map[string]Timeout   `json:"phaseTimeouts,omitempty"`
	PhaseExpirations map[string]time.Time `json:"phaseExpirations,omitempty"`
}

// handle the meta changes for next
func (instruct *Instructions) next() error {
	defer func() { instruct.Round++ }()
	if instruct.Round == 0 {
		instruct.PhaseTimeouts = map[string]Timeout{}
		if !instruct.GlobalTimeout.IsInfinite() && instruct.GlobalTimeout.Duration.Nanoseconds() != 0 {
			instruct.GlobalExpiration = time.Now().Add(instruct.GlobalTimeout.Duration)
		} else {
			instruct.GlobalExpiration = NoExpiration
		}
	}

	phase, err := instruct.Phase()
	if err != nil {
		return err
	}

	if instruct.PhaseExpirations == nil {
		return nil
	}

	_, exists := instruct.PhaseExpirations[phase]
	if !exists {
		to, hasPhaseTimeout := instruct.PhaseTimeouts[phase]
		if hasPhaseTimeout {
			instruct.PhaseExpirations[phase] = time.Now().Add(to.Duration)
		} else {
			instruct.PhaseExpirations[phase] = NoExpiration
		}
	}
	return nil
}

func (instruct Instruction) NeverTerminate() bool {
	return instruct.GlobalTimeout.IsInfinite()
}

func (instruct Instructions) GetTimeRemaining() (time.Duration, error) {
	out := NoTimeout
	now := time.Now()
	if instruct.GlobalExpiration.Unix() != 0 && instruct.GlobalExpiration.Sub(now) < out {
		out = instruct.GlobalExpiration.Sub(now)
	}

	phase, err := instruct.Phase()
	if err != nil {
		return out, err
	}
	exp, exists := instruct.PhaseExpirations[phase]
	if !exists {
		return out, nil
	}
	if exp.Sub(now) < out {
		out = exp.Sub(now)
	}
	return out, nil

}

// Next pops the first element off of Commands. If this results in Commands being
// empty, it returns ErrNoCommands
func (instruct *Instructions) Next() ([]Command, error) {
	if len(instruct.Commands) == 0 {
		return nil, ErrNoCommands
	}
	if len(instruct.Commands) == 1 {
		instruct.Commands = [][]Command{}
		return instruct.Commands[0], ErrDone
	}
	out := instruct.Commands[0]
	instruct.Commands = instruct.Commands[1:]
	return out, instruct.next()
}

// UnmarshalJSON creates Instructions from JSON, and also handles
// creating the links back this object
func (instruct *Instructions) UnmarshalJSON(data []byte) error {
	var tmp map[string]interface{}
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	err = mergo.Map(instruct, tmp)
	if err != nil {
		return err
	}
	for i := range instruct.Commands {
		for j := range instruct.Commands[i] {
			instruct.Commands[i][j].parent = instruct
		}
	}
	return nil
}

func (instruct Instructions) Phase() (string, error) {
	if len(instruct.Commands) == 0 || len(instruct.Commands[0]) == 0 {
		return "", ErrNoCommands
	}
	out, exists := instruct.Commands[0][0].Meta[PhaseKey]
	if !exists {
		return "", ErrNoPhase
	}
	return out, nil
}
