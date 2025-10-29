package progress

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	progressFileName = ".cronkoans_progress.json"
)

// KoanProgress represents the progress for a single koan
type KoanProgress struct {
	KoanID      string    `json:"koan_id"`
	Completed   bool      `json:"completed"`
	Attempts    int       `json:"attempts"`
	HintsUsed   int       `json:"hints_used"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Progress represents the overall progress
type Progress struct {
	Koans       map[string]*KoanProgress `json:"koans"`
	LastKoanID  string                   `json:"last_koan_id"`
	StartedAt   time.Time                `json:"started_at"`
	UpdatedAt   time.Time                `json:"updated_at"`
	Version     string                   `json:"version"`
}

// Tracker manages progress persistence
type Tracker struct {
	filePath string
	progress *Progress
}

// NewTracker creates a new progress tracker
func NewTracker() (*Tracker, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	filePath := filepath.Join(homeDir, progressFileName)
	tracker := &Tracker{
		filePath: filePath,
	}

	// Try to load existing progress
	if err := tracker.Load(); err != nil {
		// If file doesn't exist, create new progress
		if os.IsNotExist(err) {
			tracker.progress = newProgress()
		} else {
			return nil, err
		}
	}

	return tracker, nil
}

// newProgress creates a new progress structure
func newProgress() *Progress {
	return &Progress{
		Koans:     make(map[string]*KoanProgress),
		StartedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   "1.0",
	}
}

// Load loads progress from the JSON file
func (t *Tracker) Load() error {
	data, err := os.ReadFile(t.filePath)
	if err != nil {
		return err
	}

	var progress Progress
	if err := json.Unmarshal(data, &progress); err != nil {
		return fmt.Errorf("failed to parse progress file: %w", err)
	}

	// Ensure the map is initialized
	if progress.Koans == nil {
		progress.Koans = make(map[string]*KoanProgress)
	}

	t.progress = &progress
	return nil
}

// Save saves progress to the JSON file
func (t *Tracker) Save() error {
	t.progress.UpdatedAt = time.Now()

	data, err := json.MarshalIndent(t.progress, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal progress: %w", err)
	}

	if err := os.WriteFile(t.filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write progress file: %w", err)
	}

	return nil
}

// IsCompleted checks if a koan is completed
func (t *Tracker) IsCompleted(koanID string) bool {
	if kp, ok := t.progress.Koans[koanID]; ok {
		return kp.Completed
	}
	return false
}

// MarkCompleted marks a koan as completed
func (t *Tracker) MarkCompleted(koanID string, attempts int, hintsUsed int) error {
	now := time.Now()

	kp := t.getOrCreateKoanProgress(koanID)
	kp.Completed = true
	kp.Attempts = attempts
	kp.HintsUsed = hintsUsed
	kp.CompletedAt = &now

	t.progress.LastKoanID = koanID

	return t.Save()
}

// RecordAttempt records an attempt for a koan
func (t *Tracker) RecordAttempt(koanID string) error {
	kp := t.getOrCreateKoanProgress(koanID)
	kp.Attempts++

	return t.Save()
}

// RecordHint records that a hint was used
func (t *Tracker) RecordHint(koanID string) error {
	kp := t.getOrCreateKoanProgress(koanID)
	kp.HintsUsed++

	return t.Save()
}

// getOrCreateKoanProgress gets or creates progress for a koan
func (t *Tracker) getOrCreateKoanProgress(koanID string) *KoanProgress {
	if kp, ok := t.progress.Koans[koanID]; ok {
		return kp
	}

	kp := &KoanProgress{
		KoanID:    koanID,
		Completed: false,
		Attempts:  0,
		HintsUsed: 0,
	}
	t.progress.Koans[koanID] = kp

	return kp
}

// GetProgress gets the progress for a specific koan
func (t *Tracker) GetProgress(koanID string) *KoanProgress {
	if kp, ok := t.progress.Koans[koanID]; ok {
		return kp
	}
	return nil
}

// GetLastKoanID returns the ID of the last koan worked on
func (t *Tracker) GetLastKoanID() string {
	return t.progress.LastKoanID
}

// GetCompletedCount returns the number of completed koans
func (t *Tracker) GetCompletedCount() int {
	count := 0
	for _, kp := range t.progress.Koans {
		if kp.Completed {
			count++
		}
	}
	return count
}

// GetTotalAttempts returns the total number of attempts across all koans
func (t *Tracker) GetTotalAttempts() int {
	total := 0
	for _, kp := range t.progress.Koans {
		total += kp.Attempts
	}
	return total
}

// GetTotalHints returns the total number of hints used across all koans
func (t *Tracker) GetTotalHints() int {
	total := 0
	for _, kp := range t.progress.Koans {
		total += kp.HintsUsed
	}
	return total
}

// Reset clears all progress
func (t *Tracker) Reset() error {
	t.progress = newProgress()
	return t.Save()
}

// GetStats returns statistics about the progress
func (t *Tracker) GetStats(totalKoans int) Stats {
	completed := t.GetCompletedCount()
	percentage := 0.0
	if totalKoans > 0 {
		percentage = float64(completed) / float64(totalKoans) * 100
	}

	return Stats{
		TotalKoans:      totalKoans,
		CompletedKoans:  completed,
		RemainingKoans:  totalKoans - completed,
		PercentComplete: percentage,
		TotalAttempts:   t.GetTotalAttempts(),
		TotalHintsUsed:  t.GetTotalHints(),
		StartedAt:       t.progress.StartedAt,
		UpdatedAt:       t.progress.UpdatedAt,
	}
}

// Stats represents progress statistics
type Stats struct {
	TotalKoans      int
	CompletedKoans  int
	RemainingKoans  int
	PercentComplete float64
	TotalAttempts   int
	TotalHintsUsed  int
	StartedAt       time.Time
	UpdatedAt       time.Time
}

// GetFilePath returns the path to the progress file
func (t *Tracker) GetFilePath() string {
	return t.filePath
}

// Exists checks if the progress file exists
func (t *Tracker) Exists() bool {
	_, err := os.Stat(t.filePath)
	return err == nil
}
