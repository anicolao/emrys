package bootstrap

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

	// Check if we already have the Phase 1 packages
	if strings.Contains(configStr, "# Phase 1 Bootstrap Packages") {
		fmt.Println("✓ Configuration already includes Phase 1 packages")
		return nil
	}

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

	// Also add SSH configuration
	sshConfig := `
  # SSH server configuration for remote access
  services.openssh = {
    enable = true;
    settings = {
      PasswordAuthentication = false;
      PermitRootLogin = "no";
    };
  };`

	// Also add auto-login configuration (commented out by default for security)
	autoLoginConfig := `
  # Auto-login configuration (uncomment to enable)
  # WARNING: Only enable on physically secure, dedicated hardware
  # system.defaults.loginwindow = {
  #   autoLoginUser = "__EMRYS_USERNAME__";
  # };`

	// Check if SSH config already exists
	if !strings.Contains(configStr, "services.openssh") {
		// Insert SSH config before the closing brace
		configStr = strings.Replace(configStr, "\n}", sshConfig+"\n}", 1)
	}

	// Check if auto-login config already exists
	if !strings.Contains(configStr, "Auto-login configuration") {
		// Insert auto-login config before the closing brace
		configStr = strings.Replace(configStr, "\n}", autoLoginConfig+"\n}", 1)
	}

	// Replace the packages section
	configStr = strings.Replace(configStr, packagesSection, updatedPackagesSection, 1)

	// Write the updated configuration
	if err := os.WriteFile(configPath, []byte(configStr), 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Printf("✓ Updated configuration at %s\n", configPath)
	return nil
}

// ApplyConfiguration applies the updated nix-darwin configuration
func ApplyConfiguration() error {
	fmt.Println("Applying nix-darwin configuration...")
	fmt.Println("Note: This may take several minutes and will require sudo access")
	fmt.Println()

	// Source nix and run darwin-rebuild
	applyCmd := `
		set -e
		if [ -e '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh' ]; then
			. '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh'
		fi
		darwin-rebuild switch --flake ~/.nixpkgs#emrys
	`

	cmd := exec.Command("sh", "-c", applyCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}

	fmt.Println("✓ Configuration applied successfully")
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
	if err := ApplyConfiguration(); err != nil {
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
