// Copyright 2017 Kirill Danshin and Gramework contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//

package common

import (
	"database/sql/driver"
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

func TestCalcMTA(t *testing.T) {
	maxAprox := 10 * time.Millisecond

	type args struct {
		nextSample time.Duration
	}
	tests := []struct {
		name string
		args []args
		want []time.Duration
	}{
		{
			name: "basic",
			args: []args{
				{
					nextSample: 0,
				},
				{
					nextSample: 10 * time.Minute,
				},
				{
					nextSample: 30 * time.Minute,
				},
			},
			want: []time.Duration{
				0 * time.Minute,
				5 * time.Minute,
				13*time.Minute + 20*time.Second,
			},
		},
		{
			name: "basic 2",
			args: []args{
				{
					nextSample: 10 * time.Minute,
				},
				{
					nextSample: 15 * time.Minute,
				},
				{
					nextSample: 20 * time.Minute,
				},
			},
			want: []time.Duration{
				10 * time.Minute,
				12*time.Minute + 30*time.Second,
				15 * time.Minute,
			},
		},
		{
			name: "basic 3",
			args: []args{
				{
					nextSample: 10 * time.Millisecond,
				},
				{
					nextSample: 15 * time.Millisecond,
				},
				{
					nextSample: 20 * time.Millisecond,
				},
			},
			want: []time.Duration{
				10 * time.Millisecond,
				12*time.Millisecond + 30*time.Microsecond,
				15 * time.Millisecond,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mta := &MTAState{}
			for stepN, sample := range tt.args {
				res := mta.Calc(sample.nextSample)
				if !reflect.DeepEqual(res, tt.want[stepN]) && res-tt.want[stepN] > maxAprox {
					t.Errorf(
						"step %d: CalcMTA(%v) = %v (uint=%d), want %v (uint=%d)",
						stepN+1,
						sample,
						res,
						uint64(res),
						tt.want[stepN],
						uint64(tt.want[stepN]),
					)
				}
			}
		})
	}
}

func TestMTAStateSQL(t *testing.T) {
	type fields struct {
		SamplesCount uint64
		LatestAvg    time.Duration
	}
	bf := fields{
		SamplesCount: 1,
		LatestAvg:    time.Minute,
	}
	tests := []struct {
		name    string
		fields  fields
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "basic",
			fields:  bf,
			want:    func() []byte { b, _ := json.Marshal(bf); return b }(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mta := &MTAState{
				SamplesCount: tt.fields.SamplesCount,
				LatestAvg:    tt.fields.LatestAvg,
			}
			got, err := mta.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("MTAState.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MTAState.Value() = %v, want %v", got, tt.want)
			}

			newMTA := &MTAState{}
			err = newMTA.Scan(got)
			if err != nil {
				t.Errorf("MTAState.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(mta, newMTA) {
				t.Errorf("MTAState.Value() = %v, want %v", newMTA, mta)
			}
		})
	}
}

func TestMTAState_Scan(t *testing.T) {
	type fields struct {
		SamplesCount uint64
		LatestAvg    time.Duration
	}
	type args struct {
		value interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "nil check",
			args:    args{nil},
			wantErr: false,
		},
		{
			name:    "non-[]byte check",
			args:    args{0},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mta := &MTAState{
				SamplesCount: tt.fields.SamplesCount,
				LatestAvg:    tt.fields.LatestAvg,
			}
			if err := mta.Scan(tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("MTAState.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
