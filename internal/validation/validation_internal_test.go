package validation

import (
	"testing"
)

// TestValidationPackageInternal tests internal validation functions
func TestValidationPackageInternal(t *testing.T) {
	// This test verifies the validation package is accessible
	err := ValidateYYYY("2025")
	if err != nil {
		t.Errorf("ValidateYYYY(2025) should not error")
	}

	err = ValidateYYYYMMDD("20251222")
	if err != nil {
		t.Errorf("ValidateYYYYMMDD(20251222) should not error")
	}

	err = ValidateDateRange("2025", "20251222")
	if err != nil {
		t.Errorf("ValidateDateRange(2025, 20251222) should not error")
	}
}
