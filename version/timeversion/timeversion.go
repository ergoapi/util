// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package timeversion provides time-based version management functionality.
// It supports version formats like YYYY.MM.DDNN where:
//   - YYYY: 4-digit year
//   - MM: 1-2 digit month (1-12)
//   - DD: 2-digit day (01-31)
//   - NN: 2-digit daily sequence number (01-99)
//
// Examples:
//   - 2025.1.0101: First version on January 1, 2025
//   - 2025.1.0122: 22nd version on January 1, 2025
//   - 2025.12.3105: 5th version on December 31, 2025
package timeversion

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
)

// TimeVersion represents a time-based version
type TimeVersion struct {
	Year     int    // 4-digit year
	Month    int    // 1-12 month
	Day      int    // 1-31 day
	Sequence int    // 1-99 daily sequence
	raw      string // original input string
}

// ParseError represents an error that occurred during time version parsing
type ParseError struct {
	Input string
	Err   error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse time version '%s': %v", e.Input, e.Err)
}

// Parse creates a new TimeVersion from a version string
// Supports formats: YYYY.M.DDSS or YYYY.MM.DDSS
func Parse(versionStr string) (*TimeVersion, error) {
	if versionStr == "" {
		return nil, &ParseError{Input: versionStr, Err: errors.New("version string cannot be empty")}
	}

	// Remove 'v' prefix if present
	original := versionStr
	versionStr = strings.TrimPrefix(versionStr, "v")

	parts := strings.Split(versionStr, ".")
	if len(parts) != 3 {
		return nil, &ParseError{Input: original, Err: errors.New("version must have format YYYY.M.DDSS or YYYY.MM.DDSS")}
	}

	// Parse year
	year, err := strconv.Atoi(parts[0])
	if err != nil || year < 1000 || year > 9999 {
		return nil, &ParseError{Input: original, Err: errors.New("year must be a 4-digit number")}
	}

	// Parse month
	month, err := strconv.Atoi(parts[1])
	if err != nil || month < 1 || month > 12 {
		return nil, &ParseError{Input: original, Err: errors.New("month must be 1-12")}
	}

	// Parse day and sequence (DDSS format)
	daySeqStr := parts[2]
	if len(daySeqStr) != 4 {
		return nil, &ParseError{Input: original, Err: errors.New("day and sequence must be 4 digits (DDSS)")}
	}

	dayStr := daySeqStr[:2]
	seqStr := daySeqStr[2:]

	day, err := strconv.Atoi(dayStr)
	if err != nil || day < 1 || day > 31 {
		return nil, &ParseError{Input: original, Err: errors.New("day must be 01-31")}
	}

	sequence, err := strconv.Atoi(seqStr)
	if err != nil || sequence < 1 || sequence > 99 {
		return nil, &ParseError{Input: original, Err: errors.New("sequence must be 01-99")}
	}

	// Validate date
	testDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	if testDate.Year() != year || int(testDate.Month()) != month || testDate.Day() != day {
		return nil, &ParseError{Input: original, Err: errors.New("invalid date")}
	}

	return &TimeVersion{
		Year:     year,
		Month:    month,
		Day:      day,
		Sequence: sequence,
		raw:      original,
	}, nil
}

// MustParse creates a TimeVersion or panics if parsing fails
func MustParse(versionStr string) *TimeVersion {
	v, err := Parse(versionStr)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns the original string representation
func (tv *TimeVersion) String() string {
	return tv.raw
}

// Canonical returns the canonical format (without 'v' prefix)
func (tv *TimeVersion) Canonical() string {
	return fmt.Sprintf("%d.%d.%02d%02d", tv.Year, tv.Month, tv.Day, tv.Sequence)
}

// Format returns the formatted version with specified month format
func (tv *TimeVersion) Format(padMonth bool) string {
	if padMonth {
		return fmt.Sprintf("%d.%02d.%02d%02d", tv.Year, tv.Month, tv.Day, tv.Sequence)
	}
	return tv.Canonical()
}

// Date returns the date part as time.Time
func (tv *TimeVersion) Date() time.Time {
	return time.Date(tv.Year, time.Month(tv.Month), tv.Day, 0, 0, 0, 0, time.UTC)
}

// IsValid returns true if the version represents a valid date
func (tv *TimeVersion) IsValid() bool {
	testDate := time.Date(tv.Year, time.Month(tv.Month), tv.Day, 0, 0, 0, 0, time.UTC)
	return testDate.Year() == tv.Year && int(testDate.Month()) == tv.Month && testDate.Day() == tv.Day
}

// Comparison methods

// IsLessThan returns true if this version is less than the other version
func (tv *TimeVersion) IsLessThan(other *TimeVersion) bool {
	return tv.Compare(other) < 0
}

// IsLessThanOrEqual returns true if this version is less than or equal to the other version
func (tv *TimeVersion) IsLessThanOrEqual(other *TimeVersion) bool {
	return tv.Compare(other) <= 0
}

// IsEqual returns true if this version is equal to the other version
func (tv *TimeVersion) IsEqual(other *TimeVersion) bool {
	return tv.Compare(other) == 0
}

// IsGreaterThan returns true if this version is greater than the other version
func (tv *TimeVersion) IsGreaterThan(other *TimeVersion) bool {
	return tv.Compare(other) > 0
}

// IsGreaterThanOrEqual returns true if this version is greater than or equal to the other version
func (tv *TimeVersion) IsGreaterThanOrEqual(other *TimeVersion) bool {
	return tv.Compare(other) >= 0
}

// Compare returns -1 if tv < other, 0 if tv == other, 1 if tv > other
func (tv *TimeVersion) Compare(other *TimeVersion) int {
	// Compare year
	if tv.Year != other.Year {
		if tv.Year < other.Year {
			return -1
		}
		return 1
	}

	// Compare month
	if tv.Month != other.Month {
		if tv.Month < other.Month {
			return -1
		}
		return 1
	}

	// Compare day
	if tv.Day != other.Day {
		if tv.Day < other.Day {
			return -1
		}
		return 1
	}

	// Compare sequence
	if tv.Sequence != other.Sequence {
		if tv.Sequence < other.Sequence {
			return -1
		}
		return 1
	}

	return 0
}

// Version generation and increment methods

// Now creates a new TimeVersion for the current date with sequence 1
func Now() *TimeVersion {
	now := time.Now()
	return &TimeVersion{
		Year:     now.Year(),
		Month:    int(now.Month()),
		Day:      now.Day(),
		Sequence: 1,
		raw:      fmt.Sprintf("%d.%d.%02d01", now.Year(), int(now.Month()), now.Day()),
	}
}

// Today creates a new TimeVersion for today with specified sequence
func Today(sequence int) (*TimeVersion, error) {
	if sequence < 1 || sequence > 99 {
		return nil, errors.New("sequence must be 1-99")
	}

	now := time.Now()
	raw := fmt.Sprintf("%d.%d.%02d%02d", now.Year(), int(now.Month()), now.Day(), sequence)

	return &TimeVersion{
		Year:     now.Year(),
		Month:    int(now.Month()),
		Day:      now.Day(),
		Sequence: sequence,
		raw:      raw,
	}, nil
}

// FromDate creates a new TimeVersion from a date and sequence
func FromDate(date time.Time, sequence int) (*TimeVersion, error) {
	if sequence < 1 || sequence > 99 {
		return nil, errors.New("sequence must be 1-99")
	}

	raw := fmt.Sprintf("%d.%d.%02d%02d", date.Year(), int(date.Month()), date.Day(), sequence)

	return &TimeVersion{
		Year:     date.Year(),
		Month:    int(date.Month()),
		Day:      date.Day(),
		Sequence: sequence,
		raw:      raw,
	}, nil
}

// NextSequence returns a new version with the sequence incremented
func (tv *TimeVersion) NextSequence() (*TimeVersion, error) {
	if tv.Sequence >= 99 {
		return nil, errors.New("sequence cannot exceed 99")
	}

	newSeq := tv.Sequence + 1
	raw := fmt.Sprintf("%d.%d.%02d%02d", tv.Year, tv.Month, tv.Day, newSeq)

	return &TimeVersion{
		Year:     tv.Year,
		Month:    tv.Month,
		Day:      tv.Day,
		Sequence: newSeq,
		raw:      raw,
	}, nil
}

// NextDay returns a new version for the next day with sequence 1
func (tv *TimeVersion) NextDay() *TimeVersion {
	currentDate := tv.Date()
	nextDate := currentDate.AddDate(0, 0, 1)

	raw := fmt.Sprintf("%d.%d.%02d01", nextDate.Year(), int(nextDate.Month()), nextDate.Day())

	return &TimeVersion{
		Year:     nextDate.Year(),
		Month:    int(nextDate.Month()),
		Day:      nextDate.Day(),
		Sequence: 1,
		raw:      raw,
	}
}

// Package-level convenience functions

// Compare compares two time version strings
func Compare(v1, v2 string) (int, error) {
	tv1, err := Parse(v1)
	if err != nil {
		return 0, err
	}

	tv2, err := Parse(v2)
	if err != nil {
		return 0, err
	}

	return tv1.Compare(tv2), nil
}

// IsLessThan returns true if v1 < v2
func IsLessThan(v1, v2 string) (bool, error) {
	result, err := Compare(v1, v2)
	return result < 0, err
}

// IsLessThanOrEqual returns true if v1 <= v2
func IsLessThanOrEqual(v1, v2 string) (bool, error) {
	result, err := Compare(v1, v2)
	return result <= 0, err
}

// IsEqual returns true if v1 == v2
func IsEqual(v1, v2 string) (bool, error) {
	result, err := Compare(v1, v2)
	return result == 0, err
}

// IsGreaterThan returns true if v1 > v2
func IsGreaterThan(v1, v2 string) (bool, error) {
	result, err := Compare(v1, v2)
	return result > 0, err
}

// IsGreaterThanOrEqual returns true if v1 >= v2
func IsGreaterThanOrEqual(v1, v2 string) (bool, error) {
	result, err := Compare(v1, v2)
	return result >= 0, err
}

// Sort sorts a slice of time version strings in ascending order
func Sort(versions []string) error {
	// Parse all versions first to catch errors early
	parsed := make([]*TimeVersion, len(versions))
	for i, v := range versions {
		var err error
		parsed[i], err = Parse(v)
		if err != nil {
			return fmt.Errorf("invalid version at index %d: %w", i, err)
		}
	}

	// Sort indices based on version comparison
	indices := make([]int, len(versions))
	for i := range indices {
		indices[i] = i
	}

	sort.Slice(indices, func(i, j int) bool {
		return parsed[indices[i]].IsLessThan(parsed[indices[j]])
	})

	// Reorder original slice
	sorted := make([]string, len(versions))
	for i, idx := range indices {
		sorted[i] = versions[idx]
	}
	copy(versions, sorted)

	return nil
}

// Latest returns the latest version from a slice of version strings
func Latest(versions []string) (string, error) {
	if len(versions) == 0 {
		return "", errors.New("no versions provided")
	}

	latest := versions[0]
	for _, v := range versions[1:] {
		isGreater, err := IsGreaterThan(v, latest)
		if err != nil {
			return "", err
		}
		if isGreater {
			latest = v
		}
	}

	return latest, nil
}

// GetVersionsForDate returns all versions for a specific date, sorted by sequence
func GetVersionsForDate(versions []string, date time.Time) ([]string, error) {
	var result []string
	targetYear := date.Year()
	targetMonth := int(date.Month())
	targetDay := date.Day()

	for _, v := range versions {
		tv, err := Parse(v)
		if err != nil {
			return nil, err
		}

		if tv.Year == targetYear && tv.Month == targetMonth && tv.Day == targetDay {
			result = append(result, v)
		}
	}

	// Sort by sequence
	sort.Slice(result, func(i, j int) bool {
		tv1, _ := Parse(result[i])
		tv2, _ := Parse(result[j])
		return tv1.Sequence < tv2.Sequence
	})

	return result, nil
}

// GetNextVersionForToday returns the next available version for today
func GetNextVersionForToday(existingVersions []string) (string, error) {
	now := time.Now()
	todayVersions, err := GetVersionsForDate(existingVersions, now)
	if err != nil {
		return "", err
	}

	// If no versions exist for today, return sequence 1
	if len(todayVersions) == 0 {
		return fmt.Sprintf("%d.%d.%02d01", now.Year(), int(now.Month()), now.Day()), nil
	}

	// Find the highest sequence
	latest := todayVersions[len(todayVersions)-1]
	tv, err := Parse(latest)
	if err != nil {
		return "", err
	}

	nextSeq := tv.Sequence + 1
	if nextSeq > 99 {
		return "", errors.New("daily sequence limit reached (99)")
	}

	return fmt.Sprintf("%d.%d.%02d%02d", now.Year(), int(now.Month()), now.Day(), nextSeq), nil
}
