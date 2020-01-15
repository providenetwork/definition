/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package schema

type Service struct {
	Name        string            `yaml:"name,omitempty" json:"name,omitempty"`
	Description string            `yaml:"description,omitempty" json:"description,omitempty"`
	Volumes     []Volume          `yaml:"volumes,omitempty" json:"volumes,omitempty"`
	Resources   Resources         `yaml:"resources,omitempty" json:"resources,omitempty"`
	Args        []string          `yaml:"args,omitempty" json:"args,omitempty"`
	Environment map[string]string `yaml:"environment,omitempty" json:"environment,omitempty"`
	Image       string            `yaml:"image,omitempty" json:"image,omitempty"`
	Script      Script            `yaml:"script,omitempty" json:"script,omitempty"`
	InputFiles  []InputFile       `yaml:"input-files,omitempty" json:"input-files,omitempty"`
}
