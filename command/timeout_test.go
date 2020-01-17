/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	yaml "gopkg.in/yaml.v2"
)


func TestTimeoutMarshal(t *testing.T) {
	to := Timeout{Duration:0,isInfinite:true}
	data, err := yaml.Marshal(to)
	require.NoError(t, err)
	assert.Equal(t,"infinite",strings.TrimSpace(string(data)))
}


func TestTimeoutJSONUnmarshal(t *testing.T) {
	var to Timeout
	err := json.Unmarshal([]byte(`"infinite"`), &to)
	require.NoError(t, err)
	assert.True(t, to.IsInfinite())
}
