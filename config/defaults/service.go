/*
	Copyright 2019 whiteblock Inc.
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

type Service struct {
	Image string `mapstructure:"defaultServiceImage"`
}

func NewService(v *viper.Viper) (Service, error) {
	out := Service{}
	return out, v.Unmarshal(&out)
}

func setServiceBindings(v *viper.Viper) error {
	return v.BindEnv("defaultServiceImage", "DEFAULT_SERVICE_IMAGE")
}

func setServiceDefaults(v *viper.Viper) {
	v.SetDefault("defaultServiceImage", "INFO")
}
