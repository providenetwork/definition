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

package defaults

import (
	"github.com/spf13/viper"
)

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
