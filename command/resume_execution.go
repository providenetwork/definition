/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

import (
	"strings"
	"time"
)

type ResumeExecution struct {
	Tasks []string `json:"tasks,omitempty"`
}