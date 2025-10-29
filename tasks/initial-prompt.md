# Cron Koans - Initial Project Plan

## Overview
Create a Koans application to learn about Crontab configuration using Golang as the runner language.

## User Requirements
- **Language**: Golang
- **Koan Style**: Fill-in-the-blank
- **Features**:
  - Progress tracking
  - Hints system
  - Interactive mode
  - Validation mode
- **Must be easy to add new challenges/lessons**

## Project Structure

```
cronkoans/
├── main.go                    # CLI entry point
├── go.mod                     # Go module definition
├── README.md                  # Main documentation
├── CONTRIBUTING.md            # Guide for adding new lessons
├── .claude/
│   └── commands/
│       └── add-lesson.md      # Claude Code command for adding lessons
├── cmd/
│   └── runner/
│       └── runner.go          # Main runner logic
├── internal/
│   ├── koan/
│   │   ├── koan.go           # Koan data structures
│   │   ├── validator.go      # Cron expression validator
│   │   └── parser.go         # Parse lesson files
│   ├── progress/
│   │   └── tracker.go        # Progress tracking (JSON file)
│   └── ui/
│       └── display.go        # Terminal UI helpers
├── lessons/
│   ├── 01_basics.yaml
│   ├── 02_wildcards.yaml
│   ├── 03_ranges.yaml
│   ├── 04_steps.yaml
│   ├── 05_lists.yaml
│   ├── 06_special_strings.yaml
│   ├── 07_common_patterns.yaml
│   ├── 08_advanced.yaml
│   └── template.yaml         # Template for new lessons
└── tasks/
    └── initial-prompt.md     # This file
```

## Lesson Topics (25-30 koans total)

### 1. Basics (3-4 koans)
Understanding the 5 fields of cron:
- Minute (0-59)
- Hour (0-23)
- Day of month (1-31)
- Month (1-12)
- Day of week (0-7, where 0 and 7 are Sunday)

### 2. Wildcards (3-4 koans)
Using `*` to mean "every" value:
- `* * * * *` - every minute
- `0 * * * *` - every hour
- `0 0 * * *` - every day at midnight

### 3. Ranges (3-4 koans)
Using dashes to specify ranges:
- `0 9-17 * * *` - every hour from 9 AM to 5 PM
- `0 0 * * 1-5` - midnight on weekdays

### 4. Step Values (4-5 koans)
Using `*/n` for intervals:
- `*/5 * * * *` - every 5 minutes
- `0 */2 * * *` - every 2 hours
- `0 0 */2 * *` - every 2 days

### 5. Lists (3-4 koans)
Using commas to specify multiple values:
- `0 0 * * 0,6` - midnight on weekends
- `0 9,12,18 * * *` - at 9 AM, noon, and 6 PM
- `0,15,30,45 * * * *` - every 15 minutes

### 6. Special Strings (3-4 koans)
Non-standard but widely supported special strings:
- `@yearly` or `@annually` - once a year (0 0 1 1 *)
- `@monthly` - once a month (0 0 1 * *)
- `@weekly` - once a week (0 0 * * 0)
- `@daily` or `@midnight` - once a day (0 0 * * *)
- `@hourly` - once an hour (0 * * * *)
- `@reboot` - at startup

### 7. Common Patterns (4-5 koans)
Real-world examples:
- Database backups
- Log rotation
- Report generation
- Health checks
- System maintenance

### 8. Advanced (3-4 koans)
Combining operators and edge cases:
- Multiple ranges and lists
- Last day of month considerations
- Day of week vs day of month interactions
- Complex business hour patterns

## Lesson File Format (YAML)

Each lesson is a YAML file with this structure:

```yaml
title: "Lesson Title"
description: "Brief description of what this lesson teaches"
koans:
  - id: "lesson_koan_1"
    description: "What this koan teaches"
    question: "Run a task every 5 minutes"
    incomplete: "__ * * * *"
    answer: "*/5"
    hints:
      - "Think about the step value operator"
      - "Use the */n notation for intervals"
      - "The answer is */5 for every 5 minutes"
    explanation: "The */5 in the minute field means 'every 5 minutes'. The step value divides the field range into intervals."
```

This format makes it extremely easy to add new lessons - just copy the template and fill in the fields.

## Features Implementation

### Interactive Mode (Default)
- Start from first incomplete koan
- Display question and incomplete expression
- Prompt for answer
- Show immediate feedback (correct/incorrect)
- Option to see hints
- Automatically advance to next koan on success

### Validation Mode (`--validate`)
- Run through all koans
- Check all answers
- Display summary of pass/fail
- Useful for testing or reviewing

### Progress Tracking
- Save completion state to `~/.cronkoans_progress.json`
- Track which koans are completed
- Store timestamps
- Allow reset command to start over

### Hints System
- Three-level progressive hints
- Show next hint on request
- Don't penalize for using hints
- Reset hints when koan is completed

## CLI Commands

```
cronkoans                    # Start interactive mode from last position
cronkoans start              # Start from beginning (resume if in progress)
cronkoans reset              # Reset all progress
cronkoans list               # List all lessons and koans
cronkoans status             # Show progress statistics
cronkoans validate           # Validate all koans
cronkoans --help             # Show help
```

## Files to Create

1. **README.md** - Installation, usage instructions, learning path overview
2. **CONTRIBUTING.md** - Step-by-step guide for adding new lessons
3. **main.go** - CLI entry point
4. **go.mod** - Go module definition
5. **cmd/runner/runner.go** - Main runner logic
6. **internal/koan/koan.go** - Koan data structures
7. **internal/koan/validator.go** - Cron expression validator
8. **internal/koan/parser.go** - YAML parser for lessons
9. **internal/progress/tracker.go** - Progress tracking implementation
10. **internal/ui/display.go** - Terminal UI helpers (colors, formatting)
11. **lessons/01_basics.yaml** - First lesson
12. **lessons/02_wildcards.yaml** - Wildcards lesson
13. **lessons/03_ranges.yaml** - Ranges lesson
14. **lessons/04_steps.yaml** - Step values lesson
15. **lessons/05_lists.yaml** - Lists lesson
16. **lessons/06_special_strings.yaml** - Special strings lesson
17. **lessons/07_common_patterns.yaml** - Common patterns lesson
18. **lessons/08_advanced.yaml** - Advanced lesson
19. **lessons/template.yaml** - Template for new lessons
20. **.claude/commands/add-lesson.md** - Claude Code command for scaffolding new lessons

## Key Design Decisions

### YAML for Lessons
- Human-readable and easy to edit
- Well-structured for nested data
- No compilation needed to add lessons
- Can be edited in any text editor

### Home Directory for Progress
- Persistent across sessions
- Standard location (`~/.cronkoans_progress.json`)
- Won't clutter project directory

### Color-Coded Terminal Output
- Green for correct answers
- Red for incorrect answers
- Yellow for hints
- Blue for informational messages
- Enhances user experience

### Single Binary Distribution
- `go build` creates standalone executable
- No runtime dependencies
- Easy to install and distribute
- Cross-platform support

### Cron Validator
- Validates field ranges (minute 0-59, hour 0-23, etc.)
- Supports wildcards, ranges, steps, and lists
- Supports special strings (@daily, @hourly, etc.)
- Provides helpful error messages

## Success Criteria

- [ ] All 8 lesson files with 25-30 koans total
- [ ] Interactive mode works smoothly
- [ ] Validation mode runs all tests
- [ ] Progress tracking persists across sessions
- [ ] Hints system provides helpful guidance
- [ ] README clearly explains installation and usage
- [ ] CONTRIBUTING guide makes adding lessons trivial
- [ ] Claude Code command scaffolds new lessons
- [ ] Project builds with `go build`
- [ ] Clean, well-documented code

## Timeline

This is an educational project for learning cron syntax. The implementation should be:
- Well-structured and maintainable
- Easy to extend with new lessons
- Beginner-friendly with clear error messages
- Fun and engaging to use

## Next Steps

1. Initialize Go module
2. Create directory structure
3. Implement core data structures
4. Build the cron validator
5. Create the parser for YAML lessons
6. Implement progress tracking
7. Build the CLI runner
8. Write all lesson files
9. Create documentation
10. Test end-to-end
