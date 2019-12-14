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

package schema

type InputFile struct {
	SourcePath      string `yaml:"source-path,omitempty" json:"source-path,omitempty"`
	DestinationPath string `yaml:"destination-path,omitempty" json:"destination-path,omitempty"`
	Template        bool   `yaml:"template,omitempty" json:"template,omitempty"`
}

// GetSource makes it easy to change the name of the source member, as it is
// expected to change in the near future. (Also, there might be some logic behind it as well in
// the future)
func (in InputFile) GetSource() string {
	return in.SourcePath
}
