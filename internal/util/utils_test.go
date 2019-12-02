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

package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemconv(t *testing.T) {
	val, err := Memconv("20", 1000)
	assert.NoError(t, err)
	assert.Equal(t, int64(20000), val)

	val, err = Memconv("200 KB", 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(200*Kibi), val)

	val, err = Memconv("2000 KiB", 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(2000*Kibi), val)

	val, err = Memconv("2TB", 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(2*Tibi), val)

	val, err = Memconv("2GB", 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(2*Gibi), val)

	val, err = Memconv("2 MB", 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(2*Mibi), val)

	_, err = Memconv("foo", 0)
	assert.Error(t, err)

	val, err = Memconv("5mb", 0)
	assert.NoError(t, err)
	assert.Equal(t, int64(5*Mibi), val)
}
