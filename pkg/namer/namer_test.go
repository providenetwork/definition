/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package namer

import (
	"testing"
)

func TestNameWithSlashes(t *testing.T) {
	result := SanitizeVolumeName("wat/now")
	if result != "wat-now" {
		t.Fatalf("Invalid text substitution %s", result)
	}
}

func TestNameWithStartingSlash(t *testing.T) {
	result := SanitizeVolumeName("/wat/now")
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
