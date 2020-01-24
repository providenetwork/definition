/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package parser

import "github.com/whiteblock/definition/schema"

func ExtractGlobalVolumes(root schema.RootSchema) []string {
	volumes := map[string]bool{}

	for i := range root.Services {
		for j := range root.Services[i].Volumes {
			if !root.Services[i].Volumes[j].Local() {
				volumes[root.Services[i].Volumes[j].Name] = false
			}
		}
	}

	for i := range root.Sidecars {
		for j := range root.Sidecars[i].Volumes {
			if !root.Sidecars[i].Volumes[j].Local() {
				volumes[root.Sidecars[i].Volumes[j].Name] = false
			}
		}
	}

	for i := range root.TaskRunners {
		for j := range root.TaskRunners[i].Volumes {
			if !root.TaskRunners[i].Volumes[j].Local() {
				volumes[root.TaskRunners[i].Volumes[j].Name] = false
			}
		}
	}
	out := []string{}
	for vol := range volumes {
		out = append(out, vol)
	}
	return out
}
