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
