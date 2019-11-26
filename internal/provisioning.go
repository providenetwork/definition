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
package definition

import (
	"github.com/whiteblock/testexecution/pkg/definition/schema"
	"github.com/whiteblock/testexecution/pkg/util"

	"github.com/whiteblock/provisioner/commands"
)

func PrepareInstance(sysComponent schema.SystemComponent) (commands.Instance, error) {
	mem, err := util.Memconv(sysComponent.Resources.Memory)
	if err != nil {
		return commands.Instance{}, err
	}

	storage, err := util.Memconv(sysComponent.Resources.Storage) // todo not sure if memconv will return expected results for storage..?
	if err != nil {
		return commands.Instance{}, err
	}

	return commands.Instance{
		CPUs:    int64(sysComponent.Resources.Cpus),
		Memory:  mem,
		Storage: storage,
	}, nil
}

func (td *Definition) NewProvisioningRequest() ([]commands.CreateBiome, error) {
	cb := make([]commands.CreateBiome, 0)
	for _, test := range td.spec.Tests {
		instances := make([]commands.Instance, 0)
		for _, sysComponent := range test.System {
			instance, err := PrepareInstance(sysComponent)
			if err != nil {
				return cb, err
			}

			instances = append(instances, instance)
		}

		cb = append(cb, commands.CreateBiome{
			TestnetID: td.ID,
			OrgID:     0, //todo where to get OrgID
			Instances: instances,
		})
	}

	return cb, nil
}
