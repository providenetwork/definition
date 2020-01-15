/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package process

import (
	"fmt"

	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"
)

type Commands interface {
	Interpret(spec schema.RootSchema, dists []*entity.ResourceDist) ([]entity.TestCommands, error)
	Env(spec schema.RootSchema, dists []*entity.ResourceDist) ([]map[string]string, error)
}

var (
	ErrTestsDontMatchDist = fmt.Errorf("dists does not match the tests")
)

type commandProc struct {
	calc TestCalculator
}

func NewCommands(calc TestCalculator) Commands {
	return &commandProc{calc: calc}
}

func (cmdProc *commandProc) Interpret(spec schema.RootSchema,
	dists []*entity.ResourceDist) ([]entity.TestCommands, error) {

	if len(dists) != len(spec.Tests) {
		return nil, ErrTestsDontMatchDist
	}
	out := []entity.TestCommands{}
	for i, dist := range dists {
		testCommands, err := cmdProc.calc.Commands(spec, dist, i)
		if err != nil {
			return nil, err
		}
		out = append(out, testCommands)
	}
	return out, nil
}

func (cmdProc *commandProc) Env(spec schema.RootSchema, dists []*entity.ResourceDist) ([]map[string]string, error) {
	if len(dists) != len(spec.Tests) {
		return nil, ErrTestsDontMatchDist
	}
	out := []map[string]string{}
	for i, dist := range dists {
		environment, err := cmdProc.calc.Env(spec, dist, i)
		if err != nil {
			return nil, err
		}
		out = append(out, environment)
	}
	return out, nil
}
