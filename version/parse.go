package version

import (
	"strings"

	"github.com/blang/semver/v4"
)

// DEPRECATED: Use LTv2 instead
// LT returns true if v2 is less than v1
func LT(v1, v2 string) bool {
	vv1, _ := semver.Make(v1)
	vv2, _ := semver.Make(v2)
	return vv2.LT(vv1)
}

// LTv2 returns true if v1 is less than v2
func LTv2(v1, v2 string) bool {
	vv1, _ := semver.Make(v1)
	vv2, _ := semver.Make(v2)
	return vv1.LT(vv2)
}

// GTv2 returns true if v1 is greater than v2
func GTv2(v1, v2 string) bool {
	vv1, _ := semver.Make(v1)
	vv2, _ := semver.Make(v2)
	return vv1.GT(vv2)
}

func Parse(v string) (semver.Version, error) {
	return semver.Make(v)
}

func Next(now string, major, minor, patch bool) string {
	hasPrefix := strings.HasPrefix(now, "v")
	if hasPrefix {
		now = strings.TrimPrefix(now, "v")
	}
	v, err := semver.New(now)
	if err != nil {
		return now
	}
	if major {
		v.Major++
	}
	if minor {
		v.Minor++
	}
	if patch {
		v.Patch++
	}
	if hasPrefix {
		return "v" + v.String()
	}
	return v.String()
}
