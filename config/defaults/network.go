/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package defaults

import (
	"github.com/spf13/viper"
)

type Network struct {
	Name string `mapstructure:"defaultNetworkName"`
}

func NewNetwork(v *viper.Viper) (out Network, err error) {
	return out, v.Unmarshal(&out)
}

func setNetworkBindings(v *viper.Viper) error {
	return v.BindEnv("defaultNetworkName", "DEFAULT_NETWORK_NAME")
}

func setNetworkDefaults(v *viper.Viper) {
	v.SetDefault("defaultNetworkName", "default")
}
