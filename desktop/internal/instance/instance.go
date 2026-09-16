// ai-generated
package instance

import (
	"context"
)

// IsInstanceExisting checks whether ClipSync is already running.
// If an existing instance is found, it notifies that instance to open its GUI and returns true.
// If this is the primary instance, it acquires the single-instance lock and returns false.
func IsInstanceExisting() bool {
	return checkInstance()
}

// StartListener starts a background listener for wake-up signals from secondary instances.
func StartListener(ctx context.Context, showGUI func()) {
	startListener(ctx, showGUI)
}

// Cleanup releases any locks or OS resources held by the primary instance.
func Cleanup() {
	cleanup()
}
