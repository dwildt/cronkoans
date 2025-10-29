package koan

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ValidateCronExpression validates if a cron expression is syntactically correct
func ValidateCronExpression(expr string) error {
	// Check for special strings first
	specialStrings := []string{
		"@yearly", "@annually", "@monthly", "@weekly",
		"@daily", "@midnight", "@hourly", "@reboot",
	}

	exprLower := strings.ToLower(strings.TrimSpace(expr))
	for _, special := range specialStrings {
		if exprLower == special {
			return nil // Valid special string
		}
	}

	// Split into fields
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return fmt.Errorf("cron expression must have exactly 5 fields (minute hour day month weekday), got %d", len(fields))
	}

	// Validate each field
	validators := []struct {
		name     string
		min      int
		max      int
		field    string
		position int
	}{
		{"minute", 0, 59, fields[0], 0},
		{"hour", 0, 23, fields[1], 1},
		{"day", 1, 31, fields[2], 2},
		{"month", 1, 12, fields[3], 3},
		{"weekday", 0, 7, fields[4], 4}, // 0 and 7 both represent Sunday
	}

	for _, v := range validators {
		if err := validateField(v.field, v.min, v.max, v.name); err != nil {
			return fmt.Errorf("field %d (%s): %w", v.position+1, v.name, err)
		}
	}

	return nil
}

// validateField validates a single cron field
func validateField(field string, min, max int, fieldName string) error {
	// Wildcard
	if field == "*" {
		return nil
	}

	// Step values: */n or */n
	if strings.HasPrefix(field, "*/") {
		stepStr := field[2:]
		step, err := strconv.Atoi(stepStr)
		if err != nil {
			return fmt.Errorf("invalid step value: %s", stepStr)
		}
		if step <= 0 {
			return fmt.Errorf("step value must be positive, got %d", step)
		}
		return nil
	}

	// Range with step: n-m/step
	if strings.Contains(field, "/") && strings.Contains(field, "-") {
		parts := strings.Split(field, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid step syntax: %s", field)
		}

		// Validate the range part
		if err := validateRange(parts[0], min, max, fieldName); err != nil {
			return err
		}

		// Validate the step part
		step, err := strconv.Atoi(parts[1])
		if err != nil {
			return fmt.Errorf("invalid step value: %s", parts[1])
		}
		if step <= 0 {
			return fmt.Errorf("step value must be positive, got %d", step)
		}
		return nil
	}

	// List: n,m,o
	if strings.Contains(field, ",") {
		values := strings.Split(field, ",")
		for _, v := range values {
			v = strings.TrimSpace(v)
			if strings.Contains(v, "-") {
				if err := validateRange(v, min, max, fieldName); err != nil {
					return err
				}
			} else {
				if err := validateSingleValue(v, min, max, fieldName); err != nil {
					return err
				}
			}
		}
		return nil
	}

	// Range: n-m
	if strings.Contains(field, "-") {
		return validateRange(field, min, max, fieldName)
	}

	// Single value
	return validateSingleValue(field, min, max, fieldName)
}

// validateRange validates a range like "1-5"
func validateRange(rangeStr string, min, max int, fieldName string) error {
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return fmt.Errorf("invalid range syntax: %s", rangeStr)
	}

	start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return fmt.Errorf("invalid range start: %s", parts[0])
	}

	end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return fmt.Errorf("invalid range end: %s", parts[1])
	}

	if start < min || start > max {
		return fmt.Errorf("range start %d out of bounds [%d-%d]", start, min, max)
	}

	if end < min || end > max {
		return fmt.Errorf("range end %d out of bounds [%d-%d]", end, min, max)
	}

	if start > end {
		return fmt.Errorf("range start %d cannot be greater than end %d", start, end)
	}

	return nil
}

// validateSingleValue validates a single numeric value
func validateSingleValue(valueStr string, min, max int, fieldName string) error {
	value, err := strconv.Atoi(strings.TrimSpace(valueStr))
	if err != nil {
		return fmt.Errorf("invalid value: %s", valueStr)
	}

	// Special case for weekday: 0 and 7 are both valid for Sunday
	if fieldName == "weekday" && value == 7 {
		return nil
	}

	if value < min || value > max {
		return fmt.Errorf("value %d out of bounds [%d-%d]", value, min, max)
	}

	return nil
}

// DescribeCronExpression provides a human-readable description of a cron expression
func DescribeCronExpression(expr string) string {
	expr = strings.TrimSpace(expr)

	// Handle special strings
	specialDescriptions := map[string]string{
		"@yearly":   "Once a year at midnight on January 1st (0 0 1 1 *)",
		"@annually": "Once a year at midnight on January 1st (0 0 1 1 *)",
		"@monthly":  "Once a month at midnight on the 1st (0 0 1 * *)",
		"@weekly":   "Once a week at midnight on Sunday (0 0 * * 0)",
		"@daily":    "Once a day at midnight (0 0 * * *)",
		"@midnight": "Once a day at midnight (0 0 * * *)",
		"@hourly":   "Once an hour at the start of the hour (0 * * * *)",
		"@reboot":   "Once at system startup",
	}

	if desc, ok := specialDescriptions[strings.ToLower(expr)]; ok {
		return desc
	}

	// Validate first
	if err := ValidateCronExpression(expr); err != nil {
		return fmt.Sprintf("Invalid cron expression: %v", err)
	}

	fields := strings.Fields(expr)
	if len(fields) != 5 {
		return "Invalid cron expression format"
	}

	parts := []string{}

	// Minute
	if fields[0] != "*" {
		parts = append(parts, "at minute "+describeField(fields[0]))
	}

	// Hour
	if fields[1] != "*" {
		parts = append(parts, "hour "+describeField(fields[1]))
	}

	// Day of month
	if fields[2] != "*" {
		parts = append(parts, "on day "+describeField(fields[2]))
	}

	// Month
	if fields[3] != "*" {
		parts = append(parts, "in month "+describeField(fields[3]))
	}

	// Weekday
	if fields[4] != "*" {
		parts = append(parts, "on weekday "+describeField(fields[4]))
	}

	if len(parts) == 0 {
		return "Every minute"
	}

	return strings.Join(parts, ", ")
}

// describeField provides a simple description of a field value
func describeField(field string) string {
	if field == "*" {
		return "every"
	}
	if strings.HasPrefix(field, "*/") {
		return "every " + field[2:]
	}
	return field
}

// IsValidCronAnswer checks if the user's answer creates a valid cron expression
func IsValidCronAnswer(incomplete, answer string) bool {
	complete := replaceBlank(incomplete, answer)
	return ValidateCronExpression(complete) == nil
}

// GetCronValidationRegex returns a regex pattern for basic cron validation
func GetCronValidationRegex() *regexp.Regexp {
	// Basic pattern for cron expression validation
	// This is a simplified version and actual validation should use ValidateCronExpression
	pattern := `^(@(yearly|annually|monthly|weekly|daily|midnight|hourly|reboot)|` +
		`((\*|[0-9]|[1-5][0-9]|\*/[0-9]+|[0-9]+-[0-9]+|[0-9]+(,[0-9]+)*)\s+){4}` +
		`(\*|[0-7]|\*/[0-9]+|[0-7]-[0-7]|[0-7](,[0-7]+)*))$`

	return regexp.MustCompile(pattern)
}
