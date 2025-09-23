package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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

func isValidStyle(style string) bool {
	validStyles := []string{"classic", "block", "spinner", "arrow", "braille", "custom"}
	for _, s := range validStyles {
		if s == style {
			return true
		}
	}
	return false
}

func main() {
	startTime := time.Now()

	// Declare variables for flags
	var width int
	var style string
	var indeterminate bool
	var colorBarName string
	var colorTextName string
	var finished bool
	var version bool
	var quiet bool
	var customChars string
	var parallel bool // Declare parallel flag here

	// Define flags
	flag.IntVar(&width, "width", defaultWidth, "Width of the progress bar")
	flag.StringVar(&style, "style", defaultStyle, "Style of the progress bar (classic, block, spinner, arrow, braille, custom)")
	flag.BoolVar(&indeterminate, "indeterminate", false, "Render an indeterminate spinner")
	flag.StringVar(&colorBarName, "colorbar", "", fmt.Sprintf("Color for the bar. Available: %s", pbar.GetAvailableColors()))
	flag.StringVar(&colorTextName, "colortext", "", fmt.Sprintf("Color for the text. Available: %s", pbar.GetAvailableColors()))
	flag.BoolVar(&finished, "finished", false, "Render a finished state")
	flag.BoolVar(&version, "version", false, "Print version information")
	flag.BoolVar(&quiet, "quiet", false, "Output only the percentage")
	flag.StringVar(&customChars, "chars", "", "Custom characters for the progress bar (e.g., '#=')")
	flag.BoolVar(&parallel, "parallel", false, "Enable parallel progress bar rendering")

	// Custom usage function for man-page style help
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [current] [total] [flags]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n%s is a command-line tool that makes it easy to add progress bars to any Bash or Zsh script.\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Basic usage: 25%% complete out of 100\n")
		fmt.Fprintf(os.Stderr, "  %s 25 100\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  # Using a block style bar with custom width\n")
		fmt.Fprintf(os.Stderr, "  %s 50 100 --style=block --width=20\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  # Indeterminate spinner\n")
		fmt.Fprintf(os.Stderr, "  %s --indeterminate\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  # Quiet mode (output only percentage) for scripting\n")
		fmt.Fprintf(os.Stderr, "  %s 75 100 --quiet\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  # Custom characters and colors\n")
		fmt.Fprintf(os.Stderr, "  %s 60 100 --style=custom --chars='#-' --colorbar=green --colortext=yellow\n", os.Args[0])
	}

	// Custom usage function for man-page style help
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [current] [total] [flags]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n%s is a command-line tool that makes it easy to add progress bars to any Bash or Zsh script.\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nFlags:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Basic usage: 25%% complete out of 100\n")
		fmt.Fprintf(os.Stderr, "  %s 25 100\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  # Using a block style bar with custom width\n")
		fmt.Fprintf(os.Stderr, "  %s 50 100 --style=block --width=20\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  # Indeterminate spinner\n")
		fmt.Fprintf(os.Stderr, "  %s --indeterminate\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  # Quiet mode (output only percentage) for scripting\n")
		fmt.Fprintf(os.Stderr, "  %s 75 100 --quiet\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\n  # Custom characters and colors\n")
		fmt.Fprintf(os.Stderr, "  %s 60 100 --style=custom --chars='#-' --colorbar=green --colortext=yellow\n", os.Args[0])
	}

	flag.Parse()

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

	// If no positional arguments are provided, print usage and exit
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	// Handle positional arguments for current and total
	var current, total int
	if flag.NArg() == 2 {
		var err error
		current, err = strconv.Atoi(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid current value '%s'. Must be an integer.\n", flag.Arg(0))
			os.Exit(1)
		}
		total, err = strconv.Atoi(flag.Arg(1))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: Invalid total value '%s'. Must be an integer.\n", flag.Arg(1))
			os.Exit(1)
		}
	} else if flag.NArg() == 0 {
		// Default values if no positional arguments are provided
		current = 0
		total = defaultTotal
	} else {
		fmt.Fprintln(os.Stderr, "Error: When using positional arguments, provide both current and total values, or neither.")
		os.Exit(1)
	}

	// Validate style
	if !isValidStyle(style) {
		fmt.Fprintf(os.Stderr, "Error: Invalid style '%s'. Must be one of: classic, block, spinner, arrow, braille, custom\n", style)
		os.Exit(1)
	}

	// Validate and get ANSI color codes
	colorBarCode := pbar.GetColorCode(colorBarName)
	colorTextCode := pbar.GetColorCode(colorTextName)

	bar := &pbar.Bar{
		Total:         total,
		Current:       current,
		Width:         width,
		Style:         style,
		Indeterminate: indeterminate,
		ColorBar:      colorBarCode,
		ColorText:     colorTextCode,
		Finished:      finished,
		Quiet:         quiet,
		StartTime:     startTime,
		CustomChars:   customChars,
	}

	fmt.Print(bar.Render())

	// Exit with an error code if current > total (unless finished or indeterminate)
	if !finished && !indeterminate && current > total {
		fmt.Fprintln(os.Stderr, "Error: Current value cannot be greater than total.")
		os.Exit(1)
	}

	if width <= 0 {
		fmt.Fprintf(os.Stderr, "Error: Width must be positive\n")
		os.Exit(1)
	}
}
