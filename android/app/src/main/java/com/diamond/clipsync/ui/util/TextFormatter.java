package com.diamond.clipsync.ui.util;

public class TextFormatter {

    /**
     * Formats raw copied text cleanly for UI display:
     * - Converts tab characters to standard spaces (prevents '[]' rendering)
     * - Normalizes CRLF / CR to LF
     * - Cleans up invisible control characters while preserving valid formatting
     */
    public static String cleanDisplayText(String input) {
        if (input == null || input.isEmpty()) {
            return "";
        }

        // Replace tabs with 4 spaces to avoid missing glyph box '[]'
        String text = input.replace("\t", "    ");

        // Normalize carriage returns
        text = text.replace("\r\n", "\n").replace("\r", "\n");

        // Filter out non-printable ASCII control characters except \n
        StringBuilder sb = new StringBuilder(text.length());
        for (int i = 0; i < text.length(); i++) {
            char c = text.charAt(i);
            if (c == '\n' || (c >= 32 && c != 127)) {
                sb.append(c);
            } else if (Character.isWhitespace(c)) {
                sb.append(' ');
            }
        }

        return sb.toString();
    }
}
