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
	ID             string `json:"id"`
	Current        int    `json:"current"`
	Total          int    `json:"total"`
	Width          int    `json:"width"`
	Style          string `json:"style"`
	ColorBar       string `json:"colorbar"`
	ColorText      string `json:"colortext"`
	Finished       bool   `json:"finished"`
	CustomChars    string `json:"chars"`
	Message        string `json:"message"`
	ShowElapsed    *bool  `json:"showelapsed,omitempty"`
	ShowThroughput *bool  `json:"showthroughput,omitempty"`
	ShowETA        *bool  `json:"showeta,omitempty"`
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
			StartTime:      time.Now(),
			ShowElapsed:    true,
			ShowThroughput: true,
			ShowETA:        true,
			Managed:        true,
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
	if update.CustomChars != "" {
		bar.CustomChars = update.CustomChars
	}
	bar.Message = update.Message
	if update.ShowElapsed != nil {
		bar.ShowElapsed = *update.ShowElapsed
	}
	if update.ShowThroughput != nil {
		bar.ShowThroughput = *update.ShowThroughput
	}
	if update.ShowETA != nil {
		bar.ShowETA = *update.ShowETA
	}
}

// RenderAll renders all managed progress bars to the terminal.
func (m *Manager) RenderAll() {
	m.mu.Lock()
	defer m.mu.Unlock()

	var sb strings.Builder

	// Clear previous output (inlined clearLines logic)
	if m.lastLines > 0 {
		for i := 0; i < m.lastLines; i++ {
			sb.WriteString("\r\033[K") // Carriage return, clear to end of line
			if i < m.lastLines-1 {
				sb.WriteString("\033[A") // Move up one line
			}
		}
	}

	var outputLines []string
	for _, id := range m.order {
		bar := m.bars[id]
		outputLines = append(outputLines, bar.Render())
	}

	// Print new output
	sb.WriteString(strings.Join(outputLines, "\n"))
	fmt.Print(sb.String())
	m.lastLines = len(outputLines)
}

// Clear clears all rendered progress bars from the terminal.
func (m *Manager) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	// Clear previous output (inlined clearLines logic)
	if m.lastLines > 0 {
		for i := 0; i < m.lastLines; i++ {
			fmt.Print("\r\033[K") // Carriage return, clear to end of line
			if i < m.lastLines-1 {
				fmt.Print("\033[A") // Move up one line
			}
		}
	}
	// Reset lastLines after clearing
	m.lastLines = 0
}
