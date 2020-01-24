/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package distribute

import (
	"testing"

	"github.com/whiteblock/definition/config"
	mockParser "github.com/whiteblock/definition/pkg/mocks/parser"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

//  It would be too redundant to test in different funcs due to the complexity of this
//  Implementation
func TestSystemState_FullTest(t *testing.T) {
	parser := new(mockParser.Resources)
	systems := []schema.SystemComponent{
		schema.SystemComponent{
			Name: "1",
		},
		schema.SystemComponent{
			Name: "2",
		},
		schema.SystemComponent{
			Name: "3",
		},
	}
	spec := schema.RootSchema{}

	segments := []entity.Segment{}
	for i, system := range systems {
		result := make([]entity.Segment, i)
		for range result {
			result = append(result, entity.Segment{
				Name: system.Name,
			})
		}
		segments = append(segments, result...)
		parser.On("SystemComponent", spec, system).Return(result, nil).Once()
		parser.On("SystemComponentNamesOnly", system).Return(result).Once()
	}
	statePack := entity.NewStatePack(spec, config.Bucket{}, logrus.New())
	state := NewSystemState(parser)

	//Successful Add
	result, err := state.Add(statePack, spec, systems)
	require.NotNil(t, result)
	assert.NoError(t, err)
	assert.ElementsMatch(t, segments, result)
	assert.Len(t, result, len(segments))

	//Successful Remove
	result, err = state.Remove(statePack, []string{"1", "2", "3"})
	require.NotNil(t, result)
	assert.NoError(t, err)
	assert.ElementsMatch(t, segments, result)
	assert.Len(t, result, len(segments))

	parser.AssertExpectations(t)
}
