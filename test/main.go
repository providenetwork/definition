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
)

func main() {

	startingPoint := []byte(`services:
  - name: geth
    image: ethereum/client-go
    args:
      - --dev
      - --logfile=/var/log/geth/geth.log
    shared-volumes:
      - source-path: /var/log/geth
        name: geth-logs
    resources:
      cpus: 2
      memory: 4 GB
      storage: 5 GiB
sidecars:
  - name: "yes"
    sidecar-to:
      - nginx
    script:
      inline: "yes"
    resources:
      cpus: 1
      memory: 512 MB
      storage: 5 GiB
task-runners:
  - name: geth-transactions
    image: ubuntu:latest
    script:
      inline: echo hello
tests:
  - name: exercise-geth
    description: run a geth testnet and execute some simple transactions
    system:
      - type: geth
        count: 60
        port-mappings: 
        - "8888:8888"
    phases:
      - name: baseline-tps
        tasks:
          - type: geth-transactions
            timeout: 5 m
      - name: tps-with-latency
        system:
          - type: geth
            name: geth
            resources:
              networks:
                - name: default
                  latency: 100 ms
        tasks:
          - type: geth-transactions
            timeout: 5 m`)

	def, err := definition.SchemaYAML(startingPoint)
	fmt.Println(err)
	tests, err := definition.GetTests(def)

	out, _ := prettyjson.Marshal(tests)
	fmt.Println(string(out))
	fmt.Println(err)

}
