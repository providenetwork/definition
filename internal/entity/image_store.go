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
			err := fn(instance, image, meta)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
