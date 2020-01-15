/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package definition

import (
	"github.com/whiteblock/definition/schema"
)

// Definition is the top level container for
// the test definition specification.
type Definition struct {
	// ID is the test ID
	ID string

	// OrgID
	OrgID string

	Spec schema.RootSchema
}
