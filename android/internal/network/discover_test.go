package network

import (
	"testing"
)

func TestGetCurrIPS(t *testing.T) {
	sig := GetCurrIPS()
	// In test environment, it could be empty if no non-loopback active interface, or a comma-separated list of IPs.
	// Function must not panic and must return string.
	t.Logf("Detected network signature: %q", sig)
}
