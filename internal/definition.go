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
	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/command/biome"
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


//Gets the UUID which uniquely identifies this Definition
func (def Definition) GetID() string {
	//TODO
	return ""
}

//GetCommands gets the commands in dependency groups, so that
//res[n+1] is the set of commands which require the execution of the commands
//in res[n].
func (def Definition) GetCommands() ([][]command.Command, error) {
	//TODO
	return nil, nil
}

//GetProvisioningRequest calculates the biome resources necessary to support
//this test
func (def Definition) GetProvisioningRequest() (biome.CreateBiome, error) {
	//TODO
	return biome.CreateBiome{}, nil
}

//Validate returns nil if it is a valid test definition, other the returned
//error will contain an explanation for the issue
func (def Definition) Validate() error {
	//TODO
	return nil
}