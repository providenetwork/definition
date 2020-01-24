/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
