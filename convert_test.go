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

package date_test

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	testifyrequire "github.com/stretchr/testify/require"

	date "github.com/hardfinhq/go-date"
)

func TestNullDateFromPtr(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d1 := &date.Date{Year: 2000, Month: time.January, Day: 1}
	nd1 := date.NullDateFromPtr(d1)
	expected := date.NullDate{Date: *d1, Valid: true}
	assert.Equal(expected, nd1)

	var d2 *date.Date
	nd2 := date.NullDateFromPtr(d2)
	expected = date.NullDate{Valid: false}
	assert.Equal(expected, nd2)
}

func TestNullTimeFromPtr(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	var d *date.Date
	nt := date.NullTimeFromPtr(d)
	expected := sql.NullTime{Valid: false}
	assert.Equal(expected, nt)

	d = &date.Date{Year: 2000, Month: time.January, Day: 1}
	nt = date.NullTimeFromPtr(d)
	expected = sql.NullTime{Time: time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC), Valid: true}
	assert.Equal(expected, nt)

	tz, err := time.LoadLocation("America/Chicago")
	assert.Nil(err)
	nt = date.NullTimeFromPtr(d, date.OptConvertTimezone(tz))
	expected = sql.NullTime{Time: time.Date(2000, time.January, 1, 0, 0, 0, 0, tz), Valid: true}
	assert.Equal(expected, nt)

	nt = date.NullTimeFromPtr(
		d,
		date.OptConvertHour(12),
		date.OptConvertMinute(30),
		date.OptConvertSecond(35),
		date.OptConvertNanosecond(123456789),
	)
	expected = sql.NullTime{Time: time.Date(2000, time.January, 1, 12, 30, 35, 123456789, time.UTC), Valid: true}
	assert.Equal(expected, nt)
}

func TestFromTime(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Time     string
		Date     date.Date
		Timezone timezoneMetadata
		Error    string
	}

	cases := []testCase{
		{
			Time: "2022-01-31T00:00:00.000Z",
			Date: date.Date{Year: 2022, Month: time.January, Day: 31},
		},
		{
			Time:  "2020-05-11T07:10:55.209309302Z",
			Error: "timestamp contains more than just date information; 2020-05-11T07:10:55.209309302Z",
		},
		{
			Time:     "2022-01-31T00:00:00.000-05:00",
			Timezone: timezoneMetadata{Name: valueToPtr(""), Offset: valueToPtr(-18000)},
			Error:    "timestamp contains more than just date information; 2022-01-31T00:00:00-05:00",
		},
		{
			Time:     "2022-01-31T05:00:00.000Z",
			Timezone: timezoneMetadata{InTimezone: valueToPtr("America/New_York"), Name: valueToPtr("EST"), Offset: valueToPtr(-18000)},
			Error:    "timestamp contains more than just date information; 2022-01-31T00:00:00-05:00",
		},
		{
			Time:     "2024-01-11T00:00:00.000-06:00",
			Timezone: timezoneMetadata{Name: valueToPtr(""), Offset: valueToPtr(-21600)},
			Date:     date.Date{Year: 2024, Month: time.January, Day: 11},
		},
		{
			Time:     "2024-04-11T00:00:00.000-05:00",
			Timezone: timezoneMetadata{Name: valueToPtr(""), Offset: valueToPtr(-18000)},
			Date:     date.Date{Year: 2024, Month: time.April, Day: 11},
		},
		{
			Time:     "2024-04-11T05:00:00.000Z",
			Timezone: timezoneMetadata{InTimezone: valueToPtr("America/Chicago"), Name: valueToPtr("CDT"), Offset: valueToPtr(-18000)},
			Error:    "timestamp contains more than just date information; 2024-04-11T00:00:00-05:00",
		},
		{
			Time:  "2020-05-11T00:00:00.000000001Z",
			Error: "timestamp contains more than just date information; 2020-05-11T00:00:00.000000001Z",
		},
		{
			Time:  "2020-05-11T00:00:01Z",
			Error: "timestamp contains more than just date information; 2020-05-11T00:00:01Z",
		},
		{
			Time:  "2020-05-11T00:01:00Z",
			Error: "timestamp contains more than just date information; 2020-05-11T00:01:00Z",
		},
		{
			Time:  "2020-05-11T01:00:00Z",
			Error: "timestamp contains more than just date information; 2020-05-11T01:00:00Z",
		},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Time, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			timestamp, err := time.Parse(time.RFC3339Nano, tc.Time)
			assert.Nil(err)

			timestamp = tc.Timezone.In(assert, timestamp)

			name, offset := timestamp.Zone()
			assert.Equal(tc.Timezone.ExpectedName(timestamp), name)
			assert.Equal(tc.Timezone.ExpectedOffset(), offset)

			d, err := date.FromTime(timestamp)
			if tc.Error == "" {
				assert.Nil(err)
				assert.Equal(tc.Date, d)
			} else {
				assert.Equal(tc.Error, fmt.Sprintf("%v", err))
				assert.Equal(date.Date{}, d)
			}
		})
	}
}

func TestInTimezone(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Time     string
		Timezone string
		Date     string
	}

	cases := []testCase{
		{Time: "2024-02-01T06:41:35.540349Z", Timezone: "America/Los_Angeles", Date: "2024-01-31"},
		{Time: "2024-02-01T06:41:35.540349Z", Timezone: "America/Denver", Date: "2024-01-31"},
		{Time: "2024-02-01T06:41:35.540349Z", Timezone: "America/Chicago", Date: "2024-02-01"},
		{Time: "2024-02-01T06:41:35.540349Z", Timezone: "America/New_York", Date: "2024-02-01"},
		{Time: "2024-02-01T06:41:35.540349Z", Timezone: "UTC", Date: "2024-02-01"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		description := fmt.Sprintf("%s::%s", tc.Time, tc.Timezone)
		base.Run(description, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			timestamp, err := time.Parse(time.RFC3339Nano, tc.Time)
			assert.Nil(err)

			tz, err := time.LoadLocation(tc.Timezone)
			assert.Nil(err)

			expected, err := date.FromString(tc.Date)
			assert.Nil(err)

			d := date.InTimezone(timestamp, tz)
			assert.Equal(expected, d)
		})
	}
}

// timezoneMetadata is a struct that contains timezone metadata for assertions
// and translation across timezones. Intended to be used with `TestFromTime()`.
type timezoneMetadata struct {
	InTimezone *string
	Name       *string
	Offset     *int
}

// In translates a timestamp to an "in timezone" if one is set on this
// metadata struct.
func (tm timezoneMetadata) In(assert *testifyrequire.Assertions, t time.Time) time.Time {
	if tm.InTimezone == nil {
		return t
	}

	tz, err := time.LoadLocation(*tm.InTimezone)
	assert.Nil(err)
	return t.In(tz)
}

// ExpectedName returns the expected timezone name.
func (tm timezoneMetadata) ExpectedName(t time.Time) string {
	if tm.Name == nil {
		return "UTC"
	}

	name := *tm.Name
	tz := t.Location()
	if name == "" && tz == time.Local {
		name, _ = t.Zone()
	}

	return name
}

// ExpectedOffset returns the expected timezone offset in seconds.
func (tm timezoneMetadata) ExpectedOffset() int {
	if tm.Offset != nil {
		return *tm.Offset
	}
	return 0
}

// valueToPtr is a generic function that returns a pointer to the given value.
func valueToPtr[T any](v T) *T {
	return &v
}
