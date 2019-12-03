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

package definition

import (
	"testing"

	"github.com/Whiteblock/go-prettyjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAllTheThings(t *testing.T) {

	startingPoint := []byte(`services:
  - name: geth
    image: go-ethereum
    args:
      - --dev
      - --logfile=/var/log/geth/geth.log
    shared-volumes:
      - source-path: /var/log/geth
        name: geth-logs
    resources:
      cpus: 2
      memory: 4 GB
      storage: 100 GiB

sidecars:
  - name: geth-tail
    sidecar-to:
      - geth
    script:
      inline: while [ -ne /var/log/geth/geth.log ]; do sleep 1; done && tail -f /var/log/geth/geth.log
    mounted-volumes:
      - destination-path: /var/log/geth
        volume-name: geth-logs
    resources:
      cpus: 1
      memory: 512 MB
      storage: 8 GiB

task-runners:
  - name: geth-transactions
    image: nodejs
    input-files:
      - source-path: ./transaction-generator
        destination-path: /opt/transaction-generator
    script:
      inline: |
        cd /opt/transaction-generator
        npm install
        npm start -- --endpoint=http://geth-0:8545 --transaction-count=100
tests:
  - name: exercise-geth
    description: run a geth testnet and execute some simple transactions
    system:
      - type: geth
        count: 5
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

	def, err := SchemaYAML(startingPoint)
	require.NoError(t, err)
	tests, err := GetTests(def)
	assert.NoError(t, err)
	assert.NotNil(t, tests)

	out, _ := prettyjson.Marshal(*def.GetSpec())
	t.Log(string(out))
	out, _ = prettyjson.Marshal(tests)
	t.Log(string(out))
}
