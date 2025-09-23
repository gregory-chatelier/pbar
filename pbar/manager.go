package pbar

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

// Update represents an update for a single progress bar.
type Update struct {
	ID          string `json:"id"`
	Current     int    `json:"current"`
	Total       int    `json:"total"`
	Width       int    `json:"width"`
	Style       string `json:"style"`
	ColorBar    string `json:"colorbar"`
	ColorText   string `json:"colortext"`
	Finished    bool   `json:"finished"`
	Quiet       bool   `json:"quiet"`
	CustomChars string `json:"chars"`
	Message     string `json:"message"`
}

// Manager manages multiple progress bars.
type Manager struct {
	bars      map[string]*Bar
	order     []string // To maintain the order of bars
	mu        sync.Mutex
	lastLines int // Number of lines printed in the last render cycle
}

// NewManager creates a new Manager instance.
func NewManager() *Manager {
	return &Manager{
		bars: make(map[string]*Bar),
	}
}

// UpdateBar creates or updates a progress bar.
func (m *Manager) UpdateBar(update Update) {
	m.mu.Lock()
	defer m.mu.Unlock()

	bar, exists := m.bars[update.ID]
	if !exists {
		bar = &Bar{
			StartTime: time.Now(),
		}
		m.bars[update.ID] = bar
		m.order = append(m.order, update.ID)
		sort.Strings(m.order) // Keep bars sorted by ID for consistent display
	}

	// Apply updates
	bar.Current = update.Current
	bar.Total = update.Total
	if update.Width > 0 {
		bar.Width = update.Width
	} else if bar.Width == 0 { // Set default width if not provided and not already set
		bar.Width = defaultWidth
	}
	if update.Style != "" {
		bar.Style = update.Style
	} else if bar.Style == "" { // Set default style if not provided and not already set
		bar.Style = defaultStyle
	}
	if update.ColorBar != "" {
		bar.ColorBar = GetColorCode(update.ColorBar)
	}
	if update.ColorText != "" {
		bar.ColorText = GetColorCode(update.ColorText)
	}
	bar.Finished = update.Finished
	bar.Quiet = update.Quiet
	if update.CustomChars != "" {
		bar.CustomChars = update.CustomChars
	}
	bar.Message = update.Message
}

// RenderAll renders all managed progress bars to the terminal.
func (m *Manager) RenderAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear previous output
	m.clearLines()

	var outputLines []string
	for _, id := range m.order {
		bar := m.bars[id]
		outputLines = append(outputLines, bar.Render())
	}

	// Print new output
	fmt.Print(strings.Join(outputLines, "\n"))
	m.lastLines = len(outputLines)
}

// Clear clears all rendered progress bars from the terminal.
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clearLines()
}

// clearLines moves the cursor up and clears the lines.
func (m *Manager) clearLines() {
	if m.lastLines > 0 {
		// Move cursor up by lastLines
		fmt.Printf("\033[%dA", m.lastLines)
		// Clear each line
		for i := 0; i < m.lastLines; i++ {
			fmt.Print("\033[K") // Clear from cursor to end of line
			if i < m.lastLines-1 {
				fmt.Print("\n") // Move to next line, but not after the last one
			}
		}
		// Move cursor back up to the start of the first cleared line
		fmt.Printf("\033[%dA", m.lastLines-1)
	}
}
