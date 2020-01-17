/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

// StartContainer is the command for starting a container
type StartContainer struct {
	Name   string `json:"name"`
	Attach bool   `json:"attach"`
	// Timeout is the maximum amount of time to wait for the task before terminating it.
	// This is ignored if attach is false
	Timeout Timeout `json:"timeout"`
}
