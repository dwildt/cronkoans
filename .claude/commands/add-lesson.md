---
description: Create a new Cron Koans lesson file with the proper structure
---

You are helping create a new lesson for the Cron Koans project. Follow these steps:

1. **Ask for lesson details:**
   - Lesson topic and title
   - What concepts the lesson will teach
   - How many koans (3-5 recommended)

2. **Determine the lesson number:**
   - Look in the `lessons/` directory
   - Find the highest numbered lesson file (e.g., `08_advanced.yaml`)
   - Use the next number for the new lesson

3. **Create the lesson file:**
   - Filename format: `lessons/NN_topic_name.yaml` (e.g., `09_timezones.yaml`)
   - Use lowercase with underscores
   - Keep the name concise but descriptive

4. **Structure the lesson:**
   - Use the template from `lessons/template.yaml` as a guide
   - Each koan must have:
     - Unique ID (format: `topicname_N`)
     - Description
     - Question (real-world scenario)
     - Incomplete expression (with exactly one `__` placeholder)
     - Answer (that creates a valid cron expression)
     - 3 progressive hints
     - Educational explanation

5. **Ensure quality:**
   - Questions use real-world scenarios
   - Hints progress from general to specific
   - Explanations teach WHY, not just WHAT
   - All koan IDs are unique
   - Answers create valid cron expressions

6. **Validate the lesson:**
   - After creating the file, run: `go run main.go validate`
   - Fix any validation errors
   - Test interactively if possible

7. **Example of a well-structured koan:**
   ```yaml
   - id: "monitoring_1"
     description: "Health check frequency"
     question: "Check server health every 2 minutes"
     incomplete: "__/__ * * * *"
     answer: "*/2"
     hints:
       - "This requires a step value in the minute field"
       - "Use the */n notation where n is the interval"
       - "For every 2 minutes, use */2"
     explanation: "*/2 in the minute field means 'every 2 minutes'. The step value divides the 60 minutes into intervals of 2, running at minutes 0, 2, 4, 6, ..., 58."
   ```

8. **Topic ideas** (if the user needs inspiration):
   - Time zones and UTC
   - Last day of month patterns
   - Maintenance windows
   - Monitoring schedules
   - Backup strategies
   - Seasonal scheduling
   - Cron limitations and gotchas
   - Best practices

Remember: The goal is to create educational content that helps users learn cron through practice. Each koan should teach one specific concept and build on previous knowledge.
