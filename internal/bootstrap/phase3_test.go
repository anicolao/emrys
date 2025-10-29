package bootstrap

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetVoiceConfigPath(t *testing.T) {
	path := GetVoiceConfigPath()

	if path == "" {
		t.Error("Expected non-empty voice config path")
	}

	if !strings.Contains(path, ".config/emrys/voice.conf") {
		t.Errorf("Expected path to contain '.config/emrys/voice.conf', got %s", path)
	}
}

func TestCreateVoiceConfig(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create voice config
	err := CreateVoiceConfig()
	if err != nil {
		t.Fatalf("CreateVoiceConfig failed: %v", err)
	}

	// Verify the config file was created
	configPath := GetVoiceConfigPath()
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
		"# Emrys Voice Output Configuration",
		"enabled =",
		"voice = Jamie",
		"rate =",
		"volume =",
		"quiet_hours =",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(configStr, expected) {
			t.Errorf("Config file missing expected content: %s", expected)
		}
	}

	// Test idempotency - creating again should not fail
	err = CreateVoiceConfig()
	if err != nil {
		t.Errorf("CreateVoiceConfig should be idempotent, but failed on second call: %v", err)
	}
}

func TestUpdateNixDarwinConfigForVoice(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create a mock nix-darwin configuration
	nixpkgsDir := filepath.Join(tmpDir, ".nixpkgs")
	err := os.MkdirAll(nixpkgsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .nixpkgs directory: %v", err)
	}

	configPath := filepath.Join(nixpkgsDir, "darwin-configuration.nix")
	mockConfig := `{ config, pkgs, lib, ... }:

{
  system.primaryUser = "testuser";
  
  environment.systemPackages = with pkgs; [
    vim
    git
  ];
}`

	err = os.WriteFile(configPath, []byte(mockConfig), 0644)
	if err != nil {
		t.Fatalf("Failed to create mock config: %v", err)
	}

	// Update configuration for voice
	err = UpdateNixDarwinConfigForVoice()
	if err != nil {
		t.Fatalf("UpdateNixDarwinConfigForVoice failed: %v", err)
	}

	// Read the updated configuration
	content, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read updated config: %v", err)
	}

	configStr := string(content)

	// Verify voice configuration was added
	if !strings.Contains(configStr, "# Phase 3: Voice Output Configuration") {
		t.Error("Configuration missing voice output section")
	}

	if !strings.Contains(configStr, "Jamie") {
		t.Error("Configuration missing Jamie voice reference")
	}

	// Test idempotency - updating again should not fail
	err = UpdateNixDarwinConfigForVoice()
	if err != nil {
		t.Errorf("UpdateNixDarwinConfigForVoice should be idempotent, but failed on second call: %v", err)
	}
}

func TestIsPhase3Complete(t *testing.T) {
	// This test will likely fail on non-macOS systems or systems without Jamie voice
	// We're mainly testing that the function doesn't panic

	complete := IsPhase3Complete()

	// The result depends on the system state, so we just verify it returns a bool
	_ = complete

	// Log the result for debugging
	t.Logf("IsPhase3Complete returned: %v", complete)
}

func TestDefaultVoiceConstant(t *testing.T) {
	if DefaultVoice != "Jamie" {
		t.Errorf("Expected DefaultVoice to be 'Jamie', got '%s'", DefaultVoice)
	}
}

func TestListAvailableVoices(t *testing.T) {
	// This test requires macOS 'say' command
	t.Skip("Skipping voice listing test (requires macOS)")

	err := ListAvailableVoices()
	if err != nil {
		t.Errorf("ListAvailableVoices failed: %v", err)
	}
}

func TestTestVoiceOutput(t *testing.T) {
	// This test requires macOS 'say' command and Jamie voice
	t.Skip("Skipping voice output test (requires macOS and Jamie voice)")

	err := TestVoiceOutput()
	if err != nil {
		t.Errorf("TestVoiceOutput failed: %v", err)
	}
}

func TestInstallJamieVoice(t *testing.T) {
	// This test requires user interaction and macOS
	t.Skip("Skipping Jamie voice installation test (requires macOS and user interaction)")

	err := InstallJamieVoice()
	if err != nil {
		t.Logf("InstallJamieVoice returned error (expected if Jamie not installed): %v", err)
	}
}

func TestCreateVoiceConfigDirectory(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create voice config (which should create the directory)
	err := CreateVoiceConfig()
	if err != nil {
		t.Fatalf("CreateVoiceConfig failed: %v", err)
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

func TestVoiceConfigPermissions(t *testing.T) {
	// Use a temporary directory for testing
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", originalHome)

	// Create voice config
	err := CreateVoiceConfig()
	if err != nil {
		t.Fatalf("CreateVoiceConfig failed: %v", err)
	}

	// Check file permissions
	configPath := GetVoiceConfigPath()
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
