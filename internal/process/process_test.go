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

package process

import (
	"testing"

	"github.com/whiteblock/definition/internal/mocks/process"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCommands_Interpret(t *testing.T) {
	testSpec := schema.RootSchema{
		Tests: []schema.Test{
			schema.Test{},
			schema.Test{},
			schema.Test{},
		},
	}

	testDists := []*entity.ResourceDist{
		nil,
		nil,
		nil,
	}
	calcMock := new(mocks.TestCalculator)
	calcMock.On("Commands", mock.Anything, mock.Anything, mock.Anything).Return(
		entity.TestCommands{}, nil).Times(len(testDists))

	cmdProc := NewCommands(calcMock)
	res, err := cmdProc.Interpret(testSpec, testDists)
	assert.NoError(t, err)
	require.NotNil(t, res)
	assert.Len(t, res, len(testDists))
	calcMock.AssertExpectations(t)
}
