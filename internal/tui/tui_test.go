package tui

import (
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestNew(t *testing.T) {
	model := New()

	if model.width != 80 {
		t.Errorf("Expected width 80, got %d", model.width)
	}

	if model.height != 24 {
		t.Errorf("Expected height 24, got %d", model.height)
	}

	if model.ollamaStatus != "Unknown" {
		t.Errorf("Expected ollamaStatus 'Unknown', got '%s'", model.ollamaStatus)
	}

	if model.modelName != "Not loaded" {
		t.Errorf("Expected modelName 'Not loaded', got '%s'", model.modelName)
	}

	if model.voiceStatus != "Unknown" {
		t.Errorf("Expected voiceStatus 'Unknown', got '%s'", model.voiceStatus)
	}

	if model.viewMode != ViewStatus {
		t.Errorf("Expected viewMode ViewStatus, got %v", model.viewMode)
	}

	if len(model.logs) != 0 {
		t.Errorf("Expected empty logs, got %d logs", len(model.logs))
	}
}

func TestInit(t *testing.T) {
	model := New()
	cmd := model.Init()

	if cmd != nil {
		t.Error("Expected Init to return nil")
	}
}

func TestUpdate_KeyPress(t *testing.T) {
	model := New()

	// Test view switching
	testCases := []struct {
		key      string
		expected ViewMode
	}{
		{"1", ViewStatus},
		{"2", ViewLogs},
		{"3", ViewConfig},
	}

	for _, tc := range testCases {
		msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{rune(tc.key[0])}}
		updatedModel, _ := model.Update(msg)
		m := updatedModel.(Model)

		if m.viewMode != tc.expected {
			t.Errorf("Key '%s': expected viewMode %v, got %v", tc.key, tc.expected, m.viewMode)
		}
	}
}

func TestUpdate_WindowSize(t *testing.T) {
	model := New()

	msg := tea.WindowSizeMsg{
		Width:  120,
		Height: 40,
	}

	updatedModel, _ := model.Update(msg)
	m := updatedModel.(Model)

	if m.width != 120 {
		t.Errorf("Expected width 120, got %d", m.width)
	}

	if m.height != 40 {
		t.Errorf("Expected height 40, got %d", m.height)
	}
}

func TestView(t *testing.T) {
	model := New()
	view := model.View()

	if len(view) == 0 {
		t.Error("Expected non-empty view")
	}

	// Check for key elements in the view
	expectedStrings := []string{
		"Emrys TUI",
		"Status Dashboard",
		"Ollama Service",
		"Voice Output",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("View missing expected content: %s", expected)
		}
	}
}

func TestRenderStatusView(t *testing.T) {
	model := New()
	model.viewMode = ViewStatus
	
	view := model.renderStatusView()

	if len(view) == 0 {
		t.Error("Expected non-empty status view")
	}

	expectedStrings := []string{
		"Status Dashboard",
		"Ollama Service",
		"Current Model",
		"Voice Output",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("Status view missing expected content: %s", expected)
		}
	}
}

func TestRenderLogsView(t *testing.T) {
	model := New()
	model.viewMode = ViewLogs
	
	view := model.renderLogsView()

	if len(view) == 0 {
		t.Error("Expected non-empty logs view")
	}

	if !strings.Contains(view, "Activity Logs") {
		t.Error("Logs view missing 'Activity Logs'")
	}

	// Test with logs
	model.AddLog("INFO", "Test log message")
	view = model.renderLogsView()

	if !strings.Contains(view, "Test log message") {
		t.Error("Logs view not showing added log message")
	}
}

func TestRenderConfigView(t *testing.T) {
	model := New()
	model.viewMode = ViewConfig
	
	view := model.renderConfigView()

	if len(view) == 0 {
		t.Error("Expected non-empty config view")
	}

	expectedStrings := []string{
		"Configuration",
		"Voice Settings",
		"Model Settings",
		"Display Settings",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(view, expected) {
			t.Errorf("Config view missing expected content: %s", expected)
		}
	}
}

func TestAddLog(t *testing.T) {
	model := New()

	if len(model.logs) != 0 {
		t.Error("Expected empty logs initially")
	}

	model.AddLog("INFO", "First log")
	
	if len(model.logs) != 1 {
		t.Errorf("Expected 1 log, got %d", len(model.logs))
	}

	if model.logs[0].Level != "INFO" {
		t.Errorf("Expected level 'INFO', got '%s'", model.logs[0].Level)
	}

	if model.logs[0].Message != "First log" {
		t.Errorf("Expected message 'First log', got '%s'", model.logs[0].Message)
	}

	// Add multiple logs
	model.AddLog("ERROR", "Second log")
	model.AddLog("WARN", "Third log")

	if len(model.logs) != 3 {
		t.Errorf("Expected 3 logs, got %d", len(model.logs))
	}
}

func TestUpdateStatus(t *testing.T) {
	model := New()

	originalUpdate := model.lastUpdate

	// Wait a bit to ensure timestamp changes
	time.Sleep(10 * time.Millisecond)

	model.UpdateStatus("Running", "llama3.2", "Active")

	if model.ollamaStatus != "Running" {
		t.Errorf("Expected ollamaStatus 'Running', got '%s'", model.ollamaStatus)
	}

	if model.modelName != "llama3.2" {
		t.Errorf("Expected modelName 'llama3.2', got '%s'", model.modelName)
	}

	if model.voiceStatus != "Active" {
		t.Errorf("Expected voiceStatus 'Active', got '%s'", model.voiceStatus)
	}

	if !model.lastUpdate.After(originalUpdate) {
		t.Error("Expected lastUpdate to be updated")
	}
}

func TestGetViewName(t *testing.T) {
	model := New()

	testCases := []struct {
		viewMode ViewMode
		expected string
	}{
		{ViewStatus, "Status"},
		{ViewLogs, "Logs"},
		{ViewConfig, "Configuration"},
	}

	for _, tc := range testCases {
		model.viewMode = tc.viewMode
		name := model.getViewName()

		if name != tc.expected {
			t.Errorf("ViewMode %v: expected '%s', got '%s'", tc.viewMode, tc.expected, name)
		}
	}
}

func TestGetColoredStatus(t *testing.T) {
	model := New()

	testCases := []struct {
		status   string
		expected string // Should contain the status text
	}{
		{"Running", "Running"},
		{"Active", "Active"},
		{"Stopped", "Stopped"},
		{"Inactive", "Inactive"},
		{"Warning", "Warning"},
		{"Unknown", "Unknown"},
	}

	for _, tc := range testCases {
		result := model.getColoredStatus(tc.status)

		if !strings.Contains(result, tc.expected) {
			t.Errorf("Status '%s': expected result to contain '%s', got '%s'", tc.status, tc.expected, result)
		}
	}
}

func TestViewModeConstants(t *testing.T) {
	// Test that view mode constants are unique
	modes := []ViewMode{ViewStatus, ViewLogs, ViewConfig}
	
	for i, mode1 := range modes {
		for j, mode2 := range modes {
			if i != j && mode1 == mode2 {
				t.Errorf("ViewMode constants are not unique: %v == %v", mode1, mode2)
			}
		}
	}
}

func TestLogEntry(t *testing.T) {
	now := time.Now()
	entry := LogEntry{
		Timestamp: now,
		Level:     "INFO",
		Message:   "Test message",
	}

	if entry.Timestamp != now {
		t.Error("LogEntry timestamp not set correctly")
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level 'INFO', got '%s'", entry.Level)
	}

	if entry.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", entry.Message)
	}
}

func TestMultipleLogEntries(t *testing.T) {
	model := New()

	// Add more than 10 logs to test the limit in renderLogsView
	for i := 0; i < 15; i++ {
		model.AddLog("INFO", "Log message "+string(rune('A'+i)))
	}

	if len(model.logs) != 15 {
		t.Errorf("Expected 15 logs, got %d", len(model.logs))
	}

	view := model.renderLogsView()
	
	// The view should contain some of the recent logs
	if !strings.Contains(view, "Activity Logs") {
		t.Error("Logs view missing 'Activity Logs'")
	}
}
