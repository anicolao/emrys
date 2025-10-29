package bootstrap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// DefaultModel is the model to download and use for Emrys
const DefaultModel = "llama3.2"

// OllamaAPIURL is the URL of the Ollama API
const OllamaAPIURL = "http://localhost:11434"

// IsPhase2Complete checks if Phase 2 is complete
func IsPhase2Complete() bool {
	// Check if Ollama service is running
	if !IsOllamaRunning() {
		return false
	}

	// Check if the default model is installed
	if !IsModelInstalled(DefaultModel) {
		return false
	}

	return true
}

// IsOllamaRunning checks if the Ollama service is running
func IsOllamaRunning() bool {
	// Try to ping the Ollama API
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get(OllamaAPIURL)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// IsModelInstalled checks if a specific model is installed
func IsModelInstalled(modelName string) bool {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Parse the output to check if the model is listed
	return strings.Contains(string(output), modelName)
}

// GetInstalledModels returns a list of installed Ollama models
func GetInstalledModels() ([]string, error) {
	cmd := exec.Command("ollama", "list")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list models: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	var models []string

	// Skip the header line and parse model names
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}

		// Model name is the first field
		fields := strings.Fields(line)
		if len(fields) > 0 {
			models = append(models, fields[0])
		}
	}

	return models, nil
}

// StartOllamaService starts the Ollama service using launchd
func StartOllamaService() error {
	// First check if Ollama is already running
	if IsOllamaRunning() {
		fmt.Println("✓ Ollama service is already running")
		return nil
	}

	// Create the launch agent plist
	if err := CreateOllamaLaunchAgent(); err != nil {
		return fmt.Errorf("failed to create launch agent: %w", err)
	}

	// Load the launch agent
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	plistPath := filepath.Join(homeDir, "Library", "LaunchAgents", "com.ollama.service.plist")

	// Unload first in case it's already loaded but not running
	exec.Command("launchctl", "unload", plistPath).Run()

	// Load the launch agent
	cmd := exec.Command("launchctl", "load", plistPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to load launch agent: %w\nOutput: %s", err, string(output))
	}

	// Wait for the service to start
	fmt.Print("Starting Ollama service")
	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		fmt.Print(".")
		if IsOllamaRunning() {
			fmt.Println()
			fmt.Println("✓ Ollama service started successfully")
			return nil
		}
	}

	fmt.Println()
	return fmt.Errorf("ollama service failed to start within 30 seconds")
}

// CreateOllamaLaunchAgent creates a launchd plist for Ollama
func CreateOllamaLaunchAgent() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	// Ensure LaunchAgents directory exists
	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return fmt.Errorf("failed to create LaunchAgents directory: %w", err)
	}

	plistPath := filepath.Join(launchAgentsDir, "com.ollama.service.plist")

	// Check if plist already exists
	if _, err := os.Stat(plistPath); err == nil {
		fmt.Printf("✓ Launch agent already exists at %s\n", plistPath)
		return nil
	}

	// Find the ollama binary path
	ollamaPath, err := exec.LookPath("ollama")
	if err != nil {
		return fmt.Errorf("ollama binary not found in PATH: %w", err)
	}

	// Create the plist content
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.ollama.service</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
		<string>serve</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>StandardOutPath</key>
	<string>%s/Library/Logs/ollama.log</string>
	<key>StandardErrorPath</key>
	<string>%s/Library/Logs/ollama-error.log</string>
	<key>EnvironmentVariables</key>
	<dict>
		<key>PATH</key>
		<string>/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin:/run/current-system/sw/bin</string>
	</dict>
</dict>
</plist>
`, ollamaPath, homeDir, homeDir)

	// Write the plist file
	if err := os.WriteFile(plistPath, []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	fmt.Printf("✓ Created launch agent at %s\n", plistPath)
	return nil
}

// DownloadModel downloads and installs an Ollama model with progress indication
func DownloadModel(modelName string) error {
	fmt.Printf("Downloading model '%s'...\n", modelName)
	fmt.Println("Note: This may take several minutes depending on your internet connection")
	fmt.Println()

	// Start the pull command
	cmd := exec.Command("ollama", "pull", modelName)

	// Create pipes for stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start model download: %w", err)
	}

	// Read and display output in real-time using io.Copy
	go io.Copy(os.Stdout, stdout)

	// Read and display errors in real-time using io.Copy
	go io.Copy(os.Stderr, stderr)

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("model download failed: %w", err)
	}

	fmt.Println()
	fmt.Printf("✓ Model '%s' downloaded successfully\n", modelName)

	// Verify the model was installed
	if !IsModelInstalled(modelName) {
		return fmt.Errorf("model '%s' was not found after download", modelName)
	}

	return nil
}

// VerifyModelIntegrity verifies that a model can be used for inference
func VerifyModelIntegrity(modelName string) error {
	fmt.Printf("Verifying model '%s'...\n", modelName)

	// Test the model with a simple query
	requestBody := map[string]interface{}{
		"model":  modelName,
		"prompt": "Say 'test successful' and nothing else.",
		"stream": false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	client := &http.Client{
		Timeout: 60 * time.Second,
	}

	resp, err := client.Post(
		fmt.Sprintf("%s/api/generate", OllamaAPIURL),
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to test model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("model test failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	fmt.Printf("✓ Model '%s' verified successfully\n", modelName)
	return nil
}

// TestOllamaAPI tests the Ollama API connectivity
func TestOllamaAPI() error {
	fmt.Println("Testing Ollama API connectivity...")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Test the API root endpoint
	resp, err := client.Get(OllamaAPIURL)
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama API returned status %d", resp.StatusCode)
	}

	// Test the tags endpoint to list models
	resp, err = client.Get(fmt.Sprintf("%s/api/tags", OllamaAPIURL))
	if err != nil {
		return fmt.Errorf("failed to list models via API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to list models, status %d", resp.StatusCode)
	}

	fmt.Println("✓ Ollama API is accessible and responding")
	return nil
}

// RunPhase2 executes the complete Phase 2 bootstrap process
func RunPhase2() error {
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("  Phase 2: Ollama Setup")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Check if Phase 2 is already complete
	if IsPhase2Complete() {
		fmt.Println("✓ Phase 2 is already complete!")
		fmt.Println()
		return nil
	}

	// Step 1: Start Ollama service
	fmt.Println("Step 1: Starting Ollama service...")
	if err := StartOllamaService(); err != nil {
		return fmt.Errorf("failed to start Ollama service: %w", err)
	}
	fmt.Println()

	// Step 2: Test API connectivity
	fmt.Println("Step 2: Testing Ollama API...")
	if err := TestOllamaAPI(); err != nil {
		return fmt.Errorf("failed to test Ollama API: %w", err)
	}
	fmt.Println()

	// Step 3: Download default model
	fmt.Println("Step 3: Downloading default model...")
	if !IsModelInstalled(DefaultModel) {
		if err := DownloadModel(DefaultModel); err != nil {
			return fmt.Errorf("failed to download model: %w", err)
		}
	} else {
		fmt.Printf("✓ Model '%s' is already installed\n", DefaultModel)
	}
	fmt.Println()

	// Step 4: Verify model integrity
	fmt.Println("Step 4: Verifying model...")
	if err := VerifyModelIntegrity(DefaultModel); err != nil {
		return fmt.Errorf("failed to verify model: %w", err)
	}
	fmt.Println()

	fmt.Println("═══════════════════════════════════════")
	fmt.Println("✓ Phase 2 Bootstrap Complete!")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()
	fmt.Printf("Ollama is running at %s\n", OllamaAPIURL)
	fmt.Printf("Default model: %s\n", DefaultModel)
	fmt.Println()

	return nil
}
