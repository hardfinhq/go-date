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
	"fmt"
	"testing"
	"time"

	testifyrequire "github.com/stretchr/testify/require"

	date "github.com/hardfinhq/go-date"
)

func TestToday(base *testing.T) {
	base.Parallel()

	type testCase struct {
		Timezone string
		Now      string
		Date     date.Date
	}

	cases := []testCase{
		{
			Timezone: "UTC",
			Now:      "2020-05-11T07:10:55.209309302Z",
			Date:     date.Date{Year: 2020, Month: time.May, Day: 11},
		},
		{
			Timezone: "UTC",
			Now:      "2022-01-31T00:00:00.000Z",
			Date:     date.Date{Year: 2022, Month: time.January, Day: 31},
		},
		{
			Timezone: "America/Los_Angeles",
			Now:      "2022-01-31T00:00:00.000Z",
			Date:     date.Date{Year: 2022, Month: time.January, Day: 30},
		},
	}

	for i := range cases {
		// NOTE: Assign to loop-local (instead of declaring the `tc` variable in
		//       `range`) to avoid capturing reference to loop variable.
		tc := cases[i]
		description := fmt.Sprintf("%s:%s", tc.Now, tc.Timezone)
		base.Run(description, func(t *testing.T) {
			t.Parallel()
			assert := testifyrequire.New(t)

			now, err := time.Parse(time.RFC3339Nano, tc.Now)
			assert.Nil(err)
			tz, err := time.LoadLocation(tc.Timezone)
			assert.Nil(err)

			d := date.Today(
				date.OptTodayTimezone(tz),
				date.OptTodayNow(now),
			)
			assert.Equal(tc.Date, d)
		})
	}
}
