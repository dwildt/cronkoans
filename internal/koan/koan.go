package koan

import (
	"fmt"
)

// Koan represents a single learning exercise
type Koan struct {
	ID          string   `yaml:"id"`
	Description string   `yaml:"description"`
	Question    string   `yaml:"question"`
	Incomplete  string   `yaml:"incomplete"`
	Answer      string   `yaml:"answer"`
	Hints       []string `yaml:"hints"`
	Explanation string   `yaml:"explanation"`
}

// Lesson represents a collection of related koans
type Lesson struct {
	Title       string  `yaml:"title"`
	Description string  `yaml:"description"`
	Koans       []Koan  `yaml:"koans"`
	Filename    string  `yaml:"-"` // Not from YAML, set programmatically
}

// CompleteCronExpression returns the complete cron expression with the answer filled in
func (k *Koan) CompleteCronExpression() string {
	return replaceBlank(k.Incomplete, k.Answer)
}

// CheckAnswer validates if the user's answer is correct
func (k *Koan) CheckAnswer(userAnswer string) bool {
	// Trim spaces and compare
	return normalizeAnswer(userAnswer) == normalizeAnswer(k.Answer)
}

// GetHint returns the hint at the specified level (0-indexed)
// Returns empty string if level is out of bounds
func (k *Koan) GetHint(level int) string {
	if level < 0 || level >= len(k.Hints) {
		return ""
	}
	return k.Hints[level]
}

// HasMoreHints checks if there are more hints available
func (k *Koan) HasMoreHints(currentLevel int) bool {
	return currentLevel < len(k.Hints)-1
}

// KoanResult represents the result of attempting a koan
type KoanResult struct {
	Koan        *Koan
	UserAnswer  string
	IsCorrect   bool
	HintsUsed   int
	Attempts    int
}

// String provides a formatted representation of the koan
func (k *Koan) String() string {
	return fmt.Sprintf("[%s] %s\nQuestion: %s\nExpression: %s",
		k.ID, k.Description, k.Question, k.Incomplete)
}
