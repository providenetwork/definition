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
	"errors"
	"math"

	neta "github.com/dspinhirne/netaddr-go"
)

type NetworkState interface {
	GetNextGlobal() (Network, error)
	GetNextLocal(instance int) (Network, error)
}

var (
	ErrNoMoreGlobalNetworks = errors.New("no more global networks")
	ErrNoMoreLocalNetworks  = errors.New("no more local networks")
)

const (
	IPv4Len     uint   = 32
	ReservedIPs uint32 = 3
)

type networkState struct {
	prefixLen uint

	global      *neta.IPv4Net
	globalIndex uint32

	local        *neta.IPv4Net
	localIndexes map[int]uint32
}

func NewNetworkState(globalCIDR string, localCIDR string, maxNodes int) (NetworkState, error) {
	prefixLen := IPv4Len - uint(math.Ceil(math.Log2(float64(maxNodes+int(ReservedIPs)+1))))
	global, err := neta.ParseIPv4Net(globalCIDR)
	if err != nil {
		return nil, err
	}
	local, err := neta.ParseIPv4Net(localCIDR)
	return &networkState{
		global:       global,
		local:        local,
		prefixLen:    prefixLen,
		localIndexes: map[int]uint32{},
	}, err
}

func (ns *networkState) GetNextGlobal() (Network, error) {
	net := ns.global.NthSubnet(ns.prefixLen, ns.globalIndex)
	if net == nil {
		return nil, ErrNoMoreGlobalNetworks
	}
	ns.globalIndex++
	return NewNetwork(net), nil
}

func (ns *networkState) GetNextLocal(instance int) (Network, error) {
	index := ns.localIndexes[instance]

	net := ns.local.NthSubnet(ns.prefixLen, index)
	if net == nil {
		return nil, ErrNoMoreLocalNetworks
	}
	ns.localIndexes[instance]++
	return NewNetwork(net), nil
}
