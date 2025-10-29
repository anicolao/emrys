package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anicolao/emrys/internal/nixdarwin"
)

// Phase1Packages are the packages required for Phase 1 of the bootstrap
var Phase1Packages = []string{
	"ollama",
	"tmux",
	"go",
	"jq",
}

// IsPhase1Complete checks if all Phase 1 packages are installed
func IsPhase1Complete() bool {
	for _, pkg := range Phase1Packages {
		if !isPackageInstalled(pkg) {
			return false
		}
	}
	return true
}

// isPackageInstalled checks if a package is available in the system PATH
func isPackageInstalled(packageName string) bool {
	_, err := exec.LookPath(packageName)
	return err == nil
}

// GetMissingPackages returns a list of packages that are not yet installed
func GetMissingPackages() []string {
	var missing []string
	for _, pkg := range Phase1Packages {
		if !isPackageInstalled(pkg) {
			missing = append(missing, pkg)
		}
	}
	return missing
}

// UpdateNixDarwinConfiguration updates the nix-darwin configuration to include Phase 1 packages
func UpdateNixDarwinConfiguration() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".nixpkgs", "darwin-configuration.nix")

	// Read the current configuration
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	configStr := string(content)

	// Track if any changes were made
	configChanged := false

	// Find the environment.systemPackages section and add our packages
	// We'll add them after the existing packages
	packagesSection := `  # Basic system packages
  environment.systemPackages = with pkgs; [
    vim
    git
    curl
    wget
  ];`

	updatedPackagesSection := `  # Basic system packages
  environment.systemPackages = with pkgs; [
    vim
    git
    curl
    wget
    
    # Phase 1 Bootstrap Packages
    ollama
    tmux
    go
    jq
  ];`

	// Add SSH server configuration for remote access
	sshConfig := `
  # SSH server configuration for remote access
  # Enable Remote Login in macOS
  services.openssh.enable = true;`

	// Also add auto-login configuration for dedicated hardware
	// This enables automatic recovery after power outages
	autoLoginConfig := `
  # Auto-login configuration for dedicated Mac Mini
  # Emrys is designed to run on dedicated, physically secure hardware
  system.defaults.loginwindow = {
    autoLoginUser = "__EMRYS_USERNAME__";
  };`

	// Check if SSH config already exists
	if !strings.Contains(configStr, "services.openssh") {
		// Insert SSH config before the closing brace
		configStr = strings.Replace(configStr, "\n}", sshConfig+"\n}", 1)
		configChanged = true
	}

	// Check if auto-login config already exists
	if !strings.Contains(configStr, "Auto-login configuration") {
		// Insert auto-login config before the closing brace
		configStr = strings.Replace(configStr, "\n}", autoLoginConfig+"\n}", 1)
		configChanged = true
	}

	// Check if we need to add Phase 1 packages
	if !strings.Contains(configStr, "# Phase 1 Bootstrap Packages") {
		// Replace the packages section
		configStr = strings.Replace(configStr, packagesSection, updatedPackagesSection, 1)
		configChanged = true
	}

	// If no changes were made, we're already up to date
	if !configChanged {
		fmt.Println("✓ Configuration already includes Phase 1 packages")
		return nil
	}

	// Get the username from the existing configuration
	// Look for system.primaryUser = "username"; and extract it
	username := ""
	lines := strings.Split(configStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "system.primaryUser") {
			// Extract username from: system.primaryUser = "username";
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				username = parts[1]
				break
			}
		}
	}

	// If we couldn't find it in the config, get it from the environment
	if username == "" {
		username = os.Getenv("USER")
		if username == "" {
			// Fallback to getting username from home directory path
			username = filepath.Base(homeDir)
		}
	}

	// Replace the username placeholder in auto-login configuration
	configStr = strings.Replace(configStr, "__EMRYS_USERNAME__", username, -1)

	// Write the updated configuration
	if err := os.WriteFile(configPath, []byte(configStr), 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Printf("✓ Updated configuration at %s\n", configPath)
	return nil
}

// VerifyPackageInstallation verifies that all Phase 1 packages are installed
func VerifyPackageInstallation() error {
	fmt.Println("Verifying package installation...")

	missing := GetMissingPackages()
	if len(missing) > 0 {
		return fmt.Errorf("some packages are still missing: %s", strings.Join(missing, ", "))
	}

	fmt.Println("✓ All Phase 1 packages verified:")
	for _, pkg := range Phase1Packages {
		path, _ := exec.LookPath(pkg)
		fmt.Printf("  - %-10s %s\n", pkg, path)
	}

	return nil
}

// RunPhase1 executes the complete Phase 1 bootstrap process
func RunPhase1() error {
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  Phase 1: Package Installation")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Check if Phase 1 is already complete
	if IsPhase1Complete() {
		fmt.Println("✓ Phase 1 is already complete!")
		fmt.Println()
		if err := VerifyPackageInstallation(); err != nil {
			return err
		}
		return nil
	}

	// Show what packages are missing
	missing := GetMissingPackages()
	if len(missing) > 0 {
		fmt.Println("Missing packages:")
		for _, pkg := range missing {
			fmt.Printf("  - %s\n", pkg)
		}
		fmt.Println()
	}

	// Step 1: Update the nix-darwin configuration
	fmt.Println("Step 1: Updating nix-darwin configuration...")
	if err := UpdateNixDarwinConfiguration(); err != nil {
		return fmt.Errorf("failed to update configuration: %w", err)
	}
	fmt.Println()

	// Step 2: Apply the configuration
	fmt.Println("Step 2: Applying configuration...")
	if err := nixdarwin.ApplyConfiguration(); err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}
	fmt.Println()

	// Step 3: Verify installation
	fmt.Println("Step 3: Verifying installation...")
	if err := VerifyPackageInstallation(); err != nil {
		return fmt.Errorf("verification failed: %w", err)
	}
	fmt.Println()

	fmt.Println("═══════════════════════════════════════")
	fmt.Println("✓ Phase 1 Bootstrap Complete!")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	return nil
}
