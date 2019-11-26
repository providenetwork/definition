package definition

import (
	"errors"
	"fmt"

	"github.com/whiteblock/testexecution/pkg/definition/schema"

	"github.com/whiteblock/genesis/pkg/command"
)

// Commands returns a list of commands as defined in the test definition
func (p *Definition) Commands() ([]command.Command, error) {
	commands := []command.Command{}

	services := map[string]schema.Service{}
	for _, service := range p.spec.Services {
		services[service.Name] = service
	}

	for _, testModel := range p.spec.Tests {
		for _, systemComponent := range testModel.System {
			if _, ok := services[systemComponent.Type]; ok {
				return NewCommandTemplate(), nil
			} else {
				return commands, errors.New(fmt.Sprintf("Missing service %s", systemComponent.Type))
			}
		}
	}
	return commands, nil
}

func NewCommandTemplate() []command.Command {
	return []command.Command{
		command.Command{
			ID:           "",
			Timestamp:    0,
			Timeout:      0,
			Retry:        0,
			Target:       command.Target{},
			Dependencies: nil,
			Order: command.Order{
				Type:    command.Createcontainer,
				Payload: command.Container{},
			},
		},
		command.Command{
			ID:           "",
			Timestamp:    0,
			Timeout:      0,
			Retry:        0,
			Target:       command.Target{},
			Dependencies: []string{},
			Order: command.Order{
				Type:    command.Startcontainer,
				Payload: command.Container{},
			},
		},
	}
}
