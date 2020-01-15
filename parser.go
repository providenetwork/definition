/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package definition

import (
	"encoding/json"

	"gopkg.in/yaml.v2"
)

// SchemaYAML takes a raw set of YAML bytes and
// de-serializes them to a Definition structure from just the test spec.
func SchemaYAML(raw []byte) (out Definition, err error) {
	return out, yaml.UnmarshalStrict(raw, &out.Spec)
}

// YAML takes a raw set of YAML bytes and
// de-serializes them to a Definition structure.
func YAML(raw []byte) (out Definition, err error) {
	return out, yaml.UnmarshalStrict(raw, &out)
}

// SchemaJSON takes a raw set of JSON bytes and
// de-serializes them to a Definition structure from just the test spec.
func SchemaJSON(raw []byte) (out Definition, err error) {
	return out, json.Unmarshal(raw, &out.Spec)
}

// SchemaANY will first try yaml and then try JSON
func SchemaANY(raw []byte) (Definition, error) {
	def, err := SchemaYAML(raw)
	if err != nil {
		return SchemaJSON(raw)
	}
	return def, nil
}

// JSON takes a raw set of JSON bytes and
// de-serializes them to a Definition structure.
func JSON(raw []byte) (out Definition, err error) {
	return out, json.Unmarshal(raw, &out)
}
