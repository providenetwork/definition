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
	"encoding/json"
	"net"
	"strings"
	"testing"

	"github.com/whiteblock/definition/command"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/whiteblock/go-prettyjson"

	"gopkg.in/yaml.v2"
)

func TestAllTheThings(t *testing.T) {

	startingPoint := []byte(`services:
  - name: Quorum1
    image: quorumengineering/quorum:2.2.5
    volumes:
      - path: /output.log
        name: eea-logs
    resources:
      cpus: 4
      memory: 4 GB
      storage: 5 GiB
    input-files:
      - source-path: genesis.json
        destination-path: /data/genesis.json
      - source-path: permissioned-nodes.json
        destination-path: /data/permissioned-nodes.json
      - source-path: permissioned-nodes.json
        destination-path: /data/static-nodes.json
      - source-path: key1
        destination-path: /data/keystore/key1
      - source-path: nodekey1
        destination-path: /data/nodekey
      - source-path: passwords.txt
        destination-path: /data/passwords.txt
    script:
      inline: geth --datadir /data init data/genesis.json && geth --datadir /data --unlock 0 --password /data/passwords.txt --ethstats Node1:eea_testnet_secret@eea:80 --syncmode full --mine --minerthreads 1 --rpc --rpcaddr 0.0.0.0 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum > /output.log 2>&1                      
  - name: Quorum2
    image: quorumengineering/quorum:2.2.5
    volumes:
      - path: /output.log
        name: eea-logs
    resources:
      cpus: 4
      memory: 4 GB
      storage: 5 GiB
    input-files:
      - source-path: genesis.json
        destination-path: /data/genesis.json
      - source-path: permissioned-nodes.json
        destination-path: /data/permissioned-nodes.json
      - source-path: permissioned-nodes.json
        destination-path: /data/static-nodes.json
      - source-path: key2
        destination-path: /data/keystore/key2
      - source-path: nodekey2
        destination-path: /data/nodekey
      - source-path: passwords.txt
        destination-path: /data/passwords.txt
    script:
      inline: geth --datadir /data init data/genesis.json && geth --datadir /data --unlock 0 --password /data/passwords.txt --ethstats Node1:eea_testnet_secret@eea:80 --syncmode full --mine --minerthreads 1 --rpc --rpcaddr 0.0.0.0 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum > /output.log 2>&1                        
  - name: Quorum3
    image: quorumengineering/quorum:2.2.5
    volumes:
      - path: /output.log
        name: eea-logs
    resources:
      cpus: 4
      memory: 4 GB
      storage: 5 GiB
    input-files:
      - source-path: genesis.json
        destination-path: /data/genesis.json
      - source-path: permissioned-nodes.json
        destination-path: /data/permissioned-nodes.json
      - source-path: permissioned-nodes.json
        destination-path: /data/static-nodes.json
      - source-path: key3
        destination-path: /data/keystore/key3
      - source-path: nodekey3
        destination-path: /data/nodekey
      - source-path: passwords.txt
        destination-path: /data/passwords.txt
    script:
      inline: geth --datadir /data init data/genesis.json && geth --datadir /data --unlock 0 --password /data/passwords.txt --ethstats Node1:eea_testnet_secret@eea:80 --syncmode full --mine --minerthreads 1 --rpc --rpcaddr 0.0.0.0 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum > /output.log 2>&1   
  - name: Quorum4
    image: quorumengineering/quorum:2.2.5
    volumes:
      - path: /output.log
        name: eea-logs
    resources:
      cpus: 4
      memory: 4 GB
      storage: 5 GiB
    input-files:
      - source-path: genesis.json
        destination-path: /data/genesis.json
      - source-path: permissioned-nodes.json
        destination-path: /data/permissioned-nodes.json
      - source-path: permissioned-nodes.json
        destination-path: /data/static-nodes.json
      - source-path: key4
        destination-path: /data/keystore/key4
      - source-path: nodekey4
        destination-path: /data/nodekey
      - source-path: passwords.txt
        destination-path: /data/passwords.txt
    script:
      inline: geth --datadir /data init data/genesis.json && geth --datadir /data --unlock 0 --password /data/passwords.txt --ethstats Node1:eea_testnet_secret@eea:80 --syncmode full --mine --minerthreads 1 --rpc --rpcaddr 0.0.0.0 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum > /output.log 2>&1   
  - name: Quorum5
    image: quorumengineering/quorum:2.2.5
    volumes:
      - path: /output.log
        name: eea-logs
    resources:
      cpus: 4
      memory: 4 GB
      storage: 5 GiB
    input-files:
      - source-path: genesis.json
        destination-path: /data/genesis.json
      - source-path: permissioned-nodes.json
        destination-path: /data/permissioned-nodes.json
      - source-path: permissioned-nodes.json
        destination-path: /data/static-nodes.json
      - source-path: key5
        destination-path: /data/keystore/key5
      - source-path: nodekey5
        destination-path: /data/nodekey
      - source-path: passwords.txt
        destination-path: /data/passwords.txt
    script:
      inline: geth --datadir /data init data/genesis.json && geth --datadir /data --unlock 0 --password /data/passwords.txt --ethstats Node1:eea_testnet_secret@eea:80 --syncmode full --mine --minerthreads 1 --rpc --rpcaddr 0.0.0.0 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum > /output.log 2>&1
  - name: Quorum6
    image: quorumengineering/quorum:2.2.5
    volumes:
      - path: /output.log
        name: eea-logs
    resources:
      cpus: 4
      memory: 4 GB
      storage: 5 GiB
    input-files:
      - source-path: genesis.json
        destination-path: /data/genesis.json
      - source-path: permissioned-nodes.json
        destination-path: /data/permissioned-nodes.json
      - source-path: permissioned-nodes.json
        destination-path: /data/static-nodes.json
      - source-path: key6
        destination-path: /data/keystore/key6
      - source-path: nodekey6
        destination-path: /data/nodekey
      - source-path: passwords.txt
        destination-path: /data/passwords.txt
    script:
      inline: geth --datadir /data init data/genesis.json && geth --datadir /data --unlock 0 --password /data/passwords.txt --ethstats Node1:eea_testnet_secret@eea:80 --syncmode full --mine --minerthreads 1 --rpc --rpcaddr 0.0.0.0 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum > /output.log 2>&1
  - name: Quorum7
    image: quorumengineering/quorum:2.2.5
    volumes:
      - path: /output.log
        name: eea-logs
      - path: /etc/apt.d
        name: apt
    resources:
      cpus: 4
      memory: 4 GB
      storage: 5 GiB
    input-files:
      - source-path: genesis.json
        destination-path: /data/genesis.json
      - source-path: permissioned-nodes.json
        destination-path: /data/permissioned-nodes.json
      - source-path: permissioned-nodes.json
        destination-path: /data/static-nodes.json
      - source-path: key7
        destination-path: /data/keystore/key7
      - source-path: nodekey7
        destination-path: /data/nodekey
      - source-path: passwords.txt
        destination-path: /data/passwords.txt
    script:
      inline: geth --datadir /data init data/genesis.json && geth --datadir /data --unlock 0 --password /data/passwords.txt --ethstats Node1:eea_testnet_secret@eea:80 --syncmode full --mine --minerthreads 1 --rpc --rpcaddr 0.0.0.0 --rpcapi admin,db,eth,debug,miner,net,shh,txpool,personal,web3,quorum > /output.log 2>&1
  - name: ethstats
    image: gcr.io/whiteblock/ethstats:master
    environment:
      HOST: "0.0.0.0"
    input-files:
      - source-path: ws_secret.json
        destination-path: /eth-netstats/ws_secret.json     
task-runners:
  - name: testnet-expiration
    script:
      inline: sleep 7200
    volumes:
      - path: /etc/apt.d
        name: apt
tests:
  - name: testnet
    timeout: infinite
    description: run an EEA testnet and execute some simple transactions
    system:
      - type: Quorum1
        count: 1
        port-mappings:
          - "30303:30303"
          - "8545:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum2
        count: 1
        port-mappings:
          - "30304:30303"
          - "8546:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum3
        count: 1
        port-mappings:
          - "30305:30303"
          - "8547:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum4
        count: 1
        port-mappings:
          - "30306:30303"
          - "8548:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum5
        count: 1
        port-mappings:
          - "30307:30303"
          - "8549:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum6
        count: 1
        port-mappings:
          - "30308:30303"
          - "8550:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum7
        count: 1
        port-mappings:
          - "30308:30303"
          - "8551:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: ethstats
        count: 1
        port-mappings:
          - "80:3000"
        resources: 
            networks:
              - name: quorum_network`)

	def, err := SchemaYAML(startingPoint)
	require.NoError(t, err)
	tests, err := GetTests(def, Meta{})
	assert.NoError(t, err)
	assert.NotNil(t, tests)

	for _, test := range tests {
		assertNoDupNetworks(t, test)
		assertCorrectIPs(t, test)
	}
}

func assertNoDataLoss(t *testing.T, def Definition) {
	data, err := json.Marshal(def)
	require.NoError(t, err)
	require.NotNil(t, data)
	var defJSON Definition
	err = json.Unmarshal(data, &defJSON)
	require.NoError(t, err)

	data, err = yaml.Marshal(def)
	require.NoError(t, err)
	require.NotNil(t, data)
	var defYAML Definition
	err = yaml.UnmarshalStrict(data, &defJSON)
	require.NoError(t, err)

	assert.Equal(t, def, defJSON)
	assert.Equal(t, def, defYAML)

	data, err = json.Marshal(def.Spec)
	require.NoError(t, err)
	require.NotNil(t, data)
	defJSON2, err := SchemaJSON(data)
	require.NoError(t, err)

	data, err = yaml.Marshal(def.Spec)
	require.NoError(t, err)
	require.NotNil(t, data)
	defYAML2, err := SchemaYAML(data)
	require.NoError(t, err)

	assert.Equal(t, def.Spec, defJSON2.Spec)
	assert.Equal(t, def.Spec, defYAML2.Spec)
}

func assertNoDupNetworks(t *testing.T, test command.Test) {
	networks := map[string]bool{}
	volumes := map[string]bool{}
	for _, outer := range test.Commands {
		for _, inner := range outer {
			if inner.Order.Type == command.Createnetwork {
				var network command.Network
				err := inner.ParseOrderPayloadInto(&network)

				out, _ := prettyjson.Marshal(network)
				t.Log(string(out))
				require.NoError(t, err)
				_, exists := networks[network.Name]
				assert.False(t, exists, "duplicate network found "+network.Name)
				networks[network.Name] = true
			} else if inner.Order.Type == command.Createvolume {
				var vol command.Volume
				err := inner.ParseOrderPayloadInto(&vol)
				require.NoError(t, err)
				_, exists := volumes[vol.Name]
				assert.False(t, exists, "duplicate volume found "+vol.Name)
				volumes[vol.Name] = true
			}
		}
	}
	assert.True(t, len(volumes) > 0)
	assert.True(t, len(networks) > 0)

}

/*
  This tests that all assigned IP addresses are unique and the ENV vars are set correctly
*/
func assertCorrectIPs(t *testing.T, test command.Test) {
	var env map[string]string
	ips := map[string]bool{}
	for _, outer := range test.Commands {
		for _, inner := range outer {
			switch inner.Order.Type {
			case command.Createcontainer:
				var cont command.Container
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				env = cont.Environment
				t.Log(env)
				require.True(t, len(cont.Ports) > 0)
				_, exists := ips[inner.Target.IP+cont.IP] //sidecar net ips are localized
				require.False(t, exists)
				ips[inner.Target.IP+cont.IP] = true

				require.NotNil(t, net.ParseIP(cont.IP)) //ensure IP is valid

			case command.Attachnetwork:
				var cont command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				name := cont.ContainerName + "_QUORUM_NETWORK"
				name = strings.Replace(name, "-", "_", -1)
				name = strings.ToUpper(name)
				t.Log(name)
				require.Equal(t, env[name], cont.IP)

				_, exists := ips[cont.IP]
				require.False(t, exists)
				ips[cont.IP] = true

				require.NotNil(t, net.ParseIP(cont.IP)) //ensure IP is valid
			}
		}
	}
}

var test2 = []byte(`{
  "services": [
    {
      "image": "nginx:alpine",
      "input-files": [
        {
          "destination-path": "/usr/share/nginx/html/index.html",
          "path": "/f/k1"
        },
        {
          "destination-path": "/usr/share/nginx/html/Dockerfile",
          "path": "/f/k0"
        }
      ],
      "name": "nginx",
      "resources": {
        "cpus": 1,
        "memory": "500 MB",
        "storage": "1 GiB"
      },
      "script": {}
    }
  ],
  "task-runners": [
    {
      "name": "wait-5-minutes",
      "resources": {},
      "script": {
        "inline": "sleep 600"
      }
    }
  ],
  "tests": [
    {
      "description": "Serve the default static content of the nginx image.",
      "name": "serve-static-files",
      "phases": [
        {
          "name": "testnet",
          "system": [
            {
              "name": "nginx",
              "resources": {
                "networks": [
                  {
                    "name": "nginx"
                  }
                ]
              },
              "type": "nginx"
            }
          ],
          "tasks": [
            {
              "timeout": 600000000000,
              "type": "wait-5-minutes"
            }
          ],
          "timeout": 0
        }
      ],
      "system": [
        {
          "count": 1,
          "name": "nginx",
          "portMappings": [
            "80:80"
          ],
          "resources": {},
          "type": "nginx"
        }
      ],
      "timeout": 0
    }
  ]
}`)

func TestJSONSane(t *testing.T) {
	_, err := SchemaYAML(test2)
	assert.Error(t, err)
	defJSON, err := SchemaJSON(test2)
	require.NoError(t, err)
	data, err := json.Marshal(defJSON.Spec)
	require.NoError(t, err)
	require.NotNil(t, data)
	defJSON2, err := SchemaJSON(data)
	require.NoError(t, err)

	assert.Equal(t, defJSON.Spec, defJSON2.Spec)
	assert.Len(t, defJSON.Spec.Tests[0].System[0].PortMappings, 1)

	tests, err := GetTests(defJSON, Meta{})
	require.NoError(t, err)
	ensurePortMapping(t, tests[0])
}

func ensurePortMapping(t *testing.T, test command.Test) {
	for _, outer := range test.Commands {
		for _, inner := range outer {
			switch inner.Order.Type {
			case command.Createcontainer:
				var cont command.Container
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				t.Log(cont.Ports)
				if strings.Contains(cont.Name, "nginx") {
					require.True(t, len(cont.Ports) > 0)
				}

			}
		}
	}
}
