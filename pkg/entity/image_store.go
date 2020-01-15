/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package entity

type ImageStore struct {
	store map[string]map[string]map[string]string
}

func (is *ImageStore) Insert(instance, image string, meta map[string]string) {
	if is.store == nil {
		is.store = map[string]map[string]map[string]string{}
	}

	if is.store[instance] == nil {
		is.store[instance] = map[string]map[string]string{}
	}

	is.store[instance][image] = meta
}

func (is *ImageStore) ForEach(fn func(instance, image string, meta map[string]string) error) error {
	for instance, images := range is.store {
		for image, meta := range images {
			if image == "" {
				continue
			}
			err := fn(instance, image, meta)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
