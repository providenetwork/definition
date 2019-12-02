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

package parser

import (
	"strconv"
	"strings"
	"time"

	"github.com/whiteblock/definition/schema"
)

type Network interface {
	GetBandwidth(network schema.Network) string
	GetLatency(network schema.Network) (int64, error)
	GetPacketLoss(network schema.Network) (float64, error)
}

type networkParser struct {
}

func NewNetwork() Network {
	return &networkParser{}
}

func (np networkParser) GetBandwidth(network schema.Network) string {
	return network.Bandwidth
}

func (np networkParser) GetLatency(network schema.Network) (int64, error) {
	latency := strings.Replace(network.Latency, " ", "", -1)
	if latency == "" {
		return 0, nil
	}

	if !strings.ContainsAny(latency, "numsh") {
		return strconv.ParseInt(latency, 10, 32)
	}

	dur, err := time.ParseDuration(latency)
	return dur.Microseconds(), err
}

func (np networkParser) GetPacketLoss(network schema.Network) (float64, error) {
	packetLoss := strings.Trim(network.PacketLoss, " %")
	if packetLoss == "" {
		return 0, nil
	}
	return strconv.ParseFloat(packetLoss, 64)
}
