/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

import (
	"github.com/whiteblock/definition/schema"
)

type SystemDiff struct {
	Modified        []ServiceDiff
	AddedNetworks   []schema.Network
	RemovedNetworks []schema.Network
	Added           []Service
	Removed         []Service
}

func (diff *SystemDiff) Append(sys *SystemDiff) {
	if sys == nil {
		return
	}
	diff.Modified = append(diff.Modified, sys.Modified...)
	diff.Added = append(diff.Added, sys.Added...)
	diff.Removed = append(diff.Removed, sys.Removed...)
	diff.AddedNetworks = append(diff.AddedNetworks, sys.AddedNetworks...)
	diff.RemovedNetworks = append(diff.RemovedNetworks, sys.RemovedNetworks...)
}
