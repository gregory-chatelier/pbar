package pbar

import (
	"fmt"
	"strings"
	"time"
)

var spinnerChars = []string{"|", "/", "-", "\\"}

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

// Bar represents a progress bar.
type Bar struct {
	Total         int
	Current       int
	Width         int
	Style         string
	Indeterminate bool
	ColorBar      string // Now stores ANSI escape code directly
	ColorText     string // Now stores ANSI escape code directly
	Finished      bool
	Quiet         bool
	StartTime     time.Time
	LastUpdateTime time.Time
	ThroughputHistory []float64
	CustomChars   string
	spinnerState  int
}

// Render generates the string representation of the progress bar.
func (b *Bar) Render() string {
	// Update LastUpdateTime
	b.LastUpdateTime = time.Now()

	percent := float64(b.Current) / float64(b.Total)
	if percent < 0 {
		percent = 0
	}
	if percent > 1 {
		percent = 1
	}

	percentString := fmt.Sprintf("%d%%", int(percent*100))

	if b.Quiet {
		return percentString
	}

	var metadataString string
	if !b.StartTime.IsZero() {
		// Calculate elapsed time
		elapsedTime := time.Since(b.StartTime)
		elapsedTimeStr := formatDuration(elapsedTime)

		// Calculate throughput and ETA only if not indeterminate and total is greater than 0
		var throughputStr, etaStr string
		if !b.Indeterminate && b.Total > 0 && b.Current > 0 {
			// Calculate current throughput
			currentThroughput := float64(b.Current) / elapsedTime.Seconds()

			// Update throughput history (simple moving average for now)
			b.ThroughputHistory = append(b.ThroughputHistory, currentThroughput)
			if len(b.ThroughputHistory) > 10 { // Keep last 10 samples
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
			}
		}
		metadataString = fmt.Sprintf(" Elapsed %s%s%s", elapsedTimeStr, throughputStr, etaStr)
	}

	if b.Finished {
		return fmt.Sprintf("[✔] 100%%%s", metadataString)
	}

	if b.Indeterminate {
		char := spinnerChars[b.spinnerState%len(spinnerChars)]
		b.spinnerState++
		return fmt.Sprintf("[%s]%s", char, metadataString)
	}

	style := b.Style
	if style == "" {
		style = "classic"
	}

	var barString string
	switch style {
	case "spinner":
		char := spinnerChars[b.spinnerState%len(spinnerChars)]
		b.spinnerState++
		barString = fmt.Sprintf("[%s]", char)
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

	return fmt.Sprintf("%s %s%s", barString, percentString, metadataString)
}

func (b *Bar) renderBar(filledChar, emptyChar, colorCode string) string {
	percent := float64(b.Current) / float64(b.Total)
	filledWidth := int(percent * float64(b.Width))
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
	percent := float64(b.Current) / float64(b.Total)
	filledWidth := int(percent * float64(b.Width))
	emptyWidth := b.Width - filledWidth

	var filled string
	if filledWidth > 0 {
		filled = strings.Repeat("-", filledWidth-1) + ">"
	} else {
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
	percent := float64(b.Current) / float64(b.Total)
	filledWidth := int(percent * float64(b.Width))
	
	var barContentBuilder strings.Builder
	barContentBuilder.WriteString("[")

	for i := 0; i < b.Width; i++ {
		if i < filledWidth {
			barContentBuilder.WriteString("⣿") // Full block
		} else if i == filledWidth && filledWidth < b.Width {
			// Calculate fractional fill for the current character
			fraction := (percent * float64(b.Width)) - float64(filledWidth)
			brailleIndex := int(fraction * float64(len(brailleChars)-1))
			barContentBuilder.WriteString(brailleChars[brailleIndex])
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