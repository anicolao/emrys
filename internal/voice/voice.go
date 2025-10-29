package voice

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"
)

// Config holds voice output configuration
type Config struct {
	Enabled    bool    // Whether voice output is enabled
	Voice      string  // Voice name (e.g., "Jamie")
	Rate       int     // Speech rate in words per minute (default: 200)
	Volume     float64 // Volume from 0.0 to 1.0 (default: 0.7)
	QuietHours bool    // Whether quiet hours are enabled
	QuietStart int     // Quiet hours start (hour in 24h format)
	QuietEnd   int     // Quiet hours end (hour in 24h format)
}

// DefaultConfig returns the default voice configuration
func DefaultConfig() Config {
	return Config{
		Enabled:    true,
		Voice:      "Jamie",
		Rate:       200,
		Volume:     0.7,
		QuietHours: false,
		QuietStart: 22, // 10 PM
		QuietEnd:   7,  // 7 AM
	}
}

// Speaker manages voice output with message queuing
type Speaker struct {
	config     Config
	queue      chan string
	wg         sync.WaitGroup
	mu         sync.RWMutex
	stop       chan struct{}
	closeOnce  sync.Once
	closeMutex sync.Mutex
}

// NewSpeaker creates a new Speaker with the given configuration
func NewSpeaker(config Config) *Speaker {
	s := &Speaker{
		config: config,
		queue:  make(chan string, 100), // Buffer up to 100 messages
		stop:   make(chan struct{}),
	}

	// Start the message processing goroutine
	s.wg.Add(1)
	go s.processQueue()

	return s
}

// processQueue processes queued messages one at a time
func (s *Speaker) processQueue() {
	defer s.wg.Done()

	for {
		select {
		case msg := <-s.queue:
			if err := s.speakNow(msg); err != nil {
				// Log error but continue processing
				fmt.Printf("Voice output error: %v\n", err)
			}
		case <-s.stop:
			return
		}
	}
}

// Speak queues a message for voice output
// Returns immediately, message will be spoken asynchronously
func (s *Speaker) Speak(message string) {
	s.mu.RLock()
	enabled := s.config.Enabled
	s.mu.RUnlock()

	if !enabled {
		return
	}

	// Non-blocking send to queue
	select {
	case s.queue <- message:
		// Message queued successfully
	default:
		// Queue is full, drop the message
		fmt.Println("Voice queue full, message dropped")
	}
}

// SpeakSync speaks a message synchronously (waits for completion)
func (s *Speaker) SpeakSync(message string) error {
	s.mu.RLock()
	enabled := s.config.Enabled
	s.mu.RUnlock()

	if !enabled {
		return nil
	}

	return s.speakNow(message)
}

// speakNow executes the voice output using macOS 'say' command
func (s *Speaker) speakNow(message string) error {
	s.mu.RLock()
	config := s.config
	s.mu.RUnlock()

	// Check if we're in quiet hours
	if config.QuietHours && isQuietHours(config.QuietStart, config.QuietEnd) {
		return nil // Silently skip during quiet hours
	}

	// Build the say command with options
	args := []string{}

	// Add voice if specified
	if config.Voice != "" {
		args = append(args, "-v", config.Voice)
	}

	// Add rate if not default
	if config.Rate != 0 && config.Rate != 200 {
		args = append(args, "-r", fmt.Sprintf("%d", config.Rate))
	}

	// Add volume (say doesn't support volume directly, we use audio output)
	// Note: macOS 'say' doesn't have a volume flag, but we can control it via system volume
	// For now, we'll just document this limitation

	// Add the message
	args = append(args, message)

	// Execute the say command
	cmd := exec.Command("say", args...)
	return cmd.Run()
}

// UpdateConfig updates the speaker configuration
func (s *Speaker) UpdateConfig(config Config) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = config
}

// GetConfig returns a copy of the current configuration
func (s *Speaker) GetConfig() Config {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// Enable enables voice output
func (s *Speaker) Enable() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.Enabled = true
}

// Disable disables voice output
func (s *Speaker) Disable() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config.Enabled = false
}

// IsEnabled returns whether voice output is currently enabled
func (s *Speaker) IsEnabled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config.Enabled
}

// Close stops the speaker and waits for queued messages to complete
// It is safe to call Close multiple times
func (s *Speaker) Close() {
	s.closeOnce.Do(func() {
		close(s.stop)
		s.wg.Wait()
	})
}

// isQuietHours checks if the current time is within quiet hours
func isQuietHours(start, end int) bool {
	// Get current hour (0-23)
	hour := time.Now().Hour()

	// Handle case where quiet hours span midnight (e.g., 22:00 to 07:00)
	if start > end {
		// Quiet hours span midnight
		return hour >= start || hour < end
	}

	// Normal case: quiet hours don't span midnight
	return hour >= start && hour < end
}

// IsVoiceAvailable checks if a specific voice is available on the system
func IsVoiceAvailable(voiceName string) bool {
	// Run 'say -v ?' to list available voices
	cmd := exec.Command("say", "-v", "?")
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// Check if the voice name appears at the start of a line
	// Voice listing format: "VoiceName    language    # comment"
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Get the first field (voice name)
		fields := strings.Fields(line)
		if len(fields) > 0 && fields[0] == voiceName {
			return true
		}
	}

	return false
}

// ListAvailableVoices returns a list of available voices on the system
func ListAvailableVoices() ([]string, error) {
	cmd := exec.Command("say", "-v", "?")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list voices: %w", err)
	}

	var voices []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Extract voice name (first field before whitespace)
		fields := strings.Fields(line)
		if len(fields) > 0 {
			voices = append(voices, fields[0])
		}
	}

	return voices, nil
}

// Test speaks a test message to verify voice output is working
func Test(voiceName string) error {
	testMessage := "Emrys voice output is working correctly."

	args := []string{}
	if voiceName != "" {
		args = append(args, "-v", voiceName)
	}
	args = append(args, testMessage)

	cmd := exec.Command("say", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("voice test failed: %w", err)
	}

	return nil
}
