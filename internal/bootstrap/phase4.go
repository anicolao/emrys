package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/anicolao/emrys/internal/tui"
)

// IsPhase4Complete checks if Phase 4 is complete
func IsPhase4Complete() bool {
	// Check if the TUI binary exists in the expected location
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	tuiBinaryPath := filepath.Join(homeDir, ".local", "bin", "emrys-tui")
	if _, err := os.Stat(tuiBinaryPath); os.IsNotExist(err) {
		return false
	}

	// Check if TUI configuration exists
	configPath := GetTUIConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// GetTUIConfigPath returns the path to the TUI configuration file
func GetTUIConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "emrys", "tui.conf")
}

// BuildTUIBinary builds the Emrys TUI binary
func BuildTUIBinary() error {
	fmt.Println("Building Emrys TUI application...")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Create the binary directory if it doesn't exist
	binDir := filepath.Join(homeDir, ".local", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create binary directory: %w", err)
	}

	// Get the current working directory (where the source code is)
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Build the TUI binary
	tuiBinaryPath := filepath.Join(binDir, "emrys-tui")
	
	// Check if we have a cmd/emrys-tui directory, if not use cmd/emrys with a special flag
	tuiSourcePath := filepath.Join(cwd, "cmd", "emrys-tui")
	if _, err := os.Stat(tuiSourcePath); os.IsNotExist(err) {
		// For now, we'll create a simple wrapper script that launches the TUI
		// In a future implementation, this could be a separate binary
		scriptContent := fmt.Sprintf(`#!/bin/bash
# Emrys TUI Launcher
# This script launches the Emrys TUI application

echo "Launching Emrys TUI..."
exec go run %s/cmd/emrys --tui
`, cwd)
		
		if err := os.WriteFile(tuiBinaryPath, []byte(scriptContent), 0755); err != nil {
			return fmt.Errorf("failed to create TUI launcher script: %w", err)
		}
		
		fmt.Printf("✓ Created TUI launcher at %s\n", tuiBinaryPath)
		return nil
	}

	// Build the binary
	cmd := exec.Command("go", "build", "-o", tuiBinaryPath, tuiSourcePath)
	cmd.Dir = cwd
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build TUI binary: %w\nOutput: %s", err, string(output))
	}

	fmt.Printf("✓ Built TUI binary at %s\n", tuiBinaryPath)
	return nil
}

// CreateTUIConfig creates the TUI configuration file
func CreateTUIConfig() error {
	configPath := GetTUIConfigPath()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("✓ TUI configuration already exists at %s\n", configPath)
		return nil
	}

	// Create default configuration
	configContent := `# Emrys TUI Configuration
# This file contains settings for the Terminal User Interface

# Enable TUI on startup (true/false)
enabled = true

# Default view mode (status, logs, config)
default_view = status

# Color theme (auto, light, dark)
theme = auto

# Refresh interval in seconds
refresh_interval = 5

# Show system resources (true/false)
show_resources = true

# Log retention in days
log_retention = 7

# Maximum log entries to display
max_log_entries = 100
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Printf("✓ Created TUI configuration at %s\n", configPath)
	return nil
}

// TestTUI tests the TUI application by running it briefly
func TestTUI() error {
	fmt.Println("Testing TUI application...")
	fmt.Println()

	// Create a new TUI model
	model := tui.New()
	
	// Check if the model initializes correctly
	if model.Init() == nil {
		fmt.Println("✓ TUI model initialized successfully")
	} else {
		return fmt.Errorf("failed to initialize TUI model")
	}

	// Test rendering
	view := model.View()
	if len(view) > 0 {
		fmt.Println("✓ TUI rendering works")
	} else {
		return fmt.Errorf("TUI rendering failed")
	}

	// Display a preview of the TUI
	fmt.Println()
	fmt.Println("TUI Preview:")
	fmt.Println("─────────────────────────────────────────")
	fmt.Println(view)
	fmt.Println("─────────────────────────────────────────")
	fmt.Println()

	return nil
}

// LaunchTUI launches the TUI application in interactive mode
func LaunchTUI() error {
	fmt.Println("Launching Emrys TUI...")
	fmt.Println("Press 'q' or Ctrl+C to exit")
	fmt.Println()

	// Give the user a moment to read the message
	time.Sleep(2 * time.Second)

	// Create and start the TUI
	model := tui.New()
	
	// Create the program
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	// Run the program
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("TUI application error: %w", err)
	}

	fmt.Println()
	fmt.Println("TUI application exited")
	return nil
}

// VerifyTUIComponents verifies that all TUI components are working
func VerifyTUIComponents() error {
	fmt.Println("Verifying TUI components...")

	// Check if the TUI binary exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	tuiBinaryPath := filepath.Join(homeDir, ".local", "bin", "emrys-tui")
	if _, err := os.Stat(tuiBinaryPath); os.IsNotExist(err) {
		return fmt.Errorf("TUI binary not found at %s", tuiBinaryPath)
	}

	fmt.Printf("✓ TUI binary found at %s\n", tuiBinaryPath)

	// Check if the TUI configuration exists
	configPath := GetTUIConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("TUI configuration not found at %s", configPath)
	}

	fmt.Printf("✓ TUI configuration found at %s\n", configPath)

	// Test that we can create a TUI model
	model := tui.New()
	if model.Init() == nil {
		fmt.Println("✓ TUI model can be initialized")
	}

	return nil
}

// RunPhase4 executes the complete Phase 4 bootstrap process
func RunPhase4() error {
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  Phase 4: TUI Application Development")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Check if Phase 4 is already complete
	if IsPhase4Complete() {
		fmt.Println("✓ Phase 4 is already complete!")
		fmt.Println()
		if err := VerifyTUIComponents(); err != nil {
			fmt.Printf("Warning: Component verification failed: %v\n", err)
		}
		return nil
	}

	// Step 1: Create TUI configuration
	fmt.Println("Step 1: Creating TUI configuration...")
	if err := CreateTUIConfig(); err != nil {
		return fmt.Errorf("failed to create TUI configuration: %w", err)
	}
	fmt.Println()

	// Step 2: Build TUI binary
	fmt.Println("Step 2: Building TUI binary...")
	if err := BuildTUIBinary(); err != nil {
		return fmt.Errorf("failed to build TUI binary: %w", err)
	}
	fmt.Println()

	// Step 3: Test TUI application
	fmt.Println("Step 3: Testing TUI application...")
	if err := TestTUI(); err != nil {
		return fmt.Errorf("TUI test failed: %w", err)
	}
	fmt.Println()

	// Step 4: Verify all components
	fmt.Println("Step 4: Verifying TUI components...")
	if err := VerifyTUIComponents(); err != nil {
		return fmt.Errorf("component verification failed: %w", err)
	}
	fmt.Println()

	fmt.Println("═══════════════════════════════════════")
	fmt.Println("✓ Phase 4 Bootstrap Complete!")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()
	fmt.Printf("TUI configuration saved to: %s\n", GetTUIConfigPath())
	fmt.Println()
	fmt.Println("TUI features:")
	fmt.Println("  - Status dashboard with Ollama and voice status")
	fmt.Println("  - Activity log viewer")
	fmt.Println("  - Configuration interface")
	fmt.Println("  - Responsive layout with Lipgloss styling")
	fmt.Println("  - Multiple view modes (Status, Logs, Config)")
	fmt.Println()
	fmt.Println("You can launch the TUI with: emrys-tui")
	fmt.Println()

	return nil
}
