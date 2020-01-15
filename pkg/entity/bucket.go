/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

import (
	"github.com/whiteblock/definition/config"

	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
)

/**
 * Note: this is not the best or most efficient of algorithms, however,
 * it should meet the criteria of good enough. Feel free to make improvements on it.
 */

type Bucket struct {
	Resource //The real size it occupies
	conf     *config.Bucket
	segments []Segment
	usage    Resource //The combined usage of the segments
	log      logrus.Ext1FieldLogger
}

func NewBucket(conf *config.Bucket, log logrus.Ext1FieldLogger) *Bucket {
	out := &Bucket{conf: conf}
	out.CPUs = conf.MinCPU
	out.Memory = conf.MinMemory
	out.Storage = conf.MinStorage
	out.segments = []Segment{}
	out.log = log
	return out
}

func (b Bucket) GetSegments() []Segment {
	return b.segments
}

func (b Bucket) FindByName(name string) int {
	for i, segment := range b.segments {
		if segment.Name == name {
			return i
		}
	}
	return -1
}

func (b Bucket) Clone() (out Bucket) {
	copier.Copy(&out, b)
	copier.Copy(&out.segments, b.segments)
	out.log = b.log
	return
}

func (b Bucket) hasSpace(segment Segment) bool {
	return b.NoPortConflicts(segment.GetPorts()...) &&
		(roundValueAndMax(0, b.usage.CPUs+segment.CPUs, b.conf.UnitCPU) <= b.conf.MaxCPU) &&
		(roundValueAndMax(0, b.usage.Memory+segment.Memory, b.conf.UnitCPU) <= b.conf.MaxMemory) &&
		(roundValueAndMax(0, b.usage.Storage+segment.Storage, b.conf.UnitCPU) <= b.conf.MaxStorage)
}

func (b Bucket) findSegment(segment Segment) int {
	return b.FindByName(segment.Name)
}

func max(a int64, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func roundValueAndMax(old int64, new int64, unit int64) int64 {
	if new%unit == 0 {
		return max(old, new)
	}
	return max(old, new+(unit-(new%unit)))
}

func (b *Bucket) update(segment Segment, positive bool) {
	if positive {
		b.usage.CPUs = b.usage.CPUs + segment.CPUs
		b.usage.Memory = b.usage.Memory + segment.Memory
		b.usage.Storage = b.usage.Storage + segment.Storage
		b.InsertPorts(segment.GetPorts()...)
		b.usage.InsertPorts(segment.GetPorts()...)
	} else {
		b.usage.CPUs = b.usage.CPUs - segment.CPUs
		b.usage.Memory = b.usage.Memory - segment.Memory
		b.usage.Storage = b.usage.Storage - segment.Storage
		b.RemovePorts(segment.GetPorts()...)
		b.usage.RemovePorts(segment.GetPorts()...)
		return
	}
	b.CPUs = roundValueAndMax(b.CPUs, b.usage.CPUs, b.conf.UnitCPU)
	b.Memory = roundValueAndMax(b.Memory, b.usage.Memory, b.conf.UnitMemory)
	b.Storage = roundValueAndMax(b.Storage, b.usage.Storage, b.conf.UnitStorage)
	b.log.WithFields(logrus.Fields{
		"storage":    b.Storage,
		"memory":     b.Memory,
		"cpu":        b.CPUs,
		"storageUse": b.usage.Storage,
		"memoryUse":  b.usage.Memory,
		"cpuUse":     b.usage.CPUs,
	}).Trace("updated bucket size")
}

func (b *Bucket) tryAdd(segment Segment) bool {
	if !b.hasSpace(segment) {
		return false
	}
	b.log.WithField("segment", segment).Trace("adding a segment")
	b.update(segment, true)
	b.segments = append(b.segments, segment)
	return true
}

func (b *Bucket) tryRemove(segment Segment) bool {
	loc := b.findSegment(segment)
	if loc == -1 {
		return false
	}
	b.log.WithField("segment", b.segments[loc]).Trace("removing a segment")
	b.update(b.segments[loc], false)
	if loc == len(b.segments)-1 {
		b.segments = append(b.segments[:loc])
	} else {
		b.segments = append(b.segments[0:loc], b.segments[loc:]...)
	}
	return true
}
