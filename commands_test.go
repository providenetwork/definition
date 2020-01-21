/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package definition

import (
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"testing"

	"github.com/whiteblock/definition/command"
	"github.com/whiteblock/definition/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/whiteblock/go-prettyjson"
	"github.com/whiteblock/utility/utils"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

func getTestCommands(t *testing.T) Commands {
	v := viper.New()
	err := config.SetupViper(v)
	require.NoError(t, err)
	conf, err := config.New(v)
	require.NoError(t, err)
	return NewCommands(conf, utils.NewTestingLogger(t))
}

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
      - path: /var/pictures
        name: pictures
        scope: global
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
sidecars:
  - name: side
    sidecar-to:
      - Quorum1
      - Quorum2
      - Quorum3
      - Quorum4
      - Quorum5
      - Quorum6
      - Quorum7
    image: side
    script:
      inline:
        "yes"
    resources:
      cpus: 2
      memory: 4 GB
      storage: 30 GiB
    environment:
      NETWORK_NAME: quorem
tests:
  - name: testnet
    timeout: infinite
    system:
      - type: Quorum1
        port-mappings:
          - "30303:30303"
          - "8545:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum2
        port-mappings:
          - "30304:30303"
          - "8546:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum3
        port-mappings:
          - "30305:30303"
          - "8547:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum4
        port-mappings:
          - "30306:30303"
          - "8548:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum5
        port-mappings:
          - "30307:30303"
          - "8549:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum6
        port-mappings:
          - "30308:30303"
          - "8550:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: Quorum7
        port-mappings:
          - "30308:30303"
          - "8551:8545"
        resources: 
            networks:
              - name: quorum_network
      - type: ethstats
        port-mappings:
          - "80:3000"
        resources: 
            networks:
              - name: quorum_network`)

	def, err := SchemaYAML(startingPoint)
	require.NoError(t, err)
	cmds := getTestCommands(t)
	tests, err := cmds.GetTests(def, Meta{})
	assert.NoError(t, err)
	assert.NotNil(t, tests)

	for _, test := range tests {
		assertSanity(t, test)
		assertCorrectIPs(t, test)
		assertNamed(t, test)
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

func assertNamed(t *testing.T, test command.Test) {
	for _, outer := range test.Commands {
		for _, inner := range outer {
			switch inner.Order.Type {
			case command.Createcontainer:

				var cont command.Container
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				require.NotNil(t, cont.Environment)
				assert.Equal(t, cont.Name, cont.Environment["NAME"])
			}
		}
	}
}

func assertSanity(t *testing.T, test command.Test) {
	networks := map[string]bool{}
	volumes := map[string]bool{}
	attachedNetwork := map[string]map[string]bool{}

	for _, outer := range test.Commands {
		for _, inner := range outer {
			switch inner.Order.Type {
			case command.Createcontainer:

				var cont command.Container
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				attachedNetwork[cont.Name] = map[string]bool{}
				for _, mount := range cont.Volumes {
					_, exists := volumes[mount.Name]
					require.True(t, exists)
				}

			case command.Createnetwork:
				var network command.Network
				err := inner.ParseOrderPayloadInto(&network)
				require.NoError(t, err)

				out, _ := prettyjson.Marshal(network)
				t.Log(string(out))
				_, exists := networks[network.Name]
				assert.False(t, exists, "duplicate network found "+network.Name)
				networks[network.Name] = true

			case command.Attachnetwork:
				var cn command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cn)
				require.NoError(t, err)

				_, exists := attachedNetwork[cn.Container]
				require.True(t, exists, "container has not been created yet: "+cn.Container)

				_, exists = networks[cn.Network]
				require.True(t, exists, "network has not been created yet: "+cn.Network)

				_, exists = attachedNetwork[cn.Container][cn.Network]
				require.False(t, exists, "network has already been attached: "+cn.Network)

				attachedNetwork[cn.Container][cn.Network] = true

			case command.Detachnetwork:
				var cn command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cn)
				require.NoError(t, err)

				_, exists := attachedNetwork[cn.Container]
				require.True(t, exists, "container has not been created yet: "+cn.Container)

				_, exists = networks[cn.Network]
				require.True(t, exists, "network has not been created yet: "+cn.Network)

				_, exists = attachedNetwork[cn.Container][cn.Network]
				require.True(t, exists, "network is not attached: "+cn.Network)

				delete(attachedNetwork[cn.Container], cn.Network)

			case command.Createvolume:
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
	sidecarIPs := map[string]bool{}
	scCount := 0
	targets := map[string]bool{}
	for _, outer := range test.Commands {
		for _, inner := range outer {
			targets[inner.Target.IP] = true
			switch inner.Order.Type {
			case command.Createcontainer:
				var cont command.Container
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				require.NotNil(t, cont.Environment)
				env = cont.Environment
				t.Log(cont.Environment)
				require.True(t, len(cont.Ports) > 0 || strings.Contains(cont.Name, "side"))

				if strings.Contains(cont.Name, "side") {
					require.Equal(t, "quorem", cont.Environment["NETWORK_NAME"])
					scCount++
					scKey := inner.Target.IP + cont.Environment["SERVICE"]
					_, exists := sidecarIPs[scKey]
					require.False(t, exists)
					sidecarIPs[scKey] = true

					_, exists = ips[scKey]
					require.True(t, exists)
				} else {
					_, exists := ips[inner.Target.IP+cont.IP] //sidecar net ips are localized
					require.False(t, exists, fmt.Sprintf("duplicate entry of %s-%s", inner.Target.IP, cont.IP))
					ips[inner.Target.IP+cont.IP] = true
					require.NotNil(t, net.ParseIP(cont.IP)) //ensure IP is valid
				}
			case command.Attachnetwork:
				var cont command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				name := cont.Container + "_QUORUM_NETWORK"
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
	require.Equal(t, 7, scCount)
	require.Equal(t, 7, len(sidecarIPs))
	require.Equal(t, 2, len(targets), "there should only be two instances")
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

	cmds := getTestCommands(t)
	tests, err := cmds.GetTests(defJSON, Meta{})
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

var test3 = []byte(`services:
  - name: nginx
    image: nginx:alpine
    resources:
      cpus: 1
      memory: 500 MB
      storage: 1 GiB
    volumes:
      - name: test
        path: /var/hello 
        scope: global
task-runners:
  - name: wait-5-minutes
    script:
      inline: sleep 600
tests:
  - name: serve-static-files
    description: Serve the default static content of the nginx image.
    system:
      - name: nginx
        type: nginx
        count: 2
        port-mappings:
          - "80:80"
    phases:
      - name: testnet
        duration: 100s
        system:
          - type: nginx
            name: nginx
            resources:
              networks:
                - name: nginx
        tasks:
          - type: wait-5-minutes
            timeout: 100s
`)

func TestPhaseChangesSane(t *testing.T) {
	def, err := SchemaYAML(test3)
	assert.NoError(t, err)

	cmds := getTestCommands(t)
	tests, err := cmds.GetTests(def, Meta{})
	require.NoError(t, err)
	require.Len(t, tests, 1)
	test := tests[0]
	netCount := 0
	cntrCount := 0
	volCount := 0
	pauseCount := 0
	resumeCount := 0
	for _, outer := range test.Commands {
		for _, inner := range outer {
			switch inner.Order.Type {
			case command.Resumeexecution:
				resumeCount++
			case command.Pauseexecution:
				pauseCount++
			case command.Createcontainer:
				var cont command.Container
				err := inner.ParseOrderPayloadInto(&cont)
				if !strings.Contains(cont.Name, "service") {
					continue
				}
				require.NoError(t, err)
				assert.Len(t, cont.Volumes, 1, "the services should have a volume attached here")
				cntrCount++

			case command.Createvolume:
				var vol command.Volume
				err := inner.ParseOrderPayloadInto(&vol)
				require.NoError(t, err)
				assert.Len(t, vol.Hosts, 2, "there should be two hosts")
				assert.True(t, vol.Global, "the volume should be flagged as global")
				volCount++
			case command.Attachnetwork:
				var cont command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				if strings.HasPrefix(cont.Network, "net") {
					netCount++
				}

			case command.Detachnetwork:
				var cont command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				if strings.HasPrefix(cont.Network, "net") {
					netCount--
				}
			}
		}
	}
	assert.Equal(t, 1, resumeCount, "there should be two resume commands in this case")
	assert.Equal(t, 1, pauseCount, "there should only be one pause command in this case")
	assert.Equal(t, 1, volCount, "singleton volume should be just a single command")
	assert.Equal(t, 2, netCount, "The two containers should end with just one network attached")
	assert.Equal(t, 2, cntrCount, "There should only be two containers created")
}

var composeTest = []byte(`services:
- name: web
  environment:
    API_URL: /api
  image: web3labs/epirus-free-web:latest
- name: nginx
  image: nginx:latest
  input-files:
  - source-path: ./nginx.conf
    destination-path: /etc/nginx/nginx.conf
  - source-path: ./5xx.html
    destination-path: /www/error_pages/5xx.html
- name: api
  environment:
    ENABLE_PRIVATE_QUORUM: enabled
    MONGO_CLIENT_URI: mongodb://mongodb:27017
    MONGO_DB_NAME: epirus
    NODE_ENDPOINT: ${NODE_ENDPOINT}
  image: web3labs/epirus-free-api:latest
- name: mongodb
  image: mongo:latest
tests:
- name: compose
  system:
  - type: mongodb
    resources:
      networks:
      - name: epirus
    port-mappings:
    - 27017:27017
  phases:
  - name: phase1
    system:
    - type: api
      resources:
        networks:
        - name: epirus
  - name: phase2
    system:
    - type: web
      resources:
        networks:
        - name: epirus
  - name: phase3
    system:
    - type: nginx
      resources:
        networks:
        - name: epirus
      port-mappings:
      - 80:80
`)

func TestComposeLikeSpec(t *testing.T) {
	def, err := SchemaYAML(composeTest)
	require.NoError(t, err)

	cmds := getTestCommands(t)
	dists, err := cmds.GetDist(def)
	require.NoError(t, err)
	require.NotNil(t, dists)
	require.Len(t, dists, 1)
	require.NotNil(t, dists[0])

	for i := range dists {
		require.True(t, len(*dists[i]) > 0)
		for j := range *dists[i] {
			require.True(t, len((*dists[i])[j]) > 0, "distributions should not be empty")

		}
	}

	tests, err := cmds.GetTests(def, Meta{})
	require.NoError(t, err)
	require.Len(t, tests, 1)
	test := tests[0]
	netCount := 0
	cntrCount := 0
	volCount := 0
	for _, outer := range test.Commands {
		for _, inner := range outer {
			switch inner.Order.Type {
			case command.Createcontainer:
				var cont command.Container
				err := inner.ParseOrderPayloadInto(&cont)
				if !strings.Contains(cont.Name, "service") {
					continue
				}
				require.NoError(t, err)
				cntrCount++

			case command.Createvolume:
				var vol command.Volume
				err := inner.ParseOrderPayloadInto(&vol)
				require.NoError(t, err)
				volCount++
			case command.Attachnetwork:
				var cont command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				if strings.HasPrefix(cont.Network, "net") {
					netCount++
				}

			case command.Detachnetwork:
				var cont command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				if strings.HasPrefix(cont.Network, "net") {
					netCount--
				}
			}
		}
	}

	assert.Equal(t, 2, volCount, "there should be 2 volumes")
	assert.Equal(t, 4, netCount, "The 4 containers should end with just one network attached")
	assert.Equal(t, 4, cntrCount, "There should only be 4 containers created")
}

var sidecarTest = []byte(`services: 
  - name: geth1
    image: "ethereum/client-go:alltools-latest" 
    script:
      inline: geth
  - name: geth2 
    image: "ethereum/client-go:alltools-latest"
    script: 
      inline: geth
sidecars:
  - name: bash
    sidecar-to:
      - geth1
      - geth2
    image: ubuntu
    script:
      inline:
        tail -f /var/log/syslog
task-runners:
  - name: geth-staticpeers-helper
    image: "gcr.io/whiteblock/helpers/besu/staticpeers:master"
    script: 
      inline: help me
tests: 
  - name: geth_network_2_nodes 
    phases:
      - name: create
        tasks: 
        - type: geth-staticpeers-helper
      - name: start
        description: start the remaining node(s)
        system: 
        - type: geth1 
          resources: 
            networks: 
              - name: common-network
        - type: geth2 
          resources: 
            networks: 
              - name: common-network
`)

func TestSidecarSpec(t *testing.T) {
	def, err := SchemaYAML(sidecarTest)
	require.NoError(t, err)

	cmds := getTestCommands(t)
	dists, err := cmds.GetDist(def)
	require.NoError(t, err)
	require.NotNil(t, dists)
	require.Len(t, dists, 1)
	require.NotNil(t, dists[0])

	for i := range dists {
		require.True(t, len(*dists[i]) > 0)
		for j := range *dists[i] {
			require.True(t, len((*dists[i])[j]) > 0 || j == 0, fmt.Sprintf("distributions should not be empty %d", j))

		}
	}

	tests, err := cmds.GetTests(def, Meta{})
	require.NoError(t, err)
	require.Len(t, tests, 1)
	test := tests[0]
	netCount := 0
	cntrCount := 0
	scCount := 0
	volCount := 0
	ips := map[string]bool{}
	parentIPs := map[string]bool{}
	for _, outer := range test.Commands {
		for _, inner := range outer {
			switch inner.Order.Type {
			case command.Createcontainer:
				var cont command.Container
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				cntrCount++
				if strings.Contains(cont.Name, "bash") {
					scCount++
					_, exists := ips[cont.Environment["SERVICE"]]
					require.False(t, exists)
					ips[cont.Environment["SERVICE"]] = true
					_, exists = parentIPs[cont.Environment["SERVICE"]]
					require.True(t, exists)
				} else {
					parentIPs[cont.IP] = true
				}

			case command.Createvolume:
				var vol command.Volume
				err := inner.ParseOrderPayloadInto(&vol)
				require.NoError(t, err)
				volCount++
			case command.Attachnetwork:
				var cont command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				if strings.HasPrefix(cont.Network, "net") {
					netCount++
				}

			case command.Detachnetwork:
				var cont command.ContainerNetwork
				err := inner.ParseOrderPayloadInto(&cont)
				require.NoError(t, err)
				if strings.HasPrefix(cont.Network, "net") {
					netCount--
				}
			}
		}
	}

	assert.Equal(t, 0, volCount, "there should be 0 volumes")
	assert.Equal(t, 2, netCount, "The 2 containers should end with just one network attached")
	assert.Equal(t, 5, cntrCount, "There should only be 5 containers created")
	assert.Equal(t, 2, scCount, "There should be 2 sidecars created")
}
