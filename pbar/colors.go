package pbar

import (
	"fmt"
	"os"
	"strings"
)

// AnsiColors maps color names to their ANSI escape codes.
var AnsiColors = map[string]string{
	"reset":   "\x1b[0m",
	"black":   "\x1b[30m",
	"red":     "\x1b[31m",
	"green":   "\x1b[32m",
	"yellow":  "\x1b[33m",
	"blue":    "\x1b[34m",
	"magenta": "\x1b[35m",
	"cyan":    "\x1b[36m",
	"white":   "\x1b[37m",
}

// GetColorCode returns the ANSI escape code for a given color name.
// If the color name is invalid, it prints an error to os.Stderr and returns an empty string.
func GetColorCode(colorName string) string {
	if colorName == "" {
		return ""
	}

	if code, ok := AnsiColors[strings.ToLower(colorName)]; ok {
		return code
	}

	fmt.Fprintf(os.Stderr, "Warning: Invalid color name '%s'. Available colors: %s\n", colorName, GetAvailableColors())
	return "" // Fallback to no color
}

// GetAvailableColors returns a comma-separated string of available color names.
func GetAvailableColors() string {
	colors := make([]string, 0, len(AnsiColors)-1) // Exclude 'reset'
	for colorName := range AnsiColors {
		if colorName != "reset" {
			colors = append(colors, colorName)
		}
	}
	return strings.Join(colors, ", ")
}
