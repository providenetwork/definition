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
