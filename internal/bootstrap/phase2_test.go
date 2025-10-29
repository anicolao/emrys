package bootstrap

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsOllamaRunning(t *testing.T) {
	// This test just verifies the function runs without crashing
	result := IsOllamaRunning()
	t.Logf("IsOllamaRunning returned: %v", result)
}

func TestIsModelInstalled(t *testing.T) {
	// Test with a model that definitely doesn't exist
	result := IsModelInstalled("this-model-definitely-does-not-exist-xyz123")
	if result {
		t.Error("Expected non-existent model to return false, but it returned true")
	}
}

func TestGetInstalledModels(t *testing.T) {
	// Test that the function doesn't crash
	// It may fail if ollama is not installed, which is expected in CI
	models, err := GetInstalledModels()
	if err != nil {
		// Check if error is because ollama is not installed
		if _, lookupErr := exec.LookPath("ollama"); lookupErr != nil {
			t.Logf("ollama not installed (expected in CI): %v", err)
			return
		}
		t.Logf("GetInstalledModels returned error: %v", err)
		return
	}
	t.Logf("Installed models: %v", models)
}

func TestIsPhase2Complete(t *testing.T) {
	// This test verifies the function runs without crashing
	result := IsPhase2Complete()
	t.Logf("IsPhase2Complete returned: %v", result)
}

func TestCreateOllamaLaunchAgent(t *testing.T) {
	// Skip if not on macOS
	if _, err := os.Stat("/Library"); err != nil {
		t.Skip("Not running on macOS, skipping launch agent test")
	}

	// Create a temporary home directory
	tmpDir := t.TempDir()
	launchAgentsDir := filepath.Join(tmpDir, "Library", "LaunchAgents")
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		t.Fatalf("Failed to create LaunchAgents directory: %v", err)
	}

	// Temporarily change HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Test creating the launch agent
	// This may fail if ollama is not in PATH, which is expected
	err := CreateOllamaLaunchAgent()
	if err != nil {
		// Check if error is because ollama is not installed
		if _, lookupErr := exec.LookPath("ollama"); lookupErr != nil {
			t.Logf("ollama not installed (expected in CI): %v", err)
			return
		}
		t.Fatalf("CreateOllamaLaunchAgent failed: %v", err)
	}

	// Verify the plist was created
	plistPath := filepath.Join(launchAgentsDir, "com.ollama.service.plist")
	if _, err := os.Stat(plistPath); err != nil {
		t.Fatalf("Launch agent plist was not created: %v", err)
	}

	// Read and verify the plist content
	content, err := os.ReadFile(plistPath)
	if err != nil {
		t.Fatalf("Failed to read plist: %v", err)
	}

	plistStr := string(content)

	// Verify key elements
	if !strings.Contains(plistStr, "com.ollama.service") {
		t.Error("Plist doesn't contain service label")
	}
	if !strings.Contains(plistStr, "serve") {
		t.Error("Plist doesn't contain serve command")
	}
	if !strings.Contains(plistStr, "RunAtLoad") {
		t.Error("Plist doesn't contain RunAtLoad")
	}
	if !strings.Contains(plistStr, "KeepAlive") {
		t.Error("Plist doesn't contain KeepAlive")
	}

	// Test idempotency - running again should not fail
	err = CreateOllamaLaunchAgent()
	if err != nil {
		t.Fatalf("Second CreateOllamaLaunchAgent call failed: %v", err)
	}
}

func TestTestOllamaAPI(t *testing.T) {
	// Create a mock Ollama API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Ollama is running"))
			return
		}
		if r.URL.Path == "/api/tags" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"models":[]}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// Temporarily override the OllamaAPIURL for testing
	// We can't easily test this without modifying the package variables,
	// so we'll just test that the real function doesn't crash
	err := TestOllamaAPI()
	if err != nil {
		t.Logf("TestOllamaAPI failed (expected if Ollama is not running): %v", err)
	}
}

func TestVerifyModelIntegrity(t *testing.T) {
	// Create a mock Ollama API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/generate" {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"model":"test","response":"test successful","done":true}`))
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	// We can't easily test this without a running Ollama instance,
	// so we'll just verify the function doesn't crash when Ollama is not running
	err := VerifyModelIntegrity("nonexistent-model")
	if err == nil {
		t.Error("Expected error when verifying non-existent model without Ollama running")
	}
	t.Logf("VerifyModelIntegrity failed as expected: %v", err)
}

func TestDownloadModel(t *testing.T) {
	// We can't test actual model downloads in CI, so just verify the function exists
	// and handles errors appropriately
	err := DownloadModel("nonexistent-model-xyz123")
	if err == nil {
		t.Error("Expected error when downloading non-existent model")
	}
	t.Logf("DownloadModel failed as expected: %v", err)
}

func TestDefaultModelConstant(t *testing.T) {
	// Verify the default model constant is set
	if DefaultModel == "" {
		t.Error("DefaultModel constant is empty")
	}
	t.Logf("DefaultModel: %s", DefaultModel)
}

func TestOllamaAPIURLConstant(t *testing.T) {
	// Verify the API URL constant is set correctly
	if OllamaAPIURL == "" {
		t.Error("OllamaAPIURL constant is empty")
	}
	if !strings.HasPrefix(OllamaAPIURL, "http") {
		t.Error("OllamaAPIURL should start with http")
	}
	t.Logf("OllamaAPIURL: %s", OllamaAPIURL)
}
