// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package version provides version parsing and comparison functionality.
//
// DEPRECATED: This package contains confusing and hard-to-understand version comparison functions.
// Please migrate to github.com/ergoapi/util/version/semver for a cleaner, more intuitive API.
//
// Migration guide:
//   - LTv2(v1, v2) -> semver.IsLessThan(v1, v2)
//   - GTv2(v1, v2) -> semver.IsGreaterThan(v1, v2)
//   - NotGTv3(v1, v2) -> semver.IsLessThanOrEqual(v1, v2)
//   - NotLTv3(v1, v2) -> semver.IsGreaterThanOrEqual(v1, v2)
//   - IsLessOrEqualv3(v1, v2) -> semver.IsLessThanOrEqual(v1, v2)
//   - IsGreaterOrEqualv3(v1, v2) -> semver.IsGreaterThanOrEqual(v1, v2)
//   - Parse(v) -> semver.Parse(v)
//   - Next(now, major, minor, patch) -> use semver.Parse(now).IncrementMajor/Minor/Patch()
package version

import (
	"strings"

	"github.com/blang/semver/v4"
)

func parseVersions(v1, v2 string) (semver.Version, semver.Version, error) {
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	vv1, err1 := semver.Make(v1)
	if err1 != nil {
		return semver.Version{}, semver.Version{}, err1
	}

	vv2, err2 := semver.Make(v2)
	if err2 != nil {
		return semver.Version{}, semver.Version{}, err2
	}

	return vv1, vv2, nil
}

// Deprecated: Use github.com/ergoapi/util/version/semver.IsLessThan instead.
// LTv2 returns true if v1 is less than v2
func LTv2(v1, v2 string) bool {
	vv1, vv2, err := parseVersions(v1, v2)
	if err != nil {
		return false
	}
	return vv1.LT(vv2)
}

// Deprecated: Use github.com/ergoapi/util/version/semver.IsLessThanOrEqual instead.
// NotGTv3 returns true if v1 is less than or equal to v2
func NotGTv3(v1, v2 string) bool {
	vv1, vv2, err := parseVersions(v1, v2)
	if err != nil {
		return false
	}
	return vv1.LT(vv2) || vv1.EQ(vv2)
}

// Deprecated: Use github.com/ergoapi/util/version/semver.IsLessThanOrEqual instead.
// IsLessOrEqualv3 returns true if v1 is less than or equal to v2
func IsLessOrEqualv3(v1, v2 string) bool {
	vv1, vv2, err := parseVersions(v1, v2)
	if err != nil {
		return false
	}
	return vv1.LT(vv2) || vv1.EQ(vv2)
}

// Deprecated: Use github.com/ergoapi/util/version/semver.IsGreaterThan instead.
// GTv2 returns true if v1 is greater than v2
func GTv2(v1, v2 string) bool {
	vv1, vv2, err := parseVersions(v1, v2)
	if err != nil {
		return false
	}
	return vv1.GT(vv2)
}

// Deprecated: Use github.com/ergoapi/util/version/semver.IsGreaterThanOrEqual instead.
// IsGreaterOrEqualv3 returns true if v1 is greater than or equal to v2
func IsGreaterOrEqualv3(v1, v2 string) bool {
	vv1, vv2, err := parseVersions(v1, v2)
	if err != nil {
		return false
	}
	return vv1.GT(vv2) || vv1.EQ(vv2)
}

// Deprecated: Use github.com/ergoapi/util/version/semver.IsGreaterThanOrEqual instead.
// NotLTv3 returns true if v1 is greater than or equal to v2
func NotLTv3(v1, v2 string) bool {
	vv1, vv2, err := parseVersions(v1, v2)
	if err != nil {
		return false
	}
	return vv1.GT(vv2) || vv1.EQ(vv2)
}

// Deprecated: Use github.com/ergoapi/util/version/semver.Parse instead.
// Parse creates a Version from a version string
func Parse(v string) (semver.Version, error) {
	return semver.Make(v)
}

// Deprecated: Use github.com/ergoapi/util/version/semver.Parse(now).IncrementMajor/Minor/Patch() instead.
// Next returns the next version by incrementing major, minor, or patch
func Next(now string, major, minor, patch bool) string {
	hasPrefix := strings.HasPrefix(now, "v")
	vStr := strings.TrimPrefix(now, "v")
	v, err := semver.New(vStr)
	if err != nil {
		return now
	}
	if major {
		v.Major++
		v.Minor = 0
		v.Patch = 0
	}
	if minor {
		v.Minor++
		v.Patch = 0
	}
	if patch {
		v.Patch++
	}
	result := v.String()
	if hasPrefix {
		return "v" + result
	}
	return result
}
