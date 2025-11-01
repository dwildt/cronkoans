# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Cron Koans is an interactive command-line educational tool that teaches cron expression syntax through practice. Users learn by filling in blanks in incomplete cron expressions across 8 progressive lessons (37 total koans). The application is written in Go and uses YAML files for lesson content.

## Development Commands

### Building and Running
```bash
# Build the application
go build -o cronkoans

# Run directly (development)
go run main.go [command]

# Install globally
go install
```

### Testing and Validation
```bash
# Run all tests
go test ./...

# Validate all lesson files (checks YAML syntax, cron expression validity, unique IDs)
go run main.go validate

# Run the interactive learning mode for manual testing
go run main.go start
```

### Application Commands
- `cronkoans` or `cronkoans start` - Interactive learning mode (starts from last progress)
- `cronkoans list` - List all lessons and progress
- `cronkoans status` - Show progress statistics
- `cronkoans validate` - Validate all lesson YAML files
- `cronkoans reset` - Reset user progress
- `cronkoans --version` - Show version

### Adding New Lessons
Use the custom slash command to create new lessons:
```bash
/add-lesson
```

This command will guide you through creating a properly structured lesson file following the template in `lessons/template.yaml`.

## Architecture

### Core Data Flow
1. **Startup** (`main.go`) → Parses CLI flags and routes to appropriate command
2. **Runner** (`cmd/runner/runner.go`) → Orchestrates the learning session
3. **Lesson Loading** (`internal/koan/parser.go`) → Loads and validates all YAML lessons from `lessons/` directory
4. **Interactive Loop** (`cmd/runner/runner.go:RunInteractive`) → Presents koans, validates answers, manages hints
5. **Progress Tracking** (`internal/progress/tracker.go`) → Persists completion state to `~/.cronkoans_progress.json`

### Package Structure

**`internal/koan/`** - Core koan logic
- `koan.go` - Data structures (`Koan`, `Lesson`)
- `parser.go` - YAML loading and validation (`LoadLesson`, `LoadAllLessons`)
- `validator.go` - Cron expression validation
- `utils.go` - Helper functions (answer normalization, blank replacement)

**`internal/progress/`** - Progress persistence
- `tracker.go` - JSON-based progress tracking (`~/.cronkoans_progress.json`)
- Tracks: completion status, attempts, hints used, timestamps per koan

**`internal/ui/`** - Terminal display
- `display.go` - All user-facing output (koans, hints, validation results, progress stats)

**`cmd/runner/`** - Application orchestration
- `runner.go` - Main runner logic that coordinates lesson loading, progress tracking, and interactive mode
- `GetLessonsDir()` - Finds lesson files (checks executable directory, then current directory)

**`lessons/`** - YAML lesson files
- Named `NN_topic.yaml` (e.g., `01_basics.yaml`)
- `template.yaml` excluded from loading (filtered in `parser.go:LoadAllLessons`)

### Key Design Patterns

**Lesson Loading:** Lessons are loaded alphabetically from `lessons/*.yaml` (excluding `template.yaml`), ensuring proper 01→02→03 ordering. Each lesson file is validated on load to catch structural issues early.

**Interactive State:** The runner maintains no global state - each session starts fresh by loading lessons and progress. The `currentIndex` tracks position in the flattened list of all koans across all lessons.

**Progress Persistence:** Progress is saved immediately after each action (attempt, hint, completion) to `~/.cronkoans_progress.json` in the user's home directory. The tracker uses `LastKoanID` to resume from the last worked koan.

**Answer Validation:** User answers are normalized (trimmed, lowercased) and compared against the expected answer. The complete cron expression is also validated to ensure it's syntactically correct.

**Blank Placeholder System:** Lesson authors use `__` (double underscore) in the `incomplete` field to mark where students fill in answers. The system replaces `__` with the answer to form the complete cron expression.

## Lesson File Format

Each lesson YAML file must include:
```yaml
title: "Topic - Description"
description: "Longer explanation"
koans:
  - id: "unique_id"           # Must be unique across ALL lessons
    description: "Brief"
    question: "Real-world scenario"
    incomplete: "__ * * * *"  # Exactly one __ placeholder
    answer: "*/5"             # Creates valid cron when substituted
    hints:
      - "General hint"
      - "More specific"
      - "Almost the answer"
    explanation: "Educational explanation"
```

**Validation Requirements:**
- Koan IDs must be globally unique (checked across all lessons)
- `incomplete` must contain `__` placeholder
- Answer combined with incomplete must create a valid cron expression
- Each koan must have exactly 3 hints
- No template.yaml in lessons directory is loaded

## Important Development Notes

### No Test Files
This codebase currently has no `*_test.go` files. When adding new functionality, tests would go in the same package directory (e.g., `internal/koan/koan_test.go`).

### Dependencies
- `gopkg.in/yaml.v3` - YAML parsing for lesson files
- Go 1.24.1+ required (specified in `go.mod`)

### Progress File Location
User progress is always stored at `~/.cronkoans_progress.json`, not in the project directory. This allows the binary to be run from anywhere while maintaining progress.

### Lesson Numbering
Lessons use two-digit prefixes (`01_`, `02_`, etc.) to ensure alphabetical sorting matches logical progression. When adding lesson 9+, use `09_` not `9_`.
