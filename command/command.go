/*
	Copyright 2019 Whiteblock Inc.
	This file is a part of the Definition.

	Definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Definition is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package command

import (
	"bytes"
	"encoding/json"

	"github.com/google/uuid"
)

// OrderType is the type of order
type OrderType string

const (
	// Createcontainer attempts to create a docker container
	Createcontainer = OrderType("createcontainer")
	// Startcontainer attempts to start an already created docker container
	Startcontainer = OrderType("startcontainer")
	// Removecontainer attempts to remove a container
	Removecontainer = OrderType("removecontainer")
	// Createnetwork attempts to create a network
	Createnetwork = OrderType("createnetwork")
	// Attachnetwork attempts to remove a network
	Attachnetwork = OrderType("attachnetwork")
	// Detachnetwork detaches network
	Detachnetwork = OrderType("detachnetwork")
	// Removenetwork removes network
	Removenetwork = OrderType("removenetwork")
	// Createvolume creates volume
	Createvolume = OrderType("createvolume")

	// Removevolume removes volume
	Removevolume = OrderType("removevolume")

	// Putfileincontainer puts file in container
	Putfileincontainer = OrderType("putfileincontainer")

	// Emulation emulates
	Emulation = OrderType("emulation")

	// SwarmInit sets up the docker swarm
	SwarmInit = OrderType("swarminit")

	// Pullimage pre-emptively pulls the given image
	Pullimage = OrderType("pullimage")
)

// OrderPayload is a pointer interface for order payloads.
type OrderPayload interface {
}

// Order to be executed by Definition
type Order struct {
	// Type is the type of the order
	Type OrderType `json:"type"`
	// Payload is the payload object of the order
	Payload OrderPayload `json:"payload"`
}

// Target sets the target of a command - which testnet, instance to hit
type Target struct {
	IP string `json:"ip"`
}

// Command is the command sent to Definition.
type Command struct {
	// ID is the unique id of this command
	ID string `json:"id"`

	// Target represents the target of this command
	Target Target `json:"target"`

	// Order is the action of the command, it represents what needs to be done
	Order Order `json:"order"`

	// Meta is extra informative data to be passed with the command
	Meta map[string]string `json:"meta"`

	// Parent is a pointer to the Instructions object that contains this commands
	parent *Instructions `json:"-"`
}

// NewCommand properly creates a new command
func NewCommand(order Order, endpoint string) (Command, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return Command{}, err
	}
	return Command{
		ID: id.String(),
		Target: Target{
			IP: endpoint, //endpoint,
		},
		Order: order,
		Meta:  map[string]string{},
	}, nil
}

func (cmd Command) TestID() string {
	if cmd.Parent() == nil {
		return ""
	}
	return cmd.Parent().ID
}

// Parent exists to prevent parent from showing up when a command is marshal with
// any marshaller
func (cmd Command) Parent() *Instructions {
	return cmd.parent
}

// ParseOrderPayloadInto attempts to Marshal the payload into the object pointed to by out
func (cmd Command) ParseOrderPayloadInto(out interface{}) error {
	raw, err := json.Marshal(cmd.Order.Payload)
	if err != nil {
		return err
	}
	rdr := bytes.NewReader(raw)
	decoder := json.NewDecoder(rdr)
	decoder.DisallowUnknownFields()
	return decoder.Decode(out)
}
