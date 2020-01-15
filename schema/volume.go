/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
