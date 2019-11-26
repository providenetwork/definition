/*
	Copyright 2019 whiteblock Inc.
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

	"github.com/whiteblock/definition/internal/config"
	"github.com/whiteblock/definition/internal/entity"

	"github.com/stretchr/testify/assert"
	//"github.com/stretchr/testify/require"
)


func GenerateTestSegments(n int) []entity.Segment {
	out := make([]entity.Segment, n)
	for i := range out {
		out[i].Name = fmt.Sprint(i)
		out[i].CPUs = int64(i)
		out[i].Memory = int64(i * 100)
		out[i].Storage = int64(i * 10000)
	}
	return out
}

func TestNewBucket(t *testing.T) {
	testConf := config.Bucket{
		MinCPU: 1,
		MinMemory:2,
		MinStorage:3,
	}
	bucket := newBucket(&testConf)
	assert.Equal(t, bucket.CPUs,testConf.MinCPU)
	assert.Equal(t,bucket.Memory,testConf.MinMemory)
	assert.Equal(t,bucket.Storage,testConf.MinStorage)
}

func TestBucket_GetSegments(t *testing.T) {
	bucket := Bucket{segments: GenerateTestSegments(10)}
	assert.ElementsMatch(t, bucket.segments, bucket.GetSegments())
}

func TestBucket_ToResource(t *testing.T) {

}

func TestBucket_tryAdd(t *testing.T) {

}

func TestBucket_tryRemove(t *testing.T) {

}