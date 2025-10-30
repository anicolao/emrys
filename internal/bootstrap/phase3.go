package bootstrap

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anicolao/emrys/internal/nixdarwin"
	"github.com/anicolao/emrys/internal/voice"
)

// DefaultVoice is the default voice for Emrys
const DefaultVoice = "Jamie"

// IsPhase3Complete checks if Phase 3 is complete
func IsPhase3Complete() bool {
	// Check if Jamie voice is available
	if !voice.IsVoiceAvailable(DefaultVoice) {
		return false
	}

	// Check if voice configuration exists
	configPath := GetVoiceConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return false
	}

	return true
}

// GetVoiceConfigPath returns the path to the voice configuration file
func GetVoiceConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".config", "emrys", "voice.conf")
}

// UpdateNixDarwinConfigForVoice updates the nix-darwin configuration to install Jamie voice
func UpdateNixDarwinConfigForVoice() error {
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

	// Check if voice configuration already exists
	if strings.Contains(configStr, "# Phase 3: Voice Output Configuration") {
		fmt.Println("✓ Configuration already includes voice setup")
		return nil
	}

	// Add voice configuration before the closing brace
	voiceConfig := `
  # Phase 3: Voice Output Configuration
  # Jamie (Premium) voice installation is facilitated during Phase 3 bootstrap.
  # AppleScript opens VoiceOver Utility to the voice download section where users
  # can download the Jamie voice. The bootstrap process guides users through this.
`

	// Insert voice config before the closing brace
	configStr = strings.Replace(configStr, "\n}", voiceConfig+"\n}", 1)

	// Write the updated configuration
	if err := os.WriteFile(configPath, []byte(configStr), 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Printf("✓ Updated configuration at %s\n", configPath)
	return nil
}

// InstallJamieVoice checks if Jamie voice is installed and installs it programmatically if not
func InstallJamieVoice() error {
	fmt.Println("Checking Jamie voice installation...")

	// Check if Jamie voice is available
	if voice.IsVoiceAvailable(DefaultVoice) {
		fmt.Printf("✓ Jamie voice is already installed\n")
		return nil
	}

	// Jamie voice is not installed, install it using AppleScript
	fmt.Println()
	fmt.Println("⚠ Jamie voice is not installed on this system")
	fmt.Println("Opening VoiceOver Utility to install Jamie voice...")
	fmt.Println()

	// Try to install the voice using AppleScript
	if err := installVoiceUsingAppleScript(); err != nil {
		// If AppleScript fails, provide manual instructions
		fmt.Println()
		fmt.Printf("⚠ Could not open VoiceOver Utility automatically: %v\n", err)
		fmt.Println()
		fmt.Println("Please install Jamie (Premium) voice manually:")
		fmt.Println()
		fmt.Println("To install Jamie voice:")
		fmt.Println("  1. Open VoiceOver Utility (in /System/Applications/Utilities/)")
		fmt.Println("  2. Go to the 'Speech' section")
		fmt.Println("  3. Click on the 'Voices' tab")
		fmt.Println("  4. Find 'Jamie' in the voice list (under English (United Kingdom))")
		fmt.Println("  5. Click the download icon (cloud with down arrow) next to Jamie")
		fmt.Println("  6. Wait for the download to complete (may take several minutes)")
		fmt.Println()
	}

	// Ask if user has completed the installation
	if !confirmVoiceInstallation() {
		return fmt.Errorf("Jamie voice installation required for Phase 3")
	}

	// Check again after user confirms
	if !voice.IsVoiceAvailable(DefaultVoice) {
		fmt.Println()
		fmt.Println("⚠ Jamie voice is still not available")
		fmt.Println("Please install the voice and run this command again.")
		return fmt.Errorf("Jamie voice not found")
	}

	fmt.Println("✓ Jamie voice is now available")
	return nil
}

// installVoiceUsingAppleScript opens VoiceOver Utility to install Jamie voice using AppleScript
func installVoiceUsingAppleScript() error {
	fmt.Println("Opening VoiceOver Utility to download Jamie voice...")

	// The VoiceOver Utility is where voices are actually installed
	appleScriptCode := `
	tell application "VoiceOver Utility"
		activate
	end tell
	`

	// Execute the AppleScript
	cmd := exec.Command("osascript", "-e", appleScriptCode)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to open VoiceOver Utility: %w (output: %s)", err, string(output))
	}

	fmt.Println()
	fmt.Println("✓ VoiceOver Utility opened")
	fmt.Println()
	fmt.Println("To install Jamie voice:")
	fmt.Println("  1. In the VoiceOver Utility window, go to the 'Speech' section")
	fmt.Println("  2. Click on the 'Voices' tab")
	fmt.Println("  3. Find 'Jamie' in the voice list (under English (United Kingdom))")
	fmt.Println("  4. Click the download icon (cloud with down arrow) next to Jamie")
	fmt.Println("  5. Wait for the download to complete (may take several minutes)")
	fmt.Println("  6. Once downloaded, you can close the VoiceOver Utility")
	fmt.Println()

	return nil
}

// getMacOSVersion returns the major version number of macOS
func getMacOSVersion() (int, error) {
	cmd := exec.Command("sw_vers", "-productVersion")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get macOS version: %w", err)
	}

	versionStr := strings.TrimSpace(string(output))
	parts := strings.Split(versionStr, ".")
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid version string: %s", versionStr)
	}

	majorVersion := 0
	fmt.Sscanf(parts[0], "%d", &majorVersion)
	return majorVersion, nil
}

// confirmVoiceInstallation prompts the user to confirm voice installation
func confirmVoiceInstallation() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Have you installed the Jamie voice? (y/n): ")
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

// TestVoiceOutput tests the voice output with a confirmation phrase
func TestVoiceOutput() error {
	fmt.Println("Testing voice output...")
	fmt.Println()

	// Create a test message
	testMessage := "Hello! I am Emrys, your personal AI assistant. Voice output is working correctly."

	// Test the voice
	if err := voice.Test(DefaultVoice); err != nil {
		return fmt.Errorf("voice test failed: %w", err)
	}

	fmt.Println("✓ Voice output test successful")
	fmt.Println()

	// Speak the test message
	fmt.Printf("Speaking: \"%s\"\n", testMessage)
	fmt.Println()

	config := voice.DefaultConfig()
	config.Voice = DefaultVoice
	speaker := voice.NewSpeaker(config)
	defer speaker.Close()

	if err := speaker.SpeakSync(testMessage); err != nil {
		return fmt.Errorf("failed to speak message: %w", err)
	}

	return nil
}

// CreateVoiceConfig creates the default voice configuration file
func CreateVoiceConfig() error {
	configPath := GetVoiceConfigPath()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config already exists
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("✓ Voice configuration already exists at %s\n", configPath)
		return nil
	}

	// Create default configuration
	config := voice.DefaultConfig()
	config.Voice = DefaultVoice

	// Write configuration file
	configContent := fmt.Sprintf(`# Emrys Voice Output Configuration
# This file contains settings for text-to-speech output

# Enable or disable voice output (true/false)
enabled = %t

# Voice name (e.g., Jamie, Samantha, Alex)
voice = %s

# Speech rate in words per minute (typical range: 150-250)
rate = %d

# Volume from 0.0 to 1.0 (note: controlled via system volume)
volume = %.1f

# Enable quiet hours (true/false)
quiet_hours = %t

# Quiet hours start (24-hour format, 0-23)
quiet_start = %d

# Quiet hours end (24-hour format, 0-23)
quiet_end = %d
`,
		config.Enabled,
		config.Voice,
		config.Rate,
		config.Volume,
		config.QuietHours,
		config.QuietStart,
		config.QuietEnd,
	)

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return fmt.Errorf("failed to write configuration: %w", err)
	}

	fmt.Printf("✓ Created voice configuration at %s\n", configPath)
	return nil
}

// ListAvailableVoices lists all available voices on the system
func ListAvailableVoices() error {
	fmt.Println("Available voices on this system:")
	fmt.Println()

	voices, err := voice.ListAvailableVoices()
	if err != nil {
		return fmt.Errorf("failed to list voices: %w", err)
	}

	if len(voices) == 0 {
		fmt.Println("  No voices found")
		return nil
	}

	for i, v := range voices {
		if v == DefaultVoice {
			fmt.Printf("  %d. %s ✓ (default)\n", i+1, v)
		} else {
			fmt.Printf("  %d. %s\n", i+1, v)
		}
	}

	fmt.Println()
	return nil
}

// RunPhase3 executes the complete Phase 3 bootstrap process
func RunPhase3() error {
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  Phase 3: Voice Output Configuration")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Check if Phase 3 is already complete
	if IsPhase3Complete() {
		fmt.Println("✓ Phase 3 is already complete!")
		fmt.Println()
		if err := TestVoiceOutput(); err != nil {
			fmt.Printf("Warning: Voice test failed: %v\n", err)
		}
		return nil
	}

	// Step 1: Update nix-darwin configuration
	fmt.Println("Step 1: Updating nix-darwin configuration...")
	if err := UpdateNixDarwinConfigForVoice(); err != nil {
		return fmt.Errorf("failed to update configuration: %w", err)
	}
	fmt.Println()

	// Step 2: Apply the configuration
	fmt.Println("Step 2: Applying configuration...")
	if err := nixdarwin.ApplyConfiguration(); err != nil {
		return fmt.Errorf("failed to apply configuration: %w", err)
	}
	fmt.Println()

	// Step 3: Check and install Jamie voice
	fmt.Println("Step 3: Installing Jamie voice...")
	if err := InstallJamieVoice(); err != nil {
		return fmt.Errorf("failed to install Jamie voice: %w", err)
	}
	fmt.Println()

	// Step 4: List available voices
	fmt.Println("Step 4: Listing available voices...")
	if err := ListAvailableVoices(); err != nil {
		return fmt.Errorf("failed to list voices: %w", err)
	}
	fmt.Println()

	// Step 5: Create voice configuration
	fmt.Println("Step 5: Creating voice configuration...")
	if err := CreateVoiceConfig(); err != nil {
		return fmt.Errorf("failed to create voice configuration: %w", err)
	}
	fmt.Println()

	// Step 6: Test voice output
	fmt.Println("Step 6: Testing voice output...")
	if err := TestVoiceOutput(); err != nil {
		return fmt.Errorf("voice output test failed: %w", err)
	}
	fmt.Println()

	fmt.Println("═══════════════════════════════════════")
	fmt.Println("✓ Phase 3 Bootstrap Complete!")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()
	fmt.Printf("Voice configuration saved to: %s\n", GetVoiceConfigPath())
	fmt.Printf("Default voice: %s\n", DefaultVoice)
	fmt.Println()
	fmt.Println("Voice output features:")
	fmt.Println("  - Message queuing to prevent overlap")
	fmt.Println("  - Configurable speech rate and volume")
	fmt.Println("  - Quiet hours support")
	fmt.Println("  - Enable/disable voice output on demand")
	fmt.Println()

	return nil
}
