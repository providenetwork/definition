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

func TestNetwork_Next(t *testing.T) {
	state, err := NewNetworkState("10.0.0.0/8", "172.16.0.0/12", 1000)
	require.NoError(t, err)
	require.NotNil(t, state)

	net, err := state.GetNextGlobal()
	require.NoError(t, err)
	require.NotNil(t, net)

	numIps := 100
	ips := make([]string, numIps)
	for i := 0; i < 100; i++ {
		ip := net.Next()
		assert.NotNil(t, ip)
		ips[i] = ip.String()
	}
	assert.ElementsMatch(t, net.GetIPs(), ips)
}
