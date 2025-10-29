package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dwildt/cronkoans/cmd/runner"
	"github.com/dwildt/cronkoans/internal/ui"
)

const version = "1.0.0"

func main() {
	if err := run(); err != nil {
		ui.DisplayError(err)
		os.Exit(1)
	}
}

func run() error {
	// Define flags
	helpFlag := flag.Bool("help", false, "Show help message")
	versionFlag := flag.Bool("version", false, "Show version")
	lessonsDir := flag.String("lessons", runner.GetLessonsDir(), "Path to lessons directory")

	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("Cron Koans v%s\n", version)
		return nil
	}

	// Handle help flag
	if *helpFlag {
		ui.DisplayHelp()
		return nil
	}

	// Get the command (if any)
	args := flag.Args()
	command := "interactive"
	if len(args) > 0 {
		command = args[0]
	}

	// Create runner
	r, err := runner.NewRunner(*lessonsDir)
	if err != nil {
		return err
	}

	// Execute command
	switch command {
	case "interactive", "start", "":
		return r.RunInteractive()

	case "validate":
		return r.RunValidation()

	case "status":
		return r.ShowStatus()

	case "list":
		return r.ListLessons()

	case "reset":
		return r.Reset()

	case "help":
		ui.DisplayHelp()
		return nil

	default:
		return fmt.Errorf("unknown command: %s\nRun 'cronkoans help' for usage information", command)
	}
}
