package pbar

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"os"
	"path/filepath"
	"strconv"
)

const (
	maxThroughputHistorySize = 10
	defaultStyle             = "classic"
	defaultWidth             = 50
	stateFileName            = ".pbar.state"
)

func getSpinnerState() int {
	stateFile := filepath.Join(os.TempDir(), stateFileName)
	data, err := os.ReadFile(stateFile)
	if err != nil {
		// If the file doesn't exist, create it with initial state 0
		if os.IsNotExist(err) {
			os.WriteFile(stateFile, []byte("0"), 0644)
		}
		return 0
	}
	state, _ := strconv.Atoi(string(data))
	return state
}

func setSpinnerState(state int) {
	stateFile := filepath.Join(os.TempDir(), stateFileName)
	os.WriteFile(stateFile, []byte(strconv.Itoa(state)), 0644)
}

var spinnerChars = []string{"|", "/", "-", "\\"}
var brailleSpinnerChars = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

var brailleChars = []string{
	" ", // 0/8
	"⠁", // 1/8
	"⠃", // 2/8
	"⠇", // 3/8
	"⠏", // 4/8
	"⠟", // 5/8
	"⠿", // 6/8
	"⡿", // 7/8
	"⣿", // 8/8
}

var validStyles = map[string]bool{
	"classic": true,
	"block":   true,
	"spinner": true,
	"arrow":   true,
	"braille":         true,
	"braille-spinner": true,
	"custom":          true,
}

// Bar represents a progress bar.
type Bar struct {
	Total             int
	Current           int
	Width             int
	Style             string
	ColorBar          string // Now stores ANSI escape code directly
	ColorText         string // Now stores ANSI escape code directly
	Finished          bool
	StartTime         time.Time
	LastUpdateTime    time.Time
	ThroughputHistory []float64
	CustomChars       string
	Message           string
	spinnerState      int
	TestMode          bool
}

// Render generates the string representation of the progress bar.
func (b *Bar) Render() string {
	// Update LastUpdateTime
	b.LastUpdateTime = time.Now()

	// Ensure Total is not negative
	if b.Total < 0 {
		b.Total = 0
	}

	percent := float64(b.Current) / float64(b.Total)

	// Handle total being zero to prevent NaN or Inf
	if b.Total == 0 {
		if b.Current == 0 {
			percent = 0 // 0/0 is 0%
		} else {
			percent = 1 // X/0 (X>0) is 100%
		}
	}

	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	percentString := fmt.Sprintf("%d%%", int(percent*100))

	var metadataString string
	var throughputStr, etaStr string // Always initialize

	if !b.StartTime.IsZero() {
		// Calculate elapsed time
		elapsedTime := time.Since(b.StartTime)
		elapsedTimeStr := formatDuration(elapsedTime)

		// Calculate throughput and ETA only if not indeterminate and elapsed time is non-zero
		isIndeterminate := b.Style == "spinner" || b.Style == "braille-spinner"
		if !isIndeterminate && elapsedTime.Seconds() > 0 {
			var currentThroughput float64
			if b.Total > 0 && b.Current > 0 {
				currentThroughput = float64(b.Current) / elapsedTime.Seconds()
			} else {
				currentThroughput = 0
			}

			// Update throughput history (simple moving average for now)
			b.ThroughputHistory = append(b.ThroughputHistory, currentThroughput)
			if len(b.ThroughputHistory) > maxThroughputHistorySize { // Keep last 10 samples
				b.ThroughputHistory = b.ThroughputHistory[1:]
			}

			// Calculate average throughput
			var totalThroughput float64
			for _, t := range b.ThroughputHistory {
				totalThroughput += t
			}
			averageThroughput := totalThroughput / float64(len(b.ThroughputHistory))

			throughputStr = fmt.Sprintf(" %.2f it/s", averageThroughput)

			// Calculate ETA based on average throughput
			remainingItems := float64(b.Total - b.Current)
			if averageThroughput > 0 {
				eta := time.Duration(remainingItems / averageThroughput * float64(time.Second))
				etaStr = fmt.Sprintf(" ETA %s", formatDuration(eta))
			} else {
				etaStr = " ETA Inf"
			}
		}
		metadataString = fmt.Sprintf(" Elapsed %s%s%s", elapsedTimeStr, throughputStr, etaStr)
	}

	if b.Finished {
		if b.Message != "" {
			metadataString = fmt.Sprintf("%s %s", metadataString, b.Message)
		}
		result := fmt.Sprintf("[✔] 100%%%s", metadataString)
		result = "\r" + result + "\x1b[K"
		return result
	}

	isIndeterminate := b.Style == "spinner" || b.Style == "braille-spinner"
	if isIndeterminate {
		if !b.TestMode {
			b.spinnerState = getSpinnerState()
		}
		var char string
		switch b.Style {
		case "braille-spinner":
			char = brailleSpinnerChars[b.spinnerState%len(brailleSpinnerChars)]
		default:
			char = spinnerChars[b.spinnerState%len(spinnerChars)]
		}
		b.spinnerState++
		if !b.TestMode {
			setSpinnerState(b.spinnerState)
		}
		result := fmt.Sprintf("[%s]%s", char, metadataString)
		if b.ColorText != "" {
			result = fmt.Sprintf("[%s%s%s]%s", b.ColorText, char, "\x1b[0m", metadataString)
		}
		result = "\r" + result + "\x1b[K"
		return result
	}

	style := b.Style
	if _, ok := validStyles[style]; !ok {
		style = defaultStyle
	}

	var barString string
	switch style {
	case "spinner":
		if !b.TestMode {
			b.spinnerState = getSpinnerState()
		}
		char := spinnerChars[b.spinnerState%len(spinnerChars)]
		b.spinnerState++
		if !b.TestMode {
			setSpinnerState(b.spinnerState)
		}
		barString = fmt.Sprintf("[%s]", char)
		if b.ColorText != "" {
			barString = fmt.Sprintf("[%s%s%s]", b.ColorText, char, "\x1b[0m")
		}
	case "braille-spinner":
		if !b.TestMode {
			b.spinnerState = getSpinnerState()
		}
		char := brailleSpinnerChars[b.spinnerState%len(brailleSpinnerChars)]
		b.spinnerState++
		if !b.TestMode {
			setSpinnerState(b.spinnerState)
		}
		barString = fmt.Sprintf("[%s]", char)
		if b.ColorText != "" {
			barString = fmt.Sprintf("[%s%s%s]", b.ColorText, char, "\x1b[0m")
		}
	case "block":
		barString = b.renderBar("█", " ", b.ColorBar)
	case "classic":
		barString = b.renderBar("#", "-", b.ColorBar)
	case "arrow":
		barString = b.renderArrowBar(b.ColorBar)
	case "braille":
		barString = b.renderBrailleBar(b.ColorBar)
	case "custom":
		filledChar := "#" // Default
		emptyChar := "-"  // Default

		if len(b.CustomChars) > 0 {
			filledChar = string(b.CustomChars[0])
			if len(b.CustomChars) > 1 {
				emptyChar = string(b.CustomChars[1])
			} else {
				emptyChar = filledChar // If only one char, use it for both
			}
		}
		barString = b.renderBar(filledChar, emptyChar, b.ColorBar)
	}

	if b.ColorText != "" {
		percentString = fmt.Sprintf("%s%s%s", b.ColorText, percentString, "\x1b[0m") // Use reset code directly
	}

	result := fmt.Sprintf("%s %s%s", barString, percentString, metadataString)
	
	// Add carriage return for inline updates
	result = "\r" + result + "\x1b[K"
	
	return result
}

func (b *Bar) renderBar(filledChar, emptyChar, colorCode string) string {
	// Ensure Total is not negative
	if b.Total < 0 {
		b.Total = 0
	}

	// If width is 0 or negative, return an empty bar
	if b.Width <= 0 {
		return "[]"
	}

	percent := float64(b.Current) / float64(b.Total)

	// Handle total being zero to prevent NaN or Inf
	if b.Total == 0 {
		if b.Current == 0 {
			percent = 0 // 0/0 is 0%
		} else {
			percent = 1 // X/0 (X>0) is 100%
		}
	}

	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	filledWidth := int(percent * float64(b.Width))
	if b.Width == 1 && percent > 0 { // Special handling for width 1 and non-zero progress
		filledWidth = 1
	}
	emptyWidth := b.Width - filledWidth

	filled := strings.Repeat(filledChar, filledWidth)
	empty := strings.Repeat(emptyChar, emptyWidth)

	barContent := fmt.Sprintf("[%s%s]", filled, empty)

	if colorCode != "" {
		return fmt.Sprintf("%s%s%s", colorCode, barContent, "\x1b[0m") // Use reset code directly
	}
	return barContent
}

func (b *Bar) renderArrowBar(colorCode string) string {
	// Ensure Total is not negative
	if b.Total < 0 {
		b.Total = 0
	}

	// If width is 0 or negative, return an empty bar
	if b.Width <= 0 {
		return "[]"
	}

	percent := float64(b.Current) / float64(b.Total)
	// Handle total being zero to prevent NaN or Inf
	if b.Total == 0 {
		if b.Current == 0 {
			percent = 0 // 0/0 is 0%
		} else {
			percent = 1 // X/0 (X>0) is 100%
		}
	}

	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	filledWidth := int(percent * float64(b.Width))
	emptyWidth := b.Width - filledWidth

	var filled string
	if filledWidth > 0 {
		filled = strings.Repeat("-", filledWidth-1) + ">"
	} else if filledWidth == 0 { // Handle filledWidth being 0
		filled = ""
	} else { // Handle filledWidth being negative (shouldn't happen with clamping, but for safety)
		filled = ""
	}
	empty := strings.Repeat(" ", emptyWidth)

	barContent := fmt.Sprintf("[%s%s]", filled, empty)

	if colorCode != "" {
		return fmt.Sprintf("%s%s%s", colorCode, barContent, "\x1b[0m")
	}
	return barContent
}

func (b *Bar) renderBrailleBar(colorCode string) string {
	// Ensure Total is not negative
	if b.Total < 0 {
		b.Total = 0
	}

	// If width is 0 or negative, return an empty bar
	if b.Width <= 0 {
		return "[]"
	}

	percent := float64(b.Current) / float64(b.Total)
	// Handle total being zero to prevent NaN or Inf
	if b.Total == 0 {
		if b.Current == 0 {
			percent = 0 // 0/0 is 0%
		} else {
			percent = 1 // X/0 (X>0) is 100%
		}
	}

	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	// Calculate total braille units in the bar
	totalBrailleUnits := b.Width * (len(brailleChars) - 1)
	filledBrailleUnits := int(percent * float64(totalBrailleUnits))

	var barContentBuilder strings.Builder
	barContentBuilder.WriteString("[")

	for i := 0; i < b.Width; i++ {
		// Calculate units for the current character cell
		currentCellUnits := filledBrailleUnits - (i * (len(brailleChars) - 1))

		if currentCellUnits >= (len(brailleChars) - 1) { // Full block
			barContentBuilder.WriteString("⣿")
		} else if currentCellUnits > 0 { // Fractional block
			barContentBuilder.WriteString(brailleChars[currentCellUnits])
		} else {
			barContentBuilder.WriteString(" ") // Empty space
		}
	}
	barContentBuilder.WriteString("]")

	barContent := barContentBuilder.String()

	if colorCode != "" {
		return fmt.Sprintf("%s%s%s", colorCode, barContent, "\x1b[0m")
	}
	return barContent
}

func formatDuration(d time.Duration) string {
	seconds := int(d.Seconds())
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	remainingSeconds := seconds % 60

	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, remainingSeconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, remainingSeconds)
	}
	return fmt.Sprintf("%ds", remainingSeconds)
}

func (b *Bar) Validate() error {
	if b.Width <= 0 {
		return errors.New("width must be positive")
	}
	if !validStyles[b.Style] && b.Style != "braille-spinner" {
		return fmt.Errorf("invalid style: %s", b.Style)
	}
	return nil
}

type ProgressBar interface {
	Render() string
	Update(current int)
	Finish()
}
