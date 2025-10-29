package koan

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"
)

// LoadLesson loads a single lesson from a YAML file
func LoadLesson(filename string) (*Lesson, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	var lesson Lesson
	if err := yaml.Unmarshal(data, &lesson); err != nil {
		return nil, fmt.Errorf("failed to parse YAML from %s: %w", filename, err)
	}

	lesson.Filename = filename

	// Validate the lesson
	if err := validateLesson(&lesson); err != nil {
		return nil, fmt.Errorf("invalid lesson in %s: %w", filename, err)
	}

	return &lesson, nil
}

// LoadAllLessons loads all lesson files from a directory
func LoadAllLessons(lessonsDir string) ([]*Lesson, error) {
	// Find all YAML files in the lessons directory
	pattern := filepath.Join(lessonsDir, "*.yaml")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find lesson files: %w", err)
	}

	// Filter out template.yaml
	var lessonFiles []string
	for _, match := range matches {
		base := filepath.Base(match)
		if base != "template.yaml" {
			lessonFiles = append(lessonFiles, match)
		}
	}

	if len(lessonFiles) == 0 {
		return nil, fmt.Errorf("no lesson files found in %s", lessonsDir)
	}

	// Sort files alphabetically (this ensures 01, 02, 03... order)
	sort.Strings(lessonFiles)

	// Load each lesson
	var lessons []*Lesson
	for _, file := range lessonFiles {
		lesson, err := LoadLesson(file)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

// validateLesson validates the structure and content of a lesson
func validateLesson(lesson *Lesson) error {
	if lesson.Title == "" {
		return fmt.Errorf("lesson must have a title")
	}

	if len(lesson.Koans) == 0 {
		return fmt.Errorf("lesson must have at least one koan")
	}

	// Validate each koan
	seenIDs := make(map[string]bool)
	for i, koan := range lesson.Koans {
		if err := validateKoan(&koan); err != nil {
			return fmt.Errorf("koan %d: %w", i+1, err)
		}

		// Check for duplicate IDs
		if seenIDs[koan.ID] {
			return fmt.Errorf("duplicate koan ID: %s", koan.ID)
		}
		seenIDs[koan.ID] = true
	}

	return nil
}

// validateKoan validates a single koan
func validateKoan(koan *Koan) error {
	if koan.ID == "" {
		return fmt.Errorf("koan must have an ID")
	}

	if koan.Description == "" {
		return fmt.Errorf("koan must have a description")
	}

	if koan.Question == "" {
		return fmt.Errorf("koan must have a question")
	}

	if koan.Incomplete == "" {
		return fmt.Errorf("koan must have an incomplete expression")
	}

	if koan.Answer == "" {
		return fmt.Errorf("koan must have an answer")
	}

	// Check that incomplete expression contains the blank placeholder
	if !containsBlank(koan.Incomplete) {
		return fmt.Errorf("incomplete expression must contain __ placeholder")
	}

	// Validate that the answer creates a valid cron expression
	if !IsValidCronAnswer(koan.Incomplete, koan.Answer) {
		return fmt.Errorf("answer '%s' does not create a valid cron expression: %s",
			koan.Answer, replaceBlank(koan.Incomplete, koan.Answer))
	}

	return nil
}

// containsBlank checks if a string contains the __ placeholder
func containsBlank(s string) bool {
	return len(s) >= 2 && (s[0:2] == "__" || containsSubstring(s, "__"))
}

// containsSubstring is a helper to check for substring
func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// GetAllKoans returns a flat list of all koans from all lessons
func GetAllKoans(lessons []*Lesson) []Koan {
	var allKoans []Koan
	for _, lesson := range lessons {
		allKoans = append(allKoans, lesson.Koans...)
	}
	return allKoans
}

// FindKoanByID finds a koan by its ID across all lessons
func FindKoanByID(lessons []*Lesson, id string) *Koan {
	for _, lesson := range lessons {
		for i := range lesson.Koans {
			if lesson.Koans[i].ID == id {
				return &lesson.Koans[i]
			}
		}
	}
	return nil
}

// CountTotalKoans returns the total number of koans across all lessons
func CountTotalKoans(lessons []*Lesson) int {
	count := 0
	for _, lesson := range lessons {
		count += len(lesson.Koans)
	}
	return count
}
