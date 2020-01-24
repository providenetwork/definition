/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package maker

import (
	"testing"
	"time"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config/defaults"
	mockSearch "github.com/whiteblock/definition/pkg/mocks/search"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var testSystemComp = schema.SystemComponent{
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

var testService = schema.Service{
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

var testSidecar = schema.Sidecar{
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

func TestService_FromSystem(t *testing.T) {

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

func TestService_FromSystemDiff_NOOP(t *testing.T) {

	testSystemComp2 := schema.SystemComponent{
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

	searcher := new(mockSearch.Schema)
	searcher.On("FindServiceByType", mock.Anything, testSystemComp.Type).Return(
		testService, nil).Twice()
	searcher.On("FindSidecarByType", mock.Anything, mock.Anything).Return(
		testSidecar, nil).Twice()

	searcher.On("FindSidecarsByService", mock.Anything, mock.Anything).Return(nil).Twice()

	serv := NewService(defaults.Defaults{}, searcher, logrus.New())
	require.NotNil(t, serv)

	results, err := serv.FromSystemDiff(schema.RootSchema{}, testSystemComp, testSystemComp2)
	assert.NoError(t, err)
	require.NotNil(t, results)

	require.Len(t, results.Added, 0)
	t.Logf("%+v", results.Removed)
	require.Len(t, results.Removed, 0)
	t.Logf("%+v", results.Modified)
	require.Len(t, results.Modified, 0)

	searcher.AssertExpectations(t)
}

func TestService_FromTask(t *testing.T) {
	testTaskRunner := schema.TaskRunner{}

	searcher := new(mockSearch.Schema)
	searcher.On("FindTaskRunnerByType", mock.Anything, mock.Anything).Return(testTaskRunner, nil).Once()
	serv := NewService(defaults.Defaults{}, searcher, logrus.New())
	require.NotNil(t, serv)

	res, err := serv.FromTask(schema.RootSchema{}, schema.Task{
		Timeout: command.Timeout{command.Time{Duration: 10 * time.Minute}},
	}, 0)
	assert.NoError(t, err)
	assert.Equal(t, 10*time.Minute, res.Timeout.Duration)
	searcher.AssertExpectations(t)
}

func TestService_Defaults(t *testing.T) {
	testSystemComp := schema.SystemComponent{
		Count: 0,
		Type:  "foo",
		Args:  nil,
		Environment: map[string]string{
			"bar": "baz",
		},
		Sidecars: []schema.SystemComponentSidecar{
			{
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
	assert.Len(t, results, 1)
	assert.NoError(t, err)
}
