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

	flag.Parse()

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