package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/anicolao/emrys/internal/config"
	"github.com/anicolao/emrys/internal/nixdarwin"
)

func main() {
	fmt.Println("╔════════════════════════════════════════╗")
	fmt.Println("║           Emrys Setup                  ║")
	fmt.Println("║  Your Personal AI Assistant on macOS  ║")
	fmt.Println("╚════════════════════════════════════════╝")
	fmt.Println()

	// Check if nix-darwin is already installed
	if nixdarwin.IsInstalled() {
		fmt.Println("✓ nix-darwin is already installed!")
		fmt.Println()
		fmt.Println("Emrys is ready to use.")
		return
	}

	fmt.Println("⚠ nix-darwin is not installed yet.")
	fmt.Println()
	fmt.Println("Emrys requires nix-darwin for system configuration and package management.")
	fmt.Println("This setup will:")
	fmt.Println("  1. Install Nix (if not already installed)")
	fmt.Println("  2. Install nix-darwin")
	fmt.Println("  3. Apply a basic configuration")
	fmt.Println()

	// Check if we should proceed
	if !confirm("Would you like to proceed with the installation?") {
		fmt.Println("Installation cancelled.")
		return
	}

	fmt.Println()

	// Step 1: Check and install Nix if needed
	if !nixdarwin.IsNixInstalled() {
		fmt.Println("Step 1: Installing Nix...")
		fmt.Println("Note: You may be asked for your password (sudo access required)")
		fmt.Println()

		if err := nixdarwin.InstallNix(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println()
	} else {
		fmt.Println("✓ Nix is already installed")
		fmt.Println()
	}

	// Step 2: Install nix-darwin
	fmt.Println("Step 2: Installing nix-darwin...")

	// Use the embedded configuration and flake
	if err := nixdarwin.InstallNixDarwinWithFlake(config.DefaultNixDarwinConfig, config.DefaultFlakeConfig); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("════════════════════════════════════════")
	fmt.Println("✓ Setup completed successfully!")
	fmt.Println()
	fmt.Println("nix-darwin has been installed and configured.")
	fmt.Println("You may need to restart your terminal for all changes to take effect.")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  - Edit ~/.nixpkgs/darwin-configuration.nix to customize your setup")
	fmt.Println("  - Run 'darwin-rebuild switch' to apply configuration changes")
	fmt.Println("════════════════════════════════════════")
}

// confirm prompts the user for a yes/no confirmation
func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s (y/n): ", prompt)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}

		fmt.Println("Please answer 'y' or 'n'")
	}
}
