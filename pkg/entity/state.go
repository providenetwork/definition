/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

import "github.com/whiteblock/definition/schema"

type State struct {
	Tasks       []Service
	SystemState map[string]schema.SystemComponent
	Subnets     map[string]Network
	Network     NetworkState
	IPs         map[string]string
}

func NewState(net NetworkState) *State {
	return &State{
		SystemState: map[string]schema.SystemComponent{},
		Subnets:     map[string]Network{},
		Network:     net,
		IPs:         map[string]string{},
	}
}
