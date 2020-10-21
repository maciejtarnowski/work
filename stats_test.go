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

func TestLogDateRangeWithOffset(t *testing.T) {
	tests := []struct {
		now          time.Time
		offset       int
		expectedFrom time.Time
		expectedTo   time.Time
	}{
		{
			now:          time.Date(2020, time.October, 21, 8, 0, 0, 0, time.UTC),
			offset:       0,
			expectedFrom: time.Date(2020, time.October, 19, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2020, time.October, 21, 23, 59, 59, 0, time.UTC),
		},
		{
			now:          time.Date(2020, time.October, 21, 8, 0, 0, 0, time.UTC),
			offset:       1,
			expectedFrom: time.Date(2020, time.October, 12, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2020, time.October, 16, 23, 59, 59, 0, time.UTC),
		},
		{
			now:          time.Date(2020, time.October, 21, 8, 0, 0, 0, time.UTC),
			offset:       2,
			expectedFrom: time.Date(2020, time.October, 5, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2020, time.October, 9, 23, 59, 59, 0, time.UTC),
		},
		{
			now:          time.Date(2020, time.October, 21, 8, 0, 0, 0, time.UTC),
			offset:       3,
			expectedFrom: time.Date(2020, time.September, 28, 0, 0, 0, 0, time.UTC),
			expectedTo:   time.Date(2020, time.October, 2, 23, 59, 59, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		from, to := getDateRangeForLog(tt.now, tt.offset)

		if !from.Equal(tt.expectedFrom) {
			t.Errorf("invalid `from` for now=%s, offset=%d - expected=%s, got=%s", tt.now, tt.offset, tt.expectedFrom, from)
		}
		if !to.Equal(tt.expectedTo) {
			t.Errorf("invalid `to` for now=%s, offset=%d - expected=%s, got=%s", tt.now, tt.offset, tt.expectedTo, to)
		}
	}
}
