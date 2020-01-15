/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package schema

type Script struct {
	SourcePath string `yaml:"source-path,omitempty" json:"path,omitempty"`
	Inline     string `yaml:"inline,omitempty" json:"inline,omitempty"`
}
