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

package entity

import (
	"github.com/whiteblock/definition/config"

	"github.com/jinzhu/copier"
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
}

func NewBucket(conf *config.Bucket) *Bucket {
	out := &Bucket{conf: conf}
	out.CPUs = conf.MinCPU
	out.Memory = conf.MinMemory
	out.Storage = conf.MinStorage
	out.segments = []Segment{}
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
	copier.Copy(&out, &b)
	return
}

func (b Bucket) hasSpace(segment Segment) bool {
	return (b.usage.CPUs+segment.CPUs <= b.conf.MaxCPU) &&
		(b.usage.Memory+segment.Memory <= b.conf.MaxMemory) &&
		(b.usage.Storage+segment.Storage <= b.conf.MaxStorage)
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
	} else {
		b.usage.CPUs = b.usage.CPUs - segment.CPUs
		b.usage.Memory = b.usage.Memory - segment.Memory
		b.usage.Storage = b.usage.Storage - segment.Storage
		return
	}
	b.CPUs = roundValueAndMax(b.CPUs, b.usage.CPUs, b.conf.UnitCPU)
	b.Memory = roundValueAndMax(b.Memory, b.usage.Memory, b.conf.UnitMemory)
	b.Storage = roundValueAndMax(b.Storage, b.usage.Storage, b.conf.UnitStorage)
}

func (b *Bucket) tryAdd(segment Segment) bool {
	if !b.hasSpace(segment) {
		return false
	}
	b.update(segment, true)
	b.segments = append(b.segments, segment)
	return true
}

func (b *Bucket) tryRemove(segment Segment) bool {
	loc := b.findSegment(segment)
	if loc == -1 {
		return false
	}
	b.update(segment, false)
	if loc == len(b.segments)-1 {
		b.segments = append(b.segments[:loc])
	} else {
		b.segments = append(b.segments[0:loc], b.segments[loc:]...)
	}
	return true
}
