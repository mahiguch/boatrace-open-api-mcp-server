package unit

import (
	"testing"

	"github.com/boatrace/open-api-mcp-server/internal/validation"
)

// TestValidateYYYY tests year validation
func TestValidateYYYY(t *testing.T) {
	tests := []struct {
		name    string
		year    string
		wantErr bool
	}{
		{
			name:    "valid year 2025",
			year:    "2025",
			wantErr: false,
		},
		{
			name:    "valid year 1900 (min)",
			year:    "1900",
			wantErr: false,
		},
		{
			name:    "valid year 2100 (max)",
			year:    "2100",
			wantErr: false,
		},
		{
			name:    "invalid year 1899 (too low)",
			year:    "1899",
			wantErr: true,
		},
		{
			name:    "invalid year 2101 (too high)",
			year:    "2101",
			wantErr: true,
		},
		{
			name:    "invalid format 25 (2 digits)",
			year:    "25",
			wantErr: true,
		},
		{
			name:    "invalid format 20251 (5 digits)",
			year:    "20251",
			wantErr: true,
		},
		{
			name:    "invalid format non-numeric",
			year:    "abcd",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateYYYY(tt.year)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateYYYY() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateYYYYMMDD tests date validation
func TestValidateYYYYMMDD(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{
			name:    "valid date 2025-12-22",
			date:    "20251222",
			wantErr: false,
		},
		{
			name:    "valid date 2025-01-01",
			date:    "20250101",
			wantErr: false,
		},
		{
			name:    "valid date 2025-12-31",
			date:    "20251231",
			wantErr: false,
		},
		{
			name:    "invalid date 2025-13-01 (month 13)",
			date:    "20251301",
			wantErr: true,
		},
		{
			name:    "invalid date 2025-12-32 (day 32)",
			date:    "20251232",
			wantErr: true,
		},
		{
			name:    "invalid date 2025-02-30 (Feb 30)",
			date:    "20250230",
			wantErr: true,
		},
		{
			name:    "invalid format 2025-12-22 (dashes)",
			date:    "2025-12-22",
			wantErr: true,
		},
		{
			name:    "invalid format 25122 (5 digits)",
			date:    "25122",
			wantErr: true,
		},
		{
			name:    "invalid format non-numeric",
			date:    "abcdefgh",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateYYYYMMDD(tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateYYYYMMDD() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateDateRange tests combined year and date validation
func TestValidateDateRange(t *testing.T) {
	tests := []struct {
		name    string
		year    string
		date    string
		wantErr bool
	}{
		{
			name:    "valid date range",
			year:    "2025",
			date:    "20251222",
			wantErr: false,
		},
		{
			name:    "mismatched year",
			year:    "2024",
			date:    "20251222",
			wantErr: true,
		},
		{
			name:    "invalid year",
			year:    "25",
			date:    "20251222",
			wantErr: true,
		},
		{
			name:    "invalid date",
			year:    "2025",
			date:    "20251232",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateDateRange(tt.year, tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateDateRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
