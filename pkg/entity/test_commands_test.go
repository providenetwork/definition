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

	"github.com/whiteblock/definition/command"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTestCommands_Append(t *testing.T) {
	cmds := TestCommands{}
	testStuff := [][]command.Command{
		[]command.Command{},
		[]command.Command{},
		[]command.Command{},
	}
	testStuff2 := [][]command.Command{
		[]command.Command{},
		[]command.Command{},
	}
	cmds = cmds.Append(testStuff)
	require.Len(t, cmds, 3)
	cmds = cmds.Append(testStuff2)
	assert.Len(t, cmds, 5)
}

func TestTestCommands_MetaInject(t *testing.T) {
	cmds := TestCommands{}
	testStuff := [][]command.Command{
		[]command.Command{
			command.Command{},
			command.Command{},
			command.Command{},
		},
		[]command.Command{},
		[]command.Command{},
	}
	cmds = cmds.Append(testStuff)
	require.Len(t, cmds, 3)
	cmds.MetaInject("foo", "bar", "bar", "baz")
	require.NotNil(t, cmds[0][0].Meta)
	assert.Equal(t, "bar", cmds[0][0].Meta["foo"])
	assert.Equal(t, "baz", cmds[0][0].Meta["bar"])

}
