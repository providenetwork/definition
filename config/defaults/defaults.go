/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

//package defaults provides structures for the defaults of the test definition format
package defaults

import (
	"github.com/spf13/viper"
)

// Defaults is a top level contain for all of the defaults so that they all
// may be passed into this library by the user. This is not to be confused with configuration
type Defaults struct {
	Service   Service
	Resources Resources
	Network   Network
}

// New creates a Defaults by generating each field from the given viper config
// provider
func New(v *viper.Viper) (def Defaults, err error) {
	def.Service, err = NewService(v)
	if err != nil {
		return Defaults{}, err
	}

	def.Resources, err = NewResources(v)
	if err != nil {
		return Defaults{}, err
	}

	def.Network, err = NewNetwork(v)
	if err != nil {
		return Defaults{}, err
	}

	return
}

// SetViperBindings adds all of the enviroment bindings to the given
// viper config provider
func SetViperBindings(v *viper.Viper) error {
	err := setResourcesBindings(v)
	if err != nil {
		return err
	}
	err = setNetworkBindings(v)
	if err != nil {
		return err
	}
	return setServiceBindings(v)
}

// SetViperDefaults adds all of the default values to the given
// viper config provider
func SetViperDefaults(v *viper.Viper) {
	setResourcesDefaults(v)
	setServiceDefaults(v)
	setNetworkDefaults(v)
}
