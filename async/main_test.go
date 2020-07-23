package async

import (
	"os"
	"testing"
)

// TestMain is the testing entrypoint.
func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
