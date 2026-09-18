package utils

import (
	"testing"
)

func TestCleanDisplayText_NormalizeLineEndings(t *testing.T) {
	input := "line1\r\nline2\rline3\nline4"
	expected := "line1\nline2\nline3\nline4"
	got := CleanDisplayText(input)
	if got != expected {
		t.Fatalf("CleanDisplayText() mismatch:\n got: %q\nwant: %q", got, expected)
	}
}

func TestCleanDisplayText_ExpandTabs(t *testing.T) {
	input := "\tdef main():\n\t\tprint('hello')"
	expected := "    def main():\n        print('hello')"
	got := CleanDisplayText(input)
	if got != expected {
		t.Fatalf("CleanDisplayText() tab expansion mismatch:\n got: %q\nwant: %q", got, expected)
	}
}

func TestCleanDisplayText_StripControlCharsAndInvisibleGlyphs(t *testing.T) {
	// ASCII null \x00, bell \x07, BOM \uFEFF, zero-width space \u200B
	input := "Clean\x00Text\x07With\uFEFFInvisible\u200BChars"
	expected := "CleanTextWithInvisibleChars"
	got := CleanDisplayText(input)
	if got != expected {
		t.Fatalf("CleanDisplayText() control char removal mismatch:\n got: %q\nwant: %q", got, expected)
	}
}

func TestCleanDisplayText_EmptyString(t *testing.T) {
	got := CleanDisplayText("")
	if got != "" {
		t.Fatalf("Expected empty string, got: %q", got)
	}
}
