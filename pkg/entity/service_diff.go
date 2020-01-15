/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

import (
	"github.com/whiteblock/definition/schema"
)

type ServiceDiff struct {
	Name           string
	AddNetworks    []schema.Network
	UpdateNetworks []schema.Network
	DetachNetworks []schema.Network
	AddSidecars    []schema.Sidecar
	RemoveSidecars []schema.Sidecar

	Parent *Service
}
