package models

import (
	"testing"
)

// TestModelPackageInternal verifies models are accessible
func TestModelPackageInternal(t *testing.T) {
	// Test RaceProgram
	program := &RaceProgram{
		Date: "20251222",
		Races: []RaceInfo{
			{
				Number: 1,
				Venue:  "Test Venue",
				RaceID: "20251222-01",
			},
		},
	}

	if program.Date != "20251222" {
		t.Errorf("RaceProgram Date mismatch")
	}

	if len(program.Races) != 1 {
		t.Errorf("RaceProgram Races count mismatch")
	}

	// Test RaceResult
	result := &RaceResult{
		Date: "20251222",
		Races: []RaceOutcome{
			{
				Number: 1,
				RaceID: "20251222-01",
				Venue:  "Test Venue",
			},
		},
	}

	if result.Date != "20251222" {
		t.Errorf("RaceResult Date mismatch")
	}

	// Test RacePreview
	preview := &RacePreview{
		Date: "20251222",
		Races: []PreviewInfo{
			{
				Number: 1,
				RaceID: "20251222-01",
				Venue:  "Test Venue",
			},
		},
	}

	if preview.Date != "20251222" {
		t.Errorf("RacePreview Date mismatch")
	}
}
