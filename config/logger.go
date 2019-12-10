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
	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Logger struct {
	Verbosity      string `mapstructure:"verbosity"`
	FluentDLogging bool   `mapstructure:"fluentDLogging"`
}

func NewLogger(v *viper.Viper) (Logger, error) {
	out := Logger{}
	return out, v.Unmarshal(&out)
}

//  GetLogger gets a logger according to the config
func (l Logger) GetLogger() (*logrus.Logger, error) {
	logger := logrus.New()
	lvl, err := logrus.ParseLevel(l.Verbosity)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(lvl)
	logger.SetReportCaller(true)
	if l.FluentDLogging {
		logger.SetFormatter(joonix.NewFormatter())
	}

	return logger, nil
}

func setLoggerBindings(v *viper.Viper) error {
	err := v.BindEnv("fluentDLogging", "FLUENT_D_LOGGING")
	if err != nil {
		return err
	}
	return v.BindEnv("verbosity", "VERBOSITY")
}

func setLoggerDefaults(v *viper.Viper) {
	v.SetDefault("fluentDLogging", true)
	v.SetDefault("verbosity", "INFO")
}
