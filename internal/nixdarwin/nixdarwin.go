package nixdarwin

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// IsInstalled checks if nix-darwin is installed on the system
func IsInstalled() bool {
	// Only check if darwin-rebuild command exists
	// Configuration files may exist even if nix-darwin installation failed
	_, err := exec.LookPath("darwin-rebuild")
	return err == nil
}

// IsNixInstalled checks if Nix is installed on the system
func IsNixInstalled() bool {
	_, err := exec.LookPath("nix")
	return err == nil
}

// InstallNix installs Nix on the system
func InstallNix() error {
	fmt.Println("Installing Nix (Lix)...")
	fmt.Println("This will require sudo access and may take several minutes.")

	// Use the Lix installer
	cmd := exec.Command("sh", "-c", "curl -sSf -L https://install.lix.systems/lix | sh -s -- install")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install Nix: %w", err)
	}

	fmt.Println("✓ Nix installed successfully")
	return nil
}

// InstallNixDarwin installs nix-darwin with the provided configuration
// Deprecated: Use InstallNixDarwinWithFlake instead
func InstallNixDarwin(configPath string) error {
	fmt.Println("Installing nix-darwin...")

	// First, ensure the configuration is in the right place
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	nixpkgsDir := filepath.Join(homeDir, ".nixpkgs")
	if err := os.MkdirAll(nixpkgsDir, 0755); err != nil {
		return fmt.Errorf("failed to create .nixpkgs directory: %w", err)
	}

	// Copy the configuration file
	destConfig := filepath.Join(nixpkgsDir, "darwin-configuration.nix")
	if err := copyFile(configPath, destConfig); err != nil {
		return fmt.Errorf("failed to copy configuration: %w", err)
	}

	fmt.Printf("✓ Configuration copied to %s\n", destConfig)

	// Note: This function is deprecated. Flake-based installation is now required.
	return fmt.Errorf("legacy installation method no longer supported, please use flake-based installation")
}

// InstallNixDarwinWithConfig installs nix-darwin with the provided configuration content
func InstallNixDarwinWithConfig(configContent string) error {
	fmt.Println("Installing nix-darwin...")

	// First, ensure the configuration is in the right place
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	nixpkgsDir := filepath.Join(homeDir, ".nixpkgs")
	if err := os.MkdirAll(nixpkgsDir, 0755); err != nil {
		return fmt.Errorf("failed to create .nixpkgs directory: %w", err)
	}

	// Write the configuration content to file
	destConfig := filepath.Join(nixpkgsDir, "darwin-configuration.nix")
	if err := os.WriteFile(destConfig, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Printf("✓ Configuration written to %s\n", destConfig)

	// Run nix-darwin installation using the flake-based installer
	// We need to source nix before running nix commands
	installCmd := `
		set -e
		if [ -e '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh' ]; then
			. '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh'
		fi
		nix run nix-darwin -- switch --flake ~/.nixpkgs
	`

	cmd := exec.Command("sh", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = homeDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install nix-darwin: %w", err)
	}

	fmt.Println("✓ nix-darwin installed successfully")
	return nil
}

// InstallNixDarwinWithFlake installs nix-darwin with the provided configuration and flake content
func InstallNixDarwinWithFlake(configContent, flakeContent string) error {
	fmt.Println("Installing nix-darwin...")

	// First, ensure the configuration is in the right place
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	nixpkgsDir := filepath.Join(homeDir, ".nixpkgs")
	if err := os.MkdirAll(nixpkgsDir, 0755); err != nil {
		return fmt.Errorf("failed to create .nixpkgs directory: %w", err)
	}

	// Write the configuration content to file
	destConfig := filepath.Join(nixpkgsDir, "darwin-configuration.nix")
	if err := os.WriteFile(destConfig, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	// Write the flake.nix content to file
	destFlake := filepath.Join(nixpkgsDir, "flake.nix")
	if err := os.WriteFile(destFlake, []byte(flakeContent), 0644); err != nil {
		return fmt.Errorf("failed to write flake.nix: %w", err)
	}

	fmt.Printf("✓ Configuration written to %s\n", destConfig)
	fmt.Printf("✓ Flake written to %s\n", destFlake)

	// Run nix-darwin installation using the flake-based installer
	// We need to source nix before running nix commands
	installCmd := `
		set -e
		if [ -e '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh' ]; then
			. '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh'
		fi
		nix run nix-darwin -- switch --flake ~/.nixpkgs#emrys
	`

	cmd := exec.Command("sh", "-c", installCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = homeDir

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install nix-darwin: %w", err)
	}

	fmt.Println("✓ nix-darwin installed successfully")
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	err = os.WriteFile(dst, input, 0644)
	if err != nil {
		return err
	}

	return nil
}

// ApplyConfiguration applies the nix-darwin configuration
func ApplyConfiguration() error {
	fmt.Println("Applying nix-darwin configuration...")

	// Source nix and run darwin-rebuild
	applyCmd := `
		set -e
		if [ -e '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh' ]; then
			. '/nix/var/nix/profiles/default/etc/profile.d/nix-daemon.sh'
		fi
		darwin-rebuild switch
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
