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

package process

import (
	"fmt"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/internal/distribute"
	"github.com/whiteblock/definition/schema"
)

type TestCommands [][]command.Command

func (cmds TestCommands) Append(commands [][]command.Command) TestCommands {
	return TestCommands(append([][]command.Command(cmds), commands...))
}

type Commands interface {
	Interpret(spec schema.RootSchema, dists []*distribute.ResourceDist) ([]TestCommands, error)
}

type commandProc struct {
	calc TestCalculator
}

func NewCommands(calc TestCalculator) Commands {
	return &commandProc{calc: calc}
}

func (cmdProc *commandProc) Interpret(spec schema.RootSchema,
	dists []*distribute.ResourceDist) ([]TestCommands, error) {

	if len(dists) == len(spec.Tests) {
		return nil, fmt.Errorf("dists does not match the tests")
	}
	out := []TestCommands{}
	for i, dist := range dists {
		testCommands, err := cmdProc.calc.Commands(spec, dist, i)
		if err != nil {
			return nil, err
		}
		out = append(out, testCommands)
	}
	return out, nil
}
