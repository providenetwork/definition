/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/whiteblock/definition/command/biome"
)

func TestTest_UnmarshalJSON(t *testing.T) {
	var test2 Test
	test := Test{
		Instructions: Instructions{
			ID:           "bar",
			OrgID:        "baz",
			DefinitionID: "foo",
			Commands: [][]Command{
				[]Command{
					{
						Order: Order{
							Type:    "createContainer",
							Payload: map[string]interface{}{},
						},
						Target: Target{
							IP: "127.0.0.1",
						},
					},
				},
				[]Command{
					{
						Order: Order{
							Type:    "createContainer",
							Payload: map[string]interface{}{},
						},
						Target: Target{
							IP: "127.0.0.1",
						},
					},
				},
			},
		},
		ProvisionCommand: biome.CreateBiome{
			DefinitionID: "foo",
			TestID:       "bar",
			OrgID:        "baz",
			Instances: []biome.Instance{
				{
					Provider: biome.GCPProvider,
					CPUs:     1,
					Memory:   1,
					Storage:  1,
				},
			},
		},
	}
	data, err := json.Marshal(test)
	require.NoError(t, err)
	require.NotNil(t, data)

	err = json.Unmarshal(data, &test2)
	assert.NoError(t, err)
	assert.Equal(t, test.ProvisionCommand.DefinitionID, test2.ProvisionCommand.DefinitionID)
}
