// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package timeversion

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    *TimeVersion
	}{
		{
			name:        "Basic format - single digit month",
			input:       "2025.1.0101",
			expectError: false,
			expected:    &TimeVersion{Year: 2025, Month: 1, Day: 1, Sequence: 1, raw: "2025.1.0101"},
		},
		{
			name:        "Basic format - double digit month",
			input:       "2025.12.3105",
			expectError: false,
			expected:    &TimeVersion{Year: 2025, Month: 12, Day: 31, Sequence: 5, raw: "2025.12.3105"},
		},
		{
			name:        "With v prefix",
			input:       "v2025.1.0122",
			expectError: false,
			expected:    &TimeVersion{Year: 2025, Month: 1, Day: 1, Sequence: 22, raw: "v2025.1.0122"},
		},
		{
			name:        "Maximum values",
			input:       "9999.12.3199",
			expectError: false,
			expected:    &TimeVersion{Year: 9999, Month: 12, Day: 31, Sequence: 99, raw: "9999.12.3199"},
		},
		{
			name:        "February leap year",
			input:       "2024.2.2901",
			expectError: false,
			expected:    &TimeVersion{Year: 2024, Month: 2, Day: 29, Sequence: 1, raw: "2024.2.2901"},
		},
		// Error cases
		{
			name:        "Empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "Invalid format - too few parts",
			input:       "2025.1",
			expectError: true,
		},
		{
			name:        "Invalid format - too many parts",
			input:       "2025.1.01.01",
			expectError: true,
		},
		{
			name:        "Invalid year - too short",
			input:       "25.1.0101",
			expectError: true,
		},
		{
			name:        "Invalid year - non-numeric",
			input:       "abc.1.0101",
			expectError: true,
		},
		{
			name:        "Invalid month - zero",
			input:       "2025.0.0101",
			expectError: true,
		},
		{
			name:        "Invalid month - too large",
			input:       "2025.13.0101",
			expectError: true,
		},
		{
			name:        "Invalid day-sequence format - too short",
			input:       "2025.1.101",
			expectError: true,
		},
		{
			name:        "Invalid day-sequence format - too long",
			input:       "2025.1.01011",
			expectError: true,
		},
		{
			name:        "Invalid day - zero",
			input:       "2025.1.0001",
			expectError: true,
		},
		{
			name:        "Invalid day - too large",
			input:       "2025.1.3201",
			expectError: true,
		},
		{
			name:        "Invalid sequence - zero",
			input:       "2025.1.0100",
			expectError: true,
		},
		{
			name:        "Invalid date - February 30",
			input:       "2025.2.3001",
			expectError: true,
		},
		{
			name:        "Invalid date - February 29 on non-leap year",
			input:       "2025.2.2901",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Parse(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)

				// Check error type
				var parseErr *ParseError
				assert.ErrorAs(t, err, &parseErr)
				assert.Equal(t, tt.input, parseErr.Input)
			} else {
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, tt.expected.Year, result.Year)
				assert.Equal(t, tt.expected.Month, result.Month)
				assert.Equal(t, tt.expected.Day, result.Day)
				assert.Equal(t, tt.expected.Sequence, result.Sequence)
				assert.Equal(t, tt.expected.raw, result.raw)
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	// Valid version should not panic
	tv := MustParse("2025.1.0101")
	assert.Equal(t, "2025.1.0101", tv.String())

	// Invalid version should panic
	assert.Panics(t, func() {
		MustParse("invalid")
	})
}

func TestTimeVersionProperties(t *testing.T) {
	tv := MustParse("v2025.12.3122")

	assert.Equal(t, "v2025.12.3122", tv.String())
	assert.Equal(t, "2025.12.3122", tv.Canonical())
	assert.Equal(t, "2025.12.3122", tv.Format(false))
	assert.Equal(t, "2025.12.3122", tv.Format(true))

	assert.Equal(t, 2025, tv.Year)
	assert.Equal(t, 12, tv.Month)
	assert.Equal(t, 31, tv.Day)
	assert.Equal(t, 22, tv.Sequence)

	expectedDate := time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)
	assert.Equal(t, expectedDate, tv.Date())
	assert.True(t, tv.IsValid())
}

func TestTimeVersionFormat(t *testing.T) {
	tv := MustParse("2025.1.0101")

	assert.Equal(t, "2025.1.0101", tv.Format(false)) // no month padding
	assert.Equal(t, "2025.01.0101", tv.Format(true)) // with month padding
}

func TestVersionEquivalence(t *testing.T) {
	// Test that single-digit and double-digit months are equivalent
	tests := []struct {
		name     string
		v1, v2   string
		expected bool // should be equal
	}{
		{
			name:     "Single vs double digit month - January",
			v1:       "2025.1.0101",
			v2:       "2025.01.0101",
			expected: true,
		},
		{
			name:     "Single vs double digit month - September",
			v1:       "2025.9.1501",
			v2:       "2025.09.1501",
			expected: true,
		},
		{
			name:     "Both single digit months",
			v1:       "2025.1.0101",
			v2:       "2025.1.0101",
			expected: true,
		},
		{
			name:     "Both double digit months",
			v1:       "2025.01.0101",
			v2:       "2025.01.0101",
			expected: true,
		},
		{
			name:     "With v prefix equivalence",
			v1:       "v2025.1.0101",
			v2:       "2025.01.0101",
			expected: true,
		},
		{
			name:     "Different months should not be equal",
			v1:       "2025.1.0101",
			v2:       "2025.02.0101",
			expected: false,
		},
		{
			name:     "Different sequences should not be equal",
			v1:       "2025.01.0101",
			v2:       "2025.01.0102",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v1, err1 := Parse(tt.v1)
			v2, err2 := Parse(tt.v2)

			require.NoError(t, err1)
			require.NoError(t, err2)

			// Test TimeVersion.IsEqual method
			result := v1.IsEqual(v2)
			assert.Equal(t, tt.expected, result, "TimeVersion.IsEqual() failed for %s vs %s", tt.v1, tt.v2)

			// Test TimeVersion.Compare method
			if tt.expected {
				assert.Equal(t, 0, v1.Compare(v2), "TimeVersion.Compare() should return 0 for equal versions %s vs %s", tt.v1, tt.v2)
				assert.Equal(t, 0, v2.Compare(v1), "TimeVersion.Compare() should be symmetric for %s vs %s", tt.v1, tt.v2)
			} else {
				assert.NotEqual(t, 0, v1.Compare(v2), "TimeVersion.Compare() should not return 0 for different versions %s vs %s", tt.v1, tt.v2)
			}

			// Test package-level IsEqual function
			pkgResult, err := IsEqual(tt.v1, tt.v2)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, pkgResult, "IsEqual() function failed for %s vs %s", tt.v1, tt.v2)

			// Test that canonical representation is consistent
			if tt.expected {
				// Both should have the same numeric values
				assert.Equal(t, v1.Year, v2.Year, "Years should be equal")
				assert.Equal(t, v1.Month, v2.Month, "Months should be equal")
				assert.Equal(t, v1.Day, v2.Day, "Days should be equal")
				assert.Equal(t, v1.Sequence, v2.Sequence, "Sequences should be equal")

				// Canonical format should be the same (using single-digit month)
				assert.Equal(t, v1.Canonical(), v2.Canonical(), "Canonical formats should be identical")
			}
		})
	}
}

func TestTimeVersionComparisons(t *testing.T) {
	tv1 := MustParse("2025.1.0101")
	tv2 := MustParse("2025.1.0102")
	tv3 := MustParse("2025.1.0101")
	tv4 := MustParse("2025.1.0201")
	tv5 := MustParse("2025.2.0101")
	tv6 := MustParse("2026.1.0101")

	// IsLessThan
	assert.True(t, tv1.IsLessThan(tv2))  // same day, different sequence
	assert.True(t, tv1.IsLessThan(tv4))  // different day
	assert.True(t, tv1.IsLessThan(tv5))  // different month
	assert.True(t, tv1.IsLessThan(tv6))  // different year
	assert.False(t, tv2.IsLessThan(tv1)) // reverse
	assert.False(t, tv1.IsLessThan(tv3)) // equal

	// IsLessThanOrEqual
	assert.True(t, tv1.IsLessThanOrEqual(tv2))
	assert.True(t, tv1.IsLessThanOrEqual(tv3))
	assert.False(t, tv2.IsLessThanOrEqual(tv1))

	// IsEqual
	assert.True(t, tv1.IsEqual(tv3))
	assert.False(t, tv1.IsEqual(tv2))

	// IsGreaterThan
	assert.True(t, tv2.IsGreaterThan(tv1))
	assert.False(t, tv1.IsGreaterThan(tv2))
	assert.False(t, tv1.IsGreaterThan(tv3))

	// IsGreaterThanOrEqual
	assert.True(t, tv2.IsGreaterThanOrEqual(tv1))
	assert.True(t, tv1.IsGreaterThanOrEqual(tv3))
	assert.False(t, tv1.IsGreaterThanOrEqual(tv2))

	// Compare
	assert.Equal(t, -1, tv1.Compare(tv2))
	assert.Equal(t, 0, tv1.Compare(tv3))
	assert.Equal(t, 1, tv2.Compare(tv1))
}

func TestVersionGeneration(t *testing.T) {
	// Test Now()
	now := Now()
	currentTime := time.Now()
	assert.Equal(t, currentTime.Year(), now.Year)
	assert.Equal(t, int(currentTime.Month()), now.Month)
	assert.Equal(t, currentTime.Day(), now.Day)
	assert.Equal(t, 1, now.Sequence)

	// Test Today()
	today, err := Today(5)
	require.NoError(t, err)
	assert.Equal(t, currentTime.Year(), today.Year)
	assert.Equal(t, int(currentTime.Month()), today.Month)
	assert.Equal(t, currentTime.Day(), today.Day)
	assert.Equal(t, 5, today.Sequence)

	// Test Today() with invalid sequence
	_, err = Today(0)
	assert.Error(t, err)
	_, err = Today(100)
	assert.Error(t, err)

	// Test FromDate()
	testDate := time.Date(2025, 6, 15, 12, 30, 45, 0, time.UTC)
	fromDate, err := FromDate(testDate, 10)
	require.NoError(t, err)
	assert.Equal(t, 2025, fromDate.Year)
	assert.Equal(t, 6, fromDate.Month)
	assert.Equal(t, 15, fromDate.Day)
	assert.Equal(t, 10, fromDate.Sequence)

	// Test FromDate() with invalid sequence
	_, err = FromDate(testDate, 0)
	assert.Error(t, err)
}

func TestVersionIncrement(t *testing.T) {
	base := MustParse("2025.1.0105")

	// Test NextSequence()
	next, err := base.NextSequence()
	require.NoError(t, err)
	assert.Equal(t, 2025, next.Year)
	assert.Equal(t, 1, next.Month)
	assert.Equal(t, 1, next.Day)
	assert.Equal(t, 6, next.Sequence)

	// Test NextSequence() at limit
	maxSeq := MustParse("2025.1.0199")
	_, err = maxSeq.NextSequence()
	assert.Error(t, err)

	// Test NextDay()
	nextDay := base.NextDay()
	assert.Equal(t, 2025, nextDay.Year)
	assert.Equal(t, 1, nextDay.Month)
	assert.Equal(t, 2, nextDay.Day)
	assert.Equal(t, 1, nextDay.Sequence)

	// Test NextDay() crossing month boundary
	endOfMonth := MustParse("2025.1.3101")
	nextMonth := endOfMonth.NextDay()
	assert.Equal(t, 2025, nextMonth.Year)
	assert.Equal(t, 2, nextMonth.Month)
	assert.Equal(t, 1, nextMonth.Day)
	assert.Equal(t, 1, nextMonth.Sequence)
}

func TestPackageLevelComparisons(t *testing.T) {
	// Compare
	result, err := Compare("2025.1.0101", "2025.1.0102")
	assert.NoError(t, err)
	assert.Equal(t, -1, result)

	result, err = Compare("2025.1.0101", "2025.1.0101")
	assert.NoError(t, err)
	assert.Equal(t, 0, result)

	result, err = Compare("2025.1.0102", "2025.1.0101")
	assert.NoError(t, err)
	assert.Equal(t, 1, result)

	// Error case
	_, err = Compare("invalid", "2025.1.0101")
	assert.Error(t, err)

	// IsLessThan
	resultBool, err := IsLessThan("2025.1.0101", "2025.1.0102")
	assert.NoError(t, err)
	assert.True(t, resultBool)

	// IsLessThanOrEqual
	resultBool, err = IsLessThanOrEqual("2025.1.0101", "2025.1.0101")
	assert.NoError(t, err)
	assert.True(t, resultBool)

	resultBool, err = IsLessThanOrEqual("2025.1.0101", "2025.1.0102")
	assert.NoError(t, err)
	assert.True(t, resultBool)

	resultBool, err = IsLessThanOrEqual("2025.1.0102", "2025.1.0101")
	assert.NoError(t, err)
	assert.False(t, resultBool)

	// IsEqual
	resultBool, err = IsEqual("2025.1.0101", "2025.1.0101")
	assert.NoError(t, err)
	assert.True(t, resultBool)

	resultBool, err = IsEqual("v2025.1.0101", "2025.1.0101")
	assert.NoError(t, err)
	assert.True(t, resultBool)

	// IsGreaterThan
	resultBool, err = IsGreaterThan("2025.1.0102", "2025.1.0101")
	assert.NoError(t, err)
	assert.True(t, resultBool)

	// IsGreaterThanOrEqual
	resultBool, err = IsGreaterThanOrEqual("2025.1.0101", "2025.1.0101")
	assert.NoError(t, err)
	assert.True(t, resultBool)
}

func TestSort(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
		hasError bool
	}{
		{
			name:     "Basic sorting",
			input:    []string{"2025.1.0102", "2025.1.0101", "2025.1.0103"},
			expected: []string{"2025.1.0101", "2025.1.0102", "2025.1.0103"},
			hasError: false,
		},
		{
			name:     "Cross-day sorting",
			input:    []string{"2025.1.0201", "2025.1.0101", "2025.1.0102"},
			expected: []string{"2025.1.0101", "2025.1.0102", "2025.1.0201"},
			hasError: false,
		},
		{
			name:     "Cross-month sorting",
			input:    []string{"2025.2.0101", "2025.1.3199", "2025.1.0101"},
			expected: []string{"2025.1.0101", "2025.1.3199", "2025.2.0101"},
			hasError: false,
		},
		{
			name:     "Cross-year sorting",
			input:    []string{"2026.1.0101", "2025.12.3199", "2025.1.0101"},
			expected: []string{"2025.1.0101", "2025.12.3199", "2026.1.0101"},
			hasError: false,
		},
		{
			name:     "With v prefix",
			input:    []string{"v2025.1.0102", "2025.1.0101", "v2025.1.0103"},
			expected: []string{"2025.1.0101", "v2025.1.0102", "v2025.1.0103"},
			hasError: false,
		},
		{
			name:     "Empty slice",
			input:    []string{},
			expected: []string{},
			hasError: false,
		},
		{
			name:     "Invalid version",
			input:    []string{"2025.1.0101", "invalid", "2025.1.0102"},
			expected: nil,
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := make([]string, len(tt.input))
			copy(input, tt.input)

			err := Sort(input)

			if tt.hasError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, input)
			}
		})
	}
}

func TestLatest(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected string
		hasError bool
	}{
		{
			name:     "Basic latest",
			input:    []string{"2025.1.0101", "2025.1.0103", "2025.1.0102"},
			expected: "2025.1.0103",
			hasError: false,
		},
		{
			name:     "Cross-day latest",
			input:    []string{"2025.1.0101", "2025.1.0201", "2025.1.0102"},
			expected: "2025.1.0201",
			hasError: false,
		},
		{
			name:     "Cross-year latest",
			input:    []string{"2025.1.0101", "2026.1.0101", "2025.12.3199"},
			expected: "2026.1.0101",
			hasError: false,
		},
		{
			name:     "With v prefix",
			input:    []string{"2025.1.0101", "v2025.1.0102", "2025.1.0103"},
			expected: "2025.1.0103",
			hasError: false,
		},
		{
			name:     "Single version",
			input:    []string{"2025.1.0101"},
			expected: "2025.1.0101",
			hasError: false,
		},
		{
			name:     "Empty slice",
			input:    []string{},
			expected: "",
			hasError: true,
		},
		{
			name:     "Invalid version",
			input:    []string{"2025.1.0101", "invalid"},
			expected: "",
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Latest(tt.input)

			if tt.hasError {
				assert.Error(t, err)
				assert.Equal(t, tt.expected, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestGetVersionsForDate(t *testing.T) {
	versions := []string{
		"2025.1.0101",
		"2025.1.0103",
		"2025.1.0102",
		"2025.1.0201",
		"2025.2.0101",
	}

	targetDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	result, err := GetVersionsForDate(versions, targetDate)

	require.NoError(t, err)
	expected := []string{"2025.1.0101", "2025.1.0102", "2025.1.0103"}
	assert.Equal(t, expected, result)

	// Test with no versions for date
	emptyDate := time.Date(2025, 3, 1, 0, 0, 0, 0, time.UTC)
	result, err = GetVersionsForDate(versions, emptyDate)

	require.NoError(t, err)
	assert.Empty(t, result)
}

func TestGetNextVersionForToday(t *testing.T) {
	now := time.Now()

	// Test with no existing versions
	result, err := GetNextVersionForToday([]string{})
	require.NoError(t, err)
	expected := fmt.Sprintf("%d.%d.%02d01", now.Year(), int(now.Month()), now.Day())
	assert.Equal(t, expected, result)

	// Test with existing versions for today
	existingToday := fmt.Sprintf("%d.%d.%02d05", now.Year(), int(now.Month()), now.Day())
	existingYesterday := fmt.Sprintf("%d.%d.%02d01", now.Year(), int(now.Month()), now.Day()-1)

	existing := []string{existingToday, existingYesterday}
	result, err = GetNextVersionForToday(existing)
	require.NoError(t, err)
	expectedNext := fmt.Sprintf("%d.%d.%02d06", now.Year(), int(now.Month()), now.Day())
	assert.Equal(t, expectedNext, result)
}

func TestParseError(t *testing.T) {
	_, err := Parse("invalid.version")
	require.Error(t, err)

	var parseErr *ParseError
	assert.ErrorAs(t, err, &parseErr)
	assert.Equal(t, "invalid.version", parseErr.Input)
	assert.Contains(t, parseErr.Error(), "failed to parse time version")
	assert.Contains(t, parseErr.Error(), "invalid.version")
}

func TestIsValid(t *testing.T) {
	// Valid dates
	assert.True(t, MustParse("2025.1.0101").IsValid())
	assert.True(t, MustParse("2024.2.2901").IsValid()) // leap year

	// Note: Parse() already validates dates, so IsValid() should always return true
	// for successfully parsed versions. This method is mainly for defensive programming.
}

// Benchmark tests
func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("2025.1.0101")
	}
}

func BenchmarkCompare(b *testing.B) {
	tv1 := MustParse("2025.1.0101")
	tv2 := MustParse("2025.1.0102")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tv1.IsLessThan(tv2)
	}
}

func BenchmarkStringCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsLessThan("2025.1.0101", "2025.1.0102")
	}
}

// Example tests that serve as documentation
func ExampleParse() {
	tv, err := Parse("2025.1.0122")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Version: %s\n", tv.String())
	fmt.Printf("Date: %s, Sequence: %d\n", tv.Date().Format("2006-01-02"), tv.Sequence)

	// Output:
	// Version: 2025.1.0122
	// Date: 2025-01-01, Sequence: 22
}

func ExampleTimeVersion_IsLessThan() {
	tv1 := MustParse("2025.1.0101")
	tv2 := MustParse("2025.1.0102")

	fmt.Printf("%s < %s: %t\n", tv1, tv2, tv1.IsLessThan(tv2))
	fmt.Printf("%s < %s: %t\n", tv2, tv1, tv2.IsLessThan(tv1))

	// Output:
	// 2025.1.0101 < 2025.1.0102: true
	// 2025.1.0102 < 2025.1.0101: false
}

func ExampleSort() {
	versions := []string{"2025.1.0103", "2025.1.0101", "2025.1.0202", "v2025.1.0102"}

	err := Sort(versions)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Sorted: %s\n", strings.Join(versions, ", "))

	// Output:
	// Sorted: 2025.1.0101, v2025.1.0102, 2025.1.0103, 2025.1.0202
}

func ExampleLatest() {
	versions := []string{"2025.1.0101", "2025.1.0103", "2025.1.0102", "v2025.1.0105"}

	latest, err := Latest(versions)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Latest version: %s\n", latest)

	// Output:
	// Latest version: v2025.1.0105
}

func ExampleNow() {
	// This will use the current date
	tv := Now()
	fmt.Printf("Today's first version: %s\n", tv.Canonical())
	fmt.Printf("Date: %s, Sequence: %d\n", tv.Date().Format("2006-01-02"), tv.Sequence)
}

func ExampleGetNextVersionForToday() {
	// Simulate existing versions
	existing := []string{"2025.1.0101", "2025.1.0103", "2025.1.0102"}

	nextVersion, err := GetNextVersionForToday(existing)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Next version for today: %s\n", nextVersion)
}
