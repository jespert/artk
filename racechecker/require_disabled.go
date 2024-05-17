//go:build !race

package racechecker

// Require skips the test if the race condition checker is not enabled.
// It might save some cycles, but it mostly serves as live documentation.
func Require(t testingT) {
	t.Helper()
	t.Skip("This test requires the race checker. Skipping.")
}
