/*
	Copyright 2019 whiteblock Inc.
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

package distribute

import (
	"github.com/whiteblock/definition/internal/config"
	"github.com/whiteblock/definition/internal/entity"
)

/**
 * Note: this is not the best or most efficient of algorithms, however,
 * it should meet the criteria of good enough. Feel free to make improvements on it.
 */

type bucket struct {
	entity.Resource //The real size it occupies
	conf            *config.Bucket
	segments        []entity.Segment
	usage           entity.Resource //The combined usage of the segments
}

func newBucket(conf *config.Bucket) bucket {
	out := bucket{conf: conf}
	out.CPUs = conf.MinCPU
	out.Memory = conf.MinMemory
	out.Storage = conf.MinStorage
	return out
}

func (b bucket) empty() bool {
	return len(segments) == 0
}

func (b bucket) hasSpace(segment entity.Segment) bool {
	return (b.usage.CPUs+segment.CPUs <= b.conf.MaxCpus) &&
		(b.usage.Memory+segment.Memory <= b.conf.MaxMemory) &&
		(b.usage.Storage+segment.Storage <= b.conf.MaxStorage)
}

func (b bucket) findSegment(segment entity.Segment) int {
	for i, seggy := range b.segments {
		if segment.Name == seggy.Name {
			return i
		}
	}
	return -1
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

func (b bucket) update(segment entity.Segment, positive bool) {
	if positive {
		b.usage.CPUs = b.usage.CPUs + segment.CPUs
		b.usage.Memory = b.usage.Memory + segment.Memory
		b.usage.Storage = b.usage.Storage + segment.Storage
	} else {
		b.usage.CPUs = b.usage.CPUs - segment.CPUs
		b.usage.Memory = b.usage.Memory - segment.Memory
		b.usage.Storage = b.usage.Storage - segment.Storage
		return
	}
	b.CPUs = max(b.CPUs, b.usage.CPUs+(b.CPUs%b.conf.UnitCPU))
	b.Memory = max(b.Memory, b.usage.Memory+(b.Memory%b.conf.UnitMemory))
	b.Storage = max(b.Storage, b.usage.Storage+(b.Storage%b.conf.UnitStorage != 0))
}

func (b bucket) toResource() entity.Resource {
	return entity.Resource{
		CPUs:    b.CPUs,
		Memory:  b.Memory,
		Storage: b.Storage}
}

func (b bucket) tryAdd(segment entity.Segment) bool {
	if !b.hasSpace(segment) {
		return false
	}
	b.update(segment, true)
	b.segments = append(b.segments, segment)
	return true
}

func (b bucket) tryRemove(segment entity.Segment) bool {
	loc := b.findSegment(segment)
	if loc == -1 {
		return false
	}
	b.update(segment, false)
	if loc == len(b.segments)-1 {
		b.segments = append(b.segments[:loc])
	} else {
		b.segments = append(b.segments[0:loc], b.segments[loc:])
	}
}
