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

func TestNames_InputFileVolume(t *testing.T) {
	//todo insignificant
}

func TestNames_Sidecar(t *testing.T) {
	//todo insignificant
}

func TestNames_SidecarNetwork(t *testing.T) {
	//todo insignificant
}

func TestNames_SystemComponent(t *testing.T) {
	names := NewNames()
	require.NotNil(t, names)
	res := names.SystemComponent(schema.SystemComponent{
		Name: "Foobar",
		Type: "barfoo",
	})
	assert.Equal(t, "Foobar", res)
	res = names.SystemComponent(schema.SystemComponent{
		Type: "barfoo",
	})
	assert.Equal(t, "barfoo", res)
}

func TestNames_SystemService(t *testing.T) {
	//todo insignificant
}

func TestNames_Task(t *testing.T) {
	//todo insignificant
}
