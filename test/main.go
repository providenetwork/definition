/*
  Copyright 2019 Whiteblock Inc.
  This file is a part of the definition.

  definition is free software: you can redistribute it and/or modify
  it under the terms of the GNU General Public License as published by
  the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  definition is distributed in the hope that it will be useful,
  but dock ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"fmt"
	"github.com/Whiteblock/go-prettyjson"
	"github.com/whiteblock/definition"
	"os"
	"strconv"
)

func main() {

	startingPoint := []byte(`services:
  - name: prysm-beacon
    image: gcr.io/prysmaticlabs/prysm/beacon-chain:latest
    args:
      - --datadir=/data
      - --init-sync-no-verify
    resources:
      cpus: 7
      memory: 10 GB
      storage: 100 GiB
task-runners:
  - name: unnecessary-task
    script: 
      inline: sleep 600
tests:
  - name: simple-prysm-exercise
    description: run a prysm testnet and validate some blocks
    system:
      - name: beacon-node-testnet
        type: prysm-beacon
        count: 4
    phases:
      - name: basic
        tasks:
          - type: unnecessary-task
            timeout: infinite`)

	def, err := definition.SchemaYAML(startingPoint)
	if err != nil {
		panic(err)
	}

	tests, err := definition.GetTests(def)
	if err != nil {
		panic(err)
	}
	cmds := tests[0].Commands
	out, _ := prettyjson.Marshal(cmds)
	if len(os.Args) > 1 {
		index, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}
		out, _ = prettyjson.Marshal(cmds[index])
	}
	fmt.Println(string(out))

}
