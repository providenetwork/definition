/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package process

import (
	"errors"
	"fmt"

	"github.com/whiteblock/definition/internal/maker"
	"github.com/whiteblock/definition/internal/merger"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/pkg/namer"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
)

// System is for diff calculations
type System interface {
	GetAlreadyExists(state *entity.State, systems []schema.SystemComponent) (
		[]schema.SystemComponent, []schema.SystemComponent, bool)

	UpdateChanged(state *entity.State, spec schema.RootSchema,
		systems []schema.SystemComponent) (*entity.SystemDiff, error)

	// Add modifies State
	Add(state *entity.State, spec schema.RootSchema, systems []schema.SystemComponent) ([]entity.Service, error)

	// Remove modifies state
	Remove(state *entity.State, spec schema.RootSchema, systems []string) ([]entity.Service, error)

	Tasks(state *entity.State, spec schema.RootSchema, tasks []schema.Task) ([]entity.Service, error)
}

var (
	ErrSystemNotFound = errors.New("system not found")
)

type system struct {
	maker maker.Service
	log   logrus.Ext1FieldLogger
}

func NewSystem(
	maker maker.Service,
	log logrus.Ext1FieldLogger) System {
	return &system{maker: maker, log: log}
}

func (sys system) UpdateChanged(state *entity.State, spec schema.RootSchema,
	systems []schema.SystemComponent) (diff *entity.SystemDiff, err error) {

	diff = &entity.SystemDiff{}
	for _, systemUpdate := range systems {
		name := namer.SystemComponent(systemUpdate)
		old, exists := state.SystemState[name]
		if !exists {
			return nil, ErrSystemNotFound
		}
		system := merger.MergeSystemLeft(old, systemUpdate)
		sys.log.WithField("system", system).Debug("merged the systems")
		sysDiff, err := sys.maker.FromSystemDiff(spec, old, system)
		if err != nil {
			return nil, err
		}
		diff.Append(sysDiff)

		state.SystemState[name] = system
	}
	return
}

func (sys system) GetAlreadyExists(state *entity.State, systems []schema.SystemComponent) (
	exist []schema.SystemComponent, noExist []schema.SystemComponent, anyExist bool) {

	anyExist = false
	for _, s := range systems {
		name := namer.SystemComponent(s)
		_, exists := state.SystemState[name]
		if exists {
			anyExist = true
			exist = append(exist, s)
		} else {
			noExist = append(noExist, s)
		}
	}
	return
}

//  Add modifies State
func (sys system) Add(state *entity.State, spec schema.RootSchema,
	systems []schema.SystemComponent) ([]entity.Service, error) {
	out := []entity.Service{}

	for _, system := range systems {
		name := namer.SystemComponent(system)
		_, exists := state.SystemState[name]
		if exists {
			return nil, fmt.Errorf("already have a system with the name \"%s\"", name)
		}
		services, err := sys.maker.FromSystem(spec, system)
		if err != nil {
			return nil, err
		}
		out = append(out, services...)
	}

	for _, system := range systems {
		name := namer.SystemComponent(system)
		state.SystemState[name] = system
	}

	return out, nil
}

//Remove modifies state
func (sys system) Remove(state *entity.State, spec schema.RootSchema,
	systems []string) ([]entity.Service, error) {

	out := []entity.Service{}
	for _, toRemove := range systems {
		system, exists := state.SystemState[toRemove]
		if !exists {
			return nil, ErrSystemNotFound
		}
		services, err := sys.maker.FromSystem(spec, system)
		if err != nil {
			return nil, err
		}
		out = append(out, services...)
	}
	for _, toRemove := range systems {
		delete(state.SystemState, toRemove)
	}
	return out, nil
}

func (sys system) Tasks(state *entity.State, spec schema.RootSchema,
	tasks []schema.Task) ([]entity.Service, error) {

	out := make([]entity.Service, len(tasks))
	state.Tasks = make([]entity.Service, len(tasks))
	for i, task := range tasks {
		service, err := sys.maker.FromTask(spec, task, i)
		if err != nil {
			return nil, err
		}
		out[i] = service
		state.Tasks[i] = service
	}
	return out, nil
}
