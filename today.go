// Copyright 2024 Hardfin, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package date

import (
	"time"
)

// TodayConfig helps customize the behavior of `Today()`.
type TodayConfig struct {
	Timezone    *time.Location
	NowProvider func() time.Time
}

// TodayOption defines a function that will be applied to a `Today()` config.
type TodayOption func(*TodayConfig)

// Today determines the **current** `Date`, shifted to a given timezone
// if need be.
//
// Defaults to using UTC and `time.Now()` to determine the current time.
func Today(opts ...TodayOption) Date {
	tc := TodayConfig{
		Timezone:    time.UTC,
		NowProvider: time.Now,
	}
	for _, opt := range opts {
		opt(&tc)
	}

	now := tc.NowProvider().In(tc.Timezone)
	year, month, day := now.Date()
	return Date{Year: year, Month: month, Day: day}
}

// OptTodayTimezone returns an option that sets the timezone on a `Today()`
// config.
func OptTodayTimezone(tz *time.Location) TodayOption {
	return func(tc *TodayConfig) {
		tc.Timezone = tz
	}
}

// OptTodayNow returns an option that sets the now provider on a `Today()`
// config to return a **constant** `now` value.
//
// This is expected to be used in tests.
func OptTodayNow(now time.Time) TodayOption {
	return func(tc *TodayConfig) {
		tc.NowProvider = func() time.Time {
			return now
		}
	}
}
