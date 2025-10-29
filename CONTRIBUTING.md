# Contributing to Cron Koans

Thank you for your interest in contributing to Cron Koans! This guide will help you add new lessons, improve existing ones, or fix bugs.

## How to Add New Lessons

Adding a new lesson is straightforward and doesn't require deep programming knowledge. Lessons are written in YAML format and follow a simple structure.

### Step 1: Copy the Template

Start by copying the template file:

```bash
cp lessons/template.yaml lessons/09_your_lesson_name.yaml
```

Replace `09` with the next available number and give it a descriptive name (use underscores, not spaces).

### Step 2: Edit the Lesson Metadata

Open your new file and update the title and description:

```yaml
title: "Your Topic - Brief Description"
description: "A longer explanation of what this lesson teaches"
```

### Step 3: Create Your Koans

Each lesson should have 3-5 koans. Each koan follows this structure:

```yaml
koans:
  - id: "unique_id_1"
    description: "What this koan teaches"
    question: "Task description in plain English"
    incomplete: "__ * * * *"
    answer: "*/5"
    hints:
      - "General hint"
      - "More specific hint"
      - "Almost the answer"
    explanation: "Detailed explanation of why this works"
```

### Step 4: Guidelines for Good Koans

#### ID Naming
- Use lowercase with underscores
- Start with your lesson name: `myfeature_1`, `myfeature_2`, etc.
- Must be unique across ALL lessons

#### Questions
- Use clear, real-world scenarios
- Be specific about what you want to schedule
- Examples:
  - ✅ "Run database backup at 2:30 AM every day"
  - ❌ "Run something sometimes"

#### Incomplete Expressions
- Use `__` (double underscore) for exactly ONE blank
- Place the blank where the student should fill in
- The blank can be:
  - A complete field: `__ * * * *`
  - Part of a field: `__/__ * * * *` (for ranges with steps)
  - Multiple fields: `__ __ * * *`

#### Answers
- Should create a valid cron expression when combined with `incomplete`
- Can be:
  - A single value: `5`
  - A range: `1-5`
  - A step value: `*/10`
  - A list: `1,3,5`
  - A special string: `@daily`

#### Hints
- Provide exactly 3 hints
- Progression:
  1. General concept or direction
  2. More specific guidance
  3. Very close to the answer (but still requires thinking)
- Don't just repeat the question

#### Explanations
- Explain WHY the answer works
- Describe what the complete expression does
- Mention any edge cases or gotchas
- Keep it educational, not just "this is correct"

### Step 5: Validate Your Lesson

Test your lesson file:

```bash
go run main.go validate
```

This checks that:
- YAML syntax is correct
- All required fields are present
- Koan IDs are unique
- Answers create valid cron expressions
- Each incomplete expression has a `__` placeholder

### Step 6: Test Interactively

Actually try your lesson:

```bash
go run main.go start
```

Work through your koans and verify:
- Questions are clear
- Hints are helpful
- Explanations are educational
- Difficulty progresses appropriately

## Example: Creating a "Time Zones" Lesson

Here's a complete example:

```yaml
title: "Time Zones - UTC Scheduling"
description: "Learn how to schedule tasks considering UTC time zones"

koans:
  - id: "timezone_1"
    description: "Understanding UTC midnight"
    question: "Schedule a task for midnight UTC"
    incomplete: "0 __ * * *"
    answer: "0"
    hints:
      - "UTC midnight is hour 0 in 24-hour format"
      - "The hour field comes second in cron expressions"
      - "Use 0 for the hour field"
    explanation: "In cron, times are typically in UTC. Midnight UTC is hour 0. If your local time is EST (UTC-5), this runs at 7 PM local time."

  - id: "timezone_2"
    description: "Converting local to UTC"
    question: "Run at 9 AM EST (2 PM UTC)"
    incomplete: "0 __ * * *"
    answer: "14"
    hints:
      - "EST is UTC-5 hours"
      - "9 AM EST + 5 hours = 2 PM UTC"
      - "2 PM in 24-hour format is 14"
    explanation: "When scheduling in UTC for EST times, add 5 hours (or 4 during DST). 9 AM EST = 14:00 UTC."
```

## Lesson Ideas

Here are some topics that would make great lessons:

- **Month Names vs Numbers**: Using month names in some cron implementations
- **Last Day of Month**: Techniques for scheduling on the last day
- **Every N Weeks**: Patterns for bi-weekly or monthly intervals
- **Maintenance Windows**: Combining multiple restrictions
- **Holiday Scheduling**: Scheduling around specific dates
- **Monitoring Patterns**: Common intervals for health checks
- **Backup Strategies**: Various backup scheduling patterns
- **Daylight Saving Time**: Understanding DST impacts
- **Cron Best Practices**: Anti-patterns and recommendations
- **System Cron vs User Cron**: Differences and use cases

## Code Contributions

### Prerequisites

- Go 1.21 or higher
- Basic understanding of Go (for code changes)
- Git for version control

### Setting Up Development Environment

```bash
# Clone the repository
git clone https://github.com/dwildt/cronkoans.git
cd cronkoans

# Install dependencies
go mod download

# Build the project
go build

# Run tests
go test ./...
```

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Add comments for exported functions
- Keep functions focused and small
- Write tests for new features

### Project Structure

```
internal/
├── koan/        # Koan data structures and validation
├── progress/    # Progress tracking and persistence
└── ui/          # Terminal UI and display logic

cmd/
└── runner/      # Main runner and command handling

lessons/         # YAML lesson files
```

### Making Changes

1. Create a new branch:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. Make your changes

3. Test your changes:
   ```bash
   go test ./...
   go build
   ./cronkoans validate
   ```

4. Commit with a clear message:
   ```bash
   git commit -m "Add feature: brief description"
   ```

5. Push and create a pull request

## Bug Reports

When reporting bugs, please include:

1. **Description**: What went wrong?
2. **Steps to Reproduce**: How can we see the bug?
3. **Expected Behavior**: What should have happened?
4. **Actual Behavior**: What actually happened?
5. **Environment**:
   - Go version (`go version`)
   - Operating system
   - Cron Koans version

Example:

```
**Description**
Progress doesn't save after completing a koan

**Steps to Reproduce**
1. Run `cronkoans start`
2. Complete the first koan
3. Exit the program
4. Run `cronkoans status`

**Expected Behavior**
Status should show 1 completed koan

**Actual Behavior**
Status shows 0 completed koans

**Environment**
- Go 1.21.0
- macOS 14.0
- Cron Koans v1.0.0
```

## Feature Requests

We welcome feature ideas! When suggesting a feature:

1. **Use Case**: Describe why this feature would be useful
2. **Proposed Solution**: How might it work?
3. **Alternatives**: Are there other ways to achieve this?
4. **Additional Context**: Any examples or mockups?

## Pull Request Process

1. **Fork** the repository
2. **Create** a feature branch
3. **Make** your changes
4. **Test** thoroughly
5. **Document** your changes
6. **Submit** a pull request with:
   - Clear description of changes
   - Why the changes are needed
   - Any breaking changes
   - Screenshots (for UI changes)

### Pull Request Checklist

- [ ] Code builds without errors
- [ ] Tests pass (`go test ./...`)
- [ ] New features have tests
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] No unnecessary files included

## Questions?

- Check existing issues on GitHub
- Open a new issue for discussion
- Reach out to maintainers

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on the code, not the person
- Help newcomers feel welcome
- Assume good intentions

## License

By contributing, you agree that your contributions will be licensed under the same MIT License that covers the project.

## Thank You!

Every contribution, whether it's a new lesson, a bug fix, or improved documentation, makes Cron Koans better for everyone. Thank you for being part of this project!
