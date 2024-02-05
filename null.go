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
	"database/sql/driver"
)

// NullDate is a `Date` that can be null.
type NullDate struct {
	Date  Date
	Valid bool
}

// Scan implements `sql.Scanner`; it unmarshals nullable values of the type
// `time.Time` onto the current `NullDate` struct.
func (nd *NullDate) Scan(value any) error {
	if value == nil {
		nd.Date = Date{}
		nd.Valid = false
		return nil
	}

	err := nd.Date.Scan(value)
	if err != nil {
		return err
	}

	nd.Valid = true
	return nil
}

// Value implements `driver.Valuer`; it marshals the value to a `time.Time`
// (or `nil`) to be serialized into the database.
func (nd NullDate) Value() (driver.Value, error) {
	if !nd.Valid {
		return nil, nil
	}

	return nd.Date.Value()
}
