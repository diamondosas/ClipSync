package ping_test

import (
	"testing"
	"clipsync/internal/ping"
)

func TestPingIPS(t *testing.T) {
	// Table-driven test to handle multiple scenarios
	tests := []struct {
		name    string
		ips     []string
		wantAtLeastOne bool
	}{
		{
			name: "Localhost is reachable",
			ips:  []string{"127.0.0.1"},
			wantAtLeastOne: true,
		},
		{	
			name: "Nil input",
			ips:  nil,
			wantAtLeastOne: false,
		},
		{
			name: "Empty input",
			ips:  []string{},
			wantAtLeastOne: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ping.PingIPS(tt.ips)
			
			if tt.wantAtLeastOne && len(got) == 0 {
				t.Errorf("%s: expected at least one reachable IP, got 0", tt.name)
			}
			
			if !tt.wantAtLeastOne && len(got) > 0 {
				t.Errorf("%s: expected 0 reachable IPs, got %v", tt.name, got)
			}
		})
	}
}