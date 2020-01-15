/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
