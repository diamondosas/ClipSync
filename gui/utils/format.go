// ai-generated
package utils

import (
	"strings"
	"unicode"
)

// CleanDisplayText formats raw text for clean UI rendering in Gio.
// It converts tabs (\t) into 4 spaces to preserve indentation without
// rendering "tofu" boxes ([]), normalizes newlines, strips non-printable
// control characters (ASCII < 32 except newline), and removes invisible
// zero-width unicode characters.
func CleanDisplayText(raw string) string {
	if raw == "" {
		return ""
	}

	// 1. Normalize line breaks: \r\n -> \n, standalone \r -> \n
	s := strings.ReplaceAll(raw, "\r\n", "\n")
	s = strings.ReplaceAll(s, "\r", "\n")

	// 2. Expand tabs to 4 spaces for clean visual alignment
	s = strings.ReplaceAll(s, "\t", "    ")

	// 3. Filter unprintable control characters and invisible glyphs
	var b strings.Builder
	b.Grow(len(s))

	for _, r := range s {
		switch {
		case r == '\n':
			b.WriteRune('\n')
		case r == '\v' || r == '\f':
			// Vertical tab, form feed -> treat as line break
			b.WriteRune('\n')
		case r < 32:
			// Strip all other C0 ASCII control characters (null, bell, backspace, escape, etc.)
			continue
		case r == 0x7F:
			// ASCII DEL
			continue
		case r == 0xA0:
			// Non-breaking space -> regular space
			b.WriteRune(' ')
		case r == 0xFEFF || r == 0x200B || r == 0x200C || r == 0x200D:
			// Byte Order Mark (BOM), Zero-width space/non-joiner/joiner
			continue
		case unicode.IsControl(r):
			// Strip other unicode control characters (C1 controls, etc.)
			continue
		default:
			b.WriteRune(r)
		}
	}

	return b.String()
}
