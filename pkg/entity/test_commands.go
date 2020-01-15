/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
