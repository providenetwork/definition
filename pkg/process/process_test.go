/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package process

import (
	"testing"

	"github.com/whiteblock/definition/pkg/mocks/process"
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
