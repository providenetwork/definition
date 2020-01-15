/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package config

import (
	"github.com/spf13/viper"
)

// Output represents the basic configuration for a bucket
type Output struct {
	MaxParallelCommands int  `mapstructure:"maxParallelCommands"`
	AsIsCommands        bool `mapstructure:"asIsCommands"`
}

// NewOutput generates a Bucket configuration from the given viper
// Configuration
func NewOutput(v *viper.Viper) (out Output, err error) {
	err = v.Unmarshal(&out)
	if err != nil {
		return
	}
	if !out.AsIsCommands && out.MaxParallelCommands <= -1 {
		out.AsIsCommands = true //MaxParallelCommands <= 0 implies no edit
	}
	return out, v.Unmarshal(&out)
}

func setOutputBindings(v *viper.Viper) error {
	err := v.BindEnv("asIsCommands", "AS_IS_COMMANDS")
	if err != nil {
		return err
	}
	return v.BindEnv("maxParallelCommands", "MAX_PARALLEL_COMMANDS")
}

func setOutputDefaults(v *viper.Viper) {
	v.SetDefault("asIsCommands", false)
	v.SetDefault("maxParallelCommands", 10)
}
