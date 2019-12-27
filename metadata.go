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

type Metadata struct {
	Name  string
	OrgID string

	Labels map[string]string
}

func (m *Metadata) ToLabels() map[string]string {
	l := map[string]string{}

	if m.Name != "" {
		l["name"] = m.Name
	}

	if m.OrgID != "" {
		l["orgID"] = m.OrgID
	}

	for key, val := range m.Labels {
		l[key] = val
	}

	return l
}

func (m *Metadata) ToFlattenedLabels() []string {
	f := []string{}
	for key, val := range m.ToLabels() {
		f = append(f, key)
		f = append(f, val)
	}

	return f
}
