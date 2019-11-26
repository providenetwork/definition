package definition

import (
	"encoding/json"

	"github.com/whiteblock/testexecution/pkg/definition/schema"
	"github.com/whiteblock/testexecution/pkg/definition/validator"
)

// Definition is the top level container for the test definition
// specification
type Definition struct {
	ID string

	spec schema.RootSchema
}

func NewDefintion() string {
	return ""
}

func (d *Definition) Valid() (interface{}, error) {
	v, err := validator.NewValidator()
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(d.spec)
	if err != nil {
		return nil, err
	}

	err = v.Validate(data)
	if err != nil {
		return v.Errors(), err
	}

	return nil, nil
}
