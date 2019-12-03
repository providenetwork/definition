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
package internal

import (
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal/converter"
	"github.com/whiteblock/definition/internal/distribute"
	"github.com/whiteblock/definition/internal/maker"
	"github.com/whiteblock/definition/internal/merger"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/internal/process"
	"github.com/whiteblock/definition/internal/search"
)

func GetFunctionality(conf config.Config) (process.Commands, distribute.Distributor, error) {
	logger, err := conf.Logger.GetLogger()
	if err != nil {
		return nil, nil, err
	}
	//Distribute
	dist := distribute.NewDistributor(
		conf.Bucket,
		distribute.NewBiomeCalculator(
			parser.NewResources(
				parser.NewNames(),
				search.NewSchema(),
				converter.NewResource(
					conf.Defaults.Resources,
				),
			),
			distribute.NewSystemState(
				parser.NewResources(
					parser.NewNames(),
					search.NewSchema(),
					converter.NewResource(
						conf.Defaults.Resources,
					),
				),
				parser.NewNames(),
				merger.NewSystem(),
			),
			logger,
		),
	)

	//Commands
	cmds := process.NewCommands(
		process.NewTestCalculator(
			process.NewSystem(
				parser.NewNames(),
				maker.NewService(
					parser.NewNames(),
					search.NewSchema(),
					converter.NewService(),
				),
				merger.NewSystem(),
			),
			process.NewResolve(
				maker.NewCommand(
					parser.NewService(
						conf.Defaults.Service,
						parser.NewNames(),
					),
					parser.NewSidecar(
						conf.Defaults.Service,
						parser.NewNames(),
					),
					parser.NewNetwork(),
					parser.NewNames(),
				),
				process.NewDependency(
					maker.NewCommand(
						parser.NewService(
							conf.Defaults.Service,
							parser.NewNames(),
						),
						parser.NewSidecar(
							conf.Defaults.Service,
							parser.NewNames(),
						),
						parser.NewNetwork(),
						parser.NewNames(),
					),
					maker.NewService(
						parser.NewNames(),
						search.NewSchema(),
						converter.NewService(),
					),
					logger,
				),
				logger,
			),
			logger,
		),
	)
	return cmds, dist, nil
}
