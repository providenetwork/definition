/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

import (
	"github.com/whiteblock/definition/config"
	"github.com/whiteblock/definition/schema"

	"github.com/sirupsen/logrus"
)

type StatePack struct {
	State
	Buckets   ResourceBuckets
	PrevTasks []Segment
	Spec      schema.RootSchema
}

func NewStatePack(spec schema.RootSchema, conf config.Bucket, logger logrus.Ext1FieldLogger) *StatePack {
	out := &StatePack{
		Buckets:   NewResourceBuckets(conf, logger),
		PrevTasks: nil,
		Spec:      spec,
	}
	out.SystemState = map[string]schema.SystemComponent{}
	return out
}
