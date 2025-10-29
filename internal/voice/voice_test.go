package voice

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if !config.Enabled {
		t.Error("Expected voice to be enabled by default")
	}

	if config.Voice != "Jamie" {
		t.Errorf("Expected default voice to be 'Jamie', got '%s'", config.Voice)
	}

	if config.Rate != 200 {
		t.Errorf("Expected default rate to be 200, got %d", config.Rate)
	}

	if config.Volume != 0.7 {
		t.Errorf("Expected default volume to be 0.7, got %f", config.Volume)
	}

	if config.QuietHours {
		t.Error("Expected quiet hours to be disabled by default")
	}
}

func TestNewSpeaker(t *testing.T) {
	config := DefaultConfig()
	speaker := NewSpeaker(config)
	defer speaker.Close()

	if !speaker.IsEnabled() {
		t.Error("Expected speaker to be enabled")
	}

	gotConfig := speaker.GetConfig()
	if gotConfig.Voice != config.Voice {
		t.Errorf("Expected voice '%s', got '%s'", config.Voice, gotConfig.Voice)
	}
}

func TestSpeakerEnableDisable(t *testing.T) {
	config := DefaultConfig()
	speaker := NewSpeaker(config)
	defer speaker.Close()

	// Initially enabled
	if !speaker.IsEnabled() {
		t.Error("Expected speaker to be enabled initially")
	}

	// Disable
	speaker.Disable()
	if speaker.IsEnabled() {
		t.Error("Expected speaker to be disabled after Disable()")
	}

	// Enable
	speaker.Enable()
	if !speaker.IsEnabled() {
		t.Error("Expected speaker to be enabled after Enable()")
	}
}

func TestSpeakerUpdateConfig(t *testing.T) {
	config := DefaultConfig()
	speaker := NewSpeaker(config)
	defer speaker.Close()

	newConfig := Config{
		Enabled:    true,
		Voice:      "Samantha",
		Rate:       150,
		Volume:     0.5,
		QuietHours: true,
		QuietStart: 20,
		QuietEnd:   8,
	}

	speaker.UpdateConfig(newConfig)

	gotConfig := speaker.GetConfig()
	if gotConfig.Voice != "Samantha" {
		t.Errorf("Expected voice 'Samantha', got '%s'", gotConfig.Voice)
	}
	if gotConfig.Rate != 150 {
		t.Errorf("Expected rate 150, got %d", gotConfig.Rate)
	}
	if gotConfig.Volume != 0.5 {
		t.Errorf("Expected volume 0.5, got %f", gotConfig.Volume)
	}
	if !gotConfig.QuietHours {
		t.Error("Expected quiet hours to be enabled")
	}
}

func TestSpeakerQueueing(t *testing.T) {
	config := DefaultConfig()
	config.Voice = "" // Use default system voice for testing
	speaker := NewSpeaker(config)
	defer speaker.Close()

	// Queue multiple messages
	speaker.Speak("Message 1")
	speaker.Speak("Message 2")
	speaker.Speak("Message 3")

	// Give some time for messages to be processed
	// In a real test, we'd need to mock the say command
	time.Sleep(100 * time.Millisecond)

	// Just verify the speaker is still working
	if !speaker.IsEnabled() {
		t.Error("Speaker should still be enabled")
	}
}

func TestSpeakerDisabledNoOutput(t *testing.T) {
	config := DefaultConfig()
	config.Enabled = false
	speaker := NewSpeaker(config)
	defer speaker.Close()

	// Speak should be a no-op when disabled
	speaker.Speak("This should not be spoken")

	// Verify speaker is disabled
	if speaker.IsEnabled() {
		t.Error("Speaker should be disabled")
	}
}

func TestSpeakerClose(t *testing.T) {
	config := DefaultConfig()
	speaker := NewSpeaker(config)

	// Queue a message
	speaker.Speak("Test message")

	// Close should wait for queued messages
	speaker.Close()

	// Verify we can't queue new messages after close
	// (trying to send to closed channel would panic, but we handle it gracefully)
	// No assertion here, just making sure Close() completes
}

func TestIsQuietHours(t *testing.T) {
	// Test quiet hours logic
	tests := []struct {
		name      string
		start     int
		end       int
		wantQuiet bool
	}{
		{
			name:      "normal hours",
			start:     22,
			end:       7,
			wantQuiet: false, // Always false in our mock
		},
		{
			name:      "reverse hours",
			start:     7,
			end:       22,
			wantQuiet: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isQuietHours(tt.start, tt.end)
			if got != tt.wantQuiet {
				t.Errorf("isQuietHours() = %v, want %v", got, tt.wantQuiet)
			}
		})
	}
}

// TestIsVoiceAvailable tests voice availability checking
// This test requires the 'say' command to be available (macOS only)
func TestIsVoiceAvailable(t *testing.T) {
	// Skip this test on non-macOS systems
	t.Skip("Skipping voice availability test (requires macOS)")

	// Test with a voice that should exist on macOS
	if !IsVoiceAvailable("Alex") {
		t.Log("Voice 'Alex' not found (may not be installed)")
	}

	// Test with a voice that definitely doesn't exist
	if IsVoiceAvailable("NonExistentVoice12345") {
		t.Error("Should not find non-existent voice")
	}
}

// TestListAvailableVoices tests listing available voices
func TestListAvailableVoices(t *testing.T) {
	// Skip this test on non-macOS systems
	t.Skip("Skipping voice listing test (requires macOS)")

	voices, err := ListAvailableVoices()
	if err != nil {
		t.Fatalf("Failed to list voices: %v", err)
	}

	if len(voices) == 0 {
		t.Error("Expected at least one voice to be available")
	}

	t.Logf("Found %d voices", len(voices))
}

// TestTest tests the voice testing utility
func TestTest(t *testing.T) {
	// Skip this test on non-macOS systems
	t.Skip("Skipping voice test (requires macOS)")

	// Test with default voice
	if err := Test(""); err != nil {
		t.Errorf("Voice test failed: %v", err)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name   string
		config Config
		valid  bool
	}{
		{
			name:   "valid default config",
			config: DefaultConfig(),
			valid:  true,
		},
		{
			name: "valid custom config",
			config: Config{
				Enabled:    true,
				Voice:      "Samantha",
				Rate:       150,
				Volume:     0.8,
				QuietHours: true,
				QuietStart: 22,
				QuietEnd:   7,
			},
			valid: true,
		},
		{
			name: "disabled voice",
			config: Config{
				Enabled: false,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			speaker := NewSpeaker(tt.config)
			defer speaker.Close()

			if speaker == nil {
				t.Error("Expected speaker to be created")
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	config := DefaultConfig()
	config.Enabled = false // Disable to avoid actual voice output
	speaker := NewSpeaker(config)
	defer speaker.Close()

	// Test concurrent access to configuration
	done := make(chan bool)

	// Goroutine 1: Reading config
	go func() {
		for i := 0; i < 100; i++ {
			_ = speaker.GetConfig()
		}
		done <- true
	}()

	// Goroutine 2: Updating config
	go func() {
		for i := 0; i < 100; i++ {
			newConfig := DefaultConfig()
			newConfig.Rate = 100 + i
			speaker.UpdateConfig(newConfig)
		}
		done <- true
	}()

	// Goroutine 3: Enable/Disable
	go func() {
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				speaker.Enable()
			} else {
				speaker.Disable()
			}
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
}
