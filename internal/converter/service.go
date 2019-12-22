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

package converter

import (
	"strings"

	"github.com/whiteblock/definition/schema"

	"github.com/jinzhu/copier"
)

type Service interface {
	FromTaskRunner(taskRunner schema.TaskRunner) schema.Service
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (serv *service) FromTaskRunner(taskRunner schema.TaskRunner) schema.Service {
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
	for _, output := range taskRunner.OutputFiles {
		out.SharedVolumes = append(out.SharedVolumes, schema.SharedVolume{
			SourcePath: output,
			Name:       strings.Replace(output, "/", "-", -1),
		})
	}
	return out
}
