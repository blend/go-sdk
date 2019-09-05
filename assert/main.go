package assert

import (
	"os"
	"testing"
)

// Main wraps a testing.M.
func Main(m *testing.M) {
	Started()
	var statusCode int
	func() {
		defer ReportRate()
		statusCode = m.Run()
	}()
	os.Exit(statusCode)
}
