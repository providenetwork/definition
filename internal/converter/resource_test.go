/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package converter

import (
	"testing"

	"github.com/whiteblock/definition/config/defaults"

	"github.com/whiteblock/definition/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	util "github.com/whiteblock/utility/utils"
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
