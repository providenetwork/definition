/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package schema

import "github.com/whiteblock/definition/command"

type Task struct {
	Type           string            `yaml:"type,omitempty" json:"type,omitempty"`
	Description    string            `yaml:"description,omitempty" json:"description,omitempty"`
	IgnoreExitCode bool              `yaml:"ignore-exit-code,omitempty" json:"ignore-exit-code,omitempty"`
	Timeout        command.Timeout   `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	Args           []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Environment    map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Networks       []Network         `yaml:"networks,omitempty" json:"networks,omitempty"`
}
