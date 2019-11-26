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
