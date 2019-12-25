/*
	Copyright 2019 Whiteblock Inc.
	This file is a part of the definition.

	definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	definition is distributed in the hope that it will be useful,
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