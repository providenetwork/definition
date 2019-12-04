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

package config

import (
	"github.com/spf13/viper"
)

//Network is the configuration for the biome networks
type Network struct {
	GlobalNetwork      string `mapstructure:"globalNetwork"`
	SidecarNetwork     string `mapstructure:"sidecarNetwork"`
	MaxNodesPerNetwork int    `mapstructure:"maxNodesPerNetwork"`
}

//NewNetwork creates a new network config from the given viper config provider
func NewNetwork(v *viper.Viper) (out Network, err error) {
	return out, v.Unmarshal(&out)
}

func setNetworkBindings(v *viper.Viper) error {
	err := v.BindEnv("globalNetwork", "GLOBAL_NETWORK")
	if err != nil {
		return err
	}

	err = v.BindEnv("sidecarNetwork", "SIDECAR_NETWORK")
	if err != nil {
		return err
	}

	return v.BindEnv("maxNodesPerNetwork", "MAX_NODES_PER_NETWORK")
}

func setNetworkDefaults(v *viper.Viper) {
	v.SetDefault("globalNetwork", "10.0.0.0/8")
	v.SetDefault("sidecarNetwork", "172.31.0.0/13")
	v.SetDefault("maxNodesPerNetwork", 1000)
}
