/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package parser

import (
	"testing"

	"github.com/whiteblock/definition/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetwork_GetLatency(t *testing.T) {
	net := NewNetwork()
	require.NotNil(t, net)

	val, err := net.GetLatency(schema.Network{
		Latency: "80",
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(80), val)

	val, err = net.GetLatency(schema.Network{
		Latency: "",
	})
	assert.NoError(t, err)

	val, err = net.GetLatency(schema.Network{
		Latency: "foobar",
	})
	assert.Error(t, err)

	val, err = net.GetLatency(schema.Network{
		Latency: "80 ms",
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(80000), val)
}

func TestNetwork_GetPacketLoss(t *testing.T) {
	net := NewNetwork()
	require.NotNil(t, net)

	val, err := net.GetPacketLoss(schema.Network{
		PacketLoss: "80",
	})
	assert.NoError(t, err)
	assert.Equal(t, float64(80.0), val)

	val, err = net.GetPacketLoss(schema.Network{
		PacketLoss: "80%",
	})
	assert.NoError(t, err)
	assert.Equal(t, float64(80.0), val)

	val, err = net.GetPacketLoss(schema.Network{
		PacketLoss: "fruit",
	})
	assert.Error(t, err)

	val, err = net.GetPacketLoss(schema.Network{
		PacketLoss: "",
	})
	assert.NoError(t, err)
	assert.Equal(t, float64(0), val)
}
