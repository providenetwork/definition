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

package config

//Bucket represents the basic configuration for a bucket
type Bucket struct {
	MaxCPU     int64
	MaxMemory  int64
	MaxStorage int64

	MinCPU     int64
	MinMemory  int64
	MinStorage int64

	UnitCPU     int64
	UnitMemory  int64
	UnitStorage int64

	MaxBuckets int64
}
