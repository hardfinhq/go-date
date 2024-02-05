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
	"errors"
	"testing"

	testifyrequire "github.com/stretchr/testify/require"
)

// NOTE: This test file is in `package date` so that it can access
//       `mustNil()`, which is intentionally not exported.

func Test_mustNil(t *testing.T) {
	t.Parallel()
	assert := testifyrequire.New(t)

	// Happy path
	assert.NotPanics(func() {
		mustNil(nil)
	})

	// Sad path
	err := errors.New("definitely bad")
	assert.PanicsWithValue(err, func() {
		mustNil(err)
	})
}
