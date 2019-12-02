/*
	Copyright 2019 whiteblock Inc.
	This file is a part of the definition.

	definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	definition is distributed in the hope that it will be useful,
	but dock ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package definition

import (
	"encoding/json"
	"fmt"

	"github.com/whiteblock/definition/schema"
	"github.com/whiteblock/definition/validator"
)

//IDefinition is the representation of the test definition format.
type IDefinition interface {
	//Gets the UUID which uniquely identifies this definition
	GetID() string

	//GetOrgID gets the organization id
	GetOrgID() int64

	//GetSpec gets a pointer to the internal spec object. Should be used with care.
	GetSpec() *schema.RootSchema

	//Validate returns nil if it is a valid test definition, other the returned
	//error will contain an explanation for the issue
	Validate() []error
}

// Definition is the top level container for the test definition
// specification
type Definition struct {
	ID    string
	OrgID int64
	spec  schema.RootSchema
}

//Gets the UUID which uniquely identifies this definition
func (def Definition) GetID() string {
	return def.ID
}

//GetOrgID gets the organization id
func (def Definition) GetOrgID() int64 {
	return def.OrgID
}

//GetSpec gets a pointer to the internal spec object. Should be used with care.
func (def Definition) GetSpec() *schema.RootSchema {
	return &def.spec
}

//Validate returns nil if it is a valid test definition, other the returned
//error will contain an explanation for the issue
func (def Definition) Validate() []error {
	v, err := validator.NewValidator()
	if err != nil {
		return []error{err}
	}

	data, err := json.Marshal(def.spec)
	if err != nil {
		return []error{err}
	}

	err = v.Validate(data)
	if err == nil {
		return nil
	}
	out := []error{err}
	for _, schemaErr := range v.Errors() {
		out = append(out, fmt.Errorf(schemaErr.String()))
	}
	return out
}

