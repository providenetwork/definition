/*
	Copyright 2019 whiteblock Inc.
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
	mockParser "github.com/whiteblock/definition/internal/mocks/parser"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSystemState_FullTest(t *testing.T) {
	parser := new(mockParser.SchemaParser)
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
	segments := []Segment{}
	for i, system := range systems {
		result := make([]Segment, i)
		for j := range result {
			result = Segment{
				Name: system.Name,
			}
		}
		segments = append(segments, result...)
		parser.On("NameSystemComponent", system).Return(system.Name).Twice()
		parser.On("ParseSystemComponent", system).Return(result, err).Once()
	}

	state := NewSystemState(parser)
	result, err := state.Add(systems)
	require.NotNil(t, result)
	assert.NoError(t, err)
	assert.ElementsMatch(t, segments, result)
	assert.Len(t, result, len(segments))

}
