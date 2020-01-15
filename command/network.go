/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

//Network represents a logic network on which containers exist
type Network struct {
	//Name is the name of the network
	Name    string            `json:"name"`
	Subnet  string            `json:"subnet"`
	Gateway string            `json:"gateway"`
	Global  bool              `json:"global"`
	Labels  map[string]string `json:"labels"`
}
