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

const (
	LocalScope  = "local"
	GlobalScope = "singleton"
)

type Volume struct {
	Path        string `yaml:"path,omitempty" json:"path,omitempty"`
	Name        string `yaml:"name,omitempty" json:"name,omitempty"`
	Permissions string `yaml:"permissions,omitempty" json:"permissions,omitempty"`
	Scope       string `yaml:"scope" json:"scope,omitempty"`
}

func (sv Volume) GetScope() string {
	if sv.Scope == "" {
		return LocalScope
	}
	return sv.Scope
}

func (sv Volume) Local() bool {
	return sv.GetScope() == LocalScope
}
