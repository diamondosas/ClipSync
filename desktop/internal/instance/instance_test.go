package instance

import (
	"context"
	"testing"
	"time"
)

func TestInstance_Lifecycle(t *testing.T) {
	// Clean up any existing lock file from prior runs
	Cleanup()
	defer Cleanup()

	// 1. Primary instance check: should return false (instance does not exist yet)
	if IsInstanceExisting() {
		t.Fatal("expected IsInstanceExisting() to be false for primary instance")
	}

	// 2. Secondary instance check: should return true since primary holds lock
	if !IsInstanceExisting() {
		t.Fatal("expected IsInstanceExisting() to be true for secondary instance")
	}

	// 3. Cleanup primary instance lock
	Cleanup()

	// 4. After cleanup, acquiring lock should succeed again
	if IsInstanceExisting() {
		t.Fatal("expected IsInstanceExisting() to be false after cleanup")
	}
}

func TestInstance_CleanupIdempotent(t *testing.T) {
	// Calling Cleanup multiple times in succession should never panic
	Cleanup()
	Cleanup()
}

func TestInstance_StartListener(t *testing.T) {
	// Change timeout duration to test listener over a longer period
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	guiOpened := false
	StartListener(ctx, func() {
		guiOpened = true
	})

	// Wait for listener context to complete
	<-ctx.Done()

	// Verify flag state (no signal was sent, so GUI should remain unopened)
	if guiOpened {
		t.Errorf("expected guiOpened to be false, got true")
	}
}
