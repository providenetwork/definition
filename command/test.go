/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/whiteblock/definition/command/biome"

	"github.com/sirupsen/logrus"
	"github.com/whiteblock/utility/common"
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

	// TestNameKey is the location of the test name in the meta
	TestNameKey = "test"

	// DefinitionIDKey is the location of the definition in the meta
	DefinitionIDKey = "definition"
)

var (

	// NoExpiration is the placeholder for the expiration to be ignored
	NoExpiration = time.Unix(0, 0)

	// NoTimeout indicates that there is no timeout duration provided
	NoTimeout time.Duration = -1

	// NoDuration indicates that duration is not set
	NoDuration time.Duration = 0

	// ErrNoCommands is when commands are needed but there are none
	ErrNoCommands = errors.New("there are no commands")

	// ErrNoPhase means that the commands meta doesn't include phase data
	ErrNoPhase = errors.New("phase not found")

	// ErrDone is given when isntructions is finished after this round
	ErrDone = errors.New("all done")

	// ErrTooManyFailed represents when given more failures than commands
	ErrTooManyFailed = errors.New("more commands failed than currently exist")

	// ErrNotAllFound is for when one or more commands where not matched up with the given ids
	ErrNotAllFound = errors.New("one or more commands where not matched up with the given ids")
)

// Test contains the instructions necessary for the execution of a single test
type Test struct {
	Instructions
	ProvisionCommand biome.CreateBiome `json:"provisionCommand"`
}

// UnmarshalJSON prevents Instructions from taking this over
func (instruct *Test) UnmarshalJSON(data []byte) error {

	type testShim Test

	var tmp testShim
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	*instruct = Test(tmp)

	nextStep := map[string]interface{}{}
	err = json.Unmarshal(data, &nextStep)
	if err != nil {
		return err
	}

	tmp2, err := json.Marshal(nextStep["provisionCommand"])
	if err != nil {
		return err
	}
	return json.Unmarshal(tmp2, &instruct.ProvisionCommand)
}

func (t Test) Status(cnt int) common.Status {
	status := t.Instructions.Status()
	status.Phase = "setting up the environment"
	status.StepsLeft = int(t.GuessSteps()) - cnt
	return status
}

func (t Test) GuessSteps() int64 {
	return int64(len(t.Commands) + len(t.ProvisionCommand.Instances)*2)
}

// Instructions contains all of the execution based information, for use in anything that executes the
// Commands
type Instructions struct {
	ID           string      `json:"id,omitempty"`
	OrgID        string      `json:"orgID,omitempty"`
	DefinitionID string      `json:"definitionID,omitempty"`
	Timestamp    time.Time   `json:"timestamp,omitempty"`
	Commands     [][]Command `json:"commands,omitempty"`

	Round            int                    `json:"round,omitempty"`
	GlobalTimeout    Timeout                `json:"globalTimeout,omitempty"`
	GlobalExpiration time.Time              `json:"globalExpiration,omitempty"`
	PhaseTimeouts    map[string]Timeout     `json:"phaseTimeouts,omitempty"`
	PhaseExpirations map[string]time.Time   `json:"phaseExpirations,omitempty"`
	Meta             map[string]interface{} `json:"meta,omitempty"`
}

// MetaTo takes the given meta fields and puts them into the given struct ptr.
// Should be easier than having to do lots of type conversions directly on the meta.
// This uses json marshalling to accomplish it
func (instruct *Instructions) MetaTo(out interface{}) error {
	data, err := json.Marshal(instruct.Meta)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, out)
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
		if errors.Is(err, ErrNoPhase) {
			return nil
		}
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

func (instruct Instructions) NeverTerminate() bool {
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

// Peek is like Next, but does not make any state changes
func (instruct Instructions) Peek() ([]Command, error) {
	if len(instruct.Commands) == 0 {
		return nil, ErrNoCommands
	}
	if len(instruct.Commands) == 1 {
		return instruct.Commands[0], ErrDone
	}
	return instruct.Commands[0], nil
}

// Next pops the first element off of Commands. If this results in Commands being
// empty, it returns ErrDone.
func (instruct *Instructions) Next() ([]Command, error) {
	if len(instruct.Commands) == 0 {
		return nil, ErrNoCommands
	}
	if len(instruct.Commands) == 1 {
		defer func() { instruct.Commands = [][]Command{} }()
		return instruct.Commands[0], ErrDone
	}
	out := instruct.Commands[0]
	instruct.Commands = instruct.Commands[1:]
	return out, instruct.next()
}

// UnmarshalJSON creates Instructions from JSON, and also handles
// creating the links back this object
func (instruct *Instructions) UnmarshalJSON(data []byte) error {

	type instructShim Instructions

	var tmp instructShim
	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}
	*instruct = Instructions(tmp)

	for i := range instruct.Commands {
		for j := range instruct.Commands[i] {
			instruct.Commands[i][j].parent = instruct
		}
	}
	return nil
}

func (instruct *Instructions) PlaceInProperIDs(log logrus.Ext1FieldLogger, files []common.Metadata) {
	for i := range files {
		for j := range instruct.Commands {
			for k := range instruct.Commands[j] {
				if instruct.Commands[j][k].Order.Type != Putfileincontainer {
					continue
				}
				payload := instruct.Commands[j][k].Order.Payload.(FileAndContainer)
				log.WithFields(logrus.Fields{
					"sourcePath": payload.File.ID,
					"metaPath":   files[i].Path,
				}).Debug("checking if these match")
				if payload.File.ID == files[i].Path {
					payload.File.Meta = files[i]
					payload.File.ID = files[i].ID
					instruct.Commands[j][k].Order.Payload = payload

					log.WithFields(logrus.Fields{
						"payload": payload,
						"id":      payload.File.ID,
					}).Info("updating payload with file id")
				}
			}
		}
	}
}

func (instruct *Instructions) PartialCompletion(failed []string) error {
	if len(instruct.Commands) == 0 {
		return ErrNoCommands
	}

	if len(failed) > len(instruct.Commands[0]) {
		return ErrTooManyFailed
	}

	if len(failed) == len(instruct.Commands[0]) || len(failed) == 0 {
		return nil
	}

	left := []Command{}
	for _, failure := range failed {
		for _, cmd := range instruct.Commands[0] {
			if cmd.ID == failure {
				left = append(left, cmd)
				break
			}
		}
	}

	if len(left) != len(failed) {
		return ErrNotAllFound
	}

	instruct.Commands[0] = left
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

func (instruct Instructions) Status() common.Status {
	phase, err := instruct.Phase()
	if err != nil {
		phase = ""
	}

	return common.Status{
		Test:      instruct.ID,
		Org:       instruct.OrgID,
		Def:       instruct.DefinitionID,
		Phase:     phase,
		StepsLeft: len(instruct.Commands),
	}
}
