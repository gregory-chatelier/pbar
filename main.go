package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	flag "github.com/spf13/pflag"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/gregory-chatelier/pbar/pbar"
)

// Version is set at build time
var Version = "v0.0.1-dev"

const (
	defaultWidth = 40
	defaultStyle = "classic"
	defaultTotal = 100
)

func boolPtr(v bool) *bool {
	return &v
}

func isValidStyle(style string) bool {
	validStyles := []string{"classic", "block", "spinner", "arrow", "braille", "custom", "braille-spinner"}
	for _, s := range validStyles {
		if s == style {
			return true
		}
	}
	return false
}

// generateInstanceID creates a stable ID for a progress bar instance.
// It hashes relevant command-line flags to ensure uniqueness across different logical bars,
// but stability across iterations of the same logical bar.
func generateInstanceID(explicitID string) string {
	if explicitID != "" {
		return explicitID
	}

	var signatureParts []string
	flag.CommandLine.Visit(func(flag *flag.Flag) {
		if flag.Name != "current" && flag.Name != "total" && flag.Name != "message" {
			signatureParts = append(signatureParts, "--"+flag.Name, flag.Value.String())
		}
	})

	// If no explicit ID and no relevant flags, use a default signature
	if len(signatureParts) == 0 {
		return "default_pbar_instance"
	}

	// Combine parts and hash
	signature := strings.Join(signatureParts, "_")
	h := sha256.New()
	h.Write([]byte(signature))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	// Declare variables for flags
	var width int
	var style string
	var colorBarName string
	var colorTextName string
	var finishedMessage string
	var version bool
	var customChars string
	var parallel bool
	var message string // Declare message flag
	var showElapsed, showThroughput, showETA bool
	var explicitInstanceID string // New flag for explicit ID

	// Define flags
	flag.IntVar(&width, "width", defaultWidth, "Width of the progress bar")
	flag.StringVar(&style, "style", defaultStyle, "Style of the progress bar (classic, block, spinner, arrow, braille, custom)")
	flag.StringVar(&colorBarName, "colorbar", "", fmt.Sprintf("Color for the bar. Available: %s", pbar.GetAvailableColors()))
	flag.StringVar(&colorTextName, "colortext", "", fmt.Sprintf("Color for the bar. Available: %s", pbar.GetAvailableColors()))
	flag.StringVar(&finishedMessage, "finished-message", "", "Message to display when the progress bar is complete")
	flag.BoolVar(&version, "version", false, "Print version information")
	flag.StringVar(&customChars, "chars", "", "Custom characters for the progress bar (e.g., '#=')")
	flag.BoolVar(&parallel, "parallel", false, "Enable parallel progress bar rendering")
	flag.StringVar(&message, "message", "", "Optional message to display alongside the progress bar")
	flag.BoolVar(&showElapsed, "show-elapsed", true, "Show elapsed time (default: true)")
	flag.BoolVar(&showThroughput, "show-throughput", true, "Show throughput (iterations/second) (default: true)")
	flag.BoolVar(&showETA, "show-eta", true, "Show estimated time remaining (default: true)")
	flag.StringVar(&explicitInstanceID, "id", "", "Unique ID for the progress bar instance (optional)")

	flag.Parse()

	// Handle positional arguments for current and total
	positionalArgs := flag.Args()

	// If version flag is set, print version and exit
	if version {
		fmt.Println(Version)
		os.Exit(0)
	}

	// If parallel mode is enabled
	if parallel {
		manager := pbar.NewManager()

		// Hide cursor
		fmt.Print("\033[?25l")

		// Handle Ctrl+C to show cursor and clear lines before exiting
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-c
			manager.Clear()
			fmt.Print("\033[?25h") // Show cursor
			os.Exit(0)
		}()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			line := scanner.Bytes()
			var update pbar.Update
			err := json.Unmarshal(line, &update)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
				continue
			}
			if update.ShowElapsed == nil {
				update.ShowElapsed = boolPtr(showElapsed)
			}
			if update.ShowThroughput == nil {
				update.ShowThroughput = boolPtr(showThroughput)
			}
			if update.ShowETA == nil {
				update.ShowETA = boolPtr(showETA)
			}
			manager.UpdateBar(update)
			manager.RenderAll()
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		}

		manager.Clear()
		fmt.Print("\033[?25h") // Show cursor
		return
	}

	// --- Single bar mode (existing logic) ---

	// Handle positional arguments for current and total
	var current, total int
	if len(positionalArgs) == 2 {
		var err error
		current, err = strconv.Atoi(positionalArgs[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid current value '%s'. Must be an integer.\n", positionalArgs[0])
			os.Exit(1)
		}
		total, err = strconv.Atoi(positionalArgs[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid total value '%s'. Must be an integer.\n", positionalArgs[1])
			os.Exit(1)
		}
	} else if len(positionalArgs) == 0 {
		// Default values if no positional arguments are provided
		current = 0
		total = defaultTotal
	} else if len(positionalArgs) == 1 {
		fmt.Fprintf(os.Stderr, "Error: When using positional arguments, provide both current and total values. Got only: %s\n", positionalArgs[0])
		os.Exit(1)
	} else {
		fmt.Fprintf(os.Stderr, "Error: Too many positional arguments. Expected 0 or 2, got %d: %v\n", len(positionalArgs), positionalArgs)
		os.Exit(1)
	}

	// Validate style
	if !isValidStyle(style) {
		fmt.Fprintf(os.Stderr, "Error: Invalid style '%s'. Must be one of: classic, block, spinner, arrow, braille, custom, braille-spinner\n", style)
		os.Exit(1)
	}

	// Validate and get ANSI color codes
	colorBarCode := pbar.GetColorCode(colorBarName)
	colorTextCode := pbar.GetColorCode(colorTextName)

	instanceID := generateInstanceID(explicitInstanceID)

	var bar *pbar.Bar
	if current == 0 {
		// If current is 0, it's a new progress bar, so initialize a fresh state
		bar = &pbar.Bar{}
		bar.StartTime = time.Now()
		pbar.DeleteState(instanceID) // Ensure no old state interferes
	} else {
		// Otherwise, try to load the existing state
		loadedBar, err := pbar.LoadState(instanceID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Expected existing progress bar state for ID '%s' but none found: %v\n", instanceID, err)
			os.Exit(1)
		}
		bar = loadedBar
	}

	bar.Total = total
	bar.PreviousCurrent = bar.Current
	bar.Current = current
	bar.Width = width
	bar.Style = style
	bar.ColorBar = colorBarCode
	bar.ColorText = colorTextCode
	bar.Finished = current >= total
	bar.CustomChars = customChars
	bar.Message = message
	bar.CompletionMessage = finishedMessage
	bar.ShowElapsed = showElapsed
	bar.ShowThroughput = showThroughput
	bar.ShowETA = showETA

	fmt.Print(bar.Render())

	if bar.Finished {
		pbar.DeleteState(instanceID)
	} else {
		bar.LastUpdateTime = time.Now()
		pbar.SaveState(bar, instanceID)
	}

	// Exit with an error code if current > total (unless finished)
	isIndeterminate := style == "spinner" || style == "braille-spinner"
	if !(current >= total) && !isIndeterminate && current > total {
		fmt.Fprintln(os.Stderr, "Error: Current value cannot be greater than total.")
		os.Exit(1)
	}

	if width <= 0 {
		fmt.Fprintf(os.Stderr, "Error: Width must be positive\n")
		os.Exit(1)
	}
}