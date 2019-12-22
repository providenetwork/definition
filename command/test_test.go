/*
	Copyright 2019 Whiteblock Inc.
	This file is a part of the Definition.

	Definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Definition is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
