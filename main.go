package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gregory-chatelier/pbar/pbar"
)

// Version is set at build time
var Version = "v0.0.1-dev"

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

	// Define flags
	flag.IntVar(&width, "width", 40, "Width of the progress bar")
	flag.StringVar(&style, "style", "classic", "Style of the progress bar (classic, block, spinner, arrow, braille, custom)")
	flag.BoolVar(&indeterminate, "indeterminate", false, "Render an indeterminate spinner")
	flag.StringVar(&colorBarName, "colorbar", "", fmt.Sprintf("Color for the bar. Available: %s", pbar.GetAvailableColors()))
	flag.StringVar(&colorTextName, "colortext", "", fmt.Sprintf("Color for the text. Available: %s", pbar.GetAvailableColors()))
	flag.BoolVar(&finished, "finished", false, "Render a finished state")
	flag.BoolVar(&version, "version", false, "Print version information")
	flag.BoolVar(&quiet, "quiet", false, "Output only the percentage")
	flag.StringVar(&customChars, "chars", "", "Custom characters for the progress bar (e.g., '#=')")

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

	// If no positional arguments are provided, print usage and exit
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(0)
	}

	if version {
		fmt.Println(Version)
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
		total = 100
	} else {
		fmt.Fprintln(os.Stderr, "Error: When using positional arguments, provide both current and total values, or neither.")
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
}