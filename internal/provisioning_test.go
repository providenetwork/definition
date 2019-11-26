package definition

import (
	"testing"

	"github.com/whiteblock/testexecution/pkg/definition/schema"

	"github.com/stretchr/testify/assert"
	"github.com/whiteblock/provisioner/commands"
)

func TestPrepareInstance(t *testing.T) {
	sysComponent0 := schema.SystemComponent{}
	sysComponent0.Resources.Cpus = float64(3)
	sysComponent0.Resources.Memory = "5GB"
	sysComponent0.Resources.Storage = "10GB"

	expected := commands.Instance{
		CPUs:    int64(3),
		Memory:  int64(5000000000),
		Storage: int64(10000000000),
	}

	instance, err := PrepareInstance(sysComponent0)
	if err != nil {
		t.Error("PrepareInstance does not return expected result")
	}

	assert.Equal(t, expected, instance)
}
