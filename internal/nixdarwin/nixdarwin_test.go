package nixdarwin

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsNixInstalled(t *testing.T) {
	// This test will pass or fail based on whether nix is installed
	// We just verify the function doesn't crash
	result := IsNixInstalled()
	t.Logf("IsNixInstalled returned: %v", result)
}

func TestIsInstalled(t *testing.T) {
	// This test verifies the function runs without crashing
	result := IsInstalled()
	t.Logf("IsInstalled returned: %v", result)
}

func TestCopyFile(t *testing.T) {
	// Create a temporary source file
	tmpDir := t.TempDir()
	srcFile := filepath.Join(tmpDir, "source.txt")
	dstFile := filepath.Join(tmpDir, "dest.txt")

	content := "test content"
	if err := os.WriteFile(srcFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Test copying
	if err := copyFile(srcFile, dstFile); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify content
	result, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}

	if string(result) != content {
		t.Errorf("Content mismatch: got %q, want %q", string(result), content)
	}
}
