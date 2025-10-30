package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Model represents the TUI application state
type Model struct {
	width  int
	height int
	
	// Status information
	ollamaStatus string
	modelName    string
	voiceStatus  string
	
	// Command input
	commandInput string
	commandHistory []string
	
	// Logs
	logs []LogEntry
	
	// Current view mode
	viewMode ViewMode
	
	// Timestamps
	lastUpdate time.Time
}

// ViewMode represents the current view mode
type ViewMode int

const (
	ViewStatus ViewMode = iota
	ViewLogs
	ViewConfig
)

// LogEntry represents a single log entry
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

// New creates a new TUI model
func New() Model {
	return Model{
		width:        80,
		height:       24,
		ollamaStatus: "Unknown",
		modelName:    "Not loaded",
		voiceStatus:  "Unknown",
		logs:         make([]LogEntry, 0),
		viewMode:     ViewStatus,
		lastUpdate:   time.Now(),
		commandHistory: make([]string, 0),
	}
}

// Init initializes the TUI application
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles messages and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "1":
			m.viewMode = ViewStatus
		case "2":
			m.viewMode = ViewLogs
		case "3":
			m.viewMode = ViewConfig
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	
	return m, nil
}

// View renders the TUI
func (m Model) View() string {
	var b strings.Builder
	
	// Header
	b.WriteString(m.renderHeader())
	b.WriteString("\n\n")
	
	// Main content based on view mode
	switch m.viewMode {
	case ViewStatus:
		b.WriteString(m.renderStatusView())
	case ViewLogs:
		b.WriteString(m.renderLogsView())
	case ViewConfig:
		b.WriteString(m.renderConfigView())
	}
	
	b.WriteString("\n\n")
	
	// Footer
	b.WriteString(m.renderFooter())
	
	return b.String()
}

// renderHeader renders the application header
func (m Model) renderHeader() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("86")).
		Background(lipgloss.Color("236")).
		Padding(0, 1)
	
	title := titleStyle.Render("╔════════════════════════════════════════╗\n" +
		"║           Emrys TUI                    ║\n" +
		"║  Your Personal AI Assistant on macOS  ║\n" +
		"╚════════════════════════════════════════╝")
	
	return title
}

// renderStatusView renders the status dashboard
func (m Model) renderStatusView() string {
	statusStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)
	
	content := fmt.Sprintf(`Status Dashboard

Ollama Service: %s
Current Model:  %s
Voice Output:   %s

Last Update:    %s
Uptime:         %s

System Resources:
  CPU:    N/A
  Memory: N/A
  
Navigation:
  Press '1' for Status (current)
  Press '2' for Logs
  Press '3' for Configuration
  Press 'q' or Ctrl+C to quit`,
		m.getColoredStatus(m.ollamaStatus),
		m.modelName,
		m.getColoredStatus(m.voiceStatus),
		m.lastUpdate.Format("15:04:05"),
		time.Since(m.lastUpdate).Round(time.Second).String())
	
	return statusStyle.Render(content)
}

// renderLogsView renders the logs viewer
func (m Model) renderLogsView() string {
	logStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)
	
	var content strings.Builder
	content.WriteString("Activity Logs\n\n")
	
	if len(m.logs) == 0 {
		content.WriteString("No logs available yet.\n")
	} else {
		// Show last 10 logs
		start := 0
		if len(m.logs) > 10 {
			start = len(m.logs) - 10
		}
		
		for i := start; i < len(m.logs); i++ {
			log := m.logs[i]
			content.WriteString(fmt.Sprintf("[%s] %s: %s\n",
				log.Timestamp.Format("15:04:05"),
				log.Level,
				log.Message))
		}
	}
	
	content.WriteString("\nNavigation: '1' Status | '2' Logs | '3' Config | 'q' Quit")
	
	return logStyle.Render(content.String())
}

// renderConfigView renders the configuration interface
func (m Model) renderConfigView() string {
	configStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(1, 2)
	
	content := `Configuration

Voice Settings:
  Voice:         Jamie (Premium)
  Rate:          200 wpm
  Volume:        System
  Quiet Hours:   Disabled

Model Settings:
  Current Model: llama3.2
  Auto-update:   Disabled

Display Settings:
  Theme:         Auto
  Color Scheme:  Default

Navigation: '1' Status | '2' Logs | '3' Config | 'q' Quit`
	
	return configStyle.Render(content)
}

// renderFooter renders the application footer
func (m Model) renderFooter() string {
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	
	footer := fmt.Sprintf("Emrys v0.1.0 | View: %s | Width: %dx%d",
		m.getViewName(),
		m.width,
		m.height)
	
	return footerStyle.Render(footer)
}

// getViewName returns the name of the current view
func (m Model) getViewName() string {
	switch m.viewMode {
	case ViewStatus:
		return "Status"
	case ViewLogs:
		return "Logs"
	case ViewConfig:
		return "Configuration"
	default:
		return "Unknown"
	}
}

// getColoredStatus returns a colored status string
func (m Model) getColoredStatus(status string) string {
	var style lipgloss.Style
	
	switch status {
	case "Running", "Active", "Available", "OK":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("10")) // Green
		return style.Render("✓ " + status)
	case "Stopped", "Inactive", "Unavailable":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("9")) // Red
		return style.Render("✗ " + status)
	case "Warning", "Degraded":
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("11")) // Yellow
		return style.Render("⚠ " + status)
	default:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // Gray
		return style.Render("- " + status)
	}
}

// AddLog adds a log entry to the model
func (m *Model) AddLog(level, message string) {
	m.logs = append(m.logs, LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	})
}

// UpdateStatus updates the status information
func (m *Model) UpdateStatus(ollamaStatus, modelName, voiceStatus string) {
	m.ollamaStatus = ollamaStatus
	m.modelName = modelName
	m.voiceStatus = voiceStatus
	m.lastUpdate = time.Now()
}
