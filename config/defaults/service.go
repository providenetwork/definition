/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package defaults

import "github.com/spf13/viper"

type Service struct {
	Resources
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
	v.SetDefault("defaultServiceImage", "ubuntu:latest")
}
