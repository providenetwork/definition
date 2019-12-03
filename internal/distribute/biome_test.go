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
	mockDist "github.com/whiteblock/definition/internal/mocks/distribute"
	mockEntity "github.com/whiteblock/definition/internal/mocks/entity"
	mockParser "github.com/whiteblock/definition/internal/mocks/parser"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBiomeCalculator_AddNextPhase(t *testing.T) {
	testStatePack := entity.NewStatePack(schema.RootSchema{}, config.Bucket{})

	buckets := new(mockEntity.ResourceBuckets)
	buckets.On("Remove", mock.Anything).Return(nil).Times(3)
	buckets.On("Add", mock.Anything).Return(nil).Times(4)
	testStatePack.Buckets = buckets

	state := new(mockDist.SystemState)
	state.On("Add", mock.Anything, mock.Anything, mock.Anything).Return(
		[]entity.Segment{}, nil).Twice()

	state.On("Remove", mock.Anything, mock.Anything).Return(
		[]entity.Segment{}, nil).Twice()
	state.On("GetAlreadyExists", mock.Anything, mock.Anything).Return(nil, nil, false).Twice()
	parser := new(mockParser.Resources)
	parser.On("Tasks", mock.Anything, mock.Anything).Return(
		[]entity.Segment{}, nil).Twice()

	calc := NewBiomeCalculator(parser, state, logrus.New())

	err := calc.AddNextPhase(testStatePack, schema.Phase{})
	assert.NoError(t, err)

	err = calc.AddNextPhase(testStatePack, schema.Phase{})
	assert.NoError(t, err)

	buckets.AssertExpectations(t)
	state.AssertExpectations(t)
	parser.AssertExpectations(t)
}
