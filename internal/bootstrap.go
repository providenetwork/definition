/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package internal

import (
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/internal/converter"
	"github.com/whiteblock/definition/internal/distribute"
	"github.com/whiteblock/definition/internal/maker"
	"github.com/whiteblock/definition/internal/parser"
	"github.com/whiteblock/definition/internal/process"
	"github.com/whiteblock/definition/internal/search"

	"github.com/sirupsen/logrus"
)

func GetFunctionality(conf config.Config, logger logrus.Ext1FieldLogger) (process.Commands, distribute.Distributor) {
	//  Distribute
	dist := distribute.NewDistributor(
		conf.Bucket,
		distribute.NewBiomeCalculator(
			parser.NewResources(
				search.NewSchema(),
				converter.NewResource(
					conf.Defaults.Resources,
				),
			),
			distribute.NewSystemState(
				parser.NewResources(
					search.NewSchema(),
					converter.NewResource(
						conf.Defaults.Resources,
					),
				),
			),
			logger,
		),
		logger,
	)

	//  Commands
	cmds := process.NewCommands(
		process.NewTestCalculator(
			conf,
			process.NewSystem(
				maker.NewService(
					conf.Defaults,
					search.NewSchema(),
					logger,
				),
				logger,
			),
			process.NewResolve(
				maker.NewCommand(
					parser.NewService(
						conf.Defaults.Service,
						converter.NewResource(conf.Defaults.Resources),
					),
					parser.NewSidecar(
						conf.Defaults.Service,
						converter.NewResource(conf.Defaults.Resources),
					),
					parser.NewNetwork(),
				),
				process.NewDependency(
					maker.NewCommand(
						parser.NewService(
							conf.Defaults.Service,
							converter.NewResource(conf.Defaults.Resources),
						),
						parser.NewSidecar(
							conf.Defaults.Service,
							converter.NewResource(conf.Defaults.Resources),
						),
						parser.NewNetwork(),
					),
					maker.NewService(
						conf.Defaults,
						search.NewSchema(),
						logger,
					),
					logger,
				),
				logger,
			),
			logger,
		),
	)
	return cmds, dist
}
