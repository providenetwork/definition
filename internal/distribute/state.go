/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package distribute

import (
	"errors"
	"fmt"

	"github.com/whiteblock/definition/internal/merger"
	"github.com/whiteblock/definition/internal/namer"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"
)

type SystemState interface {
	Add(sp *entity.StatePack, spec schema.RootSchema,
		systems []schema.SystemComponent) ([]entity.Segment, error)

	GetAlreadyExists(sp *entity.StatePack, systems []schema.SystemComponent) (
		exist []schema.SystemComponent, noExist []schema.SystemComponent, anyExist bool)

	UpdateChanged(sp *entity.StatePack, spec schema.RootSchema,
		systems []schema.SystemComponent) ([]entity.Segment, []entity.Segment, error)

	Remove(sp *entity.StatePack, systems []string) ([]entity.Segment, error)
}

var (
	ErrSystemNotFound = errors.New("system not found")
)

type systemState struct {
	parser parser.Resources
}

func NewSystemState(
	parser parser.Resources) SystemState {

	return &systemState{
		parser: parser}
}

func (state systemState) UpdateChanged(sp *entity.StatePack, spec schema.RootSchema,
	systems []schema.SystemComponent) (toAdd []entity.Segment,
	toRemove []entity.Segment, err error) {

	for _, systemUpdate := range systems {
		name := namer.SystemComponent(systemUpdate)
		old, exists := sp.SystemState[name]
		if !exists {
			return nil, nil, ErrSystemNotFound
		}
		system := merger.MergeSystemLeft(systemUpdate, old)

		segs, err := state.parser.FromSystemDiff(spec, old, system)
		if err != nil {
			return nil, nil, err
		}

		if system.Count < old.Count {
			toRemove = append(toRemove, segs...)
		} else {
			toAdd = append(toAdd, segs...)
		}
		sp.SystemState[name] = system
	}
	return
}

func (state systemState) GetAlreadyExists(sp *entity.StatePack, systems []schema.SystemComponent) (
	exist []schema.SystemComponent, noExist []schema.SystemComponent, anyExist bool) {

	anyExist = false
	for _, s := range systems {
		name := namer.SystemComponent(s)
		_, exists := sp.SystemState[name]
		if exists {
			anyExist = true
			exist = append(exist, s)
		} else {
			noExist = append(noExist, s)
		}
	}
	return
}

func (state *systemState) Add(sp *entity.StatePack, spec schema.RootSchema,
	systems []schema.SystemComponent) ([]entity.Segment, error) {

	out := []entity.Segment{}
	for _, system := range systems {
		name := namer.SystemComponent(system)
		_, exists := sp.SystemState[name]
		if exists {
			return nil, fmt.Errorf("already have a system with the name \"%s\"", name)
		}
		segments, err := state.parser.SystemComponent(spec, system)
		if err != nil {
			return nil, err
		}
		out = append(out, segments...)
	}

	for _, system := range systems {
		name := namer.SystemComponent(system)
		sp.SystemState[name] = system
	}

	return out, nil
}

func (state *systemState) Remove(sp *entity.StatePack,
	systems []string) ([]entity.Segment, error) {

	out := []entity.Segment{}
	for _, toRemove := range systems {
		system, exists := sp.SystemState[toRemove]
		if !exists {
			return nil, ErrSystemNotFound
		}
		segments := state.parser.SystemComponentNamesOnly(system)
		out = append(out, segments...)
	}
	for _, toRemove := range systems {
		delete(sp.SystemState, toRemove)
	}
	return out, nil
}
