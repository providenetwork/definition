/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

//Netconf represents network impairments which are to be applied to a network
type Netconf struct {
	//  Container is the target container
	Container string `json:"container"`
	// Network is the target network
	Network string `json:"network"`
	// Limit is the max number of packets to hold the in queue
	Limit int `json:"limit"`
	// Loss represents packet loss % ie 100% = 100
	Loss float64 `json:"loss"`
	//  Delay represents the latency to be applied in microseconds
	Delay int `json:"delay"`
	// Rate represents the bandwidth constraint to be applied to the network
	Rate string `json:"rate"`
	// Duplication represents the percentage of packets to duplicate
	Duplication float64 `json:"duplicate"`
	// Corrupt represents the percentage of packets to corrupt
	Corrupt float64 `json:"corrupt"`
	// Reorder represents the percentage of packets that get reordered
	Reorder float64 `json:"reorder"`
}
