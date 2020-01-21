/**
 * Copyright 2019 Whiteblock Inc. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

package command

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	InfiniteTimeTerm = "infinite"
	DefaultTimeout   = 2 * time.Minute
)

// Time represents the frame of time which can be infinite
type Time struct {
	time.Duration
	isInfinite bool
}

func (to Time) IsInfinite() bool {
	return to.isInfinite
}

func (to Time) MarshalJSON() ([]byte, error) {
	if to.IsInfinite() {
		return json.Marshal(InfiniteTimeTerm)
	}
	return json.Marshal(to.Duration)
}

func (to Time) MarshalYAML() (interface{}, error) {
	if to.IsInfinite() {
		return InfiniteTimeTerm, nil
	}
	return to.Duration, nil
}

type Timeout struct {
	Time
}

func (to Timeout) SetInfinite() Timeout {
	to.isInfinite = true
	return to
}

func (to *Timeout) UnmarshalJSON(data []byte) error {
	var dat interface{}

	err := json.Unmarshal(data, &dat)
	if err != nil {
		return err
	}
	if dat == nil {
		to.Duration = DefaultTimeout
		return nil
	}
	if floatVal, ok := dat.(float64); ok {
		if floatVal == 0 {
			to.Duration = DefaultTimeout
			return nil
		} else { //if no units are specified, use seconds
			to.Duration = time.Duration(floatVal) * time.Second
		}
	}
	if strVal, ok := dat.(string); ok {
		if strVal == InfiniteTimeTerm {
			to.isInfinite = true
			return nil
		}
		strVal = strings.Replace(strVal, " ", "", -1)
		to.Duration, err = time.ParseDuration(strVal)
		return err
	}

	return json.Unmarshal(data, &to.Duration)
}

// UnmarshalYAML converts from yaml into this object
func (to *Timeout) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// This is such a small object, let's just convert to JSON to reuse that logic
	var tmp interface{}
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}
	data, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	return to.UnmarshalJSON(data)
}

type Duration struct {
	Time
}

func (to Duration) Empty() bool {
	return !to.IsInfinite() && to.Duration == NoDuration
}

func (to *Duration) UnmarshalJSON(data []byte) error {
	var dat interface{}

	err := json.Unmarshal(data, &dat)
	if err != nil {
		return err
	}
	if dat == nil {
		to.Duration = NoDuration
		return nil
	}
	if floatVal, ok := dat.(float64); ok {
		//if no units are specified, use seconds
		to.Duration = time.Duration(floatVal) * time.Second
	}
	if strVal, ok := dat.(string); ok {
		if strVal == InfiniteTimeTerm {
			to.isInfinite = true
			return nil
		}
		strVal = strings.Replace(strVal, " ", "", -1)
		to.Duration, err = time.ParseDuration(strVal)
		return err
	}

	return json.Unmarshal(data, &to.Duration)
}

// UnmarshalYAML converts from yaml into this object
func (to *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// This is such a small object, let's just convert to JSON to reuse that logic
	var tmp interface{}
	err := unmarshal(&tmp)
	if err != nil {
		return err
	}
	data, err := json.Marshal(tmp)
	if err != nil {
		return err
	}
	return to.UnmarshalJSON(data)
}
