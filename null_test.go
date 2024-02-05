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

func TestNullDate_Scan(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	// Wrong type
	nd := date.NullDate{}
	err := nd.Scan(1)
	assert.NotNil(err)
	assert.Equal("incompatible type for Date; type=int", fmt.Sprintf("%v", err))
	assert.Equal(date.NullDate{}, nd)

	// Happy path: nil
	nd = date.NullDate{}
	err = nd.Scan(nil)
	assert.Nil(err)
	assert.Equal(date.NullDate{}, nd)

	// Happy path: value
	nd = date.NullDate{}
	src := time.Date(1991, time.April, 26, 0, 0, 0, 0, time.UTC)
	err = nd.Scan(src)
	assert.Nil(err)
	expected := date.NullDate{Date: date.Date{Year: 1991, Month: time.April, Day: 26}, Valid: true}
	assert.Equal(expected, nd)
}

func TestNullDate_Value(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	// Not valid
	nd := date.NullDate{}
	v, err := nd.Value()
	assert.Nil(err)
	assert.Nil(v)

	// Valid
	nd = date.NullDate{Date: date.Date{Year: 1991, Month: time.April, Day: 26}, Valid: true}
	v, err = nd.Value()
	assert.Nil(err)
	expected := time.Date(1991, time.April, 26, 0, 0, 0, 0, time.UTC)
	assert.Equal(expected, v)
}
