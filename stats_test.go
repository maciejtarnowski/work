package work

import (
	"testing"
	"time"
)

func TestExpectedWorkTime(t *testing.T) {
	tests := []struct {
		from     time.Time
		to       time.Time
		expected time.Duration
	}{
		{
			from:     time.Date(2020, time.September, 7, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 7, 16, 0, 0, 0, time.UTC),
			expected: 8 * time.Hour,
		},
		{
			from:     time.Date(2020, time.September, 7, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 11, 16, 0, 0, 0, time.UTC),
			expected: 40 * time.Hour,
		},
		{
			from:     time.Date(2020, time.September, 7, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 13, 16, 0, 0, 0, time.UTC),
			expected: 40 * time.Hour,
		},
		{
			from:     time.Date(2020, time.September, 7, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 8, 8, 0, 0, 0, time.UTC),
			expected: 16 * time.Hour,
		},
		{
			from:     time.Date(2020, time.September, 7, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 9, 8, 0, 0, 0, time.UTC),
			expected: 24 * time.Hour,
		},
		{
			from:     time.Date(2020, time.September, 7, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 10, 8, 0, 0, 0, time.UTC),
			expected: 32 * time.Hour,
		},
		// edge cases
		{
			from:     time.Date(2020, time.August, 31, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 4, 8, 0, 0, 0, time.UTC),
			expected: 40 * time.Hour,
		},
		{
			from:     time.Date(2020, time.August, 31, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.August, 31, 8, 0, 0, 0, time.UTC),
			expected: 8 * time.Hour,
		},
		{
			from:     time.Date(2020, time.August, 31, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 1, 8, 0, 0, 0, time.UTC),
			expected: 16 * time.Hour,
		},
		{
			from:     time.Date(2020, time.August, 31, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 2, 8, 0, 0, 0, time.UTC),
			expected: 24 * time.Hour,
		},
		{
			from:     time.Date(2020, time.August, 31, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 3, 8, 0, 0, 0, time.UTC),
			expected: 32 * time.Hour,
		},
		{
			from:     time.Date(2020, time.August, 31, 8, 0, 0, 0, time.UTC),
			to:       time.Date(2020, time.September, 3, 6, 0, 0, 0, time.UTC),
			expected: 32 * time.Hour,
		},
	}

	for _, tt := range tests {
		d := calculateExpected(tt.from, tt.to)

		if d != tt.expected {
			t.Errorf("invalid duration for from=%s, to=%s - expected=%s, got=%s", tt.from, tt.to, tt.expected, d)
		}
	}
}
