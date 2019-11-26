/*
	Copyright 2019 whiteblock Inc.
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
package internal

import (
	"encoding/json"
	"github.com/whiteblock/definition/schema"
	"github.com/whiteblock/definition/validator"
	"gopkg.in/yaml.v2"
)

// Definition is the top level container for the test definition
// specification
type Definition struct {
	ID   string
	spec schema.RootSchema
}

// ParseYAML takes a raw set of bytes and
// de-serializes them to a Definition structure
func ParseYAML(raw []byte) (Definition, error) {
	data := Definition{}
	err := yaml.Unmarshal(raw, &data)

	if err != nil {
		return data, err
	}

	return data, nil
}

func (d *Definition) Valid() (interface{}, error) {
	v, err := validator.NewValidator()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(d.spec)
	if err != nil {
		return nil, err
	}

	err = v.Validate(data)
	if err != nil {
		return v.Errors(), err
	}

	return nil, nil
}
