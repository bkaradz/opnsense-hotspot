package printing

import (
	"fmt"
	"time"
)

// padString pads a string with spaces on the right to the specified width.
func padString(str string, width int, padChar string) string {
	if len(str) >= width {
		return str
	}
	padding := ""
	for i := 0; i < width-len(str); i++ {
		padding += padChar
	}
	return str + padding
}

// itemToString converts an Item to a formatted string.
func itemToString(item Item) string {
	const rightCols = 30
	const leftCols = 15

	left := padString(item.Key, leftCols, " ")
	right := padString(item.Value, rightCols, " ")
	return fmt.Sprintf("%s%s", left, right)
}

// FormatTime converts a duration in seconds into a human-readable string,
// showing only non-zero components and using correct singular/plural forms.
func FormatTime(seconds int64) string {
	// Convert seconds to a time.Duration
	duration := time.Duration(seconds) * time.Second

	// --- Custom Formatting to eliminate redundancy and use correct pluralization ---

	// Calculate individual components
	hours := duration / time.Hour
	minutes := (duration % time.Hour) / time.Minute
	remainingSeconds := (duration % time.Minute) / time.Second

	// Build the string slice with only non-zero parts
	parts := make([]string, 0)

	// 1. Handle Hours
	if hours > 0 {
		unit := "hour"
		if hours > 1 {
			unit += "s" // Add 's' for plural
		}
		parts = append(parts, fmt.Sprintf("%d %s", hours, unit))
	}

	// 2. Handle Minutes
	if minutes > 0 {
		unit := "minute"
		if minutes > 1 {
			unit += "s" // Add 's' for plural
		}
		parts = append(parts, fmt.Sprintf("%d %s", minutes, unit))
	}

	// 3. Handle Seconds
	// Always include seconds if the total duration is less than a minute
	// OR if seconds are the remaining component, AND the count is greater than 0.
	if remainingSeconds > 0 || (seconds > 0 && seconds < 60) {
		unit := "second"
		if remainingSeconds > 1 {
			unit += "s" // Add 's' for plural
		}
		// If the input is exactly 1 second, remainingSeconds will be 1, so unit is 'second'.
		// If the input is 0, this block is skipped.
		parts = append(parts, fmt.Sprintf("%d %s", remainingSeconds, unit))
	}

	// Handle the case of zero seconds input
	if seconds == 0 {
		parts = append(parts, "0 seconds")
	}

	// Join the parts with a comma and space
	customFormat := ""
	for i, part := range parts {
		// Use "and" for the last component if there are more than two components
		if i == len(parts)-1 && len(parts) > 1 {
			customFormat += " and "
		} else if i > 0 {
			customFormat += ", "
		}
		customFormat += part
	}

	return customFormat
}
