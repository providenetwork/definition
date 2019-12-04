/*
	Copyright 2019 Whiteblock Inc.
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
	"github.com/whiteblock/definition/schema"
)

// Definition is the top level container for
// the test definition specification.
type Definition struct {
	// ID is the test ID
	ID    string

	// OrgID
	OrgID int64

	Spec schema.RootSchema
}
