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
		inUse:   2,
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
	for i := uint32(2); i < n.inUse; i++ {
		net := n.network.Nth(i)
		if net != nil {
			out = append(out, net.String())
		}
	}
	return out
}
