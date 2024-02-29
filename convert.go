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
	"database/sql"
	"fmt"
	"time"
)

// ConvertConfig helps customize the behavior of conversion functions like
// `NullTimeFromPtr()`.
//
// It allows setting the fields in a `time.Time{}` **other** than year, month,
// and day (i.e. the fields that aren't present in a date). By default, these
// are:
// - hour=0
// - minute=0
// - second=0
// - nanosecond=0
// - timezone/loc=time.UTC
type ConvertConfig struct {
	Hour       int
	Minute     int
	Second     int
	Nanosecond int
	Timezone   *time.Location
}

// ConvertOption defines a function that will be applied to a convert config.
type ConvertOption func(*ConvertConfig)

// OptConvertHour returns an option that sets the hour on a convert config.
func OptConvertHour(hour int) ConvertOption {
	return func(cc *ConvertConfig) {
		cc.Hour = hour
	}
}

// OptConvertMinute returns an option that sets the minute on a convert config.
func OptConvertMinute(minute int) ConvertOption {
	return func(cc *ConvertConfig) {
		cc.Minute = minute
	}
}

// OptConvertSecond returns an option that sets the second on a convert config.
func OptConvertSecond(second int) ConvertOption {
	return func(cc *ConvertConfig) {
		cc.Second = second
	}
}

// OptConvertNanosecond returns an option that sets the nanosecond on a convert
// config.
func OptConvertNanosecond(nanosecond int) ConvertOption {
	return func(cc *ConvertConfig) {
		cc.Nanosecond = nanosecond
	}
}

// OptConvertTimezone returns an option that sets the timezone on a convert
// config.
func OptConvertTimezone(tz *time.Location) ConvertOption {
	return func(cc *ConvertConfig) {
		cc.Timezone = tz
	}
}

// NullDateFromPtr converts a `Date` pointer into a `NullDate`.
func NullDateFromPtr(d *Date) NullDate {
	if d == nil {
		return NullDate{Valid: false}
	}

	return NullDate{Date: *d, Valid: true}
}

// NullTimeFromPtr converts a date to a native Go `sql.NullTime`; the
// convention in Go is that a **date-only** is parsed (via `time.DateOnly`) as
// `time.Date(YYYY, MM, DD, 0, 0, 0, 0, time.UTC)`.
func NullTimeFromPtr(d *Date, opts ...ConvertOption) sql.NullTime {
	if d == nil {
		return sql.NullTime{Valid: false}
	}

	t := d.ToTime(opts...)
	return sql.NullTime{Time: t, Valid: true}
}

// FromString parses a string of the form YYYY-MM-DD into a `Date{}`.
func FromString(s string) (Date, error) {
	t, err := time.Parse(time.DateOnly, s)
	if err != nil {
		return Date{}, err
	}

	year, month, day := t.Date()
	d := Date{Year: year, Month: month, Day: day}
	return d, nil
}

// FromTime validates that a `time.Time{}` contains a date and converts it to a
// `Date{}`.
func FromTime(t time.Time) (Date, error) {
	if t.Hour() != 0 ||
		t.Minute() != 0 ||
		t.Second() != 0 ||
		t.Nanosecond() != 0 ||
		t.Location() != time.UTC {
		return Date{}, fmt.Errorf("timestamp contains more than just date information; %s", t.Format(time.RFC3339Nano))
	}

	year, month, day := t.Date()
	d := Date{Year: year, Month: month, Day: day}
	return d, nil
}

// InTimezone translates a timestamp into a timezone and then captures the date
// in that timezone.
func InTimezone(t time.Time, tz *time.Location) Date {
	tLocal := t.In(tz)
	year, month, day := tLocal.Date()
	return Date{Year: year, Month: month, Day: day}
}
