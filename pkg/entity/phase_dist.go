/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

type PhaseDist []Bucket

func (pd PhaseDist) FindBucket(name string) int {
	for i, bucket := range []Bucket(pd) {
		if bucket.FindByName(name) != -1 {
			return i
		}
	}
	return -1
}
