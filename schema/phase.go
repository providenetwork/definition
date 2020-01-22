/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package schema

import "github.com/whiteblock/definition/command"

type Phase struct {
	Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	System      []SystemComponent `yaml:"system,omitempty" json:"system,omitempty"`
	Remove      []string          `yaml:"remove,omitempty" json:"remove,omitempty"`
	Tasks       []Task            `yaml:"tasks,omitempty" json:"tasks,omitempty"`
	Timeout     command.Timeout   `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Duration    command.Duration  `yaml:"duration,omitempty" json:"duration,omitempty"`
}
