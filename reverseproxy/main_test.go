/*

Copyright (c) 2023 - Present. Blend Labs, Inc. All rights reserved
Use of this source code is governed by a MIT license that can be found in the LICENSE file.

*/

package reverseproxy

import (
	"os"
	"testing"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
