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
	"strconv"
	"strings"
)

const (
	_          = iota
	Kibi int64 = 1 << (10 * iota)
	Mibi
	Gibi
	Tibi
)

func Memconv(mem string, defaultMulti int64) (int64, error) {

	m := strings.ToLower(mem)
	m = strings.Replace(m, " ", "", -1)
	m = strings.Replace(m, "ib", "b", -1)

	var multiplier int64 = defaultMulti

	if strings.HasSuffix(m, "kb") || strings.HasSuffix(m, "k") {
		multiplier = Kibi
	} else if strings.HasSuffix(m, "mb") || strings.HasSuffix(m, "m") {
		multiplier = Mibi
	} else if strings.HasSuffix(m, "gb") || strings.HasSuffix(m, "g") {
		multiplier = Gibi
	} else if strings.HasSuffix(m, "tb") || strings.HasSuffix(m, "t") {
		multiplier = Tibi
	}

	i, err := strconv.ParseInt(strings.Trim(m, "bgkmti"), 10, 64)
	if err != nil {
		return -1, err
	}

	return i * multiplier, nil
}
