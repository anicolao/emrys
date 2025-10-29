package bootstrap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsPackageInstalled(t *testing.T) {
	// Test with a package that should always exist on most systems
	result := isPackageInstalled("sh")
	if !result {
		t.Error("Expected 'sh' to be installed, but it wasn't found")
	}

	// Test with a package that definitely doesn't exist
	result = isPackageInstalled("this-package-definitely-does-not-exist-xyz123")
	if result {
		t.Error("Expected non-existent package to return false, but it returned true")
	}
}

func TestGetMissingPackages(t *testing.T) {
	// This test just verifies the function runs without crashing
	missing := GetMissingPackages()
	t.Logf("Missing packages: %v", missing)

	// The result should be a valid slice (even if empty)
	if missing == nil {
		t.Error("GetMissingPackages returned nil instead of a slice")
	}
}

func TestIsPhase1Complete(t *testing.T) {
	// This test verifies the function runs without crashing
	result := IsPhase1Complete()
	t.Logf("IsPhase1Complete returned: %v", result)
}

func TestUpdateNixDarwinConfiguration(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	nixpkgsDir := filepath.Join(tmpDir, ".nixpkgs")
	if err := os.MkdirAll(nixpkgsDir, 0755); err != nil {
		t.Fatalf("Failed to create test directory: %v", err)
	}

	configPath := filepath.Join(nixpkgsDir, "darwin-configuration.nix")

	// Create a minimal test configuration
	testConfig := `{ config, pkgs, lib, ... }:

{
  system.primaryUser = "testuser";
  nixpkgs.hostPlatform = lib.mkDefault "aarch64-darwin";
  system.stateVersion = 5;
  nix.settings.experimental-features = [ "nix-command" "flakes" ];
  security.pam.services.sudo_local.touchIdAuth = true;

  # Basic system packages
  environment.systemPackages = with pkgs; [
    vim
    git
    curl
    wget
  ];

  system.defaults = {
    dock.autohide = true;
    finder.AppleShowAllExtensions = true;
    NSGlobalDomain.AppleShowAllExtensions = true;
  };

  nix.optimise.automatic = true;
  nix.gc = {
    automatic = true;
    interval = { Weekday = 0; Hour = 0; Minute = 0; };
    options = "--delete-older-than 30d";
  };
}
`

	if err := os.WriteFile(configPath, []byte(testConfig), 0644); err != nil {
		t.Fatalf("Failed to create test configuration: %v", err)
	}

	// Temporarily change HOME to our test directory
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Test updating the configuration
	err := UpdateNixDarwinConfiguration()
	if err != nil {
		t.Fatalf("UpdateNixDarwinConfiguration failed: %v", err)
	}

	// Read the updated configuration
	updatedContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read updated configuration: %v", err)
	}

	updatedStr := string(updatedContent)

	// Verify that Phase 1 packages were added
	if !contains(updatedStr, "ollama") {
		t.Error("Updated configuration doesn't contain 'ollama'")
	}
	if !contains(updatedStr, "tmux") {
		t.Error("Updated configuration doesn't contain 'tmux'")
	}
	if !contains(updatedStr, "go") {
		t.Error("Updated configuration doesn't contain 'go'")
	}
	if !contains(updatedStr, "jq") {
		t.Error("Updated configuration doesn't contain 'jq'")
	}
	if !contains(updatedStr, "# Phase 1 Bootstrap Packages") {
		t.Error("Updated configuration doesn't contain Phase 1 comment marker")
	}

	// Verify SSH configuration was added
	if !contains(updatedStr, "services.openssh") {
		t.Error("Updated configuration doesn't contain SSH configuration")
	}

	// Verify auto-login configuration was added and is enabled (not commented out)
	if !contains(updatedStr, "system.defaults.loginwindow") {
		t.Error("Updated configuration doesn't contain auto-login configuration")
	}
	if !contains(updatedStr, "autoLoginUser = \"testuser\";") {
		t.Error("Updated configuration doesn't contain enabled auto-login user")
	}
	if contains(updatedStr, "# autoLoginUser") {
		t.Error("Auto-login configuration is commented out (should be enabled)")
	}

	// Run again to test idempotency
	err = UpdateNixDarwinConfiguration()
	if err != nil {
		t.Fatalf("Second UpdateNixDarwinConfiguration failed: %v", err)
	}

	// Verify configuration hasn't been duplicated
	secondContent, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read configuration after second update: %v", err)
	}

	if string(secondContent) != updatedStr {
		t.Error("Configuration was modified on second run (should be idempotent)")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 &&
		(s == substr || len(s) >= len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
