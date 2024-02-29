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

// Package date provides tools for working with dates, extending the standard
// library `time` package.
//
// This package provides helpers for converting from a full `time.Time{}` to
// a `Date{}` and back, providing validation along the way. Many methods from
// `time.Time{}` are also provided as equivalents here (`After()`, `Before()`,
// `Sub()`, etc.). Additionally, custom serialization methods are provided both
// for JSON and SQL.
//
// The Go standard library contains no native type for dates without times.
// Instead, common convention is to use a `time.Time{}` with only the year,
// month, and day set. For example, this convention is followed when a
// timestamp of the form YYYY-MM-DD is parsed via
// `time.Parse(time.DateOnly, value)`.
package date
