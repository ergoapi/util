// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package semver

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expectedStr string
	}{
		{"Basic version", "1.2.3", false, "1.2.3"},
		{"Version with v prefix", "v1.2.3", false, "v1.2.3"},
		{"Version with pre-release", "1.2.3-alpha.1", false, "1.2.3-alpha.1"},
		{"Version with build", "1.2.3+build.1", false, "1.2.3+build.1"},
		{"Version with both", "1.2.3-alpha.1+build.1", false, "1.2.3-alpha.1+build.1"},
		{"Complex version", "v2.0.0-rc.1+20200101", false, "v2.0.0-rc.1+20200101"},
		{"Empty string", "", true, ""},
		{"Invalid format", "1.2", true, ""},
		{"Invalid characters", "1.2.3a", true, ""},
		{"Non-numeric", "a.b.c", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := Parse(tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)

				// Check error type
				var parseErr *ParseError
				assert.ErrorAs(t, err, &parseErr)
				assert.Equal(t, tt.input, parseErr.Input)
			} else {
				require.NoError(t, err)
				require.NotNil(t, v)
				assert.Equal(t, tt.expectedStr, v.String())
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	// Valid version should not panic
	v := MustParse("1.2.3")
	assert.Equal(t, "1.2.3", v.String())

	// Invalid version should panic
	assert.Panics(t, func() {
		MustParse("invalid")
	})
}

func TestVersionProperties(t *testing.T) {
	v := MustParse("v2.3.4-alpha.1+build.123")

	assert.Equal(t, "v2.3.4-alpha.1+build.123", v.String())
	assert.Equal(t, "2.3.4-alpha.1+build.123", v.Canonical())
	assert.Equal(t, uint64(2), v.Major())
	assert.Equal(t, uint64(3), v.Minor())
	assert.Equal(t, uint64(4), v.Patch())
	assert.NotEmpty(t, v.Pre())
	assert.NotEmpty(t, v.Build())
}

func TestVersionComparisons(t *testing.T) {
	v1 := MustParse("1.0.0")
	v2 := MustParse("1.0.1")
	v3 := MustParse("1.0.0")
	v4 := MustParse("2.0.0")

	// IsLessThan
	assert.True(t, v1.IsLessThan(v2))
	assert.False(t, v2.IsLessThan(v1))
	assert.False(t, v1.IsLessThan(v3))

	// IsLessThanOrEqual
	assert.True(t, v1.IsLessThanOrEqual(v2))
	assert.True(t, v1.IsLessThanOrEqual(v3))
	assert.False(t, v2.IsLessThanOrEqual(v1))

	// IsEqual
	assert.True(t, v1.IsEqual(v3))
	assert.False(t, v1.IsEqual(v2))

	// IsGreaterThan
	assert.True(t, v2.IsGreaterThan(v1))
	assert.False(t, v1.IsGreaterThan(v2))
	assert.False(t, v1.IsGreaterThan(v3))

	// IsGreaterThanOrEqual
	assert.True(t, v2.IsGreaterThanOrEqual(v1))
	assert.True(t, v1.IsGreaterThanOrEqual(v3))
	assert.False(t, v1.IsGreaterThanOrEqual(v2))

	// Compare
	assert.Equal(t, -1, v1.Compare(v2))
	assert.Equal(t, 0, v1.Compare(v3))
	assert.Equal(t, 1, v2.Compare(v1))
	assert.Equal(t, -1, v1.Compare(v4))
}

func TestVersionIncrements(t *testing.T) {
	base := MustParse("1.2.3-alpha.1+build.123")

	// IncrementMajor
	major := base.IncrementMajor()
	assert.Equal(t, uint64(2), major.Major())
	assert.Equal(t, uint64(0), major.Minor())
	assert.Equal(t, uint64(0), major.Patch())
	assert.Empty(t, major.Pre())
	assert.Empty(t, major.Build())

	// IncrementMinor
	minor := base.IncrementMinor()
	assert.Equal(t, uint64(1), minor.Major())
	assert.Equal(t, uint64(3), minor.Minor())
	assert.Equal(t, uint64(0), minor.Patch())
	assert.Empty(t, minor.Pre())
	assert.Empty(t, minor.Build())

	// IncrementPatch
	patch := base.IncrementPatch()
	assert.Equal(t, uint64(1), patch.Major())
	assert.Equal(t, uint64(2), patch.Minor())
	assert.Equal(t, uint64(4), patch.Patch())
	assert.Empty(t, patch.Pre())
	assert.Empty(t, patch.Build())
}

func TestVersionStringComparisons(t *testing.T) {
	v := MustParse("1.5.0")

	// IsLessThanString
	result, err := v.IsLessThanString("1.5.1")
	assert.NoError(t, err)
	assert.True(t, result)

	result, err = v.IsLessThanString("1.4.0")
	assert.NoError(t, err)
	assert.False(t, result)

	// IsGreaterThanString
	result, err = v.IsGreaterThanString("1.4.0")
	assert.NoError(t, err)
	assert.True(t, result)

	result, err = v.IsGreaterThanString("1.5.1")
	assert.NoError(t, err)
	assert.False(t, result)

	// IsEqualString
	result, err = v.IsEqualString("1.5.0")
	assert.NoError(t, err)
	assert.True(t, result)

	result, err = v.IsEqualString("v1.5.0")
	assert.NoError(t, err)
	assert.True(t, result)

	// Error cases
	_, err = v.IsLessThanString("invalid")
	assert.Error(t, err)
}

func TestPackageLevelComparisons(t *testing.T) {
	// Compare
	result, err := Compare("1.0.0", "1.0.1")
	assert.NoError(t, err)
	assert.Equal(t, -1, result)

	result, err = Compare("1.0.0", "1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, 0, result)

	result, err = Compare("1.0.1", "1.0.0")
	assert.NoError(t, err)
	assert.Equal(t, 1, result)

	// Error case
	_, err = Compare("invalid", "1.0.0")
	assert.Error(t, err)

	// IsLessThan
	result_bool, err := IsLessThan("1.0.0", "1.0.1")
	assert.NoError(t, err)
	assert.True(t, result_bool)

	// IsLessThanOrEqual
	result_bool, err = IsLessThanOrEqual("1.0.0", "1.0.0")
	assert.NoError(t, err)
	assert.True(t, result_bool)

	result_bool, err = IsLessThanOrEqual("1.0.0", "1.0.1")
	assert.NoError(t, err)
	assert.True(t, result_bool)

	result_bool, err = IsLessThanOrEqual("1.0.1", "1.0.0")
	assert.NoError(t, err)
	assert.False(t, result_bool)

	// IsEqual
	result_bool, err = IsEqual("1.0.0", "1.0.0")
	assert.NoError(t, err)
	assert.True(t, result_bool)

	result_bool, err = IsEqual("v1.0.0", "1.0.0")
	assert.NoError(t, err)
	assert.True(t, result_bool)

	// IsGreaterThan
	result_bool, err = IsGreaterThan("1.0.1", "1.0.0")
	assert.NoError(t, err)
	assert.True(t, result_bool)

	// IsGreaterThanOrEqual
	result_bool, err = IsGreaterThanOrEqual("1.0.0", "1.0.0")
	assert.NoError(t, err)
	assert.True(t, result_bool)

	result_bool, err = IsGreaterThanOrEqual("1.0.1", "1.0.0")
	assert.NoError(t, err)
	assert.True(t, result_bool)

	result_bool, err = IsGreaterThanOrEqual("1.0.0", "1.0.1")
	assert.NoError(t, err)
	assert.False(t, result_bool)
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
			input:    []string{"2.0.0", "1.0.0", "1.5.0"},
			expected: []string{"1.0.0", "1.5.0", "2.0.0"},
			hasError: false,
		},
		{
			name:     "With v prefix",
			input:    []string{"v2.0.0", "v1.0.0", "v1.5.0"},
			expected: []string{"v1.0.0", "v1.5.0", "v2.0.0"},
			hasError: false,
		},
		{
			name:     "Mixed prefix",
			input:    []string{"2.0.0", "v1.0.0", "1.5.0"},
			expected: []string{"v1.0.0", "1.5.0", "2.0.0"},
			hasError: false,
		},
		{
			name:     "With pre-release",
			input:    []string{"2.0.0", "2.0.0-alpha", "1.0.0"},
			expected: []string{"1.0.0", "2.0.0-alpha", "2.0.0"},
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
			input:    []string{"1.0.0", "invalid", "2.0.0"},
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
			input:    []string{"1.0.0", "2.0.0", "1.5.0"},
			expected: "2.0.0",
			hasError: false,
		},
		{
			name:     "With v prefix",
			input:    []string{"v1.0.0", "v2.0.0", "v1.5.0"},
			expected: "v2.0.0",
			hasError: false,
		},
		{
			name:     "Mixed prefix",
			input:    []string{"1.0.0", "v2.0.0", "1.5.0"},
			expected: "v2.0.0",
			hasError: false,
		},
		{
			name:     "Single version",
			input:    []string{"1.0.0"},
			expected: "1.0.0",
			hasError: false,
		},
		{
			name:     "With pre-release",
			input:    []string{"2.0.0-alpha", "2.0.0", "1.0.0"},
			expected: "2.0.0",
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
			input:    []string{"1.0.0", "invalid"},
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

func TestComplexVersions(t *testing.T) {
	// Test specific semver precedence rules
	tests := []struct {
		v1, v2   string
		expected int // -1: v1 < v2, 0: v1 == v2, 1: v1 > v2
	}{
		// Basic comparison
		{"1.0.0", "2.0.0", -1},
		{"1.0.0", "1.1.0", -1},
		{"1.0.0", "1.0.1", -1},

		// Pre-release comparisons (pre-release < normal version)
		{"1.0.0-alpha", "1.0.0", -1},
		{"1.0.0-alpha.1", "1.0.0", -1},
		{"1.0.0-rc.1", "1.0.0", -1},

		// Pre-release precedence
		{"1.0.0-alpha", "1.0.0-alpha.1", -1},
		{"1.0.0-alpha.1", "1.0.0-alpha.beta", -1},
		{"1.0.0-alpha.beta", "1.0.0-beta", -1},
		{"1.0.0-beta", "1.0.0-beta.2", -1},
		{"1.0.0-beta.2", "1.0.0-beta.11", -1},
		{"1.0.0-beta.11", "1.0.0-rc.1", -1},

		// Build metadata is ignored for precedence
		{"1.0.0", "1.0.0+build.1", 0},
		{"1.0.0+build.1", "1.0.0+build.2", 0},
		{"1.0.0-alpha.1+build.1", "1.0.0-alpha.1", 0},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s vs %s", tt.v1, tt.v2), func(t *testing.T) {
			v1, err1 := Parse(tt.v1)
			v2, err2 := Parse(tt.v2)

			require.NoError(t, err1)
			require.NoError(t, err2)

			result := v1.Compare(v2)
			assert.Equal(t, tt.expected, result, "%s vs %s should be %d", tt.v1, tt.v2, tt.expected)

			// Also test individual methods
			if tt.expected < 0 {
				assert.True(t, v1.IsLessThan(v2))
				assert.False(t, v1.IsGreaterThan(v2))
				assert.False(t, v1.IsEqual(v2))
			} else if tt.expected == 0 {
				assert.False(t, v1.IsLessThan(v2))
				assert.False(t, v1.IsGreaterThan(v2))
				assert.True(t, v1.IsEqual(v2))
			} else {
				assert.False(t, v1.IsLessThan(v2))
				assert.True(t, v1.IsGreaterThan(v2))
				assert.False(t, v1.IsEqual(v2))
			}
		})
	}
}

func TestParseError(t *testing.T) {
	_, err := Parse("invalid.version")
	require.Error(t, err)

	var parseErr *ParseError
	assert.ErrorAs(t, err, &parseErr)
	assert.Equal(t, "invalid.version", parseErr.Input)
	assert.Contains(t, parseErr.Error(), "failed to parse version")
	assert.Contains(t, parseErr.Error(), "invalid.version")
}

// Benchmark tests
func BenchmarkParse(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Parse("1.2.3-alpha.1+build.123")
	}
}

func BenchmarkCompare(b *testing.B) {
	v1 := MustParse("1.2.3")
	v2 := MustParse("1.2.4")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v1.IsLessThan(v2)
	}
}

func BenchmarkStringCompare(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsLessThan("1.2.3", "1.2.4")
	}
}

// Example tests that serve as documentation
func ExampleParse() {
	v, err := Parse("v1.2.3-alpha.1+build.123")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Version: %s\n", v.String())
	fmt.Printf("Major: %d, Minor: %d, Patch: %d\n", v.Major(), v.Minor(), v.Patch())

	// Output:
	// Version: v1.2.3-alpha.1+build.123
	// Major: 1, Minor: 2, Patch: 3
}

func ExampleVersion_IsLessThan() {
	v1 := MustParse("1.0.0")
	v2 := MustParse("1.0.1")

	fmt.Printf("%s < %s: %t\n", v1, v2, v1.IsLessThan(v2))
	fmt.Printf("%s < %s: %t\n", v2, v1, v2.IsLessThan(v1))

	// Output:
	// 1.0.0 < 1.0.1: true
	// 1.0.1 < 1.0.0: false
}

func ExampleSort() {
	versions := []string{"2.0.0", "1.0.0", "1.5.0", "v1.2.0"}

	err := Sort(versions)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Sorted: %s\n", strings.Join(versions, ", "))

	// Output:
	// Sorted: 1.0.0, v1.2.0, 1.5.0, 2.0.0
}

func ExampleLatest() {
	versions := []string{"1.0.0", "2.0.0", "1.5.0", "v1.2.0"}

	latest, err := Latest(versions)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Latest version: %s\n", latest)

	// Output:
	// Latest version: 2.0.0
}
