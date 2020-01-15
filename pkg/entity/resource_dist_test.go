/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceDist(t *testing.T) {
	testSegments := GenerateTestSegments(30, 0)
	testConf := GenerateTestConf(testSegments, 5)

	rb := NewResourceBuckets(testConf, logrus.New())
	require.NotNil(t, rb)

	resourceDist := &ResourceDist{}
	for i := 0; i < 10; i += 5 {
		segments := GenerateTestSegments(i+5, i)
		err := rb.Add(segments)
		require.NoError(t, err)

		buckets := rb.Resources()
		require.NotNil(t, buckets)
		resourceDist.Add(buckets)

		res, err := resourceDist.GetPhase(i / 5)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.ElementsMatch(t, buckets, res)
	}

}
