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

package date_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	testifyrequire "github.com/stretchr/testify/require"

	date "github.com/hardfinhq/go-date"
)

func TestNewDate(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d := date.NewDate(2020, time.May, 11)
	expected := date.Date{Year: 2020, Month: time.May, Day: 11}
	assert.Equal(expected, d)
}

func TestDate_AddDays(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     string
		Delta    int
		Expected string
	}

	cases := []testCase{
		{Date: "2020-05-11", Delta: 0, Expected: "2020-05-11"},
		{Date: "2020-05-11", Delta: 1, Expected: "2020-05-12"},
		{Date: "2020-05-11", Delta: 10, Expected: "2020-05-21"},
		{Date: "2022-01-31", Delta: -5, Expected: "2022-01-26"},
		{Date: "2022-01-31", Delta: 40, Expected: "2022-03-12"},
		{Date: "2022-01-31", Delta: 120, Expected: "2022-05-31"},
		{Date: "1999-12-24", Delta: 300, Expected: "2000-10-19"},
		{Date: "1999-12-24", Delta: 2000, Expected: "2005-06-15"},
		{Date: "1999-12-24", Delta: 571, Expected: "2001-07-17"},
		// Daylight savings time
		{Date: "2023-03-10", Delta: 1, Expected: "2023-03-11"},
		{Date: "2023-03-10", Delta: 2, Expected: "2023-03-12"},
		{Date: "2023-03-10", Delta: 3, Expected: "2023-03-13"},
		{Date: "2023-03-10", Delta: 4, Expected: "2023-03-14"},
		{Date: "2023-11-03", Delta: 1, Expected: "2023-11-04"},
		{Date: "2023-11-03", Delta: 2, Expected: "2023-11-05"},
		{Date: "2023-11-03", Delta: 3, Expected: "2023-11-06"},
		{Date: "2023-11-03", Delta: 4, Expected: "2023-11-07"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		description := fmt.Sprintf("%s + %d -> %s", tc.Date, tc.Delta, tc.Expected)
		base.Run(description, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d, err := date.FromString(tc.Date)
			assert.Nil(err)

			computed := d.AddDays(tc.Delta)
			assert.Equal(tc.Expected, computed.String())
		})
	}
}

func TestDate_AddMonths(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     string
		Delta    int
		Expected string
		Contrast string
	}

	cases := []testCase{
		{Date: "2020-05-11", Delta: 0, Expected: "2020-05-11", Contrast: "2020-05-11"},
		{Date: "2020-05-11", Delta: 1, Expected: "2020-06-11", Contrast: "2020-06-11"},
		{Date: "2020-05-11", Delta: 2, Expected: "2020-07-11", Contrast: "2020-07-11"},
		{Date: "2022-01-31", Delta: -2, Expected: "2021-11-30", Contrast: "2021-12-01"},
		{Date: "2022-01-31", Delta: -1, Expected: "2021-12-31", Contrast: "2021-12-31"},
		{Date: "2022-01-31", Delta: 1, Expected: "2022-02-28", Contrast: "2022-03-03"},
		{Date: "2024-01-31", Delta: 2, Expected: "2024-03-31", Contrast: "2024-03-31"},
		{Date: "2024-01-31", Delta: 3, Expected: "2024-04-30", Contrast: "2024-05-01"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		description := fmt.Sprintf("%s + %d -> %s", tc.Date, tc.Delta, tc.Expected)
		base.Run(description, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d, err := date.FromString(tc.Date)
			assert.Nil(err)

			computed := d.AddMonths(tc.Delta)
			assert.Equal(tc.Expected, computed.String())

			// To contrast, consider how `time.Time{}.AddDate()` works.
			contrast := d.ToTime().AddDate(0, tc.Delta, 0)
			assert.Equal(tc.Contrast, contrast.Format(time.DateOnly))
		})
	}
}

func TestDate_AddMonthsStdlib(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     string
		Delta    int
		Expected string
	}

	cases := []testCase{
		{Date: "2020-05-11", Delta: 0, Expected: "2020-05-11"},
		{Date: "2020-05-11", Delta: 1, Expected: "2020-06-11"},
		{Date: "2020-05-11", Delta: 2, Expected: "2020-07-11"},
		{Date: "2022-01-31", Delta: -2, Expected: "2021-12-01"},
		{Date: "2022-01-31", Delta: -1, Expected: "2021-12-31"},
		{Date: "2022-01-31", Delta: 1, Expected: "2022-03-03"},
		{Date: "2024-01-31", Delta: 2, Expected: "2024-03-31"},
		{Date: "2024-01-31", Delta: 3, Expected: "2024-05-01"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		description := fmt.Sprintf("%s + %d -> %s", tc.Date, tc.Delta, tc.Expected)
		base.Run(description, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d, err := date.FromString(tc.Date)
			assert.Nil(err)

			computed := d.AddMonthsStdlib(tc.Delta)
			assert.Equal(tc.Expected, computed.String())
		})
	}
}

func TestDate_AddYears(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     string
		Delta    int
		Expected string
		Contrast string
	}

	cases := []testCase{
		{Date: "2020-05-11", Delta: 0, Expected: "2020-05-11", Contrast: "2020-05-11"},
		{Date: "2020-05-11", Delta: 1, Expected: "2021-05-11", Contrast: "2021-05-11"},
		{Date: "2020-05-11", Delta: 2, Expected: "2022-05-11", Contrast: "2022-05-11"},
		{Date: "2020-02-29", Delta: 1, Expected: "2021-02-28", Contrast: "2021-03-01"},
		{Date: "2020-02-29", Delta: 2, Expected: "2022-02-28", Contrast: "2022-03-01"},
		{Date: "2020-02-29", Delta: 4, Expected: "2024-02-29", Contrast: "2024-02-29"},
		{Date: "2019-02-28", Delta: 1, Expected: "2020-02-28", Contrast: "2020-02-28"},
		{Date: "2019-02-28", Delta: -1, Expected: "2018-02-28", Contrast: "2018-02-28"},
		{Date: "2019-02-28", Delta: -3, Expected: "2016-02-28", Contrast: "2016-02-28"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		description := fmt.Sprintf("%s + %d -> %s", tc.Date, tc.Delta, tc.Expected)
		base.Run(description, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d, err := date.FromString(tc.Date)
			assert.Nil(err)

			computed := d.AddYears(tc.Delta)
			assert.Equal(tc.Expected, computed.String())

			// To contrast, consider how `time.Time{}.AddDate()` works.
			contrast := d.ToTime().AddDate(tc.Delta, 0, 0)
			assert.Equal(tc.Contrast, contrast.Format(time.DateOnly))
		})
	}
}

func TestDate_AddYearsStdlib(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     string
		Delta    int
		Expected string
	}

	cases := []testCase{
		{Date: "2020-05-11", Delta: 0, Expected: "2020-05-11"},
		{Date: "2020-05-11", Delta: 1, Expected: "2021-05-11"},
		{Date: "2020-05-11", Delta: 2, Expected: "2022-05-11"},
		{Date: "2020-02-29", Delta: 1, Expected: "2021-03-01"},
		{Date: "2020-02-29", Delta: 2, Expected: "2022-03-01"},
		{Date: "2020-02-29", Delta: 4, Expected: "2024-02-29"},
		{Date: "2019-02-28", Delta: 1, Expected: "2020-02-28"},
		{Date: "2019-02-28", Delta: -1, Expected: "2018-02-28"},
		{Date: "2019-02-28", Delta: -3, Expected: "2016-02-28"},
		{Date: "2023-02-28", Delta: 1, Expected: "2024-02-28"},
		{Date: "2010-05-01", Delta: 10, Expected: "2020-05-01"},
		{Date: "2010-05-01", Delta: -10, Expected: "2000-05-01"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		description := fmt.Sprintf("%s + %d -> %s", tc.Date, tc.Delta, tc.Expected)
		base.Run(description, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d, err := date.FromString(tc.Date)
			assert.Nil(err)

			computed := d.AddYearsStdlib(tc.Delta)
			assert.Equal(tc.Expected, computed.String())
		})
	}
}

func TestDate_Sub(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     string
		Other    string
		Expected int64
	}

	cases := []testCase{
		{Date: "2020-05-11", Other: "2020-05-11", Expected: 0},
		{Date: "2020-05-11", Other: "2020-05-12", Expected: -1},
		{Date: "2020-05-11", Other: "2020-05-10", Expected: 1},
		{Date: "2020-05-11", Other: "2002-05-11", Expected: 6575},
		{Date: "2020-05-11", Other: "2022-05-11", Expected: -730},
		{Date: "2020-05-11", Other: "2012-01-31", Expected: 3023},
		{Date: "2016-04-17", Other: "2020-05-12", Expected: -1486},
		{Date: "2020-05-22", Other: "2020-05-10", Expected: 12},
		{Date: "2023-05-03", Other: "2002-05-11", Expected: 7662},
		{Date: "2013-05-19", Other: "2022-05-11", Expected: -3279},
		{Date: "2012-02-28", Other: "2012-01-31", Expected: 28},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		description := fmt.Sprintf("%s - %s", tc.Date, tc.Other)
		base.Run(description, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d, err := date.FromString(tc.Date)
			assert.Nil(err)

			other, err := date.FromString(tc.Other)
			assert.Nil(err)

			computed := d.Sub(other)
			assert.Equal(tc.Expected, computed)
		})
	}
}

func TestDate_Sub_Panic(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d1 := date.Date{Year: 1, Month: time.January, Day: 1}
	d2 := date.Date{Year: 1_000_000, Month: time.January, Day: 1}

	assert.Panics(func() { d1.Sub(d2) })

	days, err := d1.SubErr(d2)
	assert.Equal(int64(0), days)
	assert.NotNil(err)
	assert.Equal("duration is not a whole number of days; duration=-2562047h47m16.854775808s", fmt.Sprintf("%v", err))
}

func TestDate_MonthStart(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     string
		Expected string
	}

	cases := []testCase{
		{Date: "2020-02-16", Expected: "2020-02-01"},
		{Date: "2021-02-16", Expected: "2021-02-01"},
		{Date: "2023-01-01", Expected: "2023-01-01"},
	}
	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Date, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d, err := date.FromString(tc.Date)
			assert.Nil(err)
			expected, err := date.FromString(tc.Expected)
			assert.Nil(err)

			shifted := d.MonthStart()
			assert.Equal(expected, shifted)
		})
	}
}

func TestDate_MonthEnd(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     string
		Expected string
	}

	cases := []testCase{
		{Date: "2020-02-16", Expected: "2020-02-29"},
		{Date: "2021-02-16", Expected: "2021-02-28"},
		{Date: "2023-01-01", Expected: "2023-01-31"},
	}
	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Date, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d, err := date.FromString(tc.Date)
			assert.Nil(err)
			expected, err := date.FromString(tc.Expected)
			assert.Nil(err)

			shifted := d.MonthEnd()
			assert.Equal(expected, shifted)
		})
	}
}

func TestDate_Before(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d1 := date.Date{Year: 2020, Month: time.May, Day: 11}
	d2 := date.Date{Year: 2022, Month: time.May, Day: 11}
	assert.True(d1.Before(d2))
	assert.False(d2.Before(d1))
	assert.False(d2.Before(d2))
	assert.False(d1.Before(d1))

	d1 = date.Date{Year: 2022, Month: time.April, Day: 11}
	d2 = date.Date{Year: 2022, Month: time.May, Day: 11}
	assert.True(d1.Before(d2))
	assert.False(d2.Before(d1))

	d1 = date.Date{Year: 2022, Month: time.April, Day: 11}
	d2 = date.Date{Year: 2022, Month: time.April, Day: 12}
	assert.True(d1.Before(d2))
	assert.False(d2.Before(d1))
}

func TestDate_After(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d1 := date.Date{Year: 2023, Month: time.July, Day: 27}
	d2 := date.Date{Year: 2018, Month: time.January, Day: 1}
	assert.True(d1.After(d2))
	assert.False(d2.After(d1))
	assert.False(d2.After(d2))
	assert.False(d1.After(d1))
}

func TestDate_Equal(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d1 := date.Date{Year: 2023, Month: time.July, Day: 27}
	d2 := date.Date{Year: 2018, Month: time.January, Day: 1}
	d3 := date.Date{Year: 2023, Month: time.July, Day: 27}
	assert.True(d1.Equal(d1))
	assert.True(d2.Equal(d2))
	assert.True(d1.Equal(d3))
	assert.False(d2.Equal(d1))
	assert.False(d1.Equal(d2))
}

func TestDate_Compare(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d1 := date.Date{Year: 2023, Month: time.July, Day: 27}
	d2 := date.Date{Year: 2018, Month: time.January, Day: 1}
	d3 := date.Date{Year: 2023, Month: time.July, Day: 27}
	d4 := date.Date{Year: 2023, Month: time.August, Day: 27}
	assert.Equal(0, d1.Compare(d1))
	assert.Equal(0, d2.Compare(d2))
	assert.Equal(0, d1.Compare(d3))
	assert.Equal(-1, d2.Compare(d1))
	assert.Equal(1, d1.Compare(d2))
	assert.Equal(-1, d1.Compare(d4))
}

func TestDate_IsZero(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d1 := date.Date{}
	d2 := date.Date{Year: 2006, Month: time.May, Day: 25}
	d3 := date.Date{Year: 2006}
	assert.True(d1.IsZero())
	assert.False(d2.IsZero())
	assert.False(d3.IsZero())
}

func TestDate_ToTime(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d := date.Date{Year: 2006, Month: time.February, Day: 16}
	converted := d.ToTime()
	expected := time.Time(time.Date(2006, time.February, 16, 0, 0, 0, 0, time.UTC))
	assert.Equal(expected, converted)

	tz, err := time.LoadLocation("America/Chicago")
	assert.Nil(err)
	converted = d.ToTime(date.OptConvertTimezone(tz))
	expected = time.Time(time.Date(2006, time.February, 16, 0, 0, 0, 0, tz))
	assert.Equal(expected, converted)
}

func TestDate_Date(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d := date.Date{Year: 2006, Month: time.February, Day: 16}
	year, month, day := d.Date()
	assert.Equal(2006, year)
	assert.Equal(time.February, month)
	assert.Equal(16, day)
}

func TestDate_ISOWeek(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d := date.Date{Year: 2006, Month: time.February, Day: 16}
	year, week := d.ISOWeek()
	assert.Equal(2006, year)
	assert.Equal(7, week)
}

func TestDate_Weekday(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     date.Date
		Expected time.Weekday
	}

	cases := []testCase{
		{Date: date.Date{Year: 2023, Month: time.January, Day: 1}, Expected: time.Sunday},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 2}, Expected: time.Monday},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 3}, Expected: time.Tuesday},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 4}, Expected: time.Wednesday},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 5}, Expected: time.Thursday},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 6}, Expected: time.Friday},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 7}, Expected: time.Saturday},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 8}, Expected: time.Sunday},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Date.String(), func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			weekday := tc.Date.Weekday()
			assert.Equal(tc.Expected, weekday)
		})
	}
}

func TestDate_YearDay(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     date.Date
		Expected int
	}

	cases := []testCase{
		{Date: date.Date{Year: 2022, Month: time.December, Day: 31}, Expected: 365},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 1}, Expected: 1},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 5}, Expected: 5},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 6}, Expected: 6},
		{Date: date.Date{Year: 2023, Month: time.January, Day: 8}, Expected: 8},
		{Date: date.Date{Year: 2024, Month: time.December, Day: 31}, Expected: 366},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Date.String(), func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			yearDay := tc.Date.YearDay()
			assert.Equal(tc.Expected, yearDay)
		})
	}
}

func TestDate_MarshalText(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Name     string
		Date     date.Date
		Expected string
	}

	cases := []testCase{
		{Name: "Remote past", Date: date.Date{Year: 1997, Month: time.July, Day: 15}, Expected: "1997-07-15"},
		{Name: "Recent past", Date: date.Date{Year: 2020, Month: time.February, Day: 20}, Expected: "2020-02-20"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			asBytes, err := tc.Date.MarshalText()
			assert.Nil(err)
			assert.Equal(tc.Expected, string(asBytes))
		})
	}
}

func TestDate_MarshalJSON(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Name     string
		Date     *date.Date
		Expected string
	}

	cases := []testCase{
		{Name: "Remote past", Date: &date.Date{Year: 1997, Month: time.July, Day: 15}, Expected: `"1997-07-15"`},
		{Name: "Recent past", Date: &date.Date{Year: 2020, Month: time.February, Day: 20}, Expected: `"2020-02-20"`},
		{Name: "Unset", Date: nil, Expected: "null"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			asBytes, err := json.Marshal(tc.Date)
			assert.Nil(err)
			assert.Equal(tc.Expected, string(asBytes))
		})
	}
}

func TestDate_UnmarshalText(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Input []byte
		Date  date.Date
		Error string
	}

	cases := []testCase{
		{Input: []byte(`x`), Error: `parsing time "x" as "2006-01-02": cannot parse "x" as "2006"`},
		{Input: []byte(`10`), Error: `parsing time "10" as "2006-01-02": cannot parse "10" as "2006"`},
		{Input: []byte("01/26/2018"), Error: `parsing time "01/26/2018" as "2006-01-02": cannot parse "01/26/2018" as "2006"`},
		{Input: []byte("1997-07-15"), Date: date.Date{Year: 1997, Month: time.July, Day: 15}},
		{Input: []byte("2020-02-20"), Date: date.Date{Year: 2020, Month: time.February, Day: 20}},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(string(tc.Input), func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d := date.Date{}
			err := d.UnmarshalText(tc.Input)
			if err != nil {
				assert.Equal(tc.Error, fmt.Sprintf("%v", err))
				assert.Equal(date.Date{}, d)
			} else {
				assert.Equal("", tc.Error)
				assert.Equal(tc.Date, d)
			}
		})
	}
}

func TestDate_UnmarshalJSON(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Input []byte
		Date  date.Date
		Error string
	}

	cases := []testCase{
		{Input: []byte(`x`), Error: "invalid character 'x' looking for beginning of value"},
		{Input: []byte(`10`), Error: "json: cannot unmarshal number into Go value of type string"},
		{Input: []byte(`"abc"`), Error: `parsing time "abc" as "2006-01-02": cannot parse "abc" as "2006"`},
		{Input: []byte(`"01/26/2018"`), Error: `parsing time "01/26/2018" as "2006-01-02": cannot parse "01/26/2018" as "2006"`},
		{Input: []byte(`"1997-07-15"`), Date: date.Date{Year: 1997, Month: time.July, Day: 15}},
		{Input: []byte(`"2020-02-20"`), Date: date.Date{Year: 2020, Month: time.February, Day: 20}},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(string(tc.Input), func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			d := date.Date{}
			err := json.Unmarshal(tc.Input, &d)
			if err != nil {
				assert.Equal(tc.Error, fmt.Sprintf("%v", err))
				assert.Equal(date.Date{}, d)
			} else {
				assert.Equal("", tc.Error)
				assert.Equal(tc.Date, d)
			}
		})
	}
}

func TestDate_Scan(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	// Wrong type
	d := date.Date{}
	err := d.Scan(1)
	assert.NotNil(err)
	assert.Equal("incompatible type for Date; type=int", fmt.Sprintf("%v", err))
	assert.Equal(date.Date{}, d)

	// Time but not date
	d = date.Date{}
	tz, err := time.LoadLocation("America/Los_Angeles")
	assert.Nil(err)
	src := time.Date(2001, time.August, 4, 11, 10, 55, 0, tz)
	err = d.Scan(src)
	assert.NotNil(err)
	assert.Equal("timestamp contains more than just date information; 2001-08-04T11:10:55-07:00", fmt.Sprintf("%v", err))
	assert.Equal(date.Date{}, d)

	// Happy path
	d = date.Date{}
	src = time.Date(1991, time.April, 26, 0, 0, 0, 0, time.UTC)
	err = d.Scan(src)
	assert.Nil(err)
	expected := date.Date{Year: 1991, Month: time.April, Day: 26}
	assert.Equal(expected, d)
}

func TestDate_Value(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	d := date.Date{Year: 1991, Month: time.April, Day: 26}
	v, err := d.Value()
	assert.Nil(err)
	expected := time.Date(1991, time.April, 26, 0, 0, 0, 0, time.UTC)
	assert.Equal(expected, v)
}

func TestDate_String(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     date.Date
		Expected string
	}

	cases := []testCase{
		{Date: date.Date{Year: 2020, Month: time.May, Day: 11}, Expected: "2020-05-11"},
		{Date: date.Date{Year: 2022, Month: time.January, Day: 31}, Expected: "2022-01-31"},
		{Date: date.Date{Year: 1999, Month: time.December, Day: 24}, Expected: "1999-12-24"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Expected, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			assert.Equal(tc.Expected, tc.Date.String())
		})
	}
}

func TestDate_GoString(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Date     date.Date
		Expected string
	}

	cases := []testCase{
		{Date: date.Date{Year: 2020, Month: time.May, Day: 11}, Expected: "date.NewDate(2020, time.May, 11)"},
		{Date: date.Date{Year: 2022, Month: time.January, Day: 31}, Expected: "date.NewDate(2022, time.January, 31)"},
		{Date: date.Date{Year: 1999, Month: time.December, Day: 24}, Expected: "date.NewDate(1999, time.December, 24)"},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		base.Run(tc.Expected, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			assert.Equal(tc.Expected, tc.Date.GoString())
		})
	}
}
