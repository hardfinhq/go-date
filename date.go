// Copyright 2023 Hardfin, Inc.
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
	"database/sql/driver"
	"encoding"
	"encoding/json"
	"fmt"
	"time"
)

// NOTE: Ensure that
// - `Date` satisfies `fmt.Stringer`.
// - `Date` satisfies `fmt.GoStringer`.
// - `Date` satisfies `encoding.TextMarshaler`.
// - `Date` satisfies `json.Marshaler`.
// - `*Date` satisfies `encoding.TextUnmarshaler`.
// - `*Date` satisfies `json.Unmarshaler`.
// - `*Date` satisfies `sql.Scanner`.
// - `Date` satisfies `driver.Valuer`.
var (
	_ fmt.Stringer             = Date{}
	_ fmt.GoStringer           = Date{}
	_ encoding.TextMarshaler   = Date{}
	_ json.Marshaler           = Date{}
	_ encoding.TextUnmarshaler = (*Date)(nil)
	_ json.Unmarshaler         = (*Date)(nil)
	_ sql.Scanner              = (*Date)(nil)
	_ driver.Valuer            = Date{}
)

// Date is a simple date (i.e. without timestamp). This is intended to be
// JSON serialized / deserialized as YYYY-MM-DD.
type Date struct {
	Year  int
	Month time.Month
	Day   int
}

// NewDate returns a new `Date` struct. This is a pure convenience function to
// make it more ergonomic to create a `Date` struct.
func NewDate(year int, month time.Month, day int) Date {
	return Date{Year: year, Month: month, Day: day}
}

// AddDays returns the date corresponding to adding the given number of days.
func (d Date) AddDays(days int) Date {
	t := d.ToTime().AddDate(0, 0, days)
	return Date{Year: t.Year(), Month: t.Month(), Day: t.Day()}
}

// AddMonths returns the date corresponding to adding the given number of
// months. This accounts for leap years and variable length months. Typically
// the only change is in the month and year but for changes that would exceed
// the number of days in the target month, the last day of the month is used.
//
// For example:
// - adding 1 month to 2020-05-11 results in 2020-06-11
// - adding 1 month to 2022-01-31 results in 2022-02-28
// - adding 3 months to 2024-01-31 results in 2024-04-30
// - subtracting 2 months from 2022-01-31 results in 2022-11-30
//
// NOTE: This behavior is very similar to but distinct from
// `time.Time{}.AddDate()` specialized to `months` only.
func (d Date) AddMonths(months int) Date {
	updatedMonth, yearDelta := monthsChange(d.Month, months)
	updatedYear := d.Year + yearDelta
	updatedDay := minInt(d.Day, daysIn(updatedMonth, updatedYear))
	return Date{Year: updatedYear, Month: updatedMonth, Day: updatedDay}
}

// AddMonthsStdlib returns the date corresponding to adding the given number of
// months, using `time.Time{}.AddDate()` from the standard library. This may
// "overshoot" if the target date is not a valid date in that month, e.g.
// 2020-02-31.
//
// For example:
// - adding 1 month to 2020-05-11 results in 2020-06-11
// - adding 1 month to 2022-01-31 results in 2022-03-03
// - adding 3 months to 2024-01-31 results in 2024-05-01
// - subtracting 2 months from 2022-01-31 results in 2022-12-01
func (d Date) AddMonthsStdlib(months int) Date {
	t := d.ToTime().AddDate(0, months, 0)
	return Date{Year: t.Year(), Month: t.Month(), Day: t.Day()}
}

func monthsChange(month time.Month, monthDelta int) (time.Month, int) {
	monthsTotal := int(month) + monthDelta
	monthsInYear := monthsTotal % 12
	yearDelta := (monthsTotal - monthsInYear) / 12
	if monthsInYear < 1 {
		// +12 months <==> -1 year
		return time.Month(monthsInYear + 12), yearDelta - 1
	}

	return time.Month(monthsInYear), yearDelta
}

// AddYears returns the date corresponding to adding the given number of
// years, using `time.Time{}.AddDate()` from the standard library. This may
// "overshoot" if the target date is not a valid date in that month, e.g.
// 2020-02-31.
//
// For example:
// - adding 1 year to 2020-02-29 results in 2021-03-01
// - adding 1 year to 2023-02-28 results in 2024-02-28
// - adding 10 years to 2010-05-01 results in 2020-05-01
// - subtracting 10 years from 2010-05-01 results in 2000-05-01
//
// NOTE: This behavior is very similar to but distinct from
// `time.Time{}.AddDate()` specialized to `years` only.
func (d Date) AddYears(years int) Date {
	updatedMonth := d.Month
	updatedYear := d.Year + years
	updatedDay := minInt(d.Day, daysIn(updatedMonth, updatedYear))
	return Date{Year: updatedYear, Month: updatedMonth, Day: updatedDay}
}

// AddYearsStdlib returns the date corresponding to adding the given number of
// years. This accounts for leap years and variable length months. Typically
// the only change is in the month and year but for changes that would exceed
// the number of days in the target month, the last day of the month is used.
//
// For example:
// - adding 1 year to 2020-02-29 results in 2021-02-28
// - adding 1 year to 2023-02-28 results in 2024-02-28
// - adding 10 years to 2010-05-01 results in 2020-05-01
// - subtracting 10 years from 2010-05-01 results in 2000-05-01
//
// NOTE: This behavior is very similar to but distinct from
// `time.Time{}.AddDate()` specialized to `years` only.
func (d Date) AddYearsStdlib(years int) Date {
	t := d.ToTime().AddDate(years, 0, 0)
	return Date{Year: t.Year(), Month: t.Month(), Day: t.Day()}
}

// Sub returns the number of days `d - other`; this converts both dates to
// a `time.Time{}` UTC and then dispatches to `time.Time{}.Sub()`.
func (d Date) Sub(other Date) int64 {
	days, err := d.SubErr(other)
	mustNil(err)
	return int64(days)
}

// SubErr returns the number of days `d - other`; this converts both dates to
// a `time.Time{}` UTC and then dispatches to `time.Time{}.Sub()`.
//
// If the number of days is not a whole number (due to overflow), an error is
// returned.
func (d Date) SubErr(other Date) (int64, error) {
	duration := d.ToTime().Sub(other.ToTime())

	day := 24 * time.Hour
	days := duration / day
	remainder := duration % day
	if remainder != 0 {
		return 0, fmt.Errorf("duration is not a whole number of days; duration=%s", duration)
	}

	return int64(days), nil
}

// MonthStart returns the first date in the month of the current date.
func (d Date) MonthStart() Date {
	return Date{Year: d.Year, Month: d.Month, Day: 1}
}

// MonthEnd returns the last date in the month of the current date.
func (d Date) MonthEnd() Date {
	endDay := daysIn(d.Month, d.Year)
	return Date{Year: d.Year, Month: d.Month, Day: endDay}
}

// Before returns true if the date is before the other date.
func (d Date) Before(other Date) bool {
	if d.Year != other.Year {
		return d.Year < other.Year
	}

	if d.Month != other.Month {
		return d.Month < other.Month
	}

	return d.Day < other.Day
}

// After returns true if the date is after the other date.
func (d Date) After(other Date) bool {
	return other.Before(d)
}

// Equal returns true if the date is equal to the other date.
func (d Date) Equal(other Date) bool {
	return d.Year == other.Year && d.Month == other.Month && d.Day == other.Day
}

func compareInt(i1, i2 int) int {
	if i1 < i2 {
		return -1
	}

	if i1 > i2 {
		return 1
	}

	return 0
}

// Compare compares the date d with other. If d is before other, it returns
// -1; if d is after other, it returns +1; if they're the same, it returns 0.
func (d Date) Compare(other Date) int {
	if d.Year != other.Year {
		return compareInt(d.Year, other.Year)
	}

	if d.Month != other.Month {
		return compareInt(int(d.Month), int(other.Month))
	}

	return compareInt(d.Day, other.Day)
}

// IsZero returns true if the date is the zero value.
func (d Date) IsZero() bool {
	return d.Year == 0 && d.Month == 0 && d.Day == 0
}

// ToTime converts the date to a native Go `time.Time`; the convention in Go is
// that a **date-only** is parsed (via `time.DateOnly`) as
// `time.Date(YYYY, MM, DD, 0, 0, 0, 0, time.UTC)`.
func (d Date) ToTime(opts ...ConvertOption) time.Time {
	cc := ConvertConfig{Timezone: time.UTC}
	for _, opt := range opts {
		opt(&cc)
	}

	return time.Date(d.Year, d.Month, d.Day, cc.Hour, cc.Minute, cc.Second, cc.Nanosecond, cc.Timezone)
}

// ISOWeek returns the ISO 8601 year and week number in which `d` occurs.
// Week ranges from 1 to 53. Jan 01 to Jan 03 of year `n` might belong to
// week 52 or 53 of year `n-1`, and Dec 29 to Dec 31 might belong to week 1
// of year `n+1`.
func (d Date) ISOWeek() (year, week int) {
	return d.ToTime().ISOWeek()
}

// Weekday returns the day of the week specified by `d`.
func (d Date) Weekday() time.Weekday {
	return d.ToTime().Weekday()
}

// MarshalText implements the encoding.TextMarshaler interface.
func (d Date) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// MarshalJSON implements `json.Marshaler`; formats the date as YYYY-MM-DD.
func (d Date) MarshalJSON() ([]byte, error) {
	s := d.String()
	return json.Marshal(s)
}

// UnmarshalText implements the encoding.TextUnmarshaler interface. The time
// must be in the format YYYY-MM-DD.
func (d *Date) UnmarshalText(data []byte) error {
	parsed, err := FromString(string(data))
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}

// UnmarshalJSON implements `json.Unmarshaler`; parses the date as YYYY-MM-DD.
func (d *Date) UnmarshalJSON(data []byte) error {
	s := ""
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	parsed, err := FromString(s)
	if err != nil {
		return err
	}

	*d = parsed
	return nil
}

// Scan implements `sql.Scanner`; it unmarshals values of the type `time.Time`
// onto the current `Date` struct.
func (d *Date) Scan(src any) error {
	var t time.Time

	switch srcTyped := src.(type) {
	case time.Time:
		t = srcTyped
	default:
		return fmt.Errorf("incompatible type for Date; type=%T", src)
	}

	verified, err := FromTime(t)
	if err != nil {
		return err
	}

	*d = verified
	return nil
}

// Value implements `driver.Valuer`; it marshals the value to a `time.Time`
// to be serialized into the database.
func (d Date) Value() (driver.Value, error) {
	return d.ToTime(), nil
}

// String implements `fmt.Stringer`.
func (d Date) String() string {
	return d.Format(time.DateOnly)
}

// Format returns a textual representation of the date value formatted according
// to the provided layout. This uses `time.Time{}.Format()` directly and is
// provided here for convenience.
func (d Date) Format(layout string) string {
	return d.ToTime().Format(layout)
}

// GoString implements `fmt.GoStringer`.
func (d Date) GoString() string {
	return fmt.Sprintf("date.NewDate(%d, time.%s, %d)", d.Year, d.Month, d.Day)
}
