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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNetworkState_GetNextGlobal(t *testing.T) {
	state, err := NewNetworkState("10.0.0.0/8", "172.16.0.0/12", 1000)
	require.NoError(t, err)
	require.NotNil(t, state)

	net, err := state.GetNextGlobal()
	require.NoError(t, err)
	require.NotNil(t, net.network)
	assert.Equal(t, "10.0.0.0/22", net.Subnet())

	net, err = state.GetNextGlobal()
	require.NoError(t, err)
	require.NotNil(t, net.network)
	assert.Equal(t, "10.0.4.0/22", net.Subnet())

	net, err = state.GetNextGlobal()
	require.NoError(t, err)
	require.NotNil(t, net.network)
	assert.Equal(t, "10.0.8.0/22", net.Subnet())
}

func TestNetworkState_GetNextLocal(t *testing.T) {
	state, err := NewNetworkState("10.0.0.0/8", "172.16.0.0/12", 1000)
	require.NoError(t, err)
	require.NotNil(t, state)

	for i := 0; i < 10; i++ {
		net, err := state.GetNextLocal(i)
		require.NoError(t, err)
		require.NotNil(t, net.network)
		assert.Equal(t, "172.16.0.0/22", net.Subnet())

		net, err = state.GetNextLocal(i)
		require.NoError(t, err)
		require.NotNil(t, net.network)
		assert.Equal(t, "172.16.4.0/22", net.Subnet())

		net, err = state.GetNextLocal(i)
		require.NoError(t, err)
		require.NotNil(t, net.network)
		assert.Equal(t, "172.16.8.0/22", net.Subnet())
	}

}
