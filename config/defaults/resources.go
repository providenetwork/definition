/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package defaults

import "github.com/spf13/viper"

type Resources struct {
	CPUs    int64 `mapstructure:"defaultCpus"`
	Memory  int64 `mapstructure:"defaultMemory"`
	Storage int64 `mapstructure:"defaultStorage"`
}

func NewResources(v *viper.Viper) (out Resources, err error) {
	return out, v.Unmarshal(&out)
}

func setResourcesBindings(v *viper.Viper) error {
	err := v.BindEnv("defaultCpus", "DEFAULT_CPUS")
	if err != nil {
		return err
	}

	err = v.BindEnv("defaultMemory", "DEFAULT_MEMORY")
	if err != nil {
		return err
	}

	return v.BindEnv("defaultStorage", "DEFAULT_STORAGE")
}

func setResourcesDefaults(v *viper.Viper) {
	v.SetDefault("defaultCpus", 1)
	v.SetDefault("defaultMemory", 512)
	v.SetDefault("defaultStorage", 5*1024)
}
