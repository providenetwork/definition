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

package parser

import (
	"testing"

	"github.com/whiteblock/definition/internal/entity"
	mockConverter "github.com/whiteblock/definition/internal/mocks/converter"
	mockParser "github.com/whiteblock/definition/internal/mocks/parser"
	mockSearch "github.com/whiteblock/definition/internal/mocks/search"
	"github.com/whiteblock/definition/schema"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestResources_SystemComponent(t *testing.T) {
	testSystemComp := schema.SystemComponent{
		Count: 5,
		Type:  "foo",
	}
	namer := new(mockParser.Names)
	for i := 0; i < int(testSystemComp.Count); i++ {
		namer.On("SystemService", testSystemComp, i).Return("foo").Once()
	}
	searcher := new(mockSearch.Schema)
	searcher.On("FindServiceByType", mock.Anything, testSystemComp.Type).Return(
		schema.Service{}, nil).Once()

	conv := new(mockConverter.Resource)

	conv.On("FromResources", mock.Anything).Return(
		entity.Resource{}, nil).Times(int(testSystemComp.Count))

	res := NewResources(namer, searcher, conv)

	ents, err := res.SystemComponent(schema.RootSchema{}, testSystemComp)
	assert.NoError(t, err)
	require.NotNil(t, ents)
	assert.Len(t, ents, int(testSystemComp.Count))

	namer.AssertExpectations(t)
	searcher.AssertExpectations(t)
	conv.AssertExpectations(t)
}

func TestResources_SystemComponentNamesOnly(t *testing.T) {
	testSystemComp := schema.SystemComponent{
		Count: 5,
		Type:  "foo",
	}
	namer := new(mockParser.Names)
	for i := 0; i < int(testSystemComp.Count); i++ {
		namer.On("SystemService", testSystemComp, i).Return("foo").Once()
	}
	res := NewResources(namer, nil, nil)
	ents := res.SystemComponentNamesOnly(testSystemComp)
	require.NotNil(t, ents)
	assert.Len(t, ents, int(testSystemComp.Count))
}

func TestResources_Tasks(t *testing.T) {
	testTasks := []schema.Task{
		schema.Task{},
		schema.Task{},
		schema.Task{},
		schema.Task{},
		schema.Task{},
	}
	namer := new(mockParser.Names)
	for i := 0; i < len(testTasks); i++ {
		namer.On("Task", mock.Anything, i).Return("foo").Once()
	}

	searcher := new(mockSearch.Schema)
	searcher.On("FindTaskRunnerByType", mock.Anything, mock.Anything).Return(
		schema.TaskRunner{}, nil).Times(len(testTasks))

	conv := new(mockConverter.Resource)
	conv.On("FromResources", mock.Anything).Return(
		entity.Resource{}, nil).Times(len(testTasks))

	res := NewResources(namer, searcher, conv)

	result, err := res.Tasks(schema.RootSchema{}, testTasks)
	assert.NoError(t, err)
	require.NotNil(t, result)
	assert.Len(t, result, len(testTasks))

	namer.AssertExpectations(t)
	searcher.AssertExpectations(t)
	conv.AssertExpectations(t)
}