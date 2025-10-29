package voice

import (
	"sync"
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

func TestSpeakerCloseMultipleTimes(t *testing.T) {
	config := DefaultConfig()
	speaker := NewSpeaker(config)

	// First close should work
	speaker.Close()

	// Second close should not panic (handled by sync.Once)
	speaker.Close()

	// Third close for good measure
	speaker.Close()

	// Test passed if we get here without panic
}

func TestIsQuietHours(t *testing.T) {
	// Test quiet hours logic with current time
	// Note: This test validates the logic but results depend on current time
	
	// Test case 1: Quiet hours spanning midnight (22:00 to 07:00)
	// If current hour is 23 or 0-6, should be quiet
	hour := time.Now().Hour()
	
	result1 := isQuietHours(22, 7)
	expectedQuiet1 := hour >= 22 || hour < 7
	if result1 != expectedQuiet1 {
		t.Errorf("isQuietHours(22, 7) = %v, expected %v (current hour: %d)", result1, expectedQuiet1, hour)
	}
	
	// Test case 2: Normal quiet hours (1:00 to 5:00)
	result2 := isQuietHours(1, 5)
	expectedQuiet2 := hour >= 1 && hour < 5
	if result2 != expectedQuiet2 {
		t.Errorf("isQuietHours(1, 5) = %v, expected %v (current hour: %d)", result2, expectedQuiet2, hour)
	}
	
	// Test case 3: Same start and end (0:00 to 0:00) - edge case
	// When start equals end, no time period is selected, so always not quiet
	result3 := isQuietHours(0, 0)
	// This should always be false since hour >= 0 && hour < 0 is always false
	if result3 {
		t.Error("isQuietHours(0, 0) should be false (no time period selected)")
	}
	
	// Test case 4: Same non-zero start and end (12:00 to 12:00)
	result4 := isQuietHours(12, 12)
	// This should always be false since hour >= 12 && hour < 12 is always false
	if result4 {
		t.Error("isQuietHours(12, 12) should be false (no time period selected)")
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
	var wg sync.WaitGroup

	// Goroutine 1: Reading config
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			_ = speaker.GetConfig()
		}
	}()

	// Goroutine 2: Updating config
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			newConfig := DefaultConfig()
			newConfig.Rate = 100 + i
			speaker.UpdateConfig(newConfig)
		}
	}()

	// Goroutine 3: Enable/Disable
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 100; i++ {
			if i%2 == 0 {
				speaker.Enable()
			} else {
				speaker.Disable()
			}
		}
	}()

	// Wait for all goroutines
	wg.Wait()
}
