package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/dwildt/cronkoans/internal/koan"
	"github.com/dwildt/cronkoans/internal/progress"
)

// Color codes for terminal output
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorGray   = "\033[37m"
	ColorBold   = "\033[1m"
)

// DisplayWelcome shows the welcome message
func DisplayWelcome() {
	fmt.Println(ColorBold + ColorCyan + "╔═══════════════════════════════════════════════════════╗" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "║            Welcome to Cron Koans!                    ║" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "║      Learn Crontab through practice and wisdom       ║" + ColorReset)
	fmt.Println(ColorBold + ColorCyan + "╚═══════════════════════════════════════════════════════╝" + ColorReset)
	fmt.Println()
}

// DisplayKoan displays a koan to the user
func DisplayKoan(k *koan.Koan, number, total int) {
	fmt.Println(ColorBold + fmt.Sprintf("\n[Koan %d/%d]", number, total) + ColorReset)
	fmt.Println(ColorBlue + k.Description + ColorReset)
	fmt.Println()
	fmt.Println(ColorGray + "Question: " + ColorReset + k.Question)
	fmt.Println()
	fmt.Println(ColorYellow + "Incomplete expression: " + ColorBold + k.Incomplete + ColorReset)
	fmt.Println()
}

// DisplayHint displays a hint
func DisplayHint(hint string, level int) {
	fmt.Println(ColorYellow + fmt.Sprintf("\n💡 Hint %d: ", level+1) + hint + ColorReset)
	fmt.Println()
}

// DisplayCorrect shows success message
func DisplayCorrect(k *koan.Koan) {
	fmt.Println(ColorGreen + "\n✓ Correct!" + ColorReset)
	fmt.Println(ColorGreen + "Complete expression: " + ColorBold + k.CompleteCronExpression() + ColorReset)
	if k.Explanation != "" {
		fmt.Println()
		fmt.Println(ColorCyan + "📚 " + k.Explanation + ColorReset)
	}
	fmt.Println()
}

// DisplayIncorrect shows incorrect message
func DisplayIncorrect() {
	fmt.Println(ColorRed + "\n✗ Incorrect. Try again!" + ColorReset)
	fmt.Println()
}

// DisplayProgress shows progress statistics
func DisplayProgress(stats progress.Stats) {
	percentage := fmt.Sprintf("%.1f%%", stats.PercentComplete)

	fmt.Println(ColorBold + "\n📊 Your Progress" + ColorReset)
	fmt.Println(strings.Repeat("─", 50))
	fmt.Printf("Completed: %s%d/%d%s koans (%s%s%s)\n",
		ColorGreen, stats.CompletedKoans, stats.TotalKoans, ColorReset,
		ColorBold, percentage, ColorReset)
	fmt.Printf("Remaining: %s%d%s koans\n",
		ColorYellow, stats.RemainingKoans, ColorReset)
	fmt.Printf("Total attempts: %d\n", stats.TotalAttempts)
	fmt.Printf("Hints used: %d\n", stats.TotalHintsUsed)
	fmt.Println(strings.Repeat("─", 50))
	fmt.Println()
}

// DisplayCompletion shows completion message
func DisplayCompletion(stats progress.Stats) {
	fmt.Println(ColorGreen + ColorBold)
	fmt.Println("\n╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║          🎉 Congratulations! 🎉                      ║")
	fmt.Println("║                                                       ║")
	fmt.Println("║      You have completed all Cron Koans!             ║")
	fmt.Println("║                                                       ║")
	fmt.Println("║      You are now a Crontab master!                  ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝")
	fmt.Println(ColorReset)

	fmt.Printf("Total attempts: %d\n", stats.TotalAttempts)
	fmt.Printf("Hints used: %d\n", stats.TotalHintsUsed)
	fmt.Printf("Started: %s\n", stats.StartedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Completed: %s\n", stats.UpdatedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()
}

// DisplayLessonList shows all available lessons
func DisplayLessonList(lessons []*koan.Lesson, tracker *progress.Tracker) {
	fmt.Println(ColorBold + "\n📚 Available Lessons" + ColorReset)
	fmt.Println(strings.Repeat("─", 60))

	for i, lesson := range lessons {
		completed := 0
		for _, k := range lesson.Koans {
			if tracker.IsCompleted(k.ID) {
				completed++
			}
		}

		total := len(lesson.Koans)
		status := fmt.Sprintf("%d/%d", completed, total)

		var statusColor string
		if completed == total {
			statusColor = ColorGreen
		} else if completed > 0 {
			statusColor = ColorYellow
		} else {
			statusColor = ColorGray
		}

		fmt.Printf("%d. %s (%s%s%s)\n", i+1, lesson.Title, statusColor, status, ColorReset)
		fmt.Printf("   %s\n", ColorGray+lesson.Description+ColorReset)
	}

	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
}

// DisplayValidationResults shows the results of validation mode
func DisplayValidationResults(results []ValidationResult, totalKoans int) {
	passed := 0
	for _, r := range results {
		if r.Passed {
			passed++
		}
	}

	fmt.Println(ColorBold + "\n🔍 Validation Results" + ColorReset)
	fmt.Println(strings.Repeat("─", 60))

	for _, r := range results {
		if r.Passed {
			fmt.Printf("%s✓%s %s\n", ColorGreen, ColorReset, r.KoanID)
		} else {
			fmt.Printf("%s✗%s %s - %s\n", ColorRed, ColorReset, r.KoanID, r.Error)
		}
	}

	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("\nPassed: %s%d/%d%s\n", ColorGreen, passed, totalKoans, ColorReset)
	fmt.Printf("Failed: %s%d%s\n", ColorRed, totalKoans-passed, ColorReset)
	fmt.Println()
}

// ValidationResult represents the result of validating a koan
type ValidationResult struct {
	KoanID string
	Passed bool
	Error  string
}

// PromptForAnswer prompts the user for their answer
func PromptForAnswer() string {
	fmt.Print(ColorBold + "Your answer: " + ColorReset)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	return strings.TrimSpace(answer)
}

// PromptYesNo prompts for a yes/no answer
func PromptYesNo(question string) bool {
	fmt.Print(ColorBold + question + " (y/n): " + ColorReset)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.ToLower(strings.TrimSpace(answer))
	return answer == "y" || answer == "yes"
}

// DisplayError shows an error message
func DisplayError(err error) {
	fmt.Fprintf(os.Stderr, ColorRed+"Error: %v\n"+ColorReset, err)
}

// DisplayInfo shows an informational message
func DisplayInfo(message string) {
	fmt.Println(ColorBlue + "ℹ " + message + ColorReset)
}

// DisplaySuccess shows a success message
func DisplaySuccess(message string) {
	fmt.Println(ColorGreen + "✓ " + message + ColorReset)
}

// DisplayWarning shows a warning message
func DisplayWarning(message string) {
	fmt.Println(ColorYellow + "⚠ " + message + ColorReset)
}

// ClearScreen clears the terminal screen
func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

// PressEnterToContinue waits for the user to press enter
func PressEnterToContinue() {
	fmt.Print(ColorGray + "\nPress Enter to continue..." + ColorReset)
	bufio.NewReader(os.Stdin).ReadString('\n')
}

// DisplayHelp shows help information
func DisplayHelp() {
	fmt.Println(ColorBold + "Cron Koans - Help" + ColorReset)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  cronkoans              Start interactive mode")
	fmt.Println("  cronkoans start        Start from beginning")
	fmt.Println("  cronkoans reset        Reset all progress")
	fmt.Println("  cronkoans list         List all lessons")
	fmt.Println("  cronkoans status       Show progress statistics")
	fmt.Println("  cronkoans validate     Validate all koans")
	fmt.Println("  cronkoans help         Show this help message")
	fmt.Println()
	fmt.Println("During interactive mode:")
	fmt.Println("  Type your answer and press Enter")
	fmt.Println("  Type 'hint' to get a hint")
	fmt.Println("  Type 'skip' to skip the current koan")
	fmt.Println("  Type 'quit' or 'exit' to quit")
	fmt.Println()
	fmt.Println("Learn more about cron:")
	fmt.Println("  https://crontab.guru")
	fmt.Println()
}
