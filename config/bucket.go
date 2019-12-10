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

package config

import (
	"github.com/spf13/viper"
)

//Bucket represents the basic configuration for a bucket
type Bucket struct {
	MaxCPU     int64 `mapstructure:"bucketMaxCPU"`
	MaxMemory  int64 `mapstructure:"bucketMaxMemory"`
	MaxStorage int64 `mapstructure:"bucketMaxStorage"`

	MinCPU     int64 `mapstructure:"bucketMinCPU"`
	MinMemory  int64 `mapstructure:"bucketMinMemory"`
	MinStorage int64 `mapstructure:"bucketMinStorage"`

	UnitCPU     int64 `mapstructure:"bucketUnitCPU"`
	UnitMemory  int64 `mapstructure:"bucketUnitMemory"`
	UnitStorage int64 `mapstructure:"bucketUnitStorage"`

	MaxBuckets int64 `mapstructure:"maxBuckets"`
}

//NewBucket generates a Bucket configuration from the given viper
//  Configuration
func NewBucket(v *viper.Viper) (Bucket, error) {
	out := Bucket{}
	return out, v.Unmarshal(&out)
}

func setBucketBindings(v *viper.Viper) error {
	err := v.BindEnv("bucketMaxCPU", "BUCKET_MAX_CPU")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMaxMemory", "BUCKET_MAX_MEMORY")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMaxStorage", "BUCKET_MAX_STORAGE")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMinCPU", "BUCKET_MIN_CPU")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMinMemory", "BUCKET_MIN_MEMORY")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketMinStorage", "BUCKET_MIN_STORAGE")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketUnitCPU", "BUCKET_UNIT_CPU")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketUnitMemory", "BUCKET_UNIT_MEMORY")
	if err != nil {
		return err
	}

	err = v.BindEnv("bucketUnitStorage", "BUCKET_UNIT_STORAGE")
	if err != nil {
		return err
	}

	return v.BindEnv("maxBuckets", "MAX_BUCKETS")
}

func setBucketDefaults(v *viper.Viper) {
	v.SetDefault("bucketMaxCPU", 92)
	v.SetDefault("bucketMaxMemory", 600*1024)     //624 GiB in MiB
	v.SetDefault("bucketMaxStorage", 1*1024*1024) //TiB in MiB
	v.SetDefault("bucketMinCPU", 1)
	v.SetDefault("bucketMinMemory", 1*1024)   //1 GiB in MiB
	v.SetDefault("bucketMinStorage", 10*1024) //10 GiB in MiB
	v.SetDefault("bucketUnitCPU", 1)
	v.SetDefault("bucketUnitMemory", 128)     //128MiB
	v.SetDefault("bucketUnitStorage", 1*1024) //1 GiB in MiB
	v.SetDefault("maxBuckets", 3000)          //Max 3000 instances
}
