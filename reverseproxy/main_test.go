/*

Copyright (c) 2021 - Present. Blend Labs, Inc. All rights reserved
Blend Confidential - Restricted

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
