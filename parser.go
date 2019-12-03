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

var globalParser = NewParser()

// Parser represents a parser of the test definition format
type Parser interface {
	// SchemaYAML takes a raw set of YAML bytes and
	// de-serializes them to a Definition structure from just the test spec.
	SchemaYAML(raw []byte) (Definition, error)

	// YAML takes a raw set of YAML bytes and
	// de-serializes them to a Definition structure.
	YAML(raw []byte) (Definition, error)

	// SchemaJSON takes a raw set of JSON bytes and
	// de-serializes them to a Definition structure from just the test spec.
	SchemaJSON(raw []byte) (Definition, error)

	// JSON takes a raw set of JSON bytes and
	// de-serializes them to a Definition structure.
	JSON(raw []byte) (Definition, error)
}

type parser struct {
}

// NewParser creates a new parser
func NewParser() Parser {
	return &parser{}
}

// SchemaYAML takes a raw set of YAML bytes and
// de-serializes them to a Definition structure from just the test spec.
func (p parser) SchemaYAML(raw []byte) (out Definition, err error) {
	return out, yaml.Unmarshal(raw, &out.spec)
}

// YAML takes a raw set of YAML bytes and
// de-serializes them to a Definition structure.
func (p parser) YAML(raw []byte) (out Definition, err error) {
	return out, yaml.Unmarshal(raw, &out)
}

// SchemaJSON takes a raw set of JSON bytes and
// de-serializes them to a Definition structure from just the test spec.
func (p parser) SchemaJSON(raw []byte) (out Definition, err error) {
	return out, json.Unmarshal(raw, &out.spec)
}

// JSON takes a raw set of JSON bytes and
// de-serializes them to a Definition structure.
func (p parser) JSON(raw []byte) (out Definition, err error) {
	return out, json.Unmarshal(raw, &out)
}

// SchemaYAML takes a raw set of bytes and
// de-serializes them to a Definition structure from just the test spec.
func SchemaYAML(raw []byte) (Definition, error) {
	return globalParser.SchemaYAML(raw)
}

// YAML takes a raw set of YAML bytes and
// de-serializes them to a Definition structure.
func YAML(raw []byte) (Definition, error) {
	return globalParser.YAML(raw)
}

// SchemaJSON takes a raw set of JSON bytes and
// de-serializes them to a Definition structure from just the test spec.
func SchemaJSON(raw []byte) (Definition, error) {
	return globalParser.SchemaJSON(raw)
}

// JSON takes a raw set of JSON bytes and
// de-serializes them to a Definition structure.
func JSON(raw []byte) (Definition, error) {
	return globalParser.JSON(raw)
}
