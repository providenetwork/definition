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