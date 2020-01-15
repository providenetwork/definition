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

//  It would be too redundant to test in different funcs due to the complexity of this
//  Implementation
func TestResourceBuckets(t *testing.T) {
	testSegments := GenerateTestSegments(20, 0)
	testConf := GenerateTestConf(testSegments, 5)

	rb := NewResourceBuckets(testConf, logrus.New())
	require.NotNil(t, rb)
	err := rb.Add(testSegments)
	require.NoError(t, err)

	buckets := rb.Resources()
	require.NotNil(t, buckets)
	assert.True(t, len(buckets) >= 5)

	err = rb.Remove(testSegments)
	assert.NoError(t, err)

	buckets = rb.Resources()
	require.NotNil(t, buckets)
	assert.True(t, len(buckets) >= 5)
}
