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
  - name: lighthouse
    image: sigp/lighthouse
    shared-volumes:
      - source-path: /var/log/lighthouse
        name: lighthouse-logs
    resources:
      cpus: 5
      memory: 8 GB
      storage: 100 GiB
task-runners:
  - name: run-lighthouse
    script:
      inline: lighthouse bn testnet -f quick 4 1575650550   
tests:
  - name: exercise-lighthouse
    description: run a lighthouse testnet and validate some blocks
    system:
      - name: lighthouse-partition-a
        type: lighthouse
        count: 2
        resources:
          networks:
            - name: common-net
            - name: partition-a-network
      - name: lighthouse-partition-b
        type: lighthouse
        count: 2
        resources:
          networks:
            - name: common-net
            - name: partition-b-network
    phases:
      - name: baseline
        tasks:
          - type: run-lighthouse
            timeout: 2m `)

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
