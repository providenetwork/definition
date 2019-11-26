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
	"fmt"

	"github.com/whiteblock/definition/internal/config"
	"github.com/whiteblock/definition/internal/entity"
)

type ResourceBuckets interface {
	Add(segments []entity.Segment) error
	Remove(segments []entity.Segment) error
	Resources() []Bucket
}

type resourceBuckets struct {
	conf    config.Bucket
	buckets []Bucket
}

func NewResourceBuckets(conf config.Bucket) ResourceBuckets {
	return &resourceBuckets{conf: conf}
}

func (rb *resourceBuckets) add(segment entity.Segment) error {
	for i := range rb.buckets {
		if rb.buckets[i].tryAdd(segment) {
			return nil
		}
	}
	if int64(len(rb.buckets)) == rb.conf.MaxBuckets {
		return fmt.Errorf("size limits exceeded")
	}
	bucket := newBucket(&rb.conf)
	if !bucket.tryAdd(segment) {
		return fmt.Errorf("segment size too large")
	}
	rb.buckets = append(rb.buckets, bucket)
	return nil
}

func (rb *resourceBuckets) Add(segments []entity.Segment) error {
	for _, segment := range segments {
		err := rb.add(segment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rb *resourceBuckets) remove(segment entity.Segment) error {
	for i := range rb.buckets {
		if rb.buckets[i].tryRemove(segment) {
			return nil
		}
	}
	return fmt.Errorf("couldn't remove segment, doesn't exist")
}

func (rb *resourceBuckets) Remove(segments []entity.Segment) error {
	for _, segment := range segments {
		err := rb.remove(segment)
		if err != nil {
			return err
		}
	}
	return nil
}

func (rb *resourceBuckets) Resources() []Bucket {
	return rb.buckets
}
