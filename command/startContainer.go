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
	InfiniteTimeoutTerm = "infinite"
	DefaultTimeout      = 2 * time.Minute
)

// Timeout represents the frame of time for a task runner to timeout.
type Timeout struct {
	time.Duration
	isInfinite bool
}

func (to Timeout) IsInfinite() bool {
	return to.isInfinite
}

func (to Timeout) MarshalJSON() ([]byte, error) {
	if to.IsInfinite() {
		return json.Marshal(InfiniteTimeoutTerm)
	}
	return json.Marshal(to.Duration)
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
		}
	}
	if strVal, ok := dat.(string); ok {
		if strVal == InfiniteTimeoutTerm {
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

// StartContainer is the command for starting a container
type StartContainer struct {
	Name   string `json:"name"`
	Attach bool   `json:"attach"`
	// Timeout is the maximum amount of time to wait for the task before terminating it.
	// This is ignored if attach is false
	Timeout Timeout `json:"timeout"`
}
