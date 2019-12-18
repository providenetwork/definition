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

package distribute

import (
	"testing"

	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal/entity"
	mockParser "github.com/whiteblock/definition/internal/mocks/parser"
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
