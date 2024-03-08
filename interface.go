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

// NOTE: Ensure that
// - `Date` satisfies `commonDateTimeValue`.
// - `time.Time` satisfies `commonDateTimeValue`.
// - `*Date` satisfies `commonDateTimePointer`.
// - `*time.Time` satisfies `commonDateTimePointer`.
var (
	_ commonDateTimeValue[Date]      = Date{}
	_ commonDateTimeValue[time.Time] = time.Time{}
	_ commonDateTimePointer          = (*Date)(nil)
	_ commonDateTimePointer          = (*time.Time)(nil)
)

type commonDateTimeValue[T any] interface {
	After(u T) bool
	Before(u T) bool
	Compare(u T) int
	Date() (year int, month time.Month, day int)
	// Day() int
	Equal(u T) bool
	Format(layout string) string
	GoString() string
	ISOWeek() (year, week int)
	IsZero() bool
	MarshalJSON() ([]byte, error)
	MarshalText() ([]byte, error)
	// Month() time.Month
	String() string
	Weekday() time.Weekday
	// Year() int
	YearDay() int
}

type commonDateTimePointer interface {
	UnmarshalJSON(data []byte) error
	UnmarshalText(data []byte) error
}
