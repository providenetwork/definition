/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

import (
	neta "github.com/dspinhirne/netaddr-go"
)

type Network interface {
	Subnet() string
	Gateway() string
	Next() *neta.IPv4
	GetIPs() []string
}

type network struct {
	network *neta.IPv4Net
	inUse   uint32
}

func NewNetwork(net *neta.IPv4Net) Network {
	return &network{
		network: net,
		inUse:   ReservedIPs,
	}
}

func (n network) Subnet() string {
	return n.network.String()
}

func (n network) Gateway() string {
	return n.network.Network().Next().String()
}

func (n *network) Next() *neta.IPv4 {
	defer func() { n.inUse++ }()
	return n.network.Nth(n.inUse)
}

func (n network) GetIPs() []string {
	out := []string{}
	for i := ReservedIPs; i < n.inUse; i++ {
		net := n.network.Nth(i)
		if net != nil {
			out = append(out, net.String())
		}
	}
	return out
}
