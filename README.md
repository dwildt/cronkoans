# Cron Koans

Learn Crontab configuration through practice and wisdom.

## What are Koans?

Koans are a teaching method where you learn by doing. Each koan presents you with an incomplete cron expression, and your task is to fill in the missing piece. As you progress through the koans, you'll build a deep understanding of cron syntax and scheduling patterns.

## What is Cron?

Cron is a time-based job scheduler in Unix-like operating systems. It allows you to schedule commands or scripts to run automatically at specified times and dates. A cron expression is a string of 5 fields that defines when a task should run:

```
* * * * *
│ │ │ │ │
│ │ │ │ └─── Day of week (0-7, both 0 and 7 represent Sunday)
│ │ │ └───── Month (1-12)
│ │ └─────── Day of month (1-31)
│ └───────── Hour (0-23)
└─────────── Minute (0-59)
```

## Installation

### Prerequisites

- Go 1.21 or higher

### Build from Source

```bash
# Clone the repository
git clone https://github.com/dwildt/cronkoans.git
cd cronkoans

# Build the application
go build -o cronkoans

# Run it
./cronkoans
```

### Install Globally

```bash
# Build and install to $GOPATH/bin
go install

# Run from anywhere
cronkoans
```

## Usage

### Interactive Mode (Default)

Start learning from where you left off:

```bash
cronkoans
```

Or start from the beginning:

```bash
cronkoans start
```

### Commands

- `cronkoans` or `cronkoans start` - Start interactive learning mode
- `cronkoans list` - List all available lessons and your progress
- `cronkoans status` - Show your progress statistics
- `cronkoans validate` - Validate all lesson files
- `cronkoans reset` - Reset your progress and start over
- `cronkoans help` - Show help information
- `cronkoans --version` - Show version information

### Interactive Mode Commands

While working through koans, you can use these commands:

- Type your answer and press Enter to submit
- Type `hint` or `h` to get a hint
- Type `skip` to skip the current koan
- Type `quit` or `exit` to quit

## Learning Path

The koans are organized into 8 progressive lessons:

### 1. Basics (5 koans)
Understanding the five fields of a cron expression and their valid ranges.

### 2. Wildcards (4 koans)
Mastering the use of `*` to mean "every" value in any field.

### 3. Ranges (4 koans)
Learning to specify ranges of values using the dash operator (`-`).

### 4. Step Values (5 koans)
Using step values (`*/n`) to run tasks at regular intervals.

### 5. Lists (4 koans)
Specifying multiple specific values using comma-separated lists.

### 6. Special Strings (5 koans)
Using shortcuts like `@daily`, `@hourly`, `@weekly` for common schedules.

### 7. Common Patterns (5 koans)
Applying your knowledge to real-world scheduling scenarios.

### 8. Advanced (5 koans)
Combining multiple operators to create complex schedules.

**Total: 37 koans**

## Examples

Here are some examples of what you'll learn:

```bash
# Every 5 minutes
*/5 * * * *

# Every day at midnight
0 0 * * *

# Every weekday at 9 AM
0 9 * * 1-5

# Twice a day (9 AM and 6 PM)
0 9,18 * * *

# Every 15 minutes during business hours on weekdays
*/15 9-17 * * 1-5

# First day of every month
0 0 1 * *
```

## Progress Tracking

Your progress is automatically saved to `~/.cronkoans_progress.json`. This file tracks:

- Which koans you've completed
- How many attempts you made
- How many hints you used
- When you started and last updated

You can reset your progress at any time with `cronkoans reset`.

## Hints System

Each koan comes with 3 progressive hints:

1. **First hint**: General direction
2. **Second hint**: More specific guidance
3. **Third hint**: Almost gives away the answer

Don't worry about using hints - they're there to help you learn! After 2 incorrect attempts, the system will offer you a hint automatically.

## Validation Mode

Want to verify that all lesson files are correctly formatted? Run:

```bash
cronkoans validate
```

This checks that all koan answers produce valid cron expressions.

## Contributing

We welcome contributions! Whether it's:

- Adding new lessons
- Improving existing koans
- Fixing bugs
- Improving documentation

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to add new lessons.

## Resources

- [Crontab Guru](https://crontab.guru) - Interactive cron expression explainer
- [Cron Wikipedia](https://en.wikipedia.org/wiki/Cron) - Detailed history and documentation
- [Crontab Quick Reference](https://www.adminschoice.com/crontab-quick-reference)

## Project Structure

```
cronkoans/
├── main.go                    # CLI entry point
├── go.mod                     # Go module definition
├── README.md                  # This file
├── CONTRIBUTING.md            # Guide for contributors
├── cmd/
│   └── runner/
│       └── runner.go          # Main runner logic
├── internal/
│   ├── koan/
│   │   ├── koan.go           # Koan data structures
│   │   ├── validator.go      # Cron expression validator
│   │   ├── parser.go         # YAML lesson parser
│   │   └── utils.go          # Utility functions
│   ├── progress/
│   │   └── tracker.go        # Progress tracking
│   └── ui/
│       └── display.go        # Terminal UI
└── lessons/
    ├── 01_basics.yaml
    ├── 02_wildcards.yaml
    ├── 03_ranges.yaml
    ├── 04_steps.yaml
    ├── 05_lists.yaml
    ├── 06_special_strings.yaml
    ├── 07_common_patterns.yaml
    ├── 08_advanced.yaml
    └── template.yaml          # Template for new lessons
```

## License

MIT License - feel free to use this project for learning and teaching!

## Acknowledgments

Inspired by:
- [Ruby Koans](http://www.rubykoans.com/) - The original koans for learning Ruby
- [Go Koans](https://github.com/cdarwin/go-koans) - Koans for learning Go
- The Unix cron utility and its many implementations

## Support

If you encounter any issues or have suggestions:

1. Check existing issues on GitHub
2. Create a new issue with details about your problem
3. Include your Go version and operating system

Happy learning! May you achieve cron enlightenment.
