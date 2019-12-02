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
	"fmt"
	"testing"

	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal/entity"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func GenerateTestSegments(n int, offset int) []entity.Segment {
	out := make([]entity.Segment, n)
	for i := range out {
		out[i].Name = fmt.Sprint(i + offset)
		out[i].CPUs = int64((i + offset))
		out[i].Memory = int64((i + offset) * 100)
		out[i].Storage = int64((i + offset) * 10000)
	}
	return out
}

func GenerateTestConf(entities []entity.Segment, minBuckets int64) config.Bucket {
	out := config.Bucket{

		MinCPU:      0,
		MinMemory:   0,
		MinStorage:  0,
		UnitCPU:     1,
		UnitMemory:  128,
		UnitStorage: 1000,
		MaxBuckets:  minBuckets * 2,
	}
	for _, entity := range entities {
		out.MaxCPU += entity.CPUs + 1
		out.MaxMemory += entity.Memory + 1
		out.MaxStorage += entity.Storage + 1
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
	bucket := newBucket(&testConf)
	assert.Equal(t, bucket.CPUs, testConf.MinCPU)
	assert.Equal(t, bucket.Memory, testConf.MinMemory)
	assert.Equal(t, bucket.Storage, testConf.MinStorage)
}

func TestBucket_GetSegments(t *testing.T) {
	bucket := Bucket{segments: GenerateTestSegments(10, 0)}
	assert.ElementsMatch(t, bucket.segments, bucket.GetSegments())
}

func TestBucket_Runthrough(t *testing.T) {
	segments := GenerateTestSegments(10, 0)
	conf := GenerateTestConf(segments, 1)
	bucket := newBucket(&conf)

	for _, segment := range segments {
		require.True(t, bucket.hasSpace(segment))
	}

	for _, segment := range segments {
		assert.True(t, bucket.tryAdd(segment))
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
