// Copyright 2017 Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

// Package common provides common web-related shortcuts,
// solutions and algorithms
package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

var ErrCouldNotScan = errors.New("gramework/common: could not scan value")

// MTAState holds all required data for Moving Time Average.
// Useful when using within ORM models
type MTAState struct {
	SamplesCount uint64
	LatestAvg    time.Duration
	Mu           sync.RWMutex `sql:"-" json:"-" xml:"-" csv:"-"`
}

// Scan implements the Scanner interface.
func (mta *MTAState) Scan(value interface{}) error {
	if value == nil {
		mta.SamplesCount, mta.LatestAvg = 0, 0
		return nil
	}
	b, ok := value.([]byte)
	if !ok {
		return ErrCouldNotScan
	}
	return json.Unmarshal(b, mta)
}

// Value implements the driver Valuer interface.
func (mta *MTAState) Value() (driver.Value, error) {
	mta.Mu.RLock()
	b, err := json.Marshal(mta)
	mta.Mu.RUnlock()
	return b, err
}

// Calc is the wrapper on CalcMTA which uses only current state data
// and provided next sample. This function automatically updates
// current state and it is concurrency safe
//
// Warning: this MTA will return approx. value and not the exact one.
// This is what we paid so we're not storing all values.
func (mta *MTAState) Calc(nextSample time.Duration) time.Duration {
	mta.Mu.Lock()
	newAvg := CalcMTA(mta.SamplesCount, mta.LatestAvg, nextSample)
	mta.SamplesCount++
	mta.LatestAvg = newAvg
	mta.Mu.Unlock()
	return newAvg
}

// CalcMTA calculates Moving Time Average
//
// samplesCount should not count the nextSample. given that,
// when initializing MTA, you should pass 0, not 1;
// and when calling 2nd time, you should pass 1, not 2.
//
// Warning: this MTA will return approx. value and not the exact one.
// This is what we paid so we're not storing all values.
func CalcMTA(samplesCount uint64, latestAvg time.Duration, nextSample time.Duration) time.Duration {
	if samplesCount == 0 {
		return nextSample
	}
	samplesCount++

	sc := time.Duration(samplesCount)
	avg := latestAvg*(sc-1)/sc + nextSample/sc

	return avg
}
