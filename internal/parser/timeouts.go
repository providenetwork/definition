/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package parser

import (
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/schema"
)

// Timeouts extracts the global and phase timeouts from a test
func Timeouts(test schema.Test) (phaseTOs map[string]command.Timeout, global command.Timeout) {
	global = test.Timeout
	if len(test.Phases) > 0 {
		phaseTOs = map[string]command.Timeout{}
	}
	for _, phase := range test.Phases {
		phaseTOs[phase.Name] = phase.Timeout
	}
	return
}
