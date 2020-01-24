/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package merger

import (
	"github.com/whiteblock/definition/schema"

	"github.com/imdario/mergo"
	"github.com/jinzhu/copier"
)

func MergeSystemLeft(sys schema.SystemComponent, systems ...schema.SystemComponent) (system schema.SystemComponent) {
	copier.Copy(&system, sys)
	for _, merging := range systems {
		cnt := system.Count
		mergo.Map(&system, merging, mergo.WithOverride)
		if system.Count == 0 {
			system.Count = cnt
		}
	}
	return
}
