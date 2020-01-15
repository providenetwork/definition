/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

import "github.com/whiteblock/utility/common"

// File represents a file which will be placed inside either a docker container or volume
type File struct {
	// Mode is permission and mode bits
	Mode int64 `json:"mode"`
	// Destination is the mount point of the file
	Destination string `json:"destination"`
	// ID is the UUID of the file
	ID string `json:"id"`

	// Meta data on the file
	Meta common.Metadata `json:"meta"`
}
