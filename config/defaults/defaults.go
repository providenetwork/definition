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
