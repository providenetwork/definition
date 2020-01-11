/*
	Copyright 2019 Whiteblock Inc.
	This file is a part of the Definition.

	Definition is free software: you can redistribute it and/or modify
	it under the terms of the GNU General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	Definition is distributed in the hope that it will be useful,
	but dock ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU General Public License for more details.

	You should have received a copy of the GNU General Public License
	along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package parser

import "github.com/whiteblock/definition/schema"

func ExtractAllVolumes(root schema.RootSchema) []string {
	volumes := map[string]bool{}

	for i := range root.Services {
		for j := range root.Services[i].SharedVolumes {
			volumes[root.Services[i].SharedVolumes[j].Name] = false
		}
	}

	for i := range root.Sidecars {
		for j := range root.Sidecars[i].MountedVolumes {
			volumes[root.Sidecars[i].MountedVolumes[j].VolumeName] = false
		}
	}

	for i := range root.TaskRunners {
		for j := range root.TaskRunners[i].SharedVolumes {
			volumes[root.TaskRunners[i].SharedVolumes[j].Name] = false
		}
	}
	out := []string{}
	for vol := range volumes {
		out = append(out, vol)
	}
	return out
}
