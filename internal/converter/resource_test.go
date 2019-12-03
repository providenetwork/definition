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

package converter

import (
	"testing"

	"github.com/whiteblock/definition/config/defaults"
	"github.com/whiteblock/definition/internal/util"
	"github.com/whiteblock/definition/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResource_FromResources(t *testing.T) {
	testConfig := defaults.Resources{
		CPUs:    2,
		Memory:  5,
		Storage: 6,
	}
	conv := NewResource(testConfig)
	require.NotNil(t, conv)

	res, err := conv.FromResources(schema.Resources{
		Cpus:    0,
		Memory:  "",
		Storage: "",
	})
	assert.NoError(t, err)
	assert.Equal(t, testConfig.CPUs, res.CPUs)
	assert.Equal(t, testConfig.Memory, res.Memory)
	assert.Equal(t, testConfig.Storage, res.Storage)

	res, err = conv.FromResources(schema.Resources{
		Cpus:    2,
		Memory:  "10g",
		Storage: "20g",
	})
	assert.NoError(t, err)
	assert.Equal(t, int64(2), res.CPUs)
	assert.Equal(t, 10*util.Gibi/util.Mibi, res.Memory)
	assert.Equal(t, 20*util.Gibi/util.Mibi, res.Storage)

	res, err = conv.FromResources(schema.Resources{
		Cpus:    0,
		Memory:  "40",
		Storage: "30",
	})
	assert.NoError(t, err)
	assert.Equal(t, testConfig.CPUs, res.CPUs)
	assert.Equal(t, 40*util.Mibi/util.Mibi, res.Memory)
	assert.Equal(t, 30*util.Gibi/util.Mibi, res.Storage)

	_, err = conv.FromResources(schema.Resources{
		Cpus:    0,
		Memory:  "foo",
		Storage: "",
	})
	assert.Error(t, err)

	_, err = conv.FromResources(schema.Resources{
		Cpus:    0,
		Memory:  "",
		Storage: "bar",
	})
	assert.Error(t, err)

}
