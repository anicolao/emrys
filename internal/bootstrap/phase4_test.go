package bootstrap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetTUIConfigPath(t *testing.T) {
	path := GetTUIConfigPath()

	if path == "" {
		t.Error("Expected non-empty TUI config path")
	}

	if !strings.Contains(path, ".config/emrys/tui.conf") {
		t.Errorf("Expected path to contain '.config/emrys/tui.conf', got %s", path)
	}
}

func TestCreateTUIConfig(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create TUI config
	err := CreateTUIConfig()
	if err != nil {
		t.Fatalf("CreateTUIConfig failed: %v", err)
	}

	// Verify the config file was created
	configPath := GetTUIConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("Config file was not created at %s", configPath)
	}

	// Read the config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	configStr := string(content)

	// Verify essential configuration elements are present
	expectedStrings := []string{
		"# Emrys TUI Configuration",
		"enabled =",
		"default_view =",
		"theme =",
		"refresh_interval =",
		"show_resources =",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(configStr, expected) {
			t.Errorf("Config file missing expected content: %s", expected)
		}
	}

	// Test idempotency - creating again should not fail
	err = CreateTUIConfig()
	if err != nil {
		t.Errorf("CreateTUIConfig should be idempotent, but failed on second call: %v", err)
	}
}

func TestIsPhase4Complete(t *testing.T) {
	// This test checks that the function doesn't panic
	complete := IsPhase4Complete()

	// The result depends on the system state, so we just verify it returns a bool
	_ = complete

	// Log the result for debugging
	t.Logf("IsPhase4Complete returned: %v", complete)
}

func TestCreateTUIConfigDirectory(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create TUI config (which should create the directory)
	err := CreateTUIConfig()
	if err != nil {
		t.Fatalf("CreateTUIConfig failed: %v", err)
	}

	// Verify the directory exists
	configDir := filepath.Join(tmpDir, ".config", "emrys")
	info, err := os.Stat(configDir)
	if err != nil {
		t.Errorf("Config directory was not created: %v", err)
	} else if !info.IsDir() {
		t.Error("Config path exists but is not a directory")
	}
}

func TestTUIConfigPermissions(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create TUI config
	err := CreateTUIConfig()
	if err != nil {
		t.Fatalf("CreateTUIConfig failed: %v", err)
	}

	// Check file permissions
	configPath := GetTUIConfigPath()
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}

	// Check that file is readable and writable by owner
	mode := info.Mode()
	if mode.Perm()&0600 != 0600 {
		t.Errorf("Config file has incorrect permissions: %v", mode.Perm())
	}
}

func TestBuildTUIBinary(t *testing.T) {
	// Skip this test in CI as it requires the full source tree
	t.Skip("Skipping TUI binary build test (requires full source tree)")

	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	err := BuildTUIBinary()
	if err != nil {
		t.Logf("BuildTUIBinary failed (expected in test environment): %v", err)
	}
}

func TestVerifyTUIComponents(t *testing.T) {
	// This test will fail if Phase 4 is not complete, which is expected
	t.Skip("Skipping component verification test (depends on Phase 4 completion)")

	err := VerifyTUIComponents()
	if err != nil {
		t.Logf("VerifyTUIComponents failed (expected if Phase 4 not complete): %v", err)
	}
}

func TestTestTUI(t *testing.T) {
	// This test verifies that the TestTUI function works
	// We'll skip actual execution in automated tests
	t.Skip("Skipping interactive TUI test")

	err := TestTUI()
	if err != nil {
		t.Errorf("TestTUI failed: %v", err)
	}
}

func TestLaunchTUI(t *testing.T) {
	// This test requires user interaction
	t.Skip("Skipping interactive TUI launch test")

	err := LaunchTUI()
	if err != nil {
		t.Errorf("LaunchTUI failed: %v", err)
	}
}

func TestTUIConfigContent(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create TUI config
	err := CreateTUIConfig()
	if err != nil {
		t.Fatalf("CreateTUIConfig failed: %v", err)
	}

	// Read the config file
	configPath := GetTUIConfigPath()
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	configStr := string(content)

	// Verify specific default values
	expectedValues := map[string]bool{
		"enabled = true":           false,
		"default_view = status":    false,
		"theme = auto":             false,
		"refresh_interval = 5":     false,
		"show_resources = true":    false,
		"log_retention = 7":        false,
		"max_log_entries = 100":    false,
	}

	for expected := range expectedValues {
		if strings.Contains(configStr, expected) {
			expectedValues[expected] = true
		}
	}

	// Check that all expected values were found
	for expected, found := range expectedValues {
		if !found {
			t.Errorf("Config file missing expected value: %s", expected)
		}
	}
}
