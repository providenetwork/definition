/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
	require.NotNil(t, net)
	assert.Equal(t, "10.0.0.0/22", net.Subnet())

	net, err = state.GetNextGlobal()
	require.NoError(t, err)
	require.NotNil(t, net)
	assert.Equal(t, "10.0.4.0/22", net.Subnet())

	net, err = state.GetNextGlobal()
	require.NoError(t, err)
	require.NotNil(t, net)
	assert.Equal(t, "10.0.8.0/22", net.Subnet())
}

func TestNetworkState_GetNextLocal(t *testing.T) {
	state, err := NewNetworkState("10.0.0.0/8", "172.16.0.0/12", 1000)
	require.NoError(t, err)
	require.NotNil(t, state)

	for i := 0; i < 10; i++ {
		net, err := state.GetNextLocal(i)
		require.NoError(t, err)
		require.NotNil(t, net)
		assert.Equal(t, "172.16.0.0/22", net.Subnet())

		net, err = state.GetNextLocal(i)
		require.NoError(t, err)
		require.NotNil(t, net)
		assert.Equal(t, "172.16.4.0/22", net.Subnet())

		net, err = state.GetNextLocal(i)
		require.NoError(t, err)
		require.NotNil(t, net)
		assert.Equal(t, "172.16.8.0/22", net.Subnet())
	}

}
