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

package entity

import (
	"github.com/whiteblock/definition/command"
)

type TestCommands [][]command.Command

func (cmds TestCommands) Append(commands [][]command.Command) TestCommands {
	return TestCommands(append([][]command.Command(cmds), commands...))
}

func (cmds TestCommands) MetaInject(kv ...string) {
	if len(kv)%2 != 0 {
		return
	}
	for i := range cmds {
		for j := range cmds[i] {
			if cmds[i][j].Meta == nil {
				cmds[i][j].Meta = map[string]string{}
			}
			for k := 0; k < len(kv)/2; k++ {
				cmds[i][j].Meta[kv[2*k]] = kv[2*k+1]
			}
		}
	}
}
