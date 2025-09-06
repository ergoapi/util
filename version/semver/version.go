// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package semver provides clear and intuitive semantic version comparison functionality.
// This package is designed to replace the confusing version comparison functions
// in the parent package with a more user-friendly API.
package semver

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"

	"github.com/blang/semver/v4"
)

// Version represents a semantic version with comparison methods
type Version struct {
	version semver.Version
	raw     string
}

// ParseError represents an error that occurred during version parsing
type ParseError struct {
	Input string
	Err   error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("failed to parse version '%s': %v", e.Input, e.Err)
}

// Parse creates a new Version from a version string
// Supports versions with or without 'v' prefix (e.g., "1.2.3" or "v1.2.3")
func Parse(versionStr string) (*Version, error) {
	if versionStr == "" {
		return nil, &ParseError{Input: versionStr, Err: errors.New("version string cannot be empty")}
	}

	// Store original input for error reporting
	original := versionStr

	// Remove 'v' prefix if present
	versionStr = strings.TrimPrefix(versionStr, "v")

	parsed, err := semver.Make(versionStr)
	if err != nil {
		return nil, &ParseError{Input: original, Err: err}
	}

	return &Version{
		version: parsed,
		raw:     original,
	}, nil
}

// MustParse creates a Version or panics if parsing fails
// Only use this when you're certain the input is valid
func MustParse(versionStr string) *Version {
	v, err := Parse(versionStr)
	if err != nil {
		panic(err)
	}
	return v
}

// String returns the string representation of the version
func (v *Version) String() string {
	return v.raw
}

// Canonical returns the canonical semver string (without 'v' prefix)
func (v *Version) Canonical() string {
	return v.version.String()
}

// Major returns the major version number
func (v *Version) Major() uint64 {
	return v.version.Major
}

// Minor returns the minor version number
func (v *Version) Minor() uint64 {
	return v.version.Minor
}

// Patch returns the patch version number
func (v *Version) Patch() uint64 {
	return v.version.Patch
}

// Pre returns the pre-release version identifiers
func (v *Version) Pre() []semver.PRVersion {
	return v.version.Pre
}

// Build returns the build metadata identifiers
func (v *Version) Build() []string {
	return v.version.Build
}

// Comparison methods with clear, intuitive names

// IsLessThan returns true if this version is less than the other version
func (v *Version) IsLessThan(other *Version) bool {
	return v.version.LT(other.version)
}

// IsLessThanOrEqual returns true if this version is less than or equal to the other version
func (v *Version) IsLessThanOrEqual(other *Version) bool {
	return v.version.LT(other.version) || v.version.EQ(other.version)
}

// IsEqual returns true if this version is equal to the other version
func (v *Version) IsEqual(other *Version) bool {
	return v.version.EQ(other.version)
}

// IsGreaterThan returns true if this version is greater than the other version
func (v *Version) IsGreaterThan(other *Version) bool {
	return v.version.GT(other.version)
}

// IsGreaterThanOrEqual returns true if this version is greater than or equal to the other version
func (v *Version) IsGreaterThanOrEqual(other *Version) bool {
	return v.version.GT(other.version) || v.version.EQ(other.version)
}

// Compare returns -1 if v < other, 0 if v == other, 1 if v > other
func (v *Version) Compare(other *Version) int {
	return v.version.Compare(other.version)
}

// Version increment methods

// IncrementMajor returns a new version with the major version incremented
// Minor and patch versions are reset to 0, pre-release and build metadata are cleared
func (v *Version) IncrementMajor() *Version {
	newVer := v.version
	newVer.Major++
	newVer.Minor = 0
	newVer.Patch = 0
	newVer.Pre = nil
	newVer.Build = nil

	return &Version{
		version: newVer,
		raw:     newVer.String(),
	}
}

// IncrementMinor returns a new version with the minor version incremented
// Patch version is reset to 0, pre-release and build metadata are cleared
func (v *Version) IncrementMinor() *Version {
	newVer := v.version
	newVer.Minor++
	newVer.Patch = 0
	newVer.Pre = nil
	newVer.Build = nil

	return &Version{
		version: newVer,
		raw:     newVer.String(),
	}
}

// IncrementPatch returns a new version with the patch version incremented
// Pre-release and build metadata are cleared
func (v *Version) IncrementPatch() *Version {
	newVer := v.version
	newVer.Patch++
	newVer.Pre = nil
	newVer.Build = nil

	return &Version{
		version: newVer,
		raw:     newVer.String(),
	}
}

// Convenience functions for string comparisons

// IsLessThanString compares this version with a version string
func (v *Version) IsLessThanString(versionStr string) (bool, error) {
	other, err := Parse(versionStr)
	if err != nil {
		return false, err
	}
	return v.IsLessThan(other), nil
}

// IsGreaterThanString compares this version with a version string
func (v *Version) IsGreaterThanString(versionStr string) (bool, error) {
	other, err := Parse(versionStr)
	if err != nil {
		return false, err
	}
	return v.IsGreaterThan(other), nil
}

// IsEqualString compares this version with a version string
func (v *Version) IsEqualString(versionStr string) (bool, error) {
	other, err := Parse(versionStr)
	if err != nil {
		return false, err
	}
	return v.IsEqual(other), nil
}

// Package-level convenience functions for quick comparisons

// Compare compares two version strings and returns -1, 0, or 1
// Returns an error if either version string is invalid
func Compare(v1, v2 string) (int, error) {
	ver1, err := Parse(v1)
	if err != nil {
		return 0, err
	}

	ver2, err := Parse(v2)
	if err != nil {
		return 0, err
	}

	return ver1.Compare(ver2), nil
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

// Sort sorts a slice of version strings in ascending order
// Invalid version strings are moved to the end and an error is returned
func Sort(versions []string) error {
	// Parse all versions first to catch errors early
	parsed := make([]*Version, len(versions))
	for i, v := range versions {
		var err error
		parsed[i], err = Parse(v)
		if err != nil {
			return fmt.Errorf("invalid version at index %d: %w", i, err)
		}
	}

	// Sort using semver comparison
	for i := 0; i < len(parsed); i++ {
		for j := i + 1; j < len(parsed); j++ {
			if parsed[i].IsGreaterThan(parsed[j]) {
				parsed[i], parsed[j] = parsed[j], parsed[i]
				versions[i], versions[j] = versions[j], versions[i]
			}
		}
	}

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
