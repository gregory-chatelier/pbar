package pbar

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestClassicBar(t *testing.T) {
	t.Run("renders a classic bar at 50%", func(t *testing.T) {
		bar := &Bar{
			Total:   100,
			Current: 50,
			Width:   10,
		}

		expected := "\r[#####-----] 50%\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestBlockBar(t *testing.T) {
	t.Run("renders a block bar at 50%", func(t *testing.T) {
		bar := &Bar{
			Total:   100,
			Current: 50,
			Width:   10,
			Style:   "block",
		}

		expected := "\r[█████     ] 50%\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestSpinner(t *testing.T) {
	t.Run("renders a spinner that cycles through characters", func(t *testing.T) {
		bar := &Bar{
			Style:   "spinner",
			TestMode: true,
		}

		expectations := []string{
			"\r[|]\x1b[K",
			"\r[/]\x1b[K",
			"\r[-]\x1b[K",
			"\r[\\]\x1b[K",
			"\r[|]\x1b[K", // Check for wrap around
		}

		for i, expected := range expectations {
			bar.spinnerState = i // Manually set state for predictability
			actual := bar.Render()
			if actual != expected {
				t.Errorf("Expected '%s', but got '%s'", expected, actual)
			}
		}
	})
}

func TestBrailleSpinner(t *testing.T) {
	t.Run("renders a braille spinner that cycles through characters", func(t *testing.T) {
		bar := &Bar{
			Style:   "braille-spinner",
			TestMode: true,
		}

		expectations := []string{
			"\r[⠋]\x1b[K",
			"\r[⠙]\x1b[K",
			"\r[⠹]\x1b[K",
			"\r[⠸]\x1b[K",
			"\r[⠼]\x1b[K",
			"\r[⠴]\x1b[K",
			"\r[⠦]\x1b[K",
			"\r[⠧]\x1b[K",
			"\r[⠇]\x1b[K",
			"\r[⠏]\x1b[K",
		}

		for i, expected := range expectations {
			bar.spinnerState = i // Manually set state for predictability
			actual := bar.Render()
			if actual != expected {
				t.Errorf("Expected '%s', but got '%s'", expected, actual)
			}
		}
	})
}

func TestIndeterminateMode(t *testing.T) {
	t.Run("renders a spinner without percentage", func(t *testing.T) {
		bar := &Bar{
			Style:    "spinner",
			TestMode: true,
		}

		expectations := []string{
			"\r[|]\x1b[K",
			"\r[/]\x1b[K",
			"\r[-]\x1b[K",
			"\r[\\]\x1b[K",
		}

		for i, expected := range expectations {
			bar.spinnerState = i // Manually set state for predictability
			actual := bar.Render()
			if actual != expected {
				t.Errorf("Expected '%s', but got '%s'", expected, actual)
			}
		}
	})
}

func TestColorSupport(t *testing.T) {
	t.Run("renders a bar with color", func(t *testing.T) {
		// ANSI escape codes for green and reset
		green := "\x1b[32m"
		reset := "\x1b[0m"

		bar := &Bar{
			Total:    100,
			Current:  50,
			Width:    10,
			ColorBar: green, // Use the ANSI code directly, as the Render function expects
		}

		expected := fmt.Sprintf("\r%s[#####-----]%s 50%%\x1b[K", green, reset)
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestFinishedState(t *testing.T) {
	t.Run("renders a finished bar with default message", func(t *testing.T) {
		bar := &Bar{
			Total:    100,
			Current:  100,
			Width:    10,
			Finished: true,
		}

		expected := "\r[✔] 100% Task Complete!\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})

	t.Run("renders a finished bar with custom message", func(t *testing.T) {
		bar := &Bar{
			Total:             100,
			Current:           100,
			Width:             10,
			Finished:          true,
			CompletionMessage: "Done!",
		}

		expected := "\r[✔] 100% Done!\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}



func TestArrowBar(t *testing.T) {
	t.Run("renders an arrow bar at 50%", func(t *testing.T) {
		bar := &Bar{
			Total:   100,
			Current: 50,
			Width:   10,
			Style:   "arrow",
		}

		expected := "\r[---->     ] 50%\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestBrailleBar(t *testing.T) {
	t.Run("renders a braille bar at 50%", func(t *testing.T) {
		bar := &Bar{
			Total:   100,
			Current: 50,
			Width:   10,
			Style:   "braille",
		}

		expected := "\r[⣿⣿⣿⣿⣿     ] 50%\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestBrailleBarFractional(t *testing.T) {
	t.Run("renders a braille bar with fractional fill", func(t *testing.T) {
		bar := &Bar{
			Total:   100,
			Current: 55,
			Width:   10,
			Style:   "braille",
		}

		expected := "\r[⣿⣿⣿⣿⣿⠏    ] 55%\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestElapsedTime(t *testing.T) {
	t.Run("displays elapsed time", func(t *testing.T) {
		startTime := time.Now().Add(-5 * time.Second) // 5 seconds ago
		bar := &Bar{
			Total:     100,
			Current:   50,
			Width:     10,
			StartTime: startTime,
		}

		actual := bar.Render()

		// We can't compare exact time, so we check if it contains "Elapsed X" where X is close to 5s
		if !strings.Contains(actual, "Elapsed 5s") && !strings.Contains(actual, "Elapsed 4s") && !strings.Contains(actual, "Elapsed 6s") {
			t.Errorf("Expected elapsed time to be around 5s, but got '%s'", actual)
		}
	})
}

func TestAverageThroughput(t *testing.T) {
	t.Run("calculates average throughput and ETA", func(t *testing.T) {
		startTime := time.Now()
		bar := &Bar{
			Total:     100,
			Current:   50,
			Width:     10,
			StartTime: startTime,
		}

		// Introduce a small delay to ensure elapsed time is non-zero
		time.Sleep(10 * time.Millisecond)

		actual := bar.Render()

		if !strings.Contains(actual, "it/s") {
			t.Errorf("Expected output to contain throughput (it/s), but got '%s'", actual)
		}
		if !strings.Contains(actual, "ETA") {
			t.Errorf("Expected output to contain ETA, but got '%s'", actual)
		}
	})
}

func TestCustomBar(t *testing.T) {
	t.Run("renders a custom bar with two characters", func(t *testing.T) {
		bar := &Bar{
			Total:       100,
			Current:     50,
			Width:       10,
			Style:       "custom",
			CustomChars: "#=",
		}

		expected := "\r[#####=====] 50%\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})

	t.Run("renders a custom bar with one character", func(t *testing.T) {
		bar := &Bar{
			Total:       100,
			Current:     50,
			Width:       10,
			Style:       "custom",
			CustomChars: "*",
		}

		expected := "\r[**********] 50%\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})

	t.Run("renders a custom bar with no characters (defaults to classic)", func(t *testing.T) {
		bar := &Bar{
			Total:       100,
			Current:     50,
			Width:       10,
			Style:       "custom",
			CustomChars: "",
		}

		expected := "\r[#####-----] 50%\x1b[K"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestBarEdgeCases(t *testing.T) {
	t.Run("classic bar at 0%", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 0, Width: 10, Style: "classic"}
		expected := "\r[----------] 0%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("classic bar at 100%", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 100, Width: 10, Style: "classic"}
		expected := "\r[##########] 100%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("block bar at 0%", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 0, Width: 10, Style: "block"}
		expected := "\r[          ] 0%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("block bar at 100%", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 100, Width: 10, Style: "block"}
		expected := "\r[██████████] 100%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("total is 0, current is 0", func(t *testing.T) {
		bar := &Bar{Total: 0, Current: 0, Width: 10, Style: "classic"}
		expected := "\r[----------] 0%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("total is 0, current is non-zero", func(t *testing.T) {
		bar := &Bar{Total: 0, Current: 50, Width: 10, Style: "classic"}
		expected := "\r[##########] 100%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("width is 0", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 50, Width: 0, Style: "classic"}
		expected := "\r[] 50%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("width is 1", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 50, Width: 1, Style: "classic"}
		expected := "\r[#] 50%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("arrow bar at 0%", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 0, Width: 10, Style: "arrow"}
		expected := "\r[          ] 0%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("arrow bar at 100%", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 100, Width: 10, Style: "arrow"}
		expected := "\r[--------->] 100%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("braille bar at 0%", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 0, Width: 10, Style: "braille"}
		expected := "\r[          ] 0%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("braille bar at 100%", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 100, Width: 10, Style: "braille"}
		expected := "\r[⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿] 100%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("braille bar fractional 1/8", func(t *testing.T) {
		bar := &Bar{Total: 8, Current: 1, Width: 1, Style: "braille"}
		expected := "\r[⠁] 12%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("braille bar fractional 7/8", func(t *testing.T) {
		bar := &Bar{Total: 8, Current: 7, Width: 1, Style: "braille"}
		expected := "\r[⡿] 87%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestSmartMetadataEdgeCases(t *testing.T) {
	t.Run("elapsed time with zero start time", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 50, Width: 10}
		actual := bar.Render()
		if strings.Contains(actual, "Elapsed") {
			t.Errorf("Expected no Elapsed time, got '%s'", actual)
		}
	})

	t.Run("throughput with zero current", func(t *testing.T) {
		startTime := time.Now().Add(-time.Second) // Ensure non-zero elapsed time
		bar := &Bar{Total: 100, Current: 0, Width: 10, StartTime: startTime}
		time.Sleep(10 * time.Millisecond) // Ensure elapsed time is non-zero
		actual := bar.Render()
		if !strings.Contains(actual, "0.00 it/s") {
			t.Errorf("Expected 0.00 it/s, got '%s'", actual)
		}
	})

	t.Run("throughput with zero total", func(t *testing.T) {
		startTime := time.Now().Add(-time.Second) // Ensure non-zero elapsed time
		bar := &Bar{Total: 0, Current: 50, Width: 10, StartTime: startTime}
		time.Sleep(10 * time.Millisecond) // Ensure elapsed time is non-zero
		actual := bar.Render()
		if !strings.Contains(actual, "0.00 it/s") {
			t.Errorf("Expected 0.00 it/s, got '%s'", actual)
		}
	})

	t.Run("ETA with current equals total", func(t *testing.T) {
		startTime := time.Now().Add(-time.Second) // Ensure non-zero elapsed time
		bar := &Bar{Total: 100, Current: 100, Width: 10, StartTime: startTime}
		time.Sleep(10 * time.Millisecond) // Ensure elapsed time is non-zero
		actual := bar.Render()
		if !strings.Contains(actual, "ETA 0s") {
			t.Errorf("Expected ETA 0s, got '%s'", actual)
		}
	})
}

func TestInvalidInputs(t *testing.T) {
	t.Run("negative current", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: -10, Width: 10}
		expected := "\r[----------] 0%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("negative total", func(t *testing.T) {
		bar := &Bar{Total: -100, Current: 50, Width: 10}
		expected := "\r[##########] 100%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("negative width", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 50, Width: -10}
		expected := "\r[] 50%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})

	t.Run("invalid style", func(t *testing.T) {
		bar := &Bar{Total: 100, Current: 50, Width: 10, Style: "invalid"}
		expected := "\r[#####-----] 50%\x1b[K" // devrait utiliser le style par défaut
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestLargeNumbers(t *testing.T) {
	t.Run("very large numbers", func(t *testing.T) {
		bar := &Bar{
			Total:   1<<31 - 1, // max int32
			Current: 1 << 30,   // half of max int32
			Width:   10,
		}
		expected := "\r[#####-----] 50%\x1b[K"
		actual := bar.Render()
		if actual != expected {
			t.Errorf("Expected '%s', got '%s'", expected, actual)
		}
	})
}

func TestThroughputHistory(t *testing.T) {
	t.Run("throughput history size limit", func(t *testing.T) {
		bar := &Bar{
			Total:     100,
			Current:   50,
			Width:     10,
			StartTime: time.Now().Add(-10 * time.Second),
		}

		// Remplir l'historique avec plus de 10 valeurs
		for i := 0; i < 15; i++ {
			bar.Render()
			time.Sleep(100 * time.Millisecond)
		}

		if len(bar.ThroughputHistory) > 10 {
			t.Errorf("Throughput history exceeded limit of 10, got %d", len(bar.ThroughputHistory))
		}
	})
}

func TestLongDurations(t *testing.T) {
	t.Run("very long elapsed time", func(t *testing.T) {
		bar := &Bar{
			Total:     100,
			Current:   50,
			Width:     10,
			StartTime: time.Now().Add(-24 * 30 * time.Hour), // ~1 month ago
		}
		actual := bar.Render()
		if !strings.Contains(actual, "720h") {
			t.Errorf("Expected to show hours for long duration, got '%s'", actual)
		}
	})
}
