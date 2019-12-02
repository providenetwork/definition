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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResourceDist(t *testing.T) {
	testSegments := GenerateTestSegments(30, 0)
	testConf := GenerateTestConf(testSegments, 5)

	rb := NewResourceBuckets(testConf)
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
