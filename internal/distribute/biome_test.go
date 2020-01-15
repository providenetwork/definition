/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package distribute

import (
	"testing"

	"github.com/whiteblock/definition/config"
	mockDist "github.com/whiteblock/definition/internal/mocks/distribute"
	mockEntity "github.com/whiteblock/definition/internal/mocks/entity"
	mockParser "github.com/whiteblock/definition/internal/mocks/parser"
	"github.com/whiteblock/definition/pkg/entity"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBiomeCalculator_AddNextPhase(t *testing.T) {
	testStatePack := entity.NewStatePack(schema.RootSchema{}, config.Bucket{}, logrus.New())

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
