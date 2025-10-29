package runner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dwildt/cronkoans/internal/koan"
	"github.com/dwildt/cronkoans/internal/progress"
	"github.com/dwildt/cronkoans/internal/ui"
)

// Runner manages the koan learning session
type Runner struct {
	lessons      []*koan.Lesson
	tracker      *progress.Tracker
	lessonsDir   string
	currentIndex int
}

// NewRunner creates a new runner
func NewRunner(lessonsDir string) (*Runner, error) {
	// Load all lessons
	lessons, err := koan.LoadAllLessons(lessonsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load lessons: %w", err)
	}

	// Create progress tracker
	tracker, err := progress.NewTracker()
	if err != nil {
		return nil, fmt.Errorf("failed to create progress tracker: %w", err)
	}

	return &Runner{
		lessons:      lessons,
		tracker:      tracker,
		lessonsDir:   lessonsDir,
		currentIndex: 0,
	}, nil
}

// RunInteractive starts the interactive learning mode
func (r *Runner) RunInteractive() error {
	ui.DisplayWelcome()

	allKoans := koan.GetAllKoans(r.lessons)
	if len(allKoans) == 0 {
		return fmt.Errorf("no koans found")
	}

	// Find the first incomplete koan or start from beginning
	startIndex := 0
	lastKoanID := r.tracker.GetLastKoanID()
	if lastKoanID != "" {
		// Find the koan after the last completed one
		for i, k := range allKoans {
			if k.ID == lastKoanID {
				if i+1 < len(allKoans) {
					startIndex = i + 1
				} else {
					startIndex = i // Last koan
				}
				break
			}
		}
	}

	// Start from the first incomplete koan
	for i := startIndex; i < len(allKoans); i++ {
		if !r.tracker.IsCompleted(allKoans[i].ID) {
			startIndex = i
			break
		}
	}

	// Check if all koans are completed
	if startIndex < len(allKoans) && r.tracker.IsCompleted(allKoans[startIndex].ID) {
		allCompleted := true
		for _, k := range allKoans {
			if !r.tracker.IsCompleted(k.ID) {
				allCompleted = false
				break
			}
		}

		if allCompleted {
			stats := r.tracker.GetStats(len(allKoans))
			ui.DisplayCompletion(stats)
			return nil
		}
	}

	// Run through koans
	for i := startIndex; i < len(allKoans); i++ {
		k := &allKoans[i]

		// Skip if already completed (in case user is reviewing)
		if r.tracker.IsCompleted(k.ID) {
			continue
		}

		if err := r.runKoan(k, i+1, len(allKoans)); err != nil {
			if err.Error() == "quit" {
				return nil
			}
			return err
		}
	}

	// Show completion message
	stats := r.tracker.GetStats(len(allKoans))
	ui.DisplayCompletion(stats)

	return nil
}

// runKoan runs a single koan
func (r *Runner) runKoan(k *koan.Koan, number, total int) error {
	ui.DisplayKoan(k, number, total)

	attempts := 0
	hintsUsed := 0
	currentHintLevel := 0

	for {
		answer := ui.PromptForAnswer()
		answer = strings.TrimSpace(strings.ToLower(answer))

		// Handle special commands
		switch answer {
		case "quit", "exit":
			return fmt.Errorf("quit")
		case "skip":
			ui.DisplayWarning("Skipping this koan...")
			return nil
		case "hint", "h":
			if k.HasMoreHints(currentHintLevel - 1) {
				hint := k.GetHint(currentHintLevel)
				if hint != "" {
					ui.DisplayHint(hint, currentHintLevel)
					currentHintLevel++
					hintsUsed++
					r.tracker.RecordHint(k.ID)
				}
			} else {
				if currentHintLevel == 0 {
					hint := k.GetHint(0)
					if hint != "" {
						ui.DisplayHint(hint, 0)
						currentHintLevel++
						hintsUsed++
						r.tracker.RecordHint(k.ID)
					}
				} else {
					ui.DisplayInfo("No more hints available for this koan.")
				}
			}
			continue
		}

		attempts++
		r.tracker.RecordAttempt(k.ID)

		// Check the answer
		if k.CheckAnswer(answer) {
			ui.DisplayCorrect(k)

			// Mark as completed
			if err := r.tracker.MarkCompleted(k.ID, attempts, hintsUsed); err != nil {
				return fmt.Errorf("failed to save progress: %w", err)
			}

			// Small pause before continuing
			ui.PressEnterToContinue()
			return nil
		}

		ui.DisplayIncorrect()

		// Offer hint after 2 failed attempts
		if attempts >= 2 && currentHintLevel < len(k.Hints) {
			if ui.PromptYesNo("Would you like a hint?") {
				hint := k.GetHint(currentHintLevel)
				if hint != "" {
					ui.DisplayHint(hint, currentHintLevel)
					currentHintLevel++
					hintsUsed++
					r.tracker.RecordHint(k.ID)
				}
			}
		}
	}
}

// RunValidation validates all koans
func (r *Runner) RunValidation() error {
	allKoans := koan.GetAllKoans(r.lessons)
	var results []ui.ValidationResult

	for _, k := range allKoans {
		result := ui.ValidationResult{
			KoanID: k.ID,
			Passed: true,
		}

		// Validate that the answer creates a valid cron expression
		complete := k.CompleteCronExpression()
		if err := koan.ValidateCronExpression(complete); err != nil {
			result.Passed = false
			result.Error = err.Error()
		}

		results = append(results, result)
	}

	ui.DisplayValidationResults(results, len(allKoans))
	return nil
}

// ShowStatus displays progress statistics
func (r *Runner) ShowStatus() error {
	allKoans := koan.GetAllKoans(r.lessons)
	stats := r.tracker.GetStats(len(allKoans))
	ui.DisplayProgress(stats)

	if r.tracker.Exists() {
		ui.DisplayInfo(fmt.Sprintf("Progress file: %s", r.tracker.GetFilePath()))
	} else {
		ui.DisplayInfo("No progress file yet. Start learning to create one!")
	}

	return nil
}

// ListLessons displays all available lessons
func (r *Runner) ListLessons() error {
	ui.DisplayLessonList(r.lessons, r.tracker)
	return nil
}

// Reset resets all progress
func (r *Runner) Reset() error {
	if !r.tracker.Exists() {
		ui.DisplayInfo("No progress to reset.")
		return nil
	}

	if ui.PromptYesNo("Are you sure you want to reset all progress?") {
		if err := r.tracker.Reset(); err != nil {
			return fmt.Errorf("failed to reset progress: %w", err)
		}
		ui.DisplaySuccess("Progress has been reset.")
	} else {
		ui.DisplayInfo("Reset cancelled.")
	}

	return nil
}

// GetLessonsDir returns a default lessons directory path
func GetLessonsDir() string {
	// Try to find lessons directory relative to executable
	exePath, err := os.Executable()
	if err != nil {
		return "lessons"
	}

	exeDir := filepath.Dir(exePath)

	// Check if lessons directory exists relative to executable
	lessonsPath := filepath.Join(exeDir, "lessons")
	if _, err := os.Stat(lessonsPath); err == nil {
		return lessonsPath
	}

	// Check if we're in development (running with go run)
	// Look for lessons in current directory
	if _, err := os.Stat("lessons"); err == nil {
		return "lessons"
	}

	// Default fallback
	return "lessons"
}
