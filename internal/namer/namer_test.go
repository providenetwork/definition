/*
	Copyright 2019 Whiteblock Inc.
	This file is a part of the Definitio

	Definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later versio

	Definition is distributed in the hope that it will be useful,
	but dock ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package namer

import (
	"github.com/whiteblock/definition/schema"
	"testing"
)

func TestNameWithSlashes(t *testing.T) {
	result := InputFileVolume(schema.InputFile{
		SourcePath:      "",
		DestinationPath: "wat/now",
		Template:        false,
	})
	if result != "wat-now" {
		t.Fatalf("Invalid text substitution %s", result)
	}
}

func TestNameWithStartingSlash(t *testing.T) {
	result := InputFileVolume(schema.InputFile{
		SourcePath:      "",
		DestinationPath: "/wat/now",
		Template:        false,
	})
	if result != "wb_wat-now" {
		t.Fatalf("Invalid text substitution %s", result)
	}
}

func TestSanitizeAnyName(t *testing.T) {
	result := SanitizeVolumeName("whatever_you//say")
	if result != "whatever_you--say" {
		t.Fatalf("Invalid text substitution %s", result)
	}
}