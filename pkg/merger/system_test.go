/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package merger

import (
	"testing"

	"github.com/whiteblock/definition/schema"

	"github.com/stretchr/testify/assert"
)

func TestSystem_MergeLeft_Count(t *testing.T) {
	base := schema.SystemComponent{
		Count: 5,
	}
	add := schema.SystemComponent{}

	res := MergeSystemLeft(base, add)

	assert.Equal(t, base.Count, res.Count)
}
