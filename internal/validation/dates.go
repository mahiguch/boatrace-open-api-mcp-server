package validation

import (
	"fmt"
	"strconv"
	"time"
)

// ValidateYYYY validates a year string (YYYY format)
// Valid range: 1900-2100
func ValidateYYYY(year string) error {
	if len(year) != 4 {
		return fmt.Errorf("invalid year '%s': expected 4-digit year", year)
	}

	y, err := strconv.Atoi(year)
	if err != nil {
		return fmt.Errorf("invalid year '%s': expected numeric value", year)
	}

	if y < 1900 || y > 2100 {
		return fmt.Errorf("invalid year %d: expected year between 1900 and 2100", y)
	}

	return nil
}

// ValidateYYYYMMDD validates a date string (YYYYMMDD format)
// Validates format and that the date is a valid calendar date
func ValidateYYYYMMDD(date string) error {
	if len(date) != 8 {
		return fmt.Errorf("invalid date '%s': expected 8-digit date in YYYYMMDD format (e.g., 20251222)", date)
	}

	// Try to parse as a valid date
	_, err := time.Parse("20060102", date)
	if err != nil {
		return fmt.Errorf("invalid date '%s': expected valid calendar date in YYYYMMDD format (e.g., 20251222)", date)
	}

	return nil
}

// ValidateDateRange validates both year and date together
func ValidateDateRange(year, date string) error {
	if err := ValidateYYYY(year); err != nil {
		return err
	}

	if err := ValidateYYYYMMDD(date); err != nil {
		return err
	}

	// Verify year matches the year in the date
	dateYear := date[:4]
	if dateYear != year {
		return fmt.Errorf("date year %s does not match provided year %s", dateYear, year)
	}

	return nil
}
