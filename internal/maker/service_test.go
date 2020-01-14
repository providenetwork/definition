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

package maker

import (
	"testing"
	"time"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config/defaults"
	mockSearch "github.com/whiteblock/definition/internal/mocks/search"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_FromSystem(t *testing.T) {
	testSystemComp := schema.SystemComponent{
		Count: 5,
		Type:  "foo",
		Args:  nil,
		Environment: map[string]string{
			"bar": "baz",
		},
		Sidecars: []schema.SystemComponentSidecar{
			schema.SystemComponentSidecar{
				Type: "sidecar",
				Name: "sidecar",
				Resources: schema.Resources{
					Cpus:    13,
					Memory:  "",
					Storage: "3233",
				},
				Args: []string{"barfoo"},
				Environment: map[string]string{
					"bar": "baz",
				},
			},
		},
		Resources: schema.SystemComponentResources{
			Cpus:     1,
			Memory:   "",
			Storage:  "20GB",
			Networks: nil,
		},
	}

	testService := schema.Service{
		Args: []string{"hello"},
		Environment: map[string]string{
			"foo": "bar",
			"bar": "foo",
		},
		Resources: schema.Resources{
			Cpus:    12,
			Memory:  "10GB",
			Storage: "2",
		},
	}

	testSidecar := schema.Sidecar{
		Resources: schema.Resources{
			Cpus:    0,
			Memory:  "foobar",
			Storage: "",
		},
		Environment: map[string]string{
			"foo": "bar",
			"bar": "foo",
		},
	}

	searcher := new(mockSearch.Schema)
	searcher.On("FindServiceByType", mock.Anything, testSystemComp.Type).Return(
		testService, nil).Once()
	searcher.On("FindSidecarByType", mock.Anything, mock.Anything).Return(
		testSidecar, nil).Once()

	searcher.On("FindSidecarsByService", mock.Anything, mock.Anything).Return(nil).Once()

	serv := NewService(defaults.Defaults{}, searcher, logrus.New())
	require.NotNil(t, serv)

	results, err := serv.FromSystem(schema.RootSchema{}, testSystemComp)
	assert.NoError(t, err)
	require.NotNil(t, results)
	require.Len(t, results, int(testSystemComp.Count))

	//Merging tests on all of the results
	for _, result := range results {
		_, exists := result.SquashedService.Environment["bar"]
		assert.True(t, exists)

		assert.Equal(t, "baz", result.SquashedService.Environment["bar"])
		assert.Equal(t, testSystemComp.Resources.Storage, result.SquashedService.Resources.Storage)
		assert.Equal(t, testSystemComp.Resources.Cpus, result.SquashedService.Resources.Cpus)
		assert.Equal(t, testService.Resources.Memory, result.SquashedService.Resources.Memory)
		assert.ElementsMatch(t, testService.Args, result.SquashedService.Args)

		require.Len(t, result.Sidecars, len(testSystemComp.Sidecars))
		require.True(t, len(result.Sidecars) > 0)
		sidecar := result.Sidecars[0]
		sysSidecar := testSystemComp.Sidecars[0]
		assert.Equal(t, "baz", result.SquashedService.Environment["bar"])
		assert.Equal(t, sysSidecar.Resources.Storage, sidecar.Resources.Storage)
		assert.Equal(t, sysSidecar.Resources.Cpus, sidecar.Resources.Cpus)
		assert.Equal(t, testSidecar.Resources.Memory, sidecar.Resources.Memory)
		assert.ElementsMatch(t, sysSidecar.Args, sidecar.Args)
	}

	searcher.AssertExpectations(t)
}

func TestService_FromTask(t *testing.T) {
	testTaskRunner := schema.TaskRunner{}

	searcher := new(mockSearch.Schema)
	searcher.On("FindTaskRunnerByType", mock.Anything, mock.Anything).Return(testTaskRunner, nil).Once()
	serv := NewService(defaults.Defaults{}, searcher, logrus.New())
	require.NotNil(t, serv)

	res, err := serv.FromTask(schema.RootSchema{}, schema.Task{
		Timeout: command.Timeout{Duration: 10 * time.Minute},
	}, 0)
	assert.NoError(t, err)
	assert.Equal(t, 10*time.Minute, res.Timeout.Duration)
	searcher.AssertExpectations(t)
}

func TestService_NOOP(t *testing.T) {
	testSystemComp := schema.SystemComponent{
		Count: 0,
		Type:  "foo",
		Args:  nil,
		Environment: map[string]string{
			"bar": "baz",
		},
		Sidecars: []schema.SystemComponentSidecar{
			schema.SystemComponentSidecar{
				Type: "sidecar",
				Name: "sidecar",
				Resources: schema.Resources{
					Cpus:    13,
					Memory:  "",
					Storage: "3233",
				},
				Args: []string{"barfoo"},
				Environment: map[string]string{
					"bar": "baz",
				},
			},
		},
		Resources: schema.SystemComponentResources{
			Cpus:     1,
			Memory:   "",
			Storage:  "20GB",
			Networks: nil,
		},
	}

	testService := schema.Service{}

	testSidecar := schema.Sidecar{}

	searcher := new(mockSearch.Schema)
	searcher.On("FindServiceByType", mock.Anything, testSystemComp.Type).Return(
		testService, nil).Once()
	searcher.On("FindSidecarByType", mock.Anything, mock.Anything).Return(
		testSidecar, nil).Once()

	searcher.On("FindSidecarsByService", mock.Anything, mock.Anything).Return(nil).Once()

	serv := NewService(defaults.Defaults{}, searcher, logrus.New())
	require.NotNil(t, serv)

	results, err := serv.FromSystem(schema.RootSchema{}, testSystemComp)
	assert.Len(t, results, 0)
	assert.NoError(t, err)
}
