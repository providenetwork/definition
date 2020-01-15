/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCommand_ParseOrderPayloadInto_Success(t *testing.T) {
	containerName := "tester"
	networkName := "testnet"
	cmd := Command{
		ID:     "TEST",
		Target: Target{IP: "0.0.0.0"},
		Order: Order{
			Type: Attachnetwork,
			Payload: map[string]string{
				"container": containerName,
				"network":   networkName,
			},
		},
	}

	var cn ContainerNetwork
	err := cmd.ParseOrderPayloadInto(&cn)
	assert.NoError(t, err)
}

func TestCommand_ParseOrderPayloadInto_Failure(t *testing.T) {
	containerName := "tester"
	networkName := "testnet"
	cmd := Command{
		ID:     "TEST",
		Target: Target{IP: "0.0.0.0"},
		Order: Order{
			Type: Attachnetwork,
			Payload: map[string]string{
				"container": containerName,
				"network":   networkName,
				"i should":  "not be here",
			},
		},
	}

	var cn ContainerNetwork
	err := cmd.ParseOrderPayloadInto(&cn)
	assert.Error(t, err)
}

func TestDeserSerRoundtripCommand(t *testing.T) {
	command := Command{
		ID:     "",
		Target: Target{},
		Order: Order{
			Type:    Startcontainer,
			Payload: map[string]interface{}{"name": "test"},
		},
	}
	bytes, err := json.Marshal(command)
	if err != nil {
		t.Fatal(err)
	}
	read := Command{}
	err = json.Unmarshal(bytes, &read)
	require.NoError(t, err)

	require.Equal(t, command, read)

	payload := SimpleName{}
	err = mapstructure.Decode(read.Order.Payload, &payload)
	require.NoError(t, err)

	require.Equal(t, "test", payload.Name)
}
