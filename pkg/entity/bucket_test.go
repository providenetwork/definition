/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

import (
	"fmt"
	"testing"

	"github.com/whiteblock/definition/config"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GenerateTestSegments(n int, offset int) []Segment {
	out := make([]Segment, n)
	for i := range out {
		out[i].Name = fmt.Sprint(i + offset)
		out[i].CPUs = int64((i + offset))
		out[i].Memory = int64((i + offset) * 100)
		out[i].Storage = int64((i + offset) * 10000)
	}
	return out
}

func GenerateTestConf(entities []Segment, minBuckets int64) config.Bucket {
	out := config.Bucket{
		MinCPU:      0,
		MinMemory:   0,
		MinStorage:  0,
		UnitCPU:     1,
		UnitMemory:  128,
		UnitStorage: 1000,
		MaxBuckets:  minBuckets * 2,
	}
	for _, e := range entities {
		out.MaxCPU += e.CPUs + 1
		out.MaxMemory += e.Memory + 1
		out.MaxStorage += e.Storage + 1
	}
	out.MaxCPU = roundValueAndMax(0, out.MaxCPU/minBuckets, out.UnitCPU)
	out.MaxMemory = roundValueAndMax(0, out.MaxMemory/minBuckets, out.UnitMemory)
	out.MaxStorage = roundValueAndMax(0, out.MaxStorage/minBuckets, out.UnitStorage)
	return out
}

func TestNewBucket(t *testing.T) {
	testConf := config.Bucket{
		MinCPU:     1,
		MinMemory:  2,
		MinStorage: 3,
	}
	bucket := NewBucket(&testConf, logrus.New())
	assert.Equal(t, bucket.CPUs, testConf.MinCPU)
	assert.Equal(t, bucket.Memory, testConf.MinMemory)
	assert.Equal(t, bucket.Storage, testConf.MinStorage)
}

func TestBucket_GetSegments(t *testing.T) {
	bucket := Bucket{segments: GenerateTestSegments(10, 0)}
	assert.ElementsMatch(t, bucket.segments, bucket.GetSegments())
}

func TestBucket_Clone(t *testing.T) {
	bucket := Bucket{
		Resource: Resource{CPUs: 1, Memory: 2, Storage: 3},
		segments: GenerateTestSegments(10, 0)}
	assert.Equal(t, bucket, bucket.Clone())
}

func TestBucket_Runthrough(t *testing.T) {
	segments := GenerateTestSegments(10, 0)
	conf := GenerateTestConf(segments, 1)
	bucket := NewBucket(&conf, logrus.New())

	for _, segment := range segments {
		require.True(t, bucket.hasSpace(segment))
	}
	expectedCPU := int64(0)
	for _, segment := range segments {
		assert.True(t, bucket.tryAdd(segment))
		expectedCPU += segment.CPUs
		assert.Equal(t, expectedCPU, bucket.usage.CPUs)
		assert.True(t, expectedCPU <= bucket.CPUs)
	}

	require.ElementsMatch(t, segments, bucket.GetSegments())

	for _, segment := range segments {
		assert.True(t, bucket.findSegment(segment) >= 0)
	}

	for _, segment := range segments {
		assert.True(t, bucket.tryRemove(segment))
	}
	for _, segment := range GenerateTestSegments(10, 11) {
		assert.False(t, bucket.tryRemove(segment))
	}
}
