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

		expected := "[#####-----] 50%"
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

		expected := "[█████     ] 50%"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestSpinner(t *testing.T) {
	t.Run("renders a spinner that cycles through characters", func(t *testing.T) {
		bar := &Bar{
			Total:   100,
			Current: 25,
			Style:   "spinner",
		}

		expectations := []string{
			"[|] 25%",
			"[/] 25%",
			"[-] 25%",
			"[\\] 25%",
			"[|] 25%", // Check for wrap around
		}

		for _, expected := range expectations {
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
			Indeterminate: true,
		}

		expectations := []string{
			"[|]",
			"[/]",
			"[-]",
			"[\\]",
		}

		for _, expected := range expectations {
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

		expected := fmt.Sprintf("%s[#####-----]%s 50%%", green, reset)
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestFinishedState(t *testing.T) {
	t.Run("renders a finished bar", func(t *testing.T) {
		bar := &Bar{
			Total:    100,
			Current:  100,
			Width:    10,
			Finished: true,
		}

		expected := "[✔] 100%"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}

func TestQuietMode(t *testing.T) {
	t.Run("renders only the percentage when quiet", func(t *testing.T) {
		bar := &Bar{
			Total:   100,
			Current: 75,
			Width:   10,
			Quiet:   true,
		}

		expected := "75%"
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

		expected := "[---->     ] 50%"
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

		expected := "[⣿⣿⣿⣿⣿     ] 50%"
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

		expected := "[⣿⣿⣿⣿⣿⠏    ] 55%"
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

		expected := "[#####=====] 50%"
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

		expected := "[**********] 50%"
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

		expected := "[#####-----] 50%"
		actual := bar.Render()

		if actual != expected {
			t.Errorf("Expected '%s', but got '%s'", expected, actual)
		}
	})
}
