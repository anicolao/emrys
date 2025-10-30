# Emrys TUI Package

This package implements the Terminal User Interface (TUI) for Emrys using the [Bubbletea](https://github.com/charmbracelet/bubbletea) framework and [Lipgloss](https://github.com/charmbracelet/lipgloss) for styling.

## Overview

The Emrys TUI provides an interactive terminal interface for monitoring and interacting with the Emrys AI assistant. It follows the Elm architecture (Model-View-Update pattern) as used by Bubbletea.

## Features

### View Modes

The TUI supports three view modes that users can switch between:

1. **Status View (Press '1')**: Displays system status including:
   - Ollama service status
   - Current AI model
   - Voice output status
   - Last update timestamp
   - System resources (CPU, Memory)

2. **Logs View (Press '2')**: Shows activity logs with:
   - Timestamp for each log entry
   - Log level (INFO, WARN, ERROR)
   - Log message
   - Last 10 log entries displayed

3. **Configuration View (Press '3')**: Displays current configuration:
   - Voice settings (voice name, rate, volume, quiet hours)
   - Model settings (current model, auto-update)
   - Display settings (theme, color scheme)

### Controls

- **1**: Switch to Status view
- **2**: Switch to Logs view
- **3**: Switch to Configuration view
- **q** or **Ctrl+C**: Quit the application

## Architecture

### Model

The `Model` struct holds the application state:

```go
type Model struct {
    width  int
    height int
    
    ollamaStatus string
    modelName    string
    voiceStatus  string
    
    commandInput string
    commandHistory []string
    
    logs []LogEntry
    
    viewMode ViewMode
    
    lastUpdate time.Time
}
```

### View Modes

```go
const (
    ViewStatus ViewMode = iota
    ViewLogs
    ViewConfig
)
```

### Log Entry

```go
type LogEntry struct {
    Timestamp time.Time
    Level     string
    Message   string
}
```

## Usage

### Creating a New TUI

```go
import "github.com/anicolao/emrys/internal/tui"

model := tui.New()
```

### Running the TUI

```go
import tea "github.com/charmbracelet/bubbletea"

p := tea.NewProgram(model, tea.WithAltScreen())
if _, err := p.Run(); err != nil {
    // Handle error
}
```

### Updating Status

```go
model.UpdateStatus("Running", "llama3.2", "Active")
```

### Adding Logs

```go
model.AddLog("INFO", "System started successfully")
model.AddLog("ERROR", "Failed to connect to service")
```

## Styling

The TUI uses Lipgloss for consistent styling:

- **Header**: Bold text with colored background
- **Status Indicators**: Color-coded (Green ✓, Red ✗, Yellow ⚠)
- **Borders**: Rounded borders for content sections
- **Footer**: Dimmed text for non-critical information

### Color Scheme

- Green (✓): Active/Running/OK status
- Red (✗): Stopped/Inactive/Error status
- Yellow (⚠): Warning/Degraded status
- Gray (-): Unknown/Neutral status

## Testing

The package includes comprehensive tests covering:

- Model initialization
- View rendering
- Status updates
- Log management
- Key press handling
- Window resizing

Run tests with:

```bash
go test ./internal/tui -v
```

## Phase 4 Bootstrap

Phase 4 of the Emrys bootstrap process sets up the TUI:

1. Creates TUI configuration file at `~/.config/emrys/tui.conf`
2. Builds the TUI binary at `~/.local/bin/emrys-tui`
3. Tests TUI functionality
4. Verifies all components

## Future Enhancements

Planned features for future releases:

- Command input with history and auto-completion
- Real-time system resource monitoring
- Interactive task management
- Configuration editing within the TUI
- Search and filter for logs
- Keyboard shortcuts customization
- Theme customization
