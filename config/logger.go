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

package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Logger struct {
	Verbosity string `mapstructure:"verbosity"`
}

func NewLogger(v *viper.Viper) (Logger, error) {
	out := Logger{}
	return out, v.Unmarshal(&out)
}

//GetLogger gets a logger according to the config
func (l Logger) GetLogger() (*logrus.Logger, error) {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(l.Verbosity)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(lvl)
	return logger, nil
}

func setLoggerBindings(v *viper.Viper) error {
	return v.BindEnv("verbosity", "VERBOSITY")
}

func setLoggerDefaults(v *viper.Viper) {
	v.SetDefault("verbosity", "INFO")
}
