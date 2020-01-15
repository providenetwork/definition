/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

//package config contains the structures used for configuration for this library
package config

import (
	"github.com/whiteblock/definition/config/defaults"

	"github.com/spf13/viper"
)

//  Config is the union of all of the configuration structures of this
//library, so that it may be passed in directly by the user
type Config struct {
	// Bucket is the configuration for buckets
	Bucket Bucket

	// Output is the configuration for the command output
	Output Output

	// Logger is the configuration for the loggers
	Logger Logger

	// Network is the configuration for the networks
	Network Network

	// Defaults is the configuration of the defaults
	Defaults defaults.Defaults
}

//New generates a new config by constructing each field from viper
func New(v *viper.Viper) (conf Config, err error) {
	conf.Bucket, err = NewBucket(v)
	if err != nil {
		return Config{}, err
	}

	conf.Logger, err = NewLogger(v)
	if err != nil {
		return Config{}, err
	}

	conf.Network, err = NewNetwork(v)
	if err != nil {
		return Config{}, err
	}

	conf.Output, err = NewOutput(v)
	if err != nil {
		return Config{}, err
	}

	conf.Defaults, err = defaults.New(v)
	return
}

//SetViperBindings adds all of the environment bindings to the given
//viper config provider, for all of the configs and defaults
func SetViperBindings(v *viper.Viper) error {
	err := setLoggerBindings(v)
	if err != nil {
		return err
	}

	err = setBucketBindings(v)
	if err != nil {
		return err
	}

	err = setNetworkBindings(v)
	if err != nil {
		return err
	}

	err = setOutputBindings(v)
	if err != nil {
		return err
	}

	return defaults.SetViperBindings(v)
}

//SetViperDefaults adds all of the default values to the given
//viper config provider, for all of the configs and defaults
func SetViperDefaults(v *viper.Viper) {
	setLoggerDefaults(v)
	setBucketDefaults(v)
	setNetworkDefaults(v)
	setOutputDefaults(v)
	defaults.SetViperDefaults(v)
}

//SetupViper applies SetViperDefaults and SetViperBindings to the given
//viper provider
func SetupViper(v *viper.Viper) error {
	SetViperDefaults(v)
	return SetViperBindings(v)
}
