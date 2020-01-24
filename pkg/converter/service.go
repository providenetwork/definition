/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package converter

import (
	"github.com/whiteblock/definition/schema"

	"github.com/jinzhu/copier"
)

func FromTaskRunner(taskRunner schema.TaskRunner) schema.Service {
	out := schema.Service{
		Name:        taskRunner.Name,
		Description: taskRunner.Description,
		Resources:   taskRunner.Resources,
		Image:       taskRunner.Image,
		Script:      taskRunner.Script,
	}
	copier.Copy(&out.Args, taskRunner.Args)
	copier.Copy(&out.Environment, taskRunner.Environment)
	copier.Copy(&out.InputFiles, taskRunner.InputFiles)
	copier.Copy(&out.Volumes, taskRunner.Volumes)
	return out
}
