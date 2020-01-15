/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
