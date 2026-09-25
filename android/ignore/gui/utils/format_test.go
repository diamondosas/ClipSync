// ai-generated
package utils

import (
	"testing"
)

func TestCleanDisplayText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "converts tabs to 4 spaces",
			input:    "func main() {\n\tprintln()\n}",
			expected: "func main() {\n    println()\n}",
		},
		{
			name:     "normalizes CRLF and CR",
			input:    "line1\r\nline2\rline3\nline4",
			expected: "line1\nline2\nline3\nline4",
		},
		{
			name:     "strips non-printable ASCII control chars",
			input:    "hello\x00\x07\x08\x1bworld",
			expected: "helloworld",
		},
		{
			name:     "replaces non-breaking space",
			input:    "hello\u00a0world",
			expected: "hello world",
		},
		{
			name:     "strips zero-width and BOM characters",
			input:    "\ufeffhello\u200b\u200c\u200dworld",
			expected: "helloworld",
		},
		{
			name:     "converts vertical tabs and form feeds to newlines",
			input:    "part1\vpart2\fpart3",
			expected: "part1\npart2\npart3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CleanDisplayText(tt.input)
			if got != tt.expected {
				t.Errorf("CleanDisplayText(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
