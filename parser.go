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
