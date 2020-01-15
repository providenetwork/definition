/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
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
